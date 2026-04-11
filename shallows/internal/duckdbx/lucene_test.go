package duckdbx_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/internal/duckdbx"
	"github.com/retrovibed/retrovibed/internal/lucenex"
	"github.com/stretchr/testify/require"
)

func TestLuceneExample1(t *testing.T) {
	q, args, err := lucenex.Query(duckdbx.NewLucene(), "arch linux", lucenex.WithDefaultField("auto_description")).ToSql()
	require.NoError(t, err)
	require.Len(t, args, 2)
	require.Equal(t, "(\"auto_description\" ILIKE '%' || ? || '%') AND (\"auto_description\" ILIKE '%' || ? || '%')", q)
}

func TestLuceneExample2(t *testing.T) {
	q, args, err := lucenex.Query(duckdbx.NewLucene(), "arch AND linux", lucenex.WithDefaultField("auto_description")).ToSql()
	require.NoError(t, err)
	require.Len(t, args, 2)
	require.Equal(t, "(\"auto_description\" ILIKE '%' || ? || '%') AND (\"auto_description\" ILIKE '%' || ? || '%')", q)
}

func TestLuceneExample3(t *testing.T) {
	q, args, err := lucenex.Query(duckdbx.NewLucene(), "arch OR linux", lucenex.WithDefaultField("auto_description")).ToSql()
	require.NoError(t, err)
	require.Len(t, args, 2)
	require.Equal(t, "(\"auto_description\" ILIKE '%' || ? || '%') OR (\"auto_description\" ILIKE '%' || ? || '%')", q)
}
