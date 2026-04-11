package lucenex_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/internal/duckdbx"
	"github.com/retrovibed/retrovibed/internal/lucenex"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	t.Run("fields", func(t *testing.T) {
		q, args, err := lucenex.Query(duckdbx.NewLucene(), "(mimetype:\"video/webm\" OR mimetype:\"video/ogg\")", lucenex.WithDefaultField("auto_description")).ToSql()
		require.NoError(t, err)
		require.EqualValues(t, "(\"mimetype\" ILIKE '%' || ? || '%') OR (\"mimetype\" ILIKE '%' || ? || '%')", q)
		require.EqualValues(t, []any{"video/webm", "video/ogg"}, args)

	})
}
