package meta_test

import (
	"encoding/json"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/stretchr/testify/require"
)

// Regression: protojson encodes uint64 as a quoted string; genieql structs must
// accept the string form via the ,string tag.
func TestWireguardJSONRoundTrip(t *testing.T) {
	original := meta.Wireguard{MaximumConnections: 9876543210}

	encoded, err := json.Marshal(original)
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"maximum_connections":"`)

	var decoded meta.Wireguard
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, original.MaximumConnections, decoded.MaximumConnections)

	// also accept protojson string format directly
	require.NoError(t, json.Unmarshal([]byte(`{"maximum_connections":"9876543210"}`), &decoded))
	require.Equal(t, original.MaximumConnections, decoded.MaximumConnections)
}
