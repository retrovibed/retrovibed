package community

import (
	"log"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestCommunityMetricInsert(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	now := time.Now()
	metric := CommunityMetric{
		CommunityID: uuid.Nil.String(),
		PeriodStart: now.Add(-24 * time.Hour),
		PeriodEnd:   now,
		Subscribers: 100,
	}

	require.NoError(t, CommunityMetricInsertWithDefaults(ctx, q, metric).Scan(&metric))
	require.NotEmpty(t, metric.ID)
	require.Equal(t, uint32(100), metric.Subscribers)
}

func TestCommunityMetricUpsert(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	now := time.Now()
	metric := CommunityMetric{
		CommunityID: uuid.Nil.String(),
		PeriodStart: now.Add(-24 * time.Hour),
		PeriodEnd:   now,
		Subscribers: 100,
	}

	require.NoError(t, CommunityMetricInsertWithDefaults(ctx, q, metric).Scan(&metric))
	require.Equal(t, uint32(100), metric.Subscribers)

	metric.Subscribers = 150
	require.NoError(t, CommunityMetricInsertWithDefaults(ctx, q, metric).Scan(&metric))
	require.Equal(t, uint32(150), metric.Subscribers)
}

func TestPublishedCASMetricInsert(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	now := time.Now()
	metric := PublishedCASMetric{
		PublishedContentID: uuid.Nil.String(),
		PeriodStart:        now.Add(-24 * time.Hour),
		PeriodEnd:          now,
		Archivers:          50,
		Bytes:              1024,
		Revenue:            2048,
	}

	require.NoError(t, PublishedCASMetricInsertWithDefaults(ctx, q, metric).Scan(&metric))
	require.NotEmpty(t, metric.ID)
	require.Equal(t, uint32(50), metric.Archivers)
	require.Equal(t, int64(1024), metric.Bytes)
	require.Equal(t, int64(2048), metric.Revenue)
}

func TestCommunityInsert(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	now := time.Now()
	state := Community{
		ID:         uuid.Nil.String(),
		LastSyncAt: now,
	}

	require.NoError(t, CommunityInsertWithDefaults(ctx, q, state).Scan(&state))
	require.Equal(t, uuid.Nil.String(), state.ID)
}

func TestCommunityUpsert(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	now := time.Now()
	state := Community{
		ID:         uuid.Nil.String(),
		LastSyncAt: now,
	}

	require.NoError(t, CommunityInsertWithDefaults(ctx, q, state).Scan(&state))
	originalID := state.ID

	state.LastSyncAt = now.Add(time.Hour)
	require.NoError(t, CommunityInsertWithDefaults(ctx, q, state).Scan(&state))
	require.Equal(t, originalID, state.ID)
}

func TestCommunityFindByID(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	now := time.Now()
	state := Community{
		ID:         uuid.Nil.String(),
		LastSyncAt: now,
	}

	require.NoError(t, CommunityInsertWithDefaults(ctx, q, state).Scan(&state))

	var found Community
	require.NoError(t, CommunityFindByID(ctx, q, uuid.Nil.String()).Scan(&found))
	require.Equal(t, state.ID, found.ID)
}

func TestPublishedCASMetricPerContentSearch(t *testing.T) {
	t.Run("returns one row per content with max archivers", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		now := time.Now()
		communityID := uuid.Must(uuid.NewV4()).String()

		pcA := NewPublishedContent(PublishedContent{
			ID:           uuid.Must(uuid.NewV4()).String(),
			CommunityID:  communityID,
			KnownMediaID: uuid.Must(uuid.NewV4()).String(),
			MagnetURI:    "magnet:?xt=urn:btih:contentA",
			LibraryID:    uuid.Must(uuid.NewV4()).String(),
		})
		require.NoError(t, PublishedContentInsertWithDefaults(ctx, q, pcA).Scan(&pcA))

		metricA1 := PublishedCASMetric{
			PublishedContentID: pcA.ID,
			PeriodStart:        now.AddDate(0, 0, -3),
			PeriodEnd:          now.AddDate(0, 0, -3),
			Archivers:          10,
		}
		require.NoError(t, PublishedCASMetricInsertWithDefaults(ctx, q, metricA1).Scan(&metricA1))

		metricA2 := PublishedCASMetric{
			PublishedContentID: pcA.ID,
			PeriodStart:        now.AddDate(0, 0, -2),
			PeriodEnd:          now.AddDate(0, 0, -2),
			Archivers:          15,
		}
		require.NoError(t, PublishedCASMetricInsertWithDefaults(ctx, q, metricA2).Scan(&metricA2))

		metricA3 := PublishedCASMetric{
			PublishedContentID: pcA.ID,
			PeriodStart:        now.AddDate(0, 0, -1),
			PeriodEnd:          now.AddDate(0, 0, -1),
			Archivers:          20,
		}
		require.NoError(t, PublishedCASMetricInsertWithDefaults(ctx, q, metricA3).Scan(&metricA3))

		periodStart := now.AddDate(0, 0, -7)
		periodEnd := now

		results, err := PublishedCASMetricPerContentSearch(ctx, q, communityID, periodStart, periodEnd)
		require.NoError(t, err)
		require.Len(t, results, 1, "should return one row per content, not one per daily metric")
		require.Equal(t, pcA.ID, results[0].PublishedContentID)
		require.Equal(t, uint32(20), results[0].Archivers, "should return max archivers (20), not sum (45)")
	})
}

func TestPublishedContentMultipleCommunities(t *testing.T) {
	t.Run("same library item can be published to multiple communities", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		libraryID := uuid.Must(uuid.NewV4()).String()
		knownMediaID := uuid.Must(uuid.NewV4()).String()
		communityA := uuid.Must(uuid.NewV4()).String()
		communityB := uuid.Must(uuid.NewV4()).String()

		pcA := NewPublishedContent(PublishedContent{
			CommunityID:  communityA,
			KnownMediaID: knownMediaID,
			MagnetURI:    "magnet:?xt=urn:btih:abc123",
			LibraryID:    libraryID,
		})
		require.NoError(t, PublishedContentInsertWithDefaults(ctx, q, pcA).Scan(&pcA))
		require.NotEmpty(t, pcA.ID)

		pcB := NewPublishedContent(PublishedContent{
			CommunityID:  communityB,
			KnownMediaID: knownMediaID,
			MagnetURI:    "magnet:?xt=urn:btih:abc123",
			LibraryID:    libraryID,
		})
		require.NoError(t, PublishedContentInsertWithDefaults(ctx, q, pcB).Scan(&pcB))
		require.NotEmpty(t, pcB.ID)
		require.NotEqual(t, pcA.ID, pcB.ID)

		var foundA PublishedContent
		require.NoError(t, PublishedContentFindByID(ctx, q, pcA.ID).Scan(&foundA))
		require.Equal(t, communityA, foundA.CommunityID)
		require.Equal(t, libraryID, foundA.LibraryID)

		var foundB PublishedContent
		require.NoError(t, PublishedContentFindByID(ctx, q, pcB.ID).Scan(&foundB))
		require.Equal(t, communityB, foundB.CommunityID)
		require.Equal(t, libraryID, foundB.LibraryID)
	})
}

func TestPublishedContentRepublish(t *testing.T) {
	t.Run("republishing same library item to same community updates magnet uri", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		libraryID := uuid.Must(uuid.NewV4()).String()
		knownMediaID := uuid.Must(uuid.NewV4()).String()
		communityID := uuid.Must(uuid.NewV4()).String()

		pc := NewPublishedContent(PublishedContent{
			CommunityID:  communityID,
			KnownMediaID: knownMediaID,
			MagnetURI:    "magnet:?xt=urn:btih:original",
			LibraryID:    libraryID,
		})
		require.NoError(t, PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))
		originalID := pc.ID
		require.NotEmpty(t, originalID)

		pc2 := NewPublishedContent(PublishedContent{
			CommunityID:  communityID,
			KnownMediaID: knownMediaID,
			MagnetURI:    "magnet:?xt=urn:btih:updated",
			LibraryID:    libraryID,
		})
		require.NoError(t, PublishedContentInsertWithDefaults(ctx, q, pc2).Scan(&pc2))
		require.Equal(t, originalID, pc2.ID)
		require.Equal(t, "magnet:?xt=urn:btih:updated", pc2.MagnetURI)

		var found PublishedContent
		require.NoError(t, PublishedContentFindByID(ctx, q, originalID).Scan(&found))
		require.Equal(t, "magnet:?xt=urn:btih:updated", found.MagnetURI)
	})

	t.Run("republishing updates known_media_id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		libraryID := uuid.Must(uuid.NewV4()).String()
		knownMediaIDOriginal := uuid.Must(uuid.NewV4()).String()
		knownMediaIDUpdated := uuid.Must(uuid.NewV4()).String()
		communityID := uuid.Must(uuid.NewV4()).String()

		pc := NewPublishedContent(PublishedContent{
			CommunityID:  communityID,
			KnownMediaID: knownMediaIDOriginal,
			MagnetURI:    "magnet:?xt=urn:btih:abc123",
			LibraryID:    libraryID,
		})
		require.NoError(t, PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))
		originalID := pc.ID

		pc2 := NewPublishedContent(PublishedContent{
			CommunityID:  communityID,
			KnownMediaID: knownMediaIDUpdated,
			MagnetURI:    "magnet:?xt=urn:btih:abc123",
			LibraryID:    libraryID,
		})
		require.NoError(t, PublishedContentInsertWithDefaults(ctx, q, pc2).Scan(&pc2))
		require.Equal(t, originalID, pc2.ID)

		var found PublishedContent
		require.NoError(t, PublishedContentFindByID(ctx, q, originalID).Scan(&found))
		require.Equal(t, knownMediaIDUpdated, found.KnownMediaID)
	})
}

func TestPublishedCASMetricAggregateSearch(t *testing.T) {
	t.Run("returns max archivers per content not sum of all metrics", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		now := time.Now()
		communityID := uuid.Must(uuid.NewV4()).String()

		pcA := NewPublishedContent(PublishedContent{
			ID:           uuid.Must(uuid.NewV4()).String(),
			CommunityID:  communityID,
			KnownMediaID: uuid.Must(uuid.NewV4()).String(),
			MagnetURI:    "magnet:?xt=urn:btih:contentA",
			LibraryID:    uuid.Must(uuid.NewV4()).String(),
		})
		require.NoError(t, PublishedContentInsertWithDefaults(ctx, q, pcA).Scan(&pcA))

		pcB := NewPublishedContent(PublishedContent{
			ID:           uuid.Must(uuid.NewV4()).String(),
			CommunityID:  communityID,
			KnownMediaID: uuid.Must(uuid.NewV4()).String(),
			MagnetURI:    "magnet:?xt=urn:btih:contentB",
			LibraryID:    uuid.Must(uuid.NewV4()).String(),
		})
		require.NoError(t, PublishedContentInsertWithDefaults(ctx, q, pcB).Scan(&pcB))

		metricA1 := PublishedCASMetric{
			PublishedContentID: pcA.ID,
			PeriodStart:        now.AddDate(0, 0, -3),
			PeriodEnd:          now.AddDate(0, 0, -3),
			Archivers:          10,
		}
		require.NoError(t, PublishedCASMetricInsertWithDefaults(ctx, q, metricA1).Scan(&metricA1))

		metricA2 := PublishedCASMetric{
			PublishedContentID: pcA.ID,
			PeriodStart:        now.AddDate(0, 0, -2),
			PeriodEnd:          now.AddDate(0, 0, -2),
			Archivers:          15,
		}
		require.NoError(t, PublishedCASMetricInsertWithDefaults(ctx, q, metricA2).Scan(&metricA2))

		metricA3 := PublishedCASMetric{
			PublishedContentID: pcA.ID,
			PeriodStart:        now.AddDate(0, 0, -1),
			PeriodEnd:          now.AddDate(0, 0, -1),
			Archivers:          20,
		}
		require.NoError(t, PublishedCASMetricInsertWithDefaults(ctx, q, metricA3).Scan(&metricA3))

		metricB := PublishedCASMetric{
			PublishedContentID: pcB.ID,
			PeriodStart:        now.AddDate(0, 0, -1),
			PeriodEnd:          now.AddDate(0, 0, -1),
			Archivers:          30,
		}
		require.NoError(t, PublishedCASMetricInsertWithDefaults(ctx, q, metricB).Scan(&metricB))

		periodStart := now.AddDate(0, 0, -7)
		periodEnd := now

		totalArchivers, err := PublishedCASMetricAggregateSearch(ctx, q, communityID, periodStart, periodEnd)
		require.NoError(t, err)
		require.Equal(t, int32(50), totalArchivers, "should be max(A)=20 + max(B)=30 = 50, not sum of all metrics")
	})
}
