package jsonl_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/stretchr/testify/require"
)

func TestEncoder(t *testing.T) {
	type record struct{ Name string }

	t.Run("encodes a single value", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&buf).Encode(record{Name: "alice"}))
		require.Equal(t, `{"Name":"alice"}`+"\n", buf.String())
	})

	t.Run("encodes multiple values as separate lines", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&buf).Encode(record{Name: "alice"}, record{Name: "bob"}))
		require.Equal(t, `{"Name":"alice"}`+"\n"+`{"Name":"bob"}`+"\n", buf.String())
	})

	t.Run("encoding zero values writes nothing", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&buf).Encode())
		require.Empty(t, buf.String())
	})

	t.Run("stops at the first failing value", func(t *testing.T) {
		var buf bytes.Buffer
		require.Error(t, jsonl.NewEncoder(&buf).Encode(record{Name: "alice"}, make(chan int), record{Name: "bob"}))
		require.Equal(t, `{"Name":"alice"}`+"\n", buf.String())
	})
}

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
