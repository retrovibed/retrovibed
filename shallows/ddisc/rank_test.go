package ddisc_test

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestRankAndSelectPicksBestCandidate(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	kid := uuid.Must(uuid.NewV4()).String()

	insert := func(title string) ddisc.Discovered {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionAutoMagnet,
			ddisc.DiscoveredOptionKnownMedia(kid),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
		)
		d.Title = title
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(t.Context(), q, d).Scan(&d))
		return d
	}

	insert("Some.Movie.2024.HDCAM.x264")
	insert("Some.Movie.2024.720p.HDTV.x264")
	best := insert("Some.Movie.2024.2160p.BluRay.REMUX.HDR10")

	loc := ddisc.Locate{KnownMediaID: kid, Query: "some movie"}
	got, err := ddisc.RankAndSelect(t.Context(), q, ddisc.DefaultPolicy(), loc)
	require.NoError(t, err)
	require.Equal(t, best.ID, got.ID)
	require.Empty(t, got.PolicyRejection)

	var persisted ddisc.Discovered
	require.NoError(t, ddisc.DiscoveredFindByID(t.Context(), q, best.ID).Scan(&persisted))
	require.Equal(t, got.PolicyRank, persisted.PolicyRank)
	require.Equal(t, got.Health, persisted.Health)
}

func TestRankAndSelectNoCandidateWhenAllRejected(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	kid := uuid.Must(uuid.NewV4()).String()

	id := int160.Random()
	d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionAutoMagnet,
		ddisc.DiscoveredOptionKnownMedia(kid),
		ddisc.DiscoveredOptionMimetype(mimex.Binary),
	)
	d.Title = "Some.Movie.2024.HDCAM.x264"
	require.NoError(t, ddisc.DiscoveredInsertWithDefaults(t.Context(), q, d).Scan(&d))

	loc := ddisc.Locate{KnownMediaID: kid, Query: "some movie"}
	_, err := ddisc.RankAndSelect(t.Context(), q, ddisc.DefaultPolicy(), loc)
	require.True(t, errors.Is(err, ddisc.ErrNoCandidate))
}

func TestRankAndSelectNoCandidateWhenEmpty(t *testing.T) {
	q := sqltestx.Metadatabase(t)

	loc := ddisc.Locate{KnownMediaID: uuid.Must(uuid.NewV4()).String()}
	_, err := ddisc.RankAndSelect(t.Context(), q, ddisc.DefaultPolicy(), loc)
	require.True(t, errors.Is(err, ddisc.ErrNoCandidate))
}

func TestRankAndSelectIgnoresCandidatesNotMatchingQuery(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	kid := uuid.Must(uuid.NewV4()).String()

	insert := func(title string) ddisc.Discovered {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionAutoMagnet,
			ddisc.DiscoveredOptionKnownMedia(kid),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
		)
		d.Title = title
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(t.Context(), q, d).Scan(&d))
		return d
	}

	// higher quality, but the title has nothing to do with the query.
	insert("Nirvana.Live.At.The.Paramount.2011.1080p.BluRay")
	// lower quality, but the title matches the query.
	matching := insert("Nirvana.In.Utero.720p.HDTV")

	loc := ddisc.Locate{KnownMediaID: kid, Query: "nirvana utero"}
	got, err := ddisc.RankAndSelect(t.Context(), q, ddisc.DefaultPolicy(), loc)
	require.NoError(t, err)
	require.Equal(t, matching.ID, got.ID)
}

func TestRankAndSelectScopesNilKnownMediaIDByQuery(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	nilkid := uuid.Nil.String()

	insert := func(title string) ddisc.Discovered {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionAutoMagnet,
			ddisc.DiscoveredOptionKnownMedia(nilkid),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
		)
		d.Title = title
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(t.Context(), q, d).Scan(&d))
		return d
	}

	// two unrelated free-text searches both land on known_media_id = Nil;
	// without the title filter either could select the other's candidate.
	alien := insert("Alien.1979.2160p.BluRay.REMUX.HDR10")
	nirvana := insert("Nirvana.In.Utero.720p.HDTV")

	alienLoc := ddisc.Locate{KnownMediaID: nilkid, Query: "alien"}
	got, err := ddisc.RankAndSelect(t.Context(), q, ddisc.DefaultPolicy(), alienLoc)
	require.NoError(t, err)
	require.Equal(t, alien.ID, got.ID)

	nirvanaLoc := ddisc.Locate{KnownMediaID: nilkid, Query: "nirvana utero"}
	got, err = ddisc.RankAndSelect(t.Context(), q, ddisc.DefaultPolicy(), nirvanaLoc)
	require.NoError(t, err)
	require.Equal(t, nirvana.ID, got.ID)
}

func TestRankAndSelectNoCandidateWhenNoneMatchQuery(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	kid := uuid.Must(uuid.NewV4()).String()

	id := int160.Random()
	d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionAutoMagnet,
		ddisc.DiscoveredOptionKnownMedia(kid),
		ddisc.DiscoveredOptionMimetype(mimex.Binary),
	)
	d.Title = "Nirvana.Live.At.The.Paramount.2011.1080p.BluRay"
	require.NoError(t, ddisc.DiscoveredInsertWithDefaults(t.Context(), q, d).Scan(&d))

	loc := ddisc.Locate{KnownMediaID: kid, Query: "nirvana utero"}
	_, err := ddisc.RankAndSelect(t.Context(), q, ddisc.DefaultPolicy(), loc)
	require.True(t, errors.Is(err, ddisc.ErrNoCandidate))
}
