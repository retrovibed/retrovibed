package tracking_test

import (
	"encoding/json"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

// Regression: protojson encodes uint64 as a quoted string; genieql structs must
// accept the string form via the ,string tag.
func TestMetadataJSONRoundTrip(t *testing.T) {
	original := tracking.Metadata{
		Bytes:      1099511627776,
		Downloaded: 549755813888,
		Uploaded:   274877906944,
	}

	encoded, err := json.Marshal(original)
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"bytes":"`)
	require.Contains(t, string(encoded), `"downloaded":"`)
	require.Contains(t, string(encoded), `"uploaded":"`)

	var decoded tracking.Metadata
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, original.Bytes, decoded.Bytes)
	require.Equal(t, original.Downloaded, decoded.Downloaded)
	require.Equal(t, original.Uploaded, decoded.Uploaded)

	// also accept protojson string format directly
	require.NoError(t, json.Unmarshal([]byte(`{"bytes":"1099511627776","downloaded":"549755813888","uploaded":"274877906944"}`), &decoded))
	require.Equal(t, original.Bytes, decoded.Bytes)
	require.Equal(t, original.Downloaded, decoded.Downloaded)
	require.Equal(t, original.Uploaded, decoded.Uploaded)
}
