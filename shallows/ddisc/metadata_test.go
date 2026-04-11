//go:build !darwin && ffmpeg_enabled

package ddisc_test

import (
	"io"
	"testing"

	"github.com/retrovibed/retrovibed/ddisc"
	"github.com/retrovibed/retrovibed/internal/bytesx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestAudioMetadata(t *testing.T) {
	// example generation command for metadata.
	// ffmpeg -f lavfi -i anullsrc=r=44100:cl=mono -t 10 -q:a 0 -metadata title="vibed" -metadata artist="retrovibed" -metadata album="example 1" -metadata track="02" shallows/ddisc/.fixtures/TestAudioMetadata/example.1.flac
	t.Run("example.1.flac", func(t *testing.T) {
		extract, err := ddisc.Extract(testx.Read(".fixtures", t.Name()))
		require.NoError(t, err)
		require.Equal(t, "retrovibed - example 1 01 vibed", extract.Music.Metadata.String())
	})

	t.Run("example.2.flac", func(t *testing.T) {
		// missing track test
		extract, err := ddisc.Extract(testx.Read(".fixtures", t.Name()))
		require.NoError(t, err)
		require.Equal(t, "retrovibed - example 1 vibed", extract.Music.Metadata.String())
	})

	t.Run("example.3.flac", func(t *testing.T) {
		// missing artist test
		extract, err := ddisc.Extract(testx.Read(".fixtures", t.Name()))
		require.NoError(t, err)
		require.Equal(t, "example 1 01 vibed", extract.Music.Metadata.String())
	})

	t.Run("example.4.flac", func(t *testing.T) {
		// missing album test
		extract, err := ddisc.Extract(testx.Read(".fixtures", t.Name()))
		require.NoError(t, err)
		require.Equal(t, "retrovibed 01 vibed", extract.Music.Metadata.String())
	})
}

func TestMimetype(t *testing.T) {
	t.Run("example.1.flac", func(t *testing.T) {
		mime, err := ddisc.Mimetype(io.NewSectionReader(testx.Read(".fixtures", "TestAudioMetadata", "example.1.flac"), 0, bytesx.KiB))
		require.NoError(t, err)
		require.Equal(t, "audio/flac", mime.String())
	})
}
