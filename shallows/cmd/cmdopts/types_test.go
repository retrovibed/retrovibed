package cmdopts

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/require"
)

func TestIOOut(t *testing.T) {
	t.Run("dash opens stdout with nop close", func(t *testing.T) {
		var v IOOut
		require.NoError(t, v.UnmarshalText([]byte("-")))

		w, err := v.Open(os.Stdout)
		require.NoError(t, err)
		require.NoError(t, w.Close()) // must not close actual stdout
		require.NotNil(t, w)
	})

	t.Run("path creates file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "out.txt")

		var v IOOut
		require.NoError(t, v.UnmarshalText([]byte(path)))

		w, err := v.Open(os.Stdout)
		require.NoError(t, err)
		defer w.Close()

		_, err = w.Write([]byte("hello"))
		require.NoError(t, err)
		require.NoError(t, w.Close())

		content, err := os.ReadFile(path)
		require.NoError(t, err)
		require.Equal(t, "hello", string(content))
	})

	t.Run("path truncates existing file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "out.txt")
		require.NoError(t, os.WriteFile(path, []byte("old content longer than new"), 0600))

		var v IOOut
		require.NoError(t, v.UnmarshalText([]byte(path)))

		w, err := v.Open(os.Stdout)
		require.NoError(t, err)
		_, err = w.Write([]byte("new"))
		require.NoError(t, err)
		require.NoError(t, w.Close())

		content, err := os.ReadFile(path)
		require.NoError(t, err)
		require.Equal(t, "new", string(content))
	})

	t.Run("marshal roundtrips path", func(t *testing.T) {
		var v IOOut
		require.NoError(t, v.UnmarshalText([]byte("/some/path")))
		text, err := v.MarshalText()
		require.NoError(t, err)
		require.Equal(t, "/some/path", string(text))
	})
}

func TestFileContents(t *testing.T) {
	t.Run("plain value is stored as-is", func(t *testing.T) {
		var v FileContents
		require.NoError(t, v.UnmarshalText([]byte("hello world")))
		require.Equal(t, FileContents("hello world"), v)
	})

	t.Run("@ prefix reads file contents", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "secret.txt")
		require.NoError(t, os.WriteFile(path, []byte("file contents"), 0600))

		var v FileContents
		require.NoError(t, v.UnmarshalText([]byte("@"+path)))
		require.Equal(t, FileContents("file contents"), v)
	})

	t.Run("@ prefix with missing file returns error", func(t *testing.T) {
		var v FileContents
		err := v.UnmarshalText([]byte("@/nonexistent/path"))
		require.Error(t, err)
	})

	t.Run("marshal roundtrips value", func(t *testing.T) {
		v := FileContents("some value")
		text, err := v.MarshalText()
		require.NoError(t, err)
		require.Equal(t, "some value", string(text))
	})
}

func TestParseEnviron(t *testing.T) {
	type cli struct {
		Env []string `flag:"" name:"env" short:"e" type:"envvar"`
	}

	parse := func(t *testing.T, args ...string) ([]string, error) {
		var c cli
		parser, err := kong.New(&c, kong.NamedMapper("envvar", kong.MapperFunc(ParseEnviron)))
		require.NoError(t, err)

		if _, err := parser.Parse(args); err != nil {
			return nil, err
		}

		return c.Env, nil
	}

	t.Run("literal key value pair", func(t *testing.T) {
		env, err := parse(t, "-e", "FOO=bar")
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=bar"}, env)
	})

	t.Run("repeated flags accumulate in order", func(t *testing.T) {
		env, err := parse(t, "-e", "FOO=1", "-e", "BAR=2")
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=1", "BAR=2"}, env)
	})

	t.Run("file scheme contributes multiple pairs", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "plugin.env")
		require.NoError(t, os.WriteFile(path, []byte("FOO=1\nBAR=2\n"), 0600))

		env, err := parse(t, "-e", "file://"+path)
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=1", "BAR=2"}, env)
	})

	t.Run("file scheme mixed with literal pairs", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "plugin.env")
		require.NoError(t, os.WriteFile(path, []byte("FOO=1\n"), 0600))

		env, err := parse(t, "-e", "file://"+path, "-e", "BAR=2")
		require.NoError(t, err)
		require.Equal(t, []string{"FOO=1", "BAR=2"}, env)
	})

	t.Run("invalid token returns error", func(t *testing.T) {
		_, err := parse(t, "-e", "notapair")
		require.Error(t, err)
	})
}
