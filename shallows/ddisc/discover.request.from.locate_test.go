package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoverRequestFromLocate(t *testing.T) {
	loc := ddisc.Locate{
		KnownMediaID: "known-media-id",
		Query:        "ubuntu",
		Mimetype:     mimex.Video,
		Adult:        true,
	}

	t.Run("derives defaults from locate", func(t *testing.T) {
		r := ddisc.DiscoverRequestFromLocate(loc)
		require.Equal(t, loc.KnownMediaID, r.KnownMediaID)
		require.Equal(t, loc.Query, r.Query)
		require.Equal(t, ddisc.Category(loc.Mimetype), r.Mimetypes)
		require.True(t, r.Adult)
		require.False(t, r.Public)
	})

	t.Run("options override derived defaults", func(t *testing.T) {
		r := ddisc.DiscoverRequestFromLocate(
			loc,
			ddisc.DiscoverRequestOptionPublic(true),
			ddisc.DiscoverRequestOptionAdult(false),
		)
		require.Equal(t, loc.KnownMediaID, r.KnownMediaID)
		require.False(t, r.Adult)
		require.True(t, r.Public)
	})
}
