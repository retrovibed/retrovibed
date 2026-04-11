package squirrelx

import (
	"regexp"

	"github.com/Masterminds/squirrel"
)

// PSQL postgresql statement builder.
var PSQL = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Noop struct{}

func (t Noop) ToSql() (sql string, args []interface{}, err error) {
	return "", nil, nil
}

// Between generate a between clause.
func Between(c string, a, b interface{}) squirrel.Sqlizer {
	return squirrel.Expr("("+c+" BETWEEN ? AND ?"+")", a, b)
}

// In predicate.
func In[T any](expr string, values ...T) squirrel.Sqlizer {
	if len(values) == 0 {
		return Noop{}
	}

	r := make([]interface{}, 0, len(values))
	for _, v := range values {
		r = append(r, v)
	}
	return squirrel.Expr(expr+" IN ("+squirrel.Placeholders(len(values))+")", r...)
}

var psqlplaceholder = regexp.MustCompile(`\$\d+`)

func Sprint(q string, args ...any) string {
	return squirrel.DebugSqlizer(squirrel.Expr(psqlplaceholder.ReplaceAllString(q, "?"), args...))
}

func QueryNonZero[T comparable](expr string, s T) squirrel.Sqlizer {
	var zero T
	if s == zero {
		return Noop{}
	}

	return squirrel.Expr(expr, s)
}

func Error(cause error) squirrel.Sqlizer {
	return err{cause: cause}
}

type err struct {
	cause error
}

func (t err) ToSql() (sql string, args []interface{}, err error) {
	return "", nil, t.cause
}

type SqlizerFn func() (sql string, args []interface{}, err error)

func (t SqlizerFn) ToSql() (sql string, args []interface{}, err error) {
	return t()
}
