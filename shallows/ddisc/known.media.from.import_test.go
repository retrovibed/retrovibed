package ddisc_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestKnownMediaFromImport(t *testing.T) {
	t.Run("maps every field", func(t *testing.T) {
		kid := uuid.Must(uuid.NewV4())
		imp := &ddiscapi.Import{
			Title:      "Ubuntu Documentary",
			Overview:   "a documentary about ubuntu",
			Popularity: 4.2,
			PosterPath: "/ubuntu.jpg",
			Source:     "unit3d",
		}

		known := ddisc.KnownMediaFromImport(kid, mimex.RetrovibedDiscoveryMovies, imp)
		require.Equal(t, kid.String(), known.UID)
		require.Equal(t, "Ubuntu Documentary", known.Title)
		require.Equal(t, "a documentary about ubuntu", known.Overview)
		require.Equal(t, 4.2, known.Popularity)
		require.Equal(t, "/ubuntu.jpg", known.PosterPath)
		require.Equal(t, "unit3d", known.Source)
		require.Equal(t, mimex.Video, known.Mimetype)
		require.Equal(t, "Ubuntu Documentary", known.AutoDescription)
		require.NotEmpty(t, known.Md5)
	})

	t.Run("maps cleanly without a title or overview", func(t *testing.T) {
		kid := uuid.Must(uuid.NewV4())
		imp := &ddiscapi.Import{Source: "unit3d"}

		known := ddisc.KnownMediaFromImport(kid, mimex.RetrovibedDiscoveryMovies, imp)
		require.Equal(t, kid.String(), known.UID)
		require.Equal(t, "", known.Title)
		require.Equal(t, "", known.Overview)
		require.NotEmpty(t, known.Md5)
	})

	t.Run("distinct kids never collide on md5, even with identical blank title and overview", func(t *testing.T) {
		imp := &ddiscapi.Import{Source: "unit3d"}

		a := ddisc.KnownMediaFromImport(uuid.Must(uuid.NewV4()), mimex.RetrovibedDiscoveryMovies, imp)
		b := ddisc.KnownMediaFromImport(uuid.Must(uuid.NewV4()), mimex.RetrovibedDiscoveryMovies, imp)
		require.NotEqual(t, a.UID, b.UID)
		require.NotEqual(t, a.Md5, b.Md5)
	})
}
