package tracking_test

import (
	"encoding/json"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

// Regression: protojson encodes uint64 as a quoted string; genieql structs must
// accept the string form via the ,string tag.
func TestPeerJSONRoundTrip(t *testing.T) {
	original := tracking.Peer{Bep51Available: 9876543210}

	encoded, err := json.Marshal(original)
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"bep_51_available":"`)

	var decoded tracking.Peer
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, original.Bep51Available, decoded.Bep51Available)

	// also accept protojson string format directly
	require.NoError(t, json.Unmarshal([]byte(`{"bep_51_available":"9876543210"}`), &decoded))
	require.Equal(t, original.Bep51Available, decoded.Bep51Available)
}
