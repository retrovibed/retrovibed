package duckdbx

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/deeppool/internal/x/squirrelx"
	"github.com/james-lawrence/deeppool/internal/x/stringsx"
)

// checks if the error is a unique constraint violation.
func ErrUniqueConstraintViolation(err error) error {
	// TODO.
	return err
}

func FTSSearch(table string, q string) squirrel.Sqlizer {
	if stringsx.Blank(q) {
		return squirrelx.Noop{}
	}

	negative := []string{}
	positive := []string{}
	for _, s := range strings.Split(q, " ") {
		if strings.HasPrefix(s, "-") {
			negative = append(negative, strings.TrimPrefix(s, "-"))
		} else {
			positive = append(positive, s)
		}
	}

	pexpr := squirrel.Expr("TRUE")
	if len(positive) > 0 {
		pexpr = squirrel.Expr(
			fmt.Sprintf("COALESCE(%s.match_bm25(id, ?), 0) > 0", table),
			strings.Join(positive, " "),
		)
	}

	nexpr := squirrel.Expr("TRUE")
	if len(negative) > 0 {
		nexpr = squirrel.Expr(
			fmt.Sprintf("COALESCE(%s.match_bm25(id, ?), 0) = 0", table),
			strings.Join(negative, " "),
		)
	}

	return squirrel.And{pexpr, nexpr}
}
