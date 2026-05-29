package community_test

import (
	"encoding/json"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/stretchr/testify/require"
)

// Regression: protojson encodes uint64 as a quoted string; genieql structs must
// accept the string form via the ,string tag.
func TestPublishedContentJSONRoundTrip(t *testing.T) {
	original := community.PublishedContent{Bytes: 1099511627776}

	encoded, err := json.Marshal(original)
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"bytes":"`)

	var decoded community.PublishedContent
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, original.Bytes, decoded.Bytes)

	// also accept protojson string format directly
	require.NoError(t, json.Unmarshal([]byte(`{"bytes":"1099511627776"}`), &decoded))
	require.Equal(t, original.Bytes, decoded.Bytes)
}
