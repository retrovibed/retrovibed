package mimex_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
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
		{mimex.Video, mimex.Video},
		{mimex.Audio, mimex.Audio},
		{mimex.Image, mimex.Image},
		{mimex.Text, mimex.Text},
		{mimex.Application, mimex.Application},
	}

	for _, tc := range tests {
		t.Run(tc.mime, func(t *testing.T) {
			require.Equal(t, tc.expected, mimex.Category(tc.mime))
		})
	}
}
