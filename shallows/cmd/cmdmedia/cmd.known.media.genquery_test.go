package cmdmedia

import (
	"bytes"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/stretchr/testify/require"
)

func TestKnownGenqueryRun(t *testing.T) {
	t.Run("writes nothing for empty args", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, knowngenquery{}.run(&buf))
		require.Empty(t, buf.String())
	})

	t.Run("writes single query", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, knowngenquery{Queries: []string{"inception"}}.run(&buf))

		type output struct {
			Query string `json:"query"`
		}
		var got output
		require.NoError(t, jsonl.NewDecoder(&buf).Decode(&got))
		require.Equal(t, "inception", got.Query)
	})

	t.Run("writes one record per query", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, knowngenquery{Queries: []string{"inception", "the dark knight", "interstellar"}}.run(&buf))
		require.Equal(t, 3, strings.Count(buf.String(), "\n"))
	})
}
