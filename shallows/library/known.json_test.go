package library_test

import (
	"encoding/json"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

// Regression: protojson encodes uint64 as a quoted string; genieql structs must
// accept the string form via the ,string tag.
func TestKnownJSONRoundTrip(t *testing.T) {
	original := library.Known{Md5Lower: 9876543210}

	encoded, err := json.Marshal(original)
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"md_5_lower":"`)

	var decoded library.Known
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, original.Md5Lower, decoded.Md5Lower)

	// also accept protojson string format directly
	require.NoError(t, json.Unmarshal([]byte(`{"md_5_lower":"9876543210"}`), &decoded))
	require.Equal(t, original.Md5Lower, decoded.Md5Lower)
}
