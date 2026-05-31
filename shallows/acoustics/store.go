package acoustics

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// StoreFeatures persists a feature vector. The HNSW index updates automatically.
// HNSW does not support UPDATEs on the indexed column, so callers reindexing a
// track should DeleteFeatures first.
func StoreFeatures(ctx context.Context, q sqlx.Queryer, id uuid.UUID, vec FeatureVector, statsVersion uint32) error {
	a := AudioFeatures{
		MediaID:      id.String(),
		Features:     vec[:],
		StatsVersion: statsVersion,
	}
	return AudioFeaturesInsert(ctx, q, a).Scan(&a)
}

// DeleteFeatures removes a row from audio_features. The HNSW index updates automatically.
func DeleteFeatures(ctx context.Context, q sqlx.Queryer, id uuid.UUID) error {
	var a AudioFeatures
	return sqlx.IgnoreNoRows(AudioFeaturesDeleteByID(ctx, q, id.String()).Scan(&a))
}

// FetchFeatures retrieves a single feature vector by media_id.
func FetchFeatures(ctx context.Context, q sqlx.Queryer, id uuid.UUID) (FeatureVector, error) {
	var (
		a  AudioFeatures
		fv FeatureVector
	)
	if err := AudioFeaturesFindByID(ctx, q, id.String()).Scan(&a); err != nil {
		return fv, err
	}
	copy(fv[:], a.Features)
	return fv, nil
}

// SimilarMediaIDs returns up to n media IDs nearest to vec by cosine similarity,
// filtered by threshold (cosine sim >= threshold), excluding any IDs in `exclude`.
// Backed by the HNSW index on audio_features.features.
func SimilarMediaIDs(ctx context.Context, q sqlx.Queryer, vec FeatureVector, exclude []uuid.UUID, n int, threshold float64) ([]uuid.UUID, error) {
	excludeStrs := make([]string, len(exclude))
	for i, uid := range exclude {
		excludeStrs[i] = uid.String()
	}

	v := sqlx.Scan(AudioFeaturesSimilarByVec(ctx, q, vec[:], excludeStrs, n))

	var (
		ids    []uuid.UUID
		rowVec FeatureVector
	)
	for a := range v.Iter() {
		copy(rowVec[:], a.Features)
		if CosineSimilarity(vec, rowVec) < threshold {
			continue
		}
		id, err := uuid.FromString(a.MediaID)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, v.Err()
}

// IndexedCount returns the number of tracks with the given stats version.
func IndexedCount(ctx context.Context, q sqlx.Queryer, statsVersion uint32) (int64, error) {
	var count int64
	err := AudioFeaturesCountByVersion(ctx, q, statsVersion).Scan(&count)
	return count, err
}

// UnindexedMediaIDs returns audio media IDs without a corresponding audio_features entry.
func UnindexedMediaIDs(ctx context.Context, q sqlx.Queryer, limit int) ([]uuid.UUID, error) {
	v := sqlx.Scan(AudioFeaturesUnindexedMediaIDs(ctx, q, limit))

	var ids []uuid.UUID
	for mid := range v.Iter() {
		if id, err := uuid.FromString(mid); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, v.Err()
}
