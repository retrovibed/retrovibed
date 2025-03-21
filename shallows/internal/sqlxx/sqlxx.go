package sqlxx

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
)

type scanner[T any] interface {
	Scan(i *T) error
	Next() bool
	Close() error
	Err() error
}

// ScanEach record into the specified type automatically closes the scanner when done.
func ScanEach[T any](s scanner[T], do func(*T) error) (err error) {
	defer s.Close()
	for s.Next() {
		var (
			p T
		)

		if err = s.Scan(&p); err != nil {
			return errorsx.WithStack(err)
		}

		if err = do(&p); err != nil {
			return errorsx.WithStack(err)
		}
	}

	return errorsx.WithStack(s.Err())
}

// ScanInto a slice, automatically closes the scanner once done.
func ScanInto[T any](s scanner[T], dst *[]T) (err error) {
	defer s.Close()
	for s.Next() {
		var (
			p T
		)

		if err = s.Scan(&p); err != nil {
			return errorsx.WithStack(err)
		}

		*dst = append(*dst, p)
	}

	return errorsx.WithStack(s.Err())
}

// row interface or scanning a single row.
type row interface {
	Scan(dest ...interface{}) error
}

// NewValueScanner ...
func NewValueScanner[T any](r row) ValueScanner[T] {
	return ValueScanner[T]{row: r}
}

// ValueScanner helper for scanning single values in a type safe manner.
type ValueScanner[T any] struct {
	row
}

// Scan a result
func (t ValueScanner[T]) Scan(v *T) error {
	return t.row.Scan(v)
}

// NewValueScanner ...
func NewValuesScanner[T any](r *sql.Rows) ValuesScanner[T] {
	return ValuesScanner[T]{Rows: r}
}

// ValuesScanner helper for scanning single values in a type safe manner.
type ValuesScanner[T any] struct {
	*sql.Rows
}

// Scan a result
func (t ValuesScanner[T]) Scan(v *T) error {
	return t.Rows.Scan(v)
}

type rowsscanner[T any, X scanner[T]] func(*sql.Rows, error) X
type queryadvance[T any] func(last *T) squirrel.SelectBuilder

func NewBatched[T any, X scanner[T]](s func(*sql.Rows, error) X, adv queryadvance[T]) Batcher[T, X] {
	return Batcher[T, X]{
		limit:      100,
		rowscanner: s,
		qadv:       adv,
	}
}

type Batcher[T any, X scanner[T]] struct {
	limit       uint64
	qadv        queryadvance[T]
	rowscanner  rowsscanner[T, X]
	completesig bool
}

func (t Batcher[T, X]) Limit(l uint64) Batcher[T, X] {
	t.limit = l
	return t
}

// Next returns an empty array once all records have been processed as determined by the array set being < then the limit.
func (t *Batcher[T, X]) Next(ctx context.Context, q sqlx.Queryer, query squirrel.SelectBuilder) (_ squirrel.SelectBuilder, b []T, err error) {
	if t.completesig {
		t.completesig = false // reset batcher for reuse.
		return query, b, nil
	}

	var (
		s scanner[T] = t.rowscanner(query.Limit(t.limit).RunWith(q).QueryContext(ctx))
	)
	defer s.Close()

	if err = ScanInto(s, &b); err != nil {
		return query, b, err
	}

	if l := len(b); l > 0 {
		query = t.qadv(&b[l-1])
	}

	if uint64(len(b)) < t.limit {
		t.completesig = true
	}

	return query, b, errorsx.Compact(s.Err(), s.Close())
}

func MustCount(ctx context.Context, q sqlx.Queryer, query string) (c int64) {
	if err := Count(ctx, q, query).Scan(&c); err != nil {
		panic(err)
	}

	return c
}

func Count(ctx context.Context, q sqlx.Queryer, query string) ValueScanner[int64] {
	return NewValueScanner[int64](q.QueryRowContext(ctx, query))
}

func Exec(ctx context.Context, q sqlx.Queryer, query string, args ...any) error {
	_, err := q.ExecContext(ctx, query, args...)
	return err
}
