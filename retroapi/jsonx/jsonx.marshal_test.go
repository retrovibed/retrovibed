package jsonx_test

import (
	"math"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	type payload struct {
		Name string
		ID   uint64
	}

	t.Run("uint64 encodes as exact number", func(t *testing.T) {
		in := payload{Name: "derp", ID: math.MaxUint64}

		encoded, err := jsonx.Marshal(&in)
		require.NoError(t, err)
		require.Equal(t, `{"Name":"derp","ID":18446744073709551615}`, string(encoded))
	})

	t.Run("struct with no exported fields encodes as empty object", func(t *testing.T) {
		type empty struct {
			unexported string
		}

		encoded, err := jsonx.Marshal(&empty{unexported: "derp"})
		require.NoError(t, err)
		require.Equal(t, `{}`, string(encoded))
	})
}
