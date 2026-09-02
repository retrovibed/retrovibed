package jsonx_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalRead(t *testing.T) {
	type payload struct {
		Name string
		ID   uint64
	}

	t.Run("matches Unmarshal", func(t *testing.T) {
		in := payload{Name: "match", ID: math.MaxUint64}

		encoded, err := jsonx.Marshal(&in)
		require.NoError(t, err)

		var a, b payload
		require.NoError(t, jsonx.Unmarshal(encoded, &a))
		require.NoError(t, jsonx.UnmarshalRead(bytes.NewReader(encoded), &b))
		require.Equal(t, a, b)
	})
}
