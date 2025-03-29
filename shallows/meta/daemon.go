package meta

import (
	"context"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/squirrelx"
	"github.com/retrovibed/retrovibed/internal/stringsx"
)

type DaemonOption func(*Daemon)

func DaemonOptionMaybeID(v *Daemon) {
	v.ID = stringsx.FirstNonBlank(v.ID, md5x.String(v.Hostname))
}

func DaemonOptionEnsureDescription(v *Daemon) {
	v.Description = stringsx.FirstNonBlank(v.Description, v.Hostname, v.ID)
}

func DaemonOptionTestDefaults(v *Daemon) {
	v.ID = ""
	v.Hostname = errorsx.Must(os.Hostname())
}

func DaemonSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) DaemonScanner {
	return NewDaemonScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func DaemonSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(DaemonScannerStaticColumns)...).From("meta_daemons")
}

func DaemonFromHost() Daemon {
	return Daemon{
		CreatedAt:   time.Now(),
		Description: "",
		Hostname:    stringsx.FirstNonBlank(errorsx.Zero(os.Hostname()), "localhost"),
		ID:          uuid.Max.String(),
		UpdatedAt:   time.Now(),
	}
}
