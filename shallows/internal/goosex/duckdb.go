package goosex

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"sync"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

type DuckdbStore struct{}

func (t DuckdbStore) Tablename() string {
	return goose.DefaultTablename
}

func (t DuckdbStore) CreateVersionTable(ctx context.Context, db database.DBTxConn) error {
	q := fmt.Sprintf(`CREATE SEQUENCE %s_seq;
		CREATE TABLE %s (
		id integer PRIMARY KEY DEFAULT NEXTVAL('%s_seq'),
		version_id bigint NOT NULL,
		is_applied boolean NOT NULL,
		tstamp timestamp NOT NULL DEFAULT now()
	)`, t.Tablename(), t.Tablename(), t.Tablename())
	log.Println("store: create version table")
	_, err := db.ExecContext(ctx, q)
	return err
}

func (t DuckdbStore) Insert(ctx context.Context, db database.DBTxConn, req database.InsertRequest) error {
	q := fmt.Sprintf(`INSERT INTO %s (version_id, is_applied) VALUES ($1, $2)`, t.Tablename())
	row := db.QueryRowContext(ctx, q, req.Version, true)
	return row.Err()
}

func (t DuckdbStore) Delete(ctx context.Context, db database.DBTxConn, version int64) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE version_id=$1`, t.Tablename())
	row := db.QueryRowContext(ctx, q, version)
	return row.Err()
}

func (t DuckdbStore) GetMigration(ctx context.Context, db database.DBTxConn, version int64) (*database.GetMigrationResult, error) {
	q := fmt.Sprintf(`SELECT tstamp, is_applied FROM %s WHERE version_id=$1 ORDER BY tstamp DESC LIMIT 1`, t.Tablename())
	log.Println("store: get migration", version)
	var timestamp time.Time
	var isApplied bool
	err := db.QueryRowContext(ctx, q, version).Scan(&timestamp, &isApplied)
	if err != nil {
		return nil, err
	}

	return &database.GetMigrationResult{
		IsApplied: isApplied,
		Timestamp: timestamp,
	}, nil
}

func (t DuckdbStore) GetLatestVersion(ctx context.Context, db database.DBTxConn) (id int64, err error) {
	q := fmt.Sprintf(`SELECT version_id from %s ORDER BY id DESC LIMIT 1`, t.Tablename())
	log.Println("store: get latest version")
	err = db.QueryRowContext(ctx, q).Scan(&id)
	return id, err
}

func (t DuckdbStore) ListMigrations(ctx context.Context, db database.DBTxConn) ([]*database.ListMigrationsResult, error) {
	q := fmt.Sprintf(`SELECT version_id, is_applied from %s ORDER BY id DESC`, t.Tablename())
	log.Println("store: list migrations")
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []*database.ListMigrationsResult
	for rows.Next() {
		var version int64
		var isApplied bool
		if err := rows.Scan(&version, &isApplied); err != nil {
			return nil, err
		}
		migrations = append(migrations, &database.ListMigrationsResult{
			Version:   version,
			IsApplied: isApplied,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return migrations, nil
}

var m sync.Mutex

func InitializeDatabase(ctx context.Context, db *sql.DB, migrations fs.FS) (err error) {
	m.Lock()
	defer m.Unlock()

	if version, err := sqlx.String(ctx, db, "SELECT library_version FROM pragma_version()"); err != nil {
		return err
	} else {
		log.Println("detected duckdb", version)
	}

	log.Println("loading extensions", extensionSQL)
	if _, err := db.ExecContext(ctx, extensionSQL); err != nil {
		return errorsx.Wrap(err, "failed to load inet extension")
	}

	log.Println("enabling hnsw persistence")
	if _, err := db.ExecContext(ctx, "SET GLOBAL hnsw_enable_experimental_persistence = true;"); err != nil {
		return errorsx.Wrap(err, "failed to enable hnsw persistence")
	}
	log.Println("extensions ready")

	// debugging.
	// if _, err := db.ExecContext(ctx, "CALL enable_logging('QueryLog', storage_path = '/tmp/retrovibed.duckdb.log', storage_buffer_size = 0);"); err != nil {
	// 	return errorsx.Wrap(err, "failed to enable query logging to stderr")
	// }

	log.Println("checkpoint 1")
	if _, err := db.ExecContext(ctx, "FORCE CHECKPOINT;"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}
	log.Println("checkpoint 1 done")
	defer log.Println("DERP DERP")
	mprov, err := goose.NewProvider(
		"",
		db,
		migrations,
		goose.WithStore(DuckdbStore{}),
		goose.WithAllowOutofOrder(true),
		goose.WithVerbose(true),
	)
	if err != nil {
		return errorsx.Wrap(err, "unable to build migration provider")
	}

	log.Println("provider created, retrieving status")
	if statuses, err := mprov.Status(ctx); err != nil {
		log.Println("unable to retrieve migration status", err)
	} else {
		for _, s := range statuses {
			if s.State == goose.StatePending {
				log.Println("pending migration", s.Source.Path)
			}
		}
	}

	if _, err := mprov.Up(ctx); err != nil {
		return errorsx.Wrap(err, "unable to run migrations")
	}

	if _, err := db.ExecContext(ctx, "FORCE CHECKPOINT;"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	return nil
}
