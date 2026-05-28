package mimex_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/stretchr/testify/require"
)

func TestCategory(t *testing.T) {
	tests := []struct {
		mime     string
		expected string
	}{
		{"video/mp4", mimex.Video},
		{"video/x-msvideo", mimex.Video},
		{"audio/mpeg", mimex.Audio},
		{"audio/flac", mimex.Audio},
		{"image/jpeg", mimex.Image},
		{"image/png", mimex.Image},
		{"application/json", mimex.Application},
		{"text/plain", mimex.Text},
	}

	for _, tc := range tests {
		t.Run(tc.mime, func(t *testing.T) {
			require.Equal(t, tc.expected, mimex.Category(tc.mime))
		})
	}
}
