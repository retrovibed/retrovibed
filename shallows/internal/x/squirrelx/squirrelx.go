package squirrelx

import (
	"fmt"

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
	r := make([]interface{}, 0, len(values))
	for _, v := range values {
		r = append(r, v)
	}
	return squirrel.Expr(expr+" IN ("+squirrel.Placeholders(len(values))+")", r...)
}

func Sprint(s squirrel.Sqlizer) string {
	q, args, err := s.ToSql()
	if err != nil {
		return err.Error()
	}

	return q + " " + fmt.Sprint(args...)
}

func QueryNonZero[T comparable](expr string, s T) squirrel.Sqlizer {
	var zero T
	if s == zero {
		return Noop{}
	}

	return squirrel.Expr(expr, s)
}
