package ddisc_test

import (
	"testing"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestNewDiscoveredFromKnown(t *testing.T) {
	t.Run("derives id from the infohash and known media id", func(t *testing.T) {
		id := int160.Random()
		known := library.Known{UID: "known-id"}
		d := ddisc.NewDiscoveredFromKnown(id, known)

		require.Equal(t, md5x.FormatUUID(md5x.Digest(id.Bytes(), []byte(known.UID))), d.ID)
		require.Equal(t, id.Bytes(), d.Infohash)
		require.Equal(t, "known-id", d.KnownMediaID)
	})

	t.Run("copies title, description, released, and adult from known", func(t *testing.T) {
		id := int160.Random()
		released := time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)
		known := library.Known{
			Title:    "the title",
			Overview: "the overview",
			Released: released,
			Adult:    true,
		}
		d := ddisc.NewDiscoveredFromKnown(id, known)

		require.Equal(t, "the title", d.Title)
		require.Equal(t, "the overview", d.Description)
		require.True(t, released.Equal(d.ReleasedAt))
		require.True(t, d.Adult)
	})

	t.Run("defaults released at to neg-infinity when known has none", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscoveredFromKnown(id, library.Known{})

		require.True(t, d.ReleasedAt.Equal(timex.NegInf()))
	})

	t.Run("falls back to binary mimetype when known has none", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscoveredFromKnown(id, library.Known{})

		require.Equal(t, mimex.Binary, d.Mimetype)
		require.Equal(t, mimex.Application, d.Category)
	})

	t.Run("uses known's mimetype and derives its category when present", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscoveredFromKnown(id, library.Known{Mimetype: mimex.Video})

		require.Equal(t, mimex.Video, d.Mimetype)
		require.Equal(t, mimex.Category(mimex.Video), d.Category)
	})

	t.Run("never expires on its own: built from a confirmed infohash", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscoveredFromKnown(id, library.Known{})

		require.True(t, d.TombstonedAt.Equal(timex.Inf()))
	})

	t.Run("sets source to the media archive", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscoveredFromKnown(id, library.Known{})

		require.Equal(t, "retrovibed.media.archive", d.Source)
	})

	t.Run("applies caller options on top of the defaults", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscoveredFromKnown(id, library.Known{}, ddisc.DiscoveredOptionTitle("override"))

		require.Equal(t, "override", d.Title)
	})
}
