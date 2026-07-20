package duckdbx

import (
	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
)

const searchColumn = "val"

// Search filters an ad-hoc set of literal string values down to the ones
// matching lucene query q, projecting a boolean literal per match. The
// values never touch a real table -- DuckDB's UNNEST expands the bound
// []string parameter into rows inline; wrapped in a subquery (FromSelect)
// because a WHERE clause cannot reference a column alias computed in the
// same SELECT list. A blank q matches every value.
func Search(q string, values ...string) squirrel.SelectBuilder {
	inner := squirrelx.PSQL.
		Select().
		Column("UNNEST(?::VARCHAR[]) AS "+searchColumn, values)

	b := squirrelx.PSQL.
		Select("true").
		FromSelect(inner, "t")

	if stringsx.Blank(q) {
		return b
	}

	return b.Where(lucenex.Query(NewLucene(), q, lucenex.WithDefaultField(searchColumn)))
}
