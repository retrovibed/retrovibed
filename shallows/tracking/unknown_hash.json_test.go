package tracking_test

import (
	"encoding/json"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

// Regression: protojson encodes uint64 as a quoted string; genieql structs must
// accept the string form via the ,string tag.
func TestUnknownHashJSONRoundTrip(t *testing.T) {
	original := tracking.UnknownHash{Attempts: 9876543210}

	encoded, err := json.Marshal(original)
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"attempts":"`)

	var decoded tracking.UnknownHash
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, original.Attempts, decoded.Attempts)

	// also accept protojson string format directly
	require.NoError(t, json.Unmarshal([]byte(`{"attempts":"9876543210"}`), &decoded))
	require.Equal(t, original.Attempts, decoded.Attempts)
}
