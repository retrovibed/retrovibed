package acoustics

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"
	"sync"

	"github.com/gofrs/uuid/v5"
)

const (
	NumTables       = 20
	BitsPerHash     = 9
	BucketsPerTable = 1 << BitsPerHash // 512
	HyperplanesSeed = 42
)

type Index struct {
	buckets     [NumTables][BucketsPerTable][]uuid.UUID
	hyperplanes [NumTables][BitsPerHash][VectorDim]float64
	mu          sync.RWMutex
}

// InitHyperplanes generates random unit hyperplanes for all tables.
// Deterministic for a given seed.
func (idx *Index) InitHyperplanes(seed uint64) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	rng := rand.New(rand.NewPCG(seed, seed^0xdeadbeef))
	for t := range NumTables {
		for p := range BitsPerHash {
			norm := 0.0
			for d := range VectorDim {
				v := rng.NormFloat64()
				idx.hyperplanes[t][p][d] = v
				norm += v * v
			}
			norm = math.Sqrt(norm)
			for d := range VectorDim {
				idx.hyperplanes[t][p][d] /= norm
			}
		}
	}
}

// ClearBuckets empties every LSH bucket. Hyperplanes are preserved.
func (idx *Index) ClearBuckets() {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	for t := range NumTables {
		for b := range BucketsPerTable {
			idx.buckets[t][b] = nil
		}
	}
}

func (idx *Index) Hash(table int, vec FeatureVector) uint32 {
	h := uint32(0)
	for p := range BitsPerHash {
		dot := 0.0
		for d := range VectorDim {
			dot += vec[d] * idx.hyperplanes[table][p][d]
		}
		if dot > 0 {
			h |= 1 << p
		}
	}
	return h
}

func (idx *Index) Insert(id uuid.UUID, vec FeatureVector) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	for t := range NumTables {
		h := idx.Hash(t, vec)
		idx.buckets[t][h] = append(idx.buckets[t][h], id)
	}
}

// Candidates returns deduplicated media IDs from all matching LSH buckets,
// excluding IDs in the exclude set.
func (idx *Index) Candidates(vec FeatureVector, exclude map[uuid.UUID]struct{}) []uuid.UUID {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	seen := make(map[uuid.UUID]struct{}, 2048)
	for id := range exclude {
		seen[id] = struct{}{}
	}

	var result []uuid.UUID
	for t := range NumTables {
		for _, id := range idx.buckets[t][idx.Hash(t, vec)] {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			result = append(result, id)
		}
	}
	return result
}

type SimilarityResult struct {
	MediaID    uuid.UUID
	Similarity float64
}

func CosineSimilarity(a, b FeatureVector) float64 {
	dot, normA, normB := 0.0, 0.0, 0.0
	for i := range VectorDim {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if denom := math.Sqrt(normA) * math.Sqrt(normB); denom > 1e-12 {
		return dot / denom
	}
	return 0
}

// RankCandidates scores candidates against query, returns top-N above threshold
// sorted by similarity descending.
func RankCandidates(query FeatureVector, candidates []FeatureVector, ids []uuid.UUID, n int, threshold float64) []SimilarityResult {
	var results []SimilarityResult
	for i, cv := range candidates {
		if sim := CosineSimilarity(query, cv); sim >= threshold {
			results = append(results, SimilarityResult{MediaID: ids[i], Similarity: sim})
		}
	}

	slices.SortFunc(results, func(a, b SimilarityResult) int {
		return cmp.Compare(b.Similarity, a.Similarity) // descending
	})

	return results[:min(len(results), n)]
}

func (idx *Index) Remove(id uuid.UUID) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	for t := range NumTables {
		for b := range BucketsPerTable {
			if i := slices.Index(idx.buckets[t][b], id); i >= 0 {
				idx.buckets[t][b] = slices.Delete(idx.buckets[t][b], i, i+1)
			}
		}
	}
}
