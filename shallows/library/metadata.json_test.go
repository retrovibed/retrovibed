package library_test

import (
	"encoding/json"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

// Regression: protojson encodes uint64 as a quoted string; genieql structs must
// accept the string form via the ,string tag.
func TestMetadataJSONRoundTrip(t *testing.T) {
	original := library.Metadata{
		Bytes:      1099511627776,
		DiskOffset: 2199023255552,
		DiskUsage:  549755813888,
		QuotaUsage: 274877906944,
	}

	encoded, err := json.Marshal(original)
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"bytes":"`)
	require.Contains(t, string(encoded), `"disk_offset":"`)
	require.Contains(t, string(encoded), `"disk_usage":"`)
	require.Contains(t, string(encoded), `"quota_usage":"`)

	var decoded library.Metadata
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, original.Bytes, decoded.Bytes)
	require.Equal(t, original.DiskOffset, decoded.DiskOffset)
	require.Equal(t, original.DiskUsage, decoded.DiskUsage)
	require.Equal(t, original.QuotaUsage, decoded.QuotaUsage)

	// also accept protojson string format directly
	require.NoError(t, json.Unmarshal([]byte(`{"bytes":"1099511627776","disk_offset":"2199023255552","disk_usage":"549755813888","quota_usage":"274877906944"}`), &decoded))
	require.Equal(t, original.Bytes, decoded.Bytes)
	require.Equal(t, original.DiskOffset, decoded.DiskOffset)
	require.Equal(t, original.DiskUsage, decoded.DiskUsage)
	require.Equal(t, original.QuotaUsage, decoded.QuotaUsage)
}
