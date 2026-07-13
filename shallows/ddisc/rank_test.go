package ddisc_test

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestRankAndSelectPicksBestCandidate(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	kid := uuid.Must(uuid.NewV4()).String()

	insert := func(title string) ddisc.Discovered {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id,
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

	got, err := ddisc.RankAndSelect(t.Context(), q, ddisc.DefaultPolicy(), kid)
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
	d := ddisc.NewDiscovered(&id,
		ddisc.DiscoveredOptionKnownMedia(kid),
		ddisc.DiscoveredOptionMimetype(mimex.Binary),
	)
	d.Title = "Some.Movie.2024.HDCAM.x264"
	require.NoError(t, ddisc.DiscoveredInsertWithDefaults(t.Context(), q, d).Scan(&d))

	_, err := ddisc.RankAndSelect(t.Context(), q, ddisc.DefaultPolicy(), kid)
	require.True(t, errors.Is(err, ddisc.ErrNoCandidate))
}

func TestRankAndSelectNoCandidateWhenEmpty(t *testing.T) {
	q := sqltestx.Metadatabase(t)

	_, err := ddisc.RankAndSelect(t.Context(), q, ddisc.DefaultPolicy(), uuid.Must(uuid.NewV4()).String())
	require.True(t, errors.Is(err, ddisc.ErrNoCandidate))
}
