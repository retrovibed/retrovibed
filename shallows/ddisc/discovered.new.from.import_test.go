package ddisc_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/stretchr/testify/require"
)

func TestNewDiscoveredFromImport(t *testing.T) {
	t.Run("keys a magnet uri by its real infohash", func(t *testing.T) {
		magnetInfohash := []byte{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		}
		magnet := metainfo.NewMagnetFromInfohash(magnetInfohash).String()

		parsed, err := metainfo.ParseMagnetURI(magnet)
		require.NoError(t, err)
		expected := int160.FromByteArray(parsed.InfoHash)

		d := ddisc.NewDiscoveredFromImport(&ddiscapi.Import{Uri: magnet})

		require.Equal(t, expected.Bytes(), d.Infohash)
		require.Equal(t, torrentx.HashUID(&expected), d.ID)
	})

	t.Run("keys a non-magnet uri by a placeholder hash of the uri", func(t *testing.T) {
		uri := "https://tracker.example/download/1.torrent"
		d := ddisc.NewDiscoveredFromImport(&ddiscapi.Import{Uri: uri})

		require.Equal(t, md5x.FormatUUID(md5x.Digest(uri)), d.ID)
		require.Equal(t, int160.FromHashedBytes([]byte(uri)).Bytes(), d.Infohash)
	})

	t.Run("copies title, description, poster, health, bytes, and source", func(t *testing.T) {
		imp := &ddiscapi.Import{
			Uri:        "https://tracker.example/download/1.torrent",
			Title:      "the title",
			Overview:   "the overview",
			PosterPath: "http://example.com/poster.jpg",
			Health:     7,
			Bytes:      1024,
			Source:     "the source",
		}
		d := ddisc.NewDiscoveredFromImport(imp)

		require.Equal(t, imp.Uri, d.URI)
		require.Equal(t, "the title", d.Title)
		require.Equal(t, "the overview", d.Description)
		require.Equal(t, "http://example.com/poster.jpg", d.PosterURI)
		require.EqualValues(t, 7, d.Health)
		require.EqualValues(t, 1024, d.Bytes)
		require.Equal(t, "the source", d.Source)
	})

	t.Run("defaults known media id to the nil uuid", func(t *testing.T) {
		d := ddisc.NewDiscoveredFromImport(&ddiscapi.Import{Uri: "https://tracker.example/download/1.torrent"})
		require.Equal(t, uuid.Nil.String(), d.KnownMediaID)
	})

	t.Run("defaults mimetype to binary since the import carries none", func(t *testing.T) {
		d := ddisc.NewDiscoveredFromImport(&ddiscapi.Import{Uri: "https://tracker.example/download/1.torrent"})
		require.Equal(t, mimex.Binary, d.Mimetype)
		require.Equal(t, mimex.Application, d.Category)
	})

	t.Run("defaults released at to neg-infinity since the import carries no release date", func(t *testing.T) {
		d := ddisc.NewDiscoveredFromImport(&ddiscapi.Import{Uri: "https://tracker.example/download/1.torrent"})
		require.True(t, d.ReleasedAt.Equal(timex.NegInf()))
	})

	t.Run("expires on its own: unconfirmed infohash is only tentatively trusted", func(t *testing.T) {
		before := time.Now().Add(3 * time.Hour)
		d := ddisc.NewDiscoveredFromImport(&ddiscapi.Import{Uri: "https://tracker.example/download/1.torrent"})
		after := time.Now().Add(3 * time.Hour)

		require.False(t, d.TombstonedAt.Before(before))
		require.False(t, d.TombstonedAt.After(after))
	})

	t.Run("applies caller options on top of the defaults", func(t *testing.T) {
		d := ddisc.NewDiscoveredFromImport(
			&ddiscapi.Import{Uri: "https://tracker.example/download/1.torrent"},
			ddisc.DiscoveredOptionMimetype(mimex.Video),
		)
		require.Equal(t, mimex.Video, d.Mimetype)
	})
}
