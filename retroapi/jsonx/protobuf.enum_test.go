package jsonx_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/stretchr/testify/require"
)

func TestEnum(t *testing.T) {
	type payload struct {
		Mode authn.ProfileStatus
	}

	t.Run("marshals as the underlying number", func(t *testing.T) {
		encoded, err := jsonx.Marshal(&payload{Mode: authn.ProfileStatus_DISABLED})
		require.NoError(t, err)
		require.Equal(t, `{"Mode":1}`, string(encoded))
	})

	t.Run("unmarshals from the declared value name", func(t *testing.T) {
		const doc = `{"Mode":"PENDING"}`

		var out payload
		require.NoError(t, jsonx.Unmarshal([]byte(doc), &out))
		require.Equal(t, authn.ProfileStatus_PENDING, out.Mode)
	})

	t.Run("unmarshals from the underlying number", func(t *testing.T) {
		const doc = `{"Mode":1}`

		var out payload
		require.NoError(t, jsonx.Unmarshal([]byte(doc), &out))
		require.Equal(t, authn.ProfileStatus_DISABLED, out.Mode)
	})

	// a real client (e.g. the console app's generated protobuf-JSON
	// serializer) sends the canonical proto3 JSON wire format directly,
	// independent of however this package marshals its own Go structs -
	// exercise that literal wire format, not a round trip through jsonx
	// itself, or a client/server format mismatch like this can pass tests
	// while still breaking in production.
	t.Run("decodes a literal external client payload", func(t *testing.T) {
		const doc = `{"Mode":"DISABLED"}`

		var out payload
		require.NoError(t, jsonx.Unmarshal([]byte(doc), &out))
		require.Equal(t, authn.ProfileStatus_DISABLED, out.Mode)
	})
}
