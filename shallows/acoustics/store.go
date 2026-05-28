package acoustics

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// StoreFeatures persists a raw feature vector. Normalization happens at query
// time against in-memory RunningStats; LSH buckets and hyperplanes also live in
// memory and are regenerated from this table at daemon startup.
func StoreFeatures(ctx context.Context, q sqlx.Queryer, id uuid.UUID, vec FeatureVector, statsVersion uint32) error {
	_, err := q.ExecContext(ctx,
		`INSERT INTO audio_features (media_id, features, stats_version) VALUES ($1, $2, $3)
		 ON CONFLICT (media_id) DO UPDATE SET features = EXCLUDED.features, stats_version = EXCLUDED.stats_version, indexed_at = now()`,
		id, formatDoubleArray(vec[:]), statsVersion,
	)
	return errorsx.Wrap(err, "acoustics: insert features")
}

// Rebuild seeds hyperplanes, clears buckets, then streams audio_features once
// to populate RunningStats and the in-memory LSH with normalized vectors.
// Call at daemon startup and whenever the stats epoch shifts (cold start, 10× growth).
func Rebuild(ctx context.Context, q sqlx.Queryer, idx *Index, stats *RunningStats) error {
	vectors, ids, err := FetchAllVectors(ctx, q)
	if err != nil {
		return err
	}

	idx.InitHyperplanes(HyperplanesSeed)
	idx.ClearBuckets()

	*stats = RunningStats{}
	for _, v := range vectors {
		stats.Update(v)
	}

	for i, v := range vectors {
		idx.Insert(ids[i], stats.Normalize(v))
	}
	return nil
}

// FetchAllVectors returns every raw feature vector in audio_features.
func FetchAllVectors(ctx context.Context, q sqlx.Queryer) ([]FeatureVector, []uuid.UUID, error) {
	return fetchVectors(ctx, q, `SELECT media_id, features FROM audio_features`, nil)
}

// FetchCandidateVectors batch-fetches feature vectors for candidate media IDs from DuckDB.
func FetchCandidateVectors(ctx context.Context, q sqlx.Queryer, ids []uuid.UUID) ([]FeatureVector, []uuid.UUID, error) {
	if len(ids) == 0 {
		return nil, nil, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, uid := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = uid
	}

	return fetchVectors(ctx, q,
		fmt.Sprintf(`SELECT media_id, features FROM audio_features WHERE media_id IN (%s)`, strings.Join(placeholders, ",")),
		args,
	)
}

func fetchVectors(ctx context.Context, q sqlx.Queryer, query string, args []any) ([]FeatureVector, []uuid.UUID, error) {
	var (
		id      uuid.UUID
		vecStr  string
		parsed  []float64
		fv      FeatureVector
		vectors []FeatureVector
		ids     []uuid.UUID
	)

	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, errorsx.Wrap(err, "acoustics: fetch vectors")
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &vecStr)
		if err != nil {
			return nil, nil, err
		}
		parsed, err = parseDoubleArray(vecStr, VectorDim)
		if err != nil {
			return nil, nil, err
		}
		copy(fv[:], parsed)
		ids = append(ids, id)
		vectors = append(vectors, fv)
	}
	return vectors, ids, rows.Err()
}

// IndexedCount returns the number of tracks with the given stats version.
func IndexedCount(ctx context.Context, q sqlx.Queryer, statsVersion uint32) (int64, error) {
	var count int64
	err := q.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM audio_features WHERE stats_version = $1`, statsVersion,
	).Scan(&count)
	return count, err
}

// UnindexedMediaIDs returns audio media IDs without a corresponding audio_features entry.
func UnindexedMediaIDs(ctx context.Context, q sqlx.Queryer, limit int) ([]uuid.UUID, error) {
	var id uuid.UUID

	rows, err := q.QueryContext(ctx,
		`SELECT m.id FROM library_metadata m
		 LEFT JOIN audio_features af ON af.media_id = m.id
		 WHERE m.mimetype LIKE 'audio/%'
		   AND m.tombstoned_at = 'infinity'
		   AND af.media_id IS NULL
		 LIMIT $1`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func formatDoubleArray(vals []float64) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, v := range vals {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%g", v)
	}
	b.WriteByte(']')
	return b.String()
}

func parseDoubleArray(s string, expected int) ([]float64, error) {
	s = strings.Trim(s, "[]")
	if s == "" {
		return nil, nil
	}

	parts := strings.Split(s, ",")
	if len(parts) != expected {
		return nil, fmt.Errorf("acoustics: expected %d array values, got %d", expected, len(parts))
	}

	result := make([]float64, expected)
	for i, p := range parts {
		_, err := fmt.Sscanf(strings.TrimSpace(p), "%g", &result[i])
		if err != nil {
			return nil, errorsx.Wrap(err, fmt.Sprintf("acoustics: parse element %d", i))
		}
	}
	return result, nil
}
