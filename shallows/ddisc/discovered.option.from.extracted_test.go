package ddisc_test

import (
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionFromExtracted(t *testing.T) {
	musicDate := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	videoDate := time.Date(2002, 2, 2, 0, 0, 0, 0, time.UTC)

	t.Run("prefers music over video when both are present", func(t *testing.T) {
		ex := ddisc.Extracted{
			Music: &ddisc.MetadataAudio{Metadata: ddisc.Metadata{
				Title: "music title", Subtitle: "music subtitle", Collation: "music collation", Mimetype: mimex.Audio, Date: musicDate,
			}},
			Video: &ddisc.MetadataVideo{Metadata: ddisc.Metadata{
				Title: "video title", Subtitle: "video subtitle", Collation: "video collation", Mimetype: mimex.Video, Date: videoDate,
			}},
		}

		d := ddisc.Discovered{}
		ddisc.DiscoveredOptionFromExtracted(ex)(&d)

		require.Equal(t, "music title", d.Title)
		require.Equal(t, "music subtitle", d.Description)
		require.Equal(t, "music collation", d.Collation)
		require.Equal(t, mimex.Audio, d.Mimetype)
		require.Equal(t, mimex.Category(mimex.Audio), d.Category)
		require.True(t, musicDate.Equal(d.ReleasedAt))
	})

	t.Run("falls back to video when music is absent", func(t *testing.T) {
		ex := ddisc.Extracted{
			Video: &ddisc.MetadataVideo{Metadata: ddisc.Metadata{
				Title: "video title", Subtitle: "video subtitle", Collation: "video collation", Mimetype: mimex.Video, Date: videoDate,
			}},
		}

		d := ddisc.Discovered{}
		ddisc.DiscoveredOptionFromExtracted(ex)(&d)

		require.Equal(t, "video title", d.Title)
		require.Equal(t, "video subtitle", d.Description)
		require.Equal(t, "video collation", d.Collation)
		require.Equal(t, mimex.Video, d.Mimetype)
		require.Equal(t, mimex.Category(mimex.Video), d.Category)
		require.True(t, videoDate.Equal(d.ReleasedAt))
	})

	t.Run("keeps existing discovered fields when neither music nor video is present", func(t *testing.T) {
		existingReleased := time.Date(2000, 6, 6, 0, 0, 0, 0, time.UTC)
		d := ddisc.Discovered{
			Title:       "existing title",
			Description: "existing description",
			Collation:   "existing collation",
			Mimetype:    "video/mp4",
			ReleasedAt:  existingReleased,
		}
		ddisc.DiscoveredOptionFromExtracted(ddisc.Extracted{})(&d)

		require.Equal(t, "existing title", d.Title)
		require.Equal(t, "existing description", d.Description)
		require.Equal(t, "existing collation", d.Collation)
		require.Equal(t, "video/mp4", d.Mimetype)
		require.True(t, existingReleased.Equal(d.ReleasedAt))
	})
}
