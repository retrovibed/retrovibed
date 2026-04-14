package jsonl_test

import (
	"context"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/stretchr/testify/require"
)

func TestIter(t *testing.T) {
	t.Run("yields all values from valid jsonl input", func(t *testing.T) {
		input := strings.NewReader(`{"name":"alice"}` + "\n" + `{"name":"bob"}` + "\n")
		type record struct{ Name string }

		var results []record
		seq := jsonl.Iter[record](jsonl.NewDecoder(input))
		for v := range seq.Each(context.Background()) {
			results = append(results, v)
		}

		require.NoError(t, seq.Err())
		require.Equal(t, []record{{Name: "alice"}, {Name: "bob"}}, results)
	})

	t.Run("empty input yields no values", func(t *testing.T) {
		input := strings.NewReader("")
		type record struct{ Name string }

		var results []record
		seq := jsonl.Iter[record](jsonl.NewDecoder(input))
		for v := range seq.Each(context.Background()) {
			results = append(results, v)
		}

		require.NoError(t, seq.Err())
		require.Empty(t, results)
	})

	t.Run("stops and records error on invalid json", func(t *testing.T) {
		input := strings.NewReader(`{"name":"alice"}` + "\n" + `not-json` + "\n")
		type record struct{ Name string }

		var results []record
		seq := jsonl.Iter[record](jsonl.NewDecoder(input))
		for v := range seq.Each(context.Background()) {
			results = append(results, v)
		}

		require.Error(t, seq.Err())
		require.Equal(t, []record{{Name: "alice"}}, results)
	})

	t.Run("stops early when yield returns false", func(t *testing.T) {
		input := strings.NewReader(`{"name":"alice"}` + "\n" + `{"name":"bob"}` + "\n" + `{"name":"carol"}` + "\n")
		type record struct{ Name string }

		var results []record
		seq := jsonl.Iter[record](jsonl.NewDecoder(input))
		for v := range seq.Each(context.Background()) {
			results = append(results, v)
			break
		}

		require.NoError(t, seq.Err())
		require.Equal(t, []record{{Name: "alice"}}, results)
	})
}
