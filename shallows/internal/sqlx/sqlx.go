package sqlx

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/retrovibed/retrovibed/shallows/internal/ducktype"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

// Queryer interface for executing queries.
type Queryer interface {
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
	Exec(string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type Transactioner interface {
	Queryer
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

// Row interface or scanning a single row.
type Row interface {
	Scan(dest ...any) error
}

// Debug creates a DebuggingQueryer
func Debug(q Queryer) DebuggingQueryer {
	return DebuggingQueryer{
		Delegate: q,
	}
}

// DebuggingQueryer queryer that prints out the queries being executed.
type DebuggingQueryer struct {
	Delegate Queryer
}

// Query execute a query
func (t DebuggingQueryer) Query(q string, args ...any) (*sql.Rows, error) {
	return t.QueryContext(context.Background(), q, args...)
}

// QueryRow executes a query that returns a single row.
func (t DebuggingQueryer) QueryRow(q string, args ...any) *sql.Row {
	return t.QueryRowContext(context.Background(), q, args...)
}

// Exec executes a statement.
func (t DebuggingQueryer) Exec(q string, args ...any) (sql.Result, error) {
	return t.ExecContext(context.Background(), q, args...)
}

// QueryContext ...
func (t DebuggingQueryer) QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error) {
	log.Printf("%s:\n%#v\n", q, args)
	return t.Delegate.QueryContext(ctx, q, args...)
}

// QueryRowContext ...
func (t DebuggingQueryer) QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row {
	log.Printf("%s:\n%#v\n", q, args)
	return t.Delegate.QueryRowContext(ctx, q, args...)
}

// ExecContext ...
func (t DebuggingQueryer) ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error) {
	log.Printf("%s:\n%#v\n", q, args)
	return t.Delegate.ExecContext(ctx, q, args...)
}

func Count(ctx context.Context, q Queryer, query string, args ...any) (count int, err error) {
	err = NewIntRowScanner(q.QueryRowContext(ctx, query, args...)).Scan(&count)
	return count, err
}

func String(ctx context.Context, q Queryer, query string, args ...any) (s string, err error) {
	err = NewValueRowScanner[string](q.QueryRowContext(ctx, query, args...)).Scan(&s)
	return s, err
}

func Value[T any](ctx context.Context, q Queryer, query string, args ...any) (s T, err error) {
	err = NewValueRowScanner[T](q.QueryRowContext(ctx, query, args...)).Scan(&s)
	return s, err
}

func Timestamp(ctx context.Context, q Queryer, query string, args ...any) (s time.Time, err error) {
	var v ducktype.NullTime
	err = NewValueRowScanner[ducktype.NullTime](q.QueryRowContext(ctx, query, args...)).Scan(&v)

	switch v.InfinityModifier {
	case ducktype.Infinity:
		return timex.Inf(), err
	case ducktype.NegativeInfinity:
		return timex.NegInf(), err
	default:
		return v.Time, err
	}
}

// NewIntRowScanner ...
func NewIntRowScanner(r Row) IntRowScanner {
	return IntRowScanner{Row: r}
}

// IntRowScanner helper for scanning integers in a type safe manner.
type IntRowScanner struct {
	Row
}

// Scan an integer result
func (t IntRowScanner) Scan(v *int) error {
	return t.Row.Scan(v)
}

// NewBoolRowScanner ...
func NewBoolRowScanner(r Row) BoolRowScanner {
	return BoolRowScanner{Row: r}
}

// BoolRowScanner helper for scanning bools in a type safe manner.
type BoolRowScanner struct {
	Row
}

// Scan an bool result
func (t BoolRowScanner) Scan(v *bool) error {
	return t.Row.Scan(v)
}

func NewValueRowScanner[T any](r Row) ValueRowScanner[T] {
	return ValueRowScanner[T]{Row: r}
}

type ValueRowScanner[T any] struct {
	Row
}

func (t ValueRowScanner[T]) Scan(v *T) error {
	return t.Row.Scan(v)
}

// MaybeWriteCSV convience method for writing results into a csv.
func MaybeWriteCSV(dst io.Writer, rows *sql.Rows, err error) error {
	if err != nil {
		return err
	}

	return WriteCSV(rows, dst)
}

// WriteCSV output a csv of the rows into the io.Writer.
func WriteCSV(rows *sql.Rows, dst io.Writer) error {
	var (
		err     error
		columns []string
	)
	if columns, err = rows.Columns(); err != nil {
		return err
	}

	out := csv.NewWriter(dst)
	if err = out.Write(columns); err != nil {
		return errorsx.Wrap(err, "failed to write csv headers")
	}

	results := make([]any, len(columns))
	resultStrings := make([]sql.NullString, len(columns))
	for i := range resultStrings {
		results[i] = &resultStrings[i]
	}

	for rows.Next() {
		if err = rows.Scan(results...); err != nil {
			return errorsx.Wrap(err, "failed to scan results")
		}

		results := make([]string, 0, len(columns))
		for _, x := range resultStrings {
			results = append(results, x.String)
		}
		if err = out.Write(results); err != nil {
			return errorsx.Wrap(err, "failed to write csv row")
		}
	}

	if err = rows.Err(); err != nil {
		return errorsx.Wrap(err, "failed scanning results")
	}

	out.Flush()
	return out.Error()
}

// WithinTransaction will create a transaction if possible from the queryer.
func WithinTransaction(ctx context.Context, q Queryer, do func(context.Context, Queryer) error) (err error) {
	type transable interface {
		BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	}

	var (
		tx *sql.Tx
	)

	switch actual := q.(type) {
	case transable:
		if tx, err = actual.BeginTx(ctx, nil); err != nil {
			return err
		}

		return CompleteTx(tx, do(ctx, tx))
	case *sql.Tx:
		return do(ctx, q)
	default:
		return fmt.Errorf("not a transaction")
	}
}

// CompleteTx completes a transaction based on the presence of an error.
// if no error is present, then the transaction attempts to commit.
// if an error is present, then the transaction attempts to rollback.
// when a rollback is attempted the original error is returned.
func CompleteTx(tx *sql.Tx, err error) error {
	if err == nil {
		return tx.Commit()
	}

	// log if an error occurs rollingback, but continue on because if
	// we are rolling back it means there was another error anyways.
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		log.Println("error rolling back transaction", err)
	}

	return err
}

// IgnoreNoRows converts sql.ErrNoRows to nil.
func IgnoreNoRows(err error) error {
	switch errorsx.Cause(err) {
	case sql.ErrNoRows:
		return nil
	default:
		return err
	}
}

// IgnoreTxDone converts sql.ErrTxDone to nil.
func IgnoreTxDone(err error) error {
	switch err {
	case sql.ErrTxDone:
		return nil
	default:
		return err
	}
}

func ErrNoRows(err error) error {
	switch errorsx.Cause(err) {
	case sql.ErrNoRows:
		return err
	default:
		return nil
	}
}

// PrependTableName takes a reference set of columns and prepends the table to them.
func PrependTableName(table, reference string) string {
	columns := strings.Split(reference, ",")
	results := make([]string, 0, len(columns))
	for _, column := range columns {
		results = append(results, table+"."+column)
	}

	return strings.Join(results, ",")
}

// Columns split the gql column string constant into its constituent columns.
func Columns(reference ...string) (res []string) {
	for _, ref := range reference {
		res = append(res, strings.Split(ref, ",")...)
	}
	return res
}

// ExpandStrings into an array of interfaces to pass to squirrel.
func ExpandStrings(in []string) (results []any) {
	for _, i := range in {
		results = append(results, i)
	}

	return results
}

func SearchQuery(s string) string {
	remove := runes.Map(func(r rune) rune {
		if unicode.IsPunct(r) || unicode.IsControl(r) || unicode.IsSymbol(r) {
			return '\n'
		}

		return r
	})

	trans := transform.Chain(remove, runes.Map(unicode.ToLower), runes.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '\n'
		}
		return r
	}))

	s, _, _ = transform.String(trans, s)
	parts := slicesx.Filter(func(s string) bool {
		return !stringsx.Blank(s)
	}, slicesx.Map(strings.TrimSpace, strings.Split(s, "\n")...)...)

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, ":* | ") + ":*"
}

func NewNullString(s string) (x sql.NullString) {
	x.String = s
	x.Valid = true
	return x
}
