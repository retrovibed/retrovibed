package rss_test

import (
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/rss"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("example 1 - basic rss feed", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		_, _, parsed, err := rss.Parse(ctx, testx.Read(testx.Fixture("parsing", "example.1.xml")))
		require.NoError(t, err)
		require.Equal(t, len(parsed), 50)
	})

	t.Run("example 2 - enclosures", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		_, _, parsed, err := rss.Parse(ctx, testx.Read(testx.Fixture("parsing", "example.2.xml")))
		require.NoError(t, err)
		require.Equal(t, len(parsed), 3)
		require.Equal(
			t,
			[]rss.Enclosure{
				{URL: "https://archlinux.org//releng/releases/2025.04.01/torrent/", Mimetype: "application/x-bittorrent", Length: 0x49b08000},
				{URL: "https://archlinux.org//releng/releases/2025.03.01/torrent/", Mimetype: "application/x-bittorrent", Length: 0x4a608000},
				{URL: "https://archlinux.org//releng/releases/2025.02.01/torrent/", Mimetype: "application/x-bittorrent", Length: 0x49b08000},
			},
			rss.FindEnclosureURLByMimetype(mimex.Bittorrent, parsed...),
		)
	})

	t.Run("example 3 - retrovibed extensions", func(t *testing.T) {
		_, channel, parsed, err := rss.Parse(t.Context(), testx.Read(testx.Fixture("parsing", "example.3.xml")))
		require.NoError(t, err)
		require.Equal(t, len(parsed), 1)
		require.EqualValues(t, 1440, channel.TTL)
		require.Equal(t, "en-us", channel.Language)
		require.Equal(t, time.Date(2025, time.July, 22, 13, 2, 40, 0, time.FixedZone("+0000", 0)), channel.LastBuildDate.Timestamp(time.Now()))
		require.Equal(t, "Retrovibed Media Database", channel.Title)
		require.Equal(t, "57869e82c2684ac4881bd32581f969db", channel.Retrovibed.Entropy)
		require.Equal(t, "application/vnd.retrovibed.media.archive", channel.Retrovibed.Mimetype)
		require.Equal(t, rss.FindEnclosureURLByMimetype(mimex.Bittorrent, parsed...), []rss.Enclosure{
			{URL: "magnet:?xt=urn:btih:8665727372B28B0263690B82928399516641A1B4&dn=retrovibed.media.metadata.00.gz", Mimetype: "application/x-bittorrent", Length: 0x4a208000},
		})
	})

	t.Run("example 4 - timestamps", func(t *testing.T) {
		digest, channel, parsed, err := rss.Parse(t.Context(), testx.Read(testx.Fixture("parsing", "example.4.xml")))
		require.NoError(t, err)
		require.Equal(t, len(parsed), 1)
		require.EqualValues(t, 30, channel.TTL)
		require.Equal(t, "en-CA", channel.Language)
		require.Equal(t, time.Date(2025, time.September, 7, 20, 43, 27, 0, time.FixedZone("", -10800)), channel.LastBuildDate.Timestamp(time.Now()))
		require.Equal(t, "FOSS Torrents - RSS Feed for Torrent Files", channel.Title)
		require.Equal(t, "", channel.Retrovibed.Entropy)
		require.Equal(t, "", channel.Retrovibed.Mimetype)
		require.Equal(t, "https://fosstorrents.com/direct-files/0ad-0.27.1-win32.exe.torrent", channel.Items[0].Link)
		require.Equal(t, "45a1655e-8204-b2a4-1150-8e859d0caa7e", md5x.FormatUUID(digest))
	})

	t.Run("example 5 - item Expires", func(t *testing.T) {
		_, channel, parsed, err := rss.Parse(t.Context(), testx.Read(testx.Fixture("parsing", "example.5.xml")))
		require.NoError(t, err)
		require.Equal(t, 2, len(parsed))
		require.EqualValues(t, 1440, channel.TTL)
		require.Equal(t, "Item Expires Test", channel.Title)
		require.Equal(t, time.Date(2025, time.July, 23, 13, 2, 40, 0, time.FixedZone("+0000", 0)), parsed[0].Expires)
		require.Equal(t, timex.Inf(), parsed[1].Expires)
	})
}
