package tracking

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/stretchr/testify/require"
)

func TestResetDescription(t *testing.T) {
	t.Run("joins existing description with description derived from path", func(t *testing.T) {
		md := &Metadata{Infohash: md5x.Digest("example1").Sum(nil), Description: "derp0"}
		o, desc, auto := GenerateDescription("movie.mkv", md)
		require.Equal(t, "derp0 movie.mkv", o)
		require.Equal(t, "derp0 movie.mkv", desc)
		require.Equal(t, "derp0 movie mkv", auto)
	})

	t.Run("trims surrounding whitespace from the joined description", func(t *testing.T) {
		md := &Metadata{Infohash: md5x.Digest("example2").Sum(nil), Description: "  derp1  "}
		o, desc, auto := GenerateDescription("my-video.mp4", md)
		require.Equal(t, "  derp1   my-video.mp4", o)
		require.Equal(t, "derp1 my-video.mp4", desc)
		require.Equal(t, "derp1 my video mp4", auto)
	})

	t.Run("omits path derived description when path matches existing description", func(t *testing.T) {
		md := &Metadata{Infohash: md5x.Digest("example3").Sum(nil), Description: "my-video.mp4"}
		o, desc, auto := GenerateDescription("my-video.mp4", md)
		require.Equal(t, "my-video.mp4 ", o)
		require.Equal(t, "my-video.mp4", desc)
		require.Equal(t, "my video mp4", auto)
	})

	t.Run("omits path derived description when path basename matches infohash", func(t *testing.T) {
		digest := md5x.Digest("example4")
		md := &Metadata{Infohash: digest.Sum(nil), Description: "derp4"}
		o, desc, auto := GenerateDescription(md5x.FormatHex(digest), md)
		require.Equal(t, "derp4 ", o)
		require.Equal(t, "derp4", desc)
		require.Equal(t, "derp4", auto)
	})

	t.Run("uses path derived description alone when description is empty", func(t *testing.T) {
		md := &Metadata{Infohash: md5x.Digest("example5").Sum(nil), Description: ""}
		o, desc, auto := GenerateDescription("movie.mkv", md)
		require.Equal(t, " movie.mkv", o)
		require.Equal(t, "movie.mkv", desc)
		require.Equal(t, "movie mkv", auto)
	})
}
