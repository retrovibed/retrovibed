package envfile_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/envfile"
	"github.com/stretchr/testify/require"
)

const example = "FOO=\"bar\" # derp 0\n# derp 1\nBAR=\"baz\"\nBIZ=\"BAN\"\n# derp 2\n"

func TestParse(t *testing.T) {
	t.Run("comment block and inline comment both attach as hints", func(t *testing.T) {
		require.Equal(t, []envfile.Variable{
			{Key: "FOO", Value: "bar", Hint: "derp 0"},
			{Key: "BAR", Value: "baz", Hint: "derp 1"},
			{Key: "BIZ", Value: "BAN", Hint: ""},
		}, envfile.Parse(example))
	})

	t.Run("empty content yields no variables", func(t *testing.T) {
		require.Empty(t, envfile.Parse(""))
	})

	t.Run("unquoted value with no comment", func(t *testing.T) {
		require.Equal(t, []envfile.Variable{
			{Key: "FOO", Value: "bar", Hint: ""},
		}, envfile.Parse("FOO=bar\n"))
	})

	t.Run("blank line clears pending comment", func(t *testing.T) {
		require.Equal(t, []envfile.Variable{
			{Key: "FOO", Value: "bar", Hint: ""},
		}, envfile.Parse("# orphaned\n\nFOO=bar\n"))
	})
}

func TestApply(t *testing.T) {
	t.Run("edits an existing value, preserving its comment", func(t *testing.T) {
		got := envfile.Apply(example, []envfile.Variable{
			{Key: "FOO", Value: "updated"},
		})
		require.Equal(t, "FOO=updated # derp 0\n# derp 1\nBAR=\"baz\"\nBIZ=\"BAN\"\n# derp 2\n", got)
	})

	t.Run("appends an edit whose key has no matching line", func(t *testing.T) {
		got := envfile.Apply("FOO=bar\n", []envfile.Variable{
			{Key: "NEW", Value: "value"},
		})
		require.Equal(t, "FOO=bar\nNEW=value\n", got)
	})

	t.Run("quotes values containing whitespace or a comment character", func(t *testing.T) {
		got := envfile.Apply("FOO=bar\n", []envfile.Variable{
			{Key: "FOO", Value: "has space"},
		})
		require.Equal(t, "FOO=\"has space\"\n", got)
	})

	t.Run("round trip: parse(apply(content, edits)) reflects the edits", func(t *testing.T) {
		edits := []envfile.Variable{{Key: "FOO", Value: "changed"}}
		got := envfile.Parse(envfile.Apply(example, edits))
		require.Equal(t, []envfile.Variable{
			{Key: "FOO", Value: "changed", Hint: "derp 0"},
			{Key: "BAR", Value: "baz", Hint: "derp 1"},
			{Key: "BIZ", Value: "BAN", Hint: ""},
		}, got)
	})
}
