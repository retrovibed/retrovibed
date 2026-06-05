package duckdbx

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
)

func FTSSearch(table string, q string, columns ...string) squirrel.Sqlizer {
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

		if len(columns) > 0 {
			for _, term := range positive {
				pexpr = squirrel.ConcatExpr(pexpr, squirrel.Expr(fmt.Sprintf(" OR (%s ILIKE ?)", stringsx.Join(" || ", columns...)), "%"+term+"%"))
			}
		}
	}

	nexpr := squirrel.Expr("TRUE")
	if len(negative) > 0 {
		nexpr = squirrel.Expr(
			fmt.Sprintf("COALESCE(%s.match_bm25(id, ?), 0) = 0", table),
			strings.Join(negative, " "),
		)

		if len(columns) > 0 {
			for _, term := range negative {
				nexpr = squirrel.ConcatExpr(nexpr, squirrel.Expr(fmt.Sprintf(" AND (%s NOT ILIKE ?)", stringsx.Join(" || ", columns...)), term))
			}
		}
	}

	return squirrel.And{pexpr, nexpr}
}
