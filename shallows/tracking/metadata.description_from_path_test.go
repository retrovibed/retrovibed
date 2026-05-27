package tracking

import (
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/stretchr/testify/require"
)

func TestDescriptionFromPath(t *testing.T) {
	t.Run("returns filename when path is not infohash and not description", func(t *testing.T) {
		md := &Metadata{Infohash: md5x.Digest("example1").Sum(nil), Description: "derp0"}
		require.Equal(t, "example.mp4", DescriptionFromPath(md, "example.mp4"))
	})

	t.Run("returns empty when path basename matches hex-encoded infohash", func(t *testing.T) {
		md := &Metadata{Infohash: md5x.Digest("example2").Sum(nil), Description: "derp1"}
		require.Equal(t, "", DescriptionFromPath(md, md5x.FormatHex(md5x.Digest("example2"))))
	})

	t.Run("returns empty when path basename matches description", func(t *testing.T) {
		md := &Metadata{Infohash: md5x.Digest("example3").Sum(nil), Description: "my-video.mp4"}
		require.Equal(t, "", DescriptionFromPath(md, "my-video.mp4"))
	})

	t.Run("uses basename of path with directory components", func(t *testing.T) {
		md := &Metadata{Infohash: md5x.Digest("example4").Sum(nil), Description: "derp4"}
		require.Equal(t, "movie.mkv", DescriptionFromPath(md, filepath.Join("some", "nested", "dir", "movie.mkv")))
	})
}
