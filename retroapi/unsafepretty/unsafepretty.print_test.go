package unsafepretty_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/unsafepretty"
	"github.com/stretchr/testify/require"
)

func TestPrint(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		require.Equal(t, "", unsafepretty.Print(""))
	})

	t.Run("printable text is unchanged", func(t *testing.T) {
		require.Equal(t, "foo-bar_baz.123", unsafepretty.Print("foo-bar_baz.123"))
	})

	t.Run("printable unicode is unchanged", func(t *testing.T) {
		require.Equal(t, "héllo世界", unsafepretty.Print("héllo世界"))
	})

	t.Run("newlines and spaces are unchanged by default", func(t *testing.T) {
		require.Equal(t, "foo bar\nbaz", unsafepretty.Print("foo bar\nbaz"))
	})

	t.Run("tabs are escaped", func(t *testing.T) {
		require.Equal(t, `foo\tbar`, unsafepretty.Print("foo\tbar"))
	})

	t.Run("carriage returns are escaped", func(t *testing.T) {
		require.Equal(t, `foo\r`+"\n"+`bar`, unsafepretty.Print("foo\r\nbar"))
	})

	t.Run("non printable runes are rendered as unicode code points", func(t *testing.T) {
		require.Equal(t, "foo U+0000 bar", unsafepretty.Print("foo \x00 bar"))
		require.Equal(t, "U+007F", unsafepretty.Print("\x7f"))
	})

	t.Run("display spaces", func(t *testing.T) {
		require.Equal(t, "foo·bar··baz", unsafepretty.Print("foo bar  baz", unsafepretty.OptionDisplaySpaces()))
	})

	t.Run("display spaces leaves other whitespace handling alone", func(t *testing.T) {
		require.Equal(t, `foo·\tbar`+"\n"+`baz`, unsafepretty.Print("foo \tbar\nbaz", unsafepretty.OptionDisplaySpaces()))
	})

	t.Run("custom space runes", func(t *testing.T) {
		require.Equal(t, "foo[_]bar", unsafepretty.Print("foo bar", unsafepretty.OptionSpaceRunes('[', '_', ']')))
	})

	t.Run("custom newline runes", func(t *testing.T) {
		require.Equal(t, "foo⏎\nbar", unsafepretty.Print("foo\nbar", unsafepretty.OptionNewlineRunes('⏎', '\n')))
	})

	t.Run("newline runes can be removed", func(t *testing.T) {
		require.Equal(t, "foobar", unsafepretty.Print("foo\nbar", unsafepretty.OptionNewlineRunes()))
	})

	t.Run("later options override earlier options", func(t *testing.T) {
		require.Equal(t, "foo_bar", unsafepretty.Print("foo bar", unsafepretty.OptionDisplaySpaces(), unsafepretty.OptionSpaceRunes('_')))
	})

	t.Run("multiple options combine", func(t *testing.T) {
		require.Equal(t, "a·b¶c", unsafepretty.Print("a b\nc", unsafepretty.OptionDisplaySpaces(), unsafepretty.OptionNewlineRunes('¶')))
	})
}
