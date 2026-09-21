package lucenex_test

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	t.Run("fields", func(t *testing.T) {
		q, args, err := lucenex.Query(duckdbx.NewLucene(), "(mimetype:\"video/webm\" OR mimetype:\"video/ogg\")", lucenex.WithDefaultField("auto_description")).ToSql()
		require.NoError(t, err)
		require.EqualValues(t, "(\"mimetype\" ILIKE '%' || ? || '%') OR (\"mimetype\" ILIKE '%' || ? || '%')", q)
		require.EqualValues(t, []any{"video/webm", "video/ogg"}, args)
	})

	// the rendered query must be a self contained expression so that it can be safely
	// combined with sibling predicates. otherwise a top level OR escapes the enclosing
	// AND, i.e. `adult = ? AND (a) OR (b)` which matches every row where b is true regardless of adult.
	t.Run("composed with a sibling predicate", func(t *testing.T) {
		for _, tc := range []struct {
			name     string
			terms    string
			expected string
			args     []any
		}{
			{
				name:     "single term",
				terms:    "two",
				expected: "(adult = ? AND (\"title\" ILIKE '%' || ? || '%'))",
				args:     []any{false, "two"},
			},
			{
				name:     "implicit and",
				terms:    "two guys",
				expected: "(adult = ? AND ((\"title\" ILIKE '%' || ? || '%') AND (\"title\" ILIKE '%' || ? || '%')))",
				args:     []any{false, "two", "guys"},
			},
			{
				name:     "explicit and",
				terms:    "two AND guys",
				expected: "(adult = ? AND ((\"title\" ILIKE '%' || ? || '%') AND (\"title\" ILIKE '%' || ? || '%')))",
				args:     []any{false, "two", "guys"},
			},
			{
				name:     "or",
				terms:    "two OR guys",
				expected: "(adult = ? AND ((\"title\" ILIKE '%' || ? || '%') OR (\"title\" ILIKE '%' || ? || '%')))",
				args:     []any{false, "two", "guys"},
			},
			{
				name:     "chained or",
				terms:    "two OR guys OR a OR girl",
				expected: "(adult = ? AND ((((\"title\" ILIKE '%' || ? || '%') OR (\"title\" ILIKE '%' || ? || '%')) OR (\"title\" ILIKE '%' || ? || '%')) OR (\"title\" ILIKE '%' || ? || '%')))",
				args:     []any{false, "two", "guys", "a", "girl"},
			},
			{
				name:     "and then or",
				terms:    "two AND guys OR a",
				expected: "(adult = ? AND (((\"title\" ILIKE '%' || ? || '%') AND (\"title\" ILIKE '%' || ? || '%')) OR (\"title\" ILIKE '%' || ? || '%')))",
				args:     []any{false, "two", "guys", "a"},
			},
			{
				name:     "grouped or then and",
				terms:    "(two OR guys) AND girl",
				expected: "(adult = ? AND (((\"title\" ILIKE '%' || ? || '%') OR (\"title\" ILIKE '%' || ? || '%')) AND (\"title\" ILIKE '%' || ? || '%')))",
				args:     []any{false, "two", "guys", "girl"},
			},
			{
				name:     "or across fields",
				terms:    "title:two OR subtitle:guys",
				expected: "(adult = ? AND ((\"title\" ILIKE '%' || ? || '%') OR (\"subtitle\" ILIKE '%' || ? || '%')))",
				args:     []any{false, "two", "guys"},
			},
			{
				name:     "quoted phrase or term",
				terms:    "\"two guys\" OR girl",
				expected: "(adult = ? AND ((\"title\" ILIKE '%' || ? || '%') OR (\"title\" ILIKE '%' || ? || '%')))",
				args:     []any{false, "two guys", "girl"},
			},
			{
				name:     "not",
				terms:    "NOT two",
				expected: "(adult = ? AND (NOT(\"title\" ILIKE '%' || ? || '%')))",
				args:     []any{false, "two"},
			},
			{
				name:     "blank terms are ignored",
				terms:    "",
				expected: "(adult = ?)",
				args:     []any{false},
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				q, args, err := squirrel.And{
					squirrel.Eq{"adult": false},
					lucenex.Query(duckdbx.NewLucene(), tc.terms, lucenex.WithDefaultField("title")),
				}.ToSql()
				require.NoError(t, err)
				require.EqualValues(t, tc.expected, q)
				require.EqualValues(t, tc.args, args)
			})
		}
	})

	t.Run("composed as an or sibling", func(t *testing.T) {
		q, args, err := squirrel.Or{
			squirrel.Eq{"adult": false},
			lucenex.Query(duckdbx.NewLucene(), "two AND guys", lucenex.WithDefaultField("title")),
		}.ToSql()
		require.NoError(t, err)
		require.EqualValues(t, "(adult = ? OR ((\"title\" ILIKE '%' || ? || '%') AND (\"title\" ILIKE '%' || ? || '%')))", q)
		require.EqualValues(t, []any{false, "two", "guys"}, args)
	})
}
