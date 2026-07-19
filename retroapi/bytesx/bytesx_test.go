package bytesx_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/bytesx"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	mib, gib := float64(bytesx.MiB), float64(bytesx.GiB)
	tests := []struct {
		size     string
		expected uint64
	}{
		{"214.2 MB", uint64(214.2 * mib)},
		{"1.2 GB", uint64(1.2 * gib)},
		{"900 MB", uint64(900 * mib)},
		{"1 TB", uint64(bytesx.TiB)},
		{"512 B", 512},
		{"1.2 gb", uint64(1.2 * gib)},
		{"", 0},
		{"garbage", 0},
		{"1.2", 0},
		{"1.2 QB", 0},
		{"notanumber MB", 0},
	}

	for _, tc := range tests {
		t.Run(tc.size, func(t *testing.T) {
			require.Equal(t, tc.expected, bytesx.Parse(tc.size))
		})
	}
}
