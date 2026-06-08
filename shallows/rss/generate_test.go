package rss_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/rss"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("example 1 - 0 items", func(t *testing.T) {
		var (
			buf bytes.Buffer
		)
		path := testx.Fixture("generated", "example.1.xml")

		require.NoError(t, rss.Generator().Generate(&buf, rss.Channel{
			Title:         "Retrovibed Media Database",
			Link:          "https://media.community.retrovibe.space",
			Description:   "magnet links for media metadata archives for use by the retrovibed application",
			TTL:           1440,
			LastBuildDate: time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC),
			Language:      "en-us",
			Copyright:     "Retrovibed 2025",
		}, iterx.From[rss.Item]()))

		require.Equal(t, testx.ReadMD5(path), testx.IOMD5(bytes.NewReader(buf.Bytes())), "invalid generation:\n%s\n-----------------\n%s", buf.String(), testx.ReadString(path))
	})

	t.Run("example 2 - 1 item with magnet enclosure", func(t *testing.T) {
		var (
			buf bytes.Buffer
		)
		path := testx.Fixture("generated", "example.2.xml")

		require.NoError(t, rss.Generator().Generate(&buf, rss.Channel{
			Title:         "Retrovibed Media Database",
			Link:          "https://media.community.retrovibe.space",
			Description:   "magnet links for media metadata archives for use by the retrovibed application",
			TTL:           1440,
			LastBuildDate: time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC),
			Language:      "en-us",
			Copyright:     "Retrovibed 2025",
			Retrovibed: &rss.Retrovibed{
				Mimetype: "application/vnd.retrovibed.media.archive",
				Entropy:  "57869e82c2684ac4881bd32581f969db",
			},
		}, iterx.From(rss.Item{
			Guid:        "00000000-0000-0000-0000-000000000001",
			Title:       "Retrovibe Media Archive 00",
			Link:        "https://media.community.retrovibe.space/00000000-0000-0000-0000-000000000001",
			PublishDate: time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC),
			Enclosures: []rss.Enclosure{{
				Length:   1243643904,
				Mimetype: mimex.Bittorrent,
				URL:      "magnet:?xt=urn:btih:8665727372B28B0263690B82928399516641A1B4&dn=retrovibed.media.metadata.00.gz",
			}},
		})))

		require.Equal(t, testx.ReadMD5(path), testx.IOMD5(bytes.NewReader(buf.Bytes())), "invalid generation:\n%s\n-----------------\n%s", buf.String(), testx.ReadString(path))
	})

	t.Run("example 3 - expirations", func(t *testing.T) {
		var (
			buf bytes.Buffer
		)
		path := testx.Fixture("generated", "example.3.xml")

		require.NoError(t, rss.Generator().Generate(&buf, rss.Channel{
			Title:         "Retrovibed Media Database",
			Link:          "https://media.community.retrovibe.space",
			Description:   "magnet links for media metadata archives for use by the retrovibed application",
			TTL:           1440,
			LastBuildDate: time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC),
			Language:      "en-us",
			Copyright:     "Retrovibed 2025",
			Retrovibed: &rss.Retrovibed{
				Mimetype: "application/vnd.retrovibed.media.archive",
				Entropy:  "57869e82c2684ac4881bd32581f969db",
			},
		}, iterx.From(
			rss.Item{
				Guid:        "00000000-0000-0000-0000-000000000001",
				Title:       "Retrovibe Media Archive 00",
				Link:        "https://media.community.retrovibe.space/00000000-0000-0000-0000-000000000001",
				PublishDate: time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC),
				Enclosures: []rss.Enclosure{{
					Length:   1243643904,
					Mimetype: mimex.Bittorrent,
					URL:      "magnet:?xt=urn:btih:8665727372B28B0263690B82928399516641A1B4&dn=retrovibed.media.metadata.00.gz",
				}},
			},
			rss.Item{
				Guid:        "00000000-0000-0000-0000-000000000001",
				Title:       "Retrovibe Media Archive 00",
				Link:        "https://media.community.retrovibe.space/00000000-0000-0000-0000-000000000001",
				PublishDate: time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC),
				Expires:     time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC),
				Enclosures: []rss.Enclosure{{
					Length:   1243643904,
					Mimetype: mimex.Bittorrent,
					URL:      "magnet:?xt=urn:btih:8665727372B28B0263690B82928399516641A1B4&dn=retrovibed.media.metadata.00.gz",
				}},
			},
			rss.Item{
				Guid:        "00000000-0000-0000-0000-000000000001",
				Title:       "Retrovibe Media Archive 00",
				Link:        "https://media.community.retrovibe.space/00000000-0000-0000-0000-000000000001",
				PublishDate: time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC),
				Expires:     timex.Inf(),
				Enclosures: []rss.Enclosure{{
					Length:   1243643904,
					Mimetype: mimex.Bittorrent,
					URL:      "magnet:?xt=urn:btih:8665727372B28B0263690B82928399516641A1B4&dn=retrovibed.media.metadata.00.gz",
				}},
			},
			rss.Item{
				Guid:        "00000000-0000-0000-0000-000000000001",
				Title:       "Retrovibe Media Archive 00",
				Link:        "https://media.community.retrovibe.space/00000000-0000-0000-0000-000000000001",
				PublishDate: time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC),
				Expires:     timex.NegInf(),
				Enclosures: []rss.Enclosure{{
					Length:   1243643904,
					Mimetype: mimex.Bittorrent,
					URL:      "magnet:?xt=urn:btih:8665727372B28B0263690B82928399516641A1B4&dn=retrovibed.media.metadata.00.gz",
				}},
			},
		)))

		require.Equal(t, testx.ReadMD5(path), testx.IOMD5(bytes.NewReader(buf.Bytes())), "invalid generation:\n%s\n-----------------\n%s", buf.String(), testx.ReadString(path))
	})

	t.Run("multiple items with different mimetypes", func(t *testing.T) {
		var (
			buf bytes.Buffer
		)
		ts := time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC)

		err := rss.Generator().Generate(&buf, rss.Channel{
			Title:         "Test Community",
			Link:          "https://test.community.retrovibe.space",
			Description:   "Test community description",
			TTL:           1440,
			LastBuildDate: ts,
			Language:      "en",
		}, iterx.From(
			rss.Item{
				Guid:        "item-001",
				Title:       "Video Content",
				Link:        "magnet:?xt=urn:btih:feedtest1",
				PublishDate: ts,
				Description: "A video file",
				Enclosures: []rss.Enclosure{{
					URL:      "magnet:?xt=urn:btih:feedtest1",
					Mimetype: "video/mp4",
					Length:   1024000,
				}},
			},
			rss.Item{
				Guid:        "item-002",
				Title:       "Audio Content",
				Link:        "magnet:?xt=urn:btih:feedtest2",
				PublishDate: ts,
				Description: "An audio file",
				Enclosures: []rss.Enclosure{{
					URL:      "magnet:?xt=urn:btih:feedtest2",
					Mimetype: "audio/mpeg",
					Length:   512000,
				}},
			},
			rss.Item{
				Guid:        "item-003",
				Title:       "Matroska Video",
				Link:        "magnet:?xt=urn:btih:feedtest3",
				PublishDate: ts,
				Enclosures: []rss.Enclosure{{
					URL:      "magnet:?xt=urn:btih:feedtest3",
					Mimetype: "video/x-matroska",
					Length:   2048000,
				}},
			},
		))
		require.NoError(t, err)

		content := buf.String()
		require.Contains(t, content, "<rss")
		require.Contains(t, content, "Test Community")
		require.Contains(t, content, "item-001")
		require.Contains(t, content, "item-002")
		require.Contains(t, content, "item-003")
		require.Contains(t, content, "Video Content")
		require.Contains(t, content, "Audio Content")
		require.Contains(t, content, "Matroska Video")
		require.Contains(t, content, `type="video/mp4"`)
		require.Contains(t, content, `type="audio/mpeg"`)
		require.Contains(t, content, `type="video/x-matroska"`)
		require.Contains(t, content, "magnet:?xt=urn:btih:feedtest1")
		require.Contains(t, content, "magnet:?xt=urn:btih:feedtest2")
		require.Contains(t, content, "magnet:?xt=urn:btih:feedtest3")
	})

	t.Run("item with Expires", func(t *testing.T) {
		var (
			buf bytes.Buffer
		)
		ts := time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC)
		expires := time.Date(2025, time.July, 23, 13, 02, 40, 0, time.UTC)

		err := rss.Generator().Generate(&buf, rss.Channel{
			Title:         "Expires Test",
			Link:          "https://test.retrovibe.space",
			Description:   "Test item expires",
			TTL:           1440,
			LastBuildDate: ts,
			Language:      "en",
		}, iterx.From(rss.Item{
			Guid:        "item-expires-001",
			Title:       "Item with Expires",
			Link:        "https://test.retrovibe.space/item-expires-001",
			PublishDate: ts,
			Expires:     expires,
		}))
		require.NoError(t, err)

		content := buf.String()
		require.Contains(t, content, "item-expires-001")
		require.Contains(t, content, "<retrovibed:expires>Wed, 23 Jul 2025 13:02:40 +0000</retrovibed:expires>")
	})

	t.Run("feed content changes with new items", func(t *testing.T) {
		var (
			buf1 bytes.Buffer
			buf2 bytes.Buffer
		)
		ts := time.Date(2025, time.July, 22, 13, 02, 40, 0, time.UTC)

		channel := rss.Channel{
			Title:         "Republish Test",
			Link:          "https://republish.community.retrovibe.space",
			Description:   "Test republishing",
			TTL:           1440,
			LastBuildDate: ts,
			Language:      "en",
		}

		err := rss.Generator().Generate(&buf1, channel, iterx.From(
			rss.Item{
				Guid:        "initial-001",
				Title:       "Initial Content",
				Link:        "magnet:?xt=urn:btih:initial",
				PublishDate: ts,
			},
		))
		require.NoError(t, err)
		md5First := testx.IOMD5(bytes.NewReader(buf1.Bytes()))

		err = rss.Generator().Generate(&buf2, channel, iterx.From(
			rss.Item{
				Guid:        "initial-001",
				Title:       "Initial Content",
				Link:        "magnet:?xt=urn:btih:initial",
				PublishDate: ts,
			},
			rss.Item{
				Guid:        "new-001",
				Title:       "New Content",
				Link:        "magnet:?xt=urn:btih:newcontent",
				PublishDate: ts,
			},
		))
		require.NoError(t, err)
		md5Second := testx.IOMD5(bytes.NewReader(buf2.Bytes()))

		require.NotEqual(t, md5First, md5Second, "feed md5 should change after adding content")
		require.Contains(t, buf2.String(), "Initial Content")
		require.Contains(t, buf2.String(), "New Content")
	})
}
