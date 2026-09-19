package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestDiscoverRequestFromKnown(t *testing.T) {
	t.Run("derives defaults from known", func(t *testing.T) {
		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))

		r := ddisc.DiscoverRequestFromKnown(known)
		require.Equal(t, known.UID, r.KnownMediaID)
		require.Equal(t, known.Title, r.Query)
		require.Equal(t, ddisc.Category(known.Mimetype), r.Mimetypes)
		require.Equal(t, known.Adult, r.Adult)
		require.True(t, r.Public)
	})

	t.Run("options override derived defaults", func(t *testing.T) {
		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))

		r := ddisc.DiscoverRequestFromKnown(
			known,
			ddisc.DiscoverRequestOptionPublic(false),
			ddisc.DiscoverRequestOptionAdult(!known.Adult),
			ddisc.DiscoverRequestOptionQuery("override"),
			ddisc.DiscoverRequestOptionMimetypes(mimex.RetrovibedDiscoveryMusic),
		)
		require.Equal(t, known.UID, r.KnownMediaID)
		require.Equal(t, "override", r.Query)
		require.Equal(t, []string{mimex.RetrovibedDiscoveryMusic}, r.Mimetypes)
		require.Equal(t, !known.Adult, r.Adult)
		require.False(t, r.Public)
	})
}
