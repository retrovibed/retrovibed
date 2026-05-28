package acoustics

import (
	"math"
	"slices"
	"testing"

	"github.com/gofrs/uuid/v5"
)

func TestIndex(t *testing.T) {
	newVec := func(fill float64) FeatureVector {
		var fv FeatureVector
		for i := range fv {
			fv[i] = fill + float64(i)*0.01
		}
		return fv
	}

	t.Run("insert and candidates", func(t *testing.T) {
		var idx Index
		idx.InitHyperplanes(42)

		id1 := uuid.Must(uuid.NewV4())
		vec1 := newVec(1.0)
		idx.Insert(id1, vec1)

		if !slices.Contains(idx.Candidates(vec1, nil), id1) {
			t.Fatal("inserted ID not found in candidates for same vector")
		}
	})

	t.Run("similar vectors share buckets", func(t *testing.T) {
		var idx Index
		idx.InitHyperplanes(42)

		id1 := uuid.Must(uuid.NewV4())
		id2 := uuid.Must(uuid.NewV4())

		idx.Insert(id1, newVec(1.0))
		idx.Insert(id2, newVec(1.001))

		if !slices.Contains(idx.Candidates(newVec(1.0), nil), id2) {
			t.Fatal("similar vector not found in candidates")
		}
	})

	t.Run("exclude filters", func(t *testing.T) {
		var idx Index
		idx.InitHyperplanes(42)

		id1 := uuid.Must(uuid.NewV4())
		id2 := uuid.Must(uuid.NewV4())

		vec := newVec(1.0)
		idx.Insert(id1, vec)
		idx.Insert(id2, vec)

		exclude := map[uuid.UUID]struct{}{id1: {}}
		if slices.Contains(idx.Candidates(vec, exclude), id1) {
			t.Fatal("excluded ID appeared in candidates")
		}
	})

	t.Run("remove", func(t *testing.T) {
		var idx Index
		idx.InitHyperplanes(42)

		id1 := uuid.Must(uuid.NewV4())
		vec := newVec(1.0)
		idx.Insert(id1, vec)
		idx.Remove(id1)

		if slices.Contains(idx.Candidates(vec, nil), id1) {
			t.Fatal("removed ID still in candidates")
		}
	})

	t.Run("deterministic hyperplanes", func(t *testing.T) {
		var a, b Index
		a.InitHyperplanes(99)
		b.InitHyperplanes(99)

		vec := newVec(2.5)
		for tab := range NumTables {
			ha := a.Hash(tab, vec)
			hb := b.Hash(tab, vec)
			if ha != hb {
				t.Fatalf("table %d: hash mismatch for same seed", tab)
			}
		}
	})
}

func TestCosineSimilarity(t *testing.T) {
	t.Run("identical vectors", func(t *testing.T) {
		var a FeatureVector
		for i := range a {
			a[i] = float64(i + 1)
		}
		sim := CosineSimilarity(a, a)
		if math.Abs(sim-1.0) > 1e-10 {
			t.Fatalf("expected 1.0, got %f", sim)
		}
	})

	t.Run("orthogonal vectors", func(t *testing.T) {
		var a, b FeatureVector
		a[0] = 1.0
		b[1] = 1.0
		sim := CosineSimilarity(a, b)
		if math.Abs(sim) > 1e-10 {
			t.Fatalf("expected 0, got %f", sim)
		}
	})

	t.Run("opposite vectors", func(t *testing.T) {
		var a, b FeatureVector
		for i := range a {
			a[i] = float64(i + 1)
			b[i] = -float64(i + 1)
		}
		sim := CosineSimilarity(a, b)
		if math.Abs(sim-(-1.0)) > 1e-10 {
			t.Fatalf("expected -1.0, got %f", sim)
		}
	})

	t.Run("zero vector", func(t *testing.T) {
		var a, b FeatureVector
		a[0] = 1.0
		sim := CosineSimilarity(a, b)
		if sim != 0 {
			t.Fatalf("expected 0, got %f", sim)
		}
	})
}

func TestRankCandidates(t *testing.T) {
	t.Run("returns top N above threshold", func(t *testing.T) {
		var query FeatureVector
		for i := range query {
			query[i] = 1.0
		}

		ids := make([]uuid.UUID, 5)
		candidates := make([]FeatureVector, 5)
		for i := range 5 {
			ids[i] = uuid.Must(uuid.NewV4())
			for d := range VectorDim {
				candidates[i][d] = 1.0 + float64(i)*0.1*float64(d)/float64(VectorDim)
			}
		}

		results := RankCandidates(query, candidates, ids, 3, 0.0)

		if len(results) > 3 {
			t.Fatalf("expected at most 3 results, got %d", len(results))
		}

		for i := 1; i < len(results); i++ {
			if results[i].Similarity > results[i-1].Similarity {
				t.Fatal("results not sorted descending")
			}
		}
	})

	t.Run("threshold filters", func(t *testing.T) {
		var query FeatureVector
		query[0] = 1.0

		var dissimilar FeatureVector
		dissimilar[1] = 1.0 // orthogonal

		ids := []uuid.UUID{uuid.Must(uuid.NewV4())}
		candidates := []FeatureVector{dissimilar}

		results := RankCandidates(query, candidates, ids, 10, 0.5)
		if len(results) != 0 {
			t.Fatalf("expected 0 results above threshold, got %d", len(results))
		}
	})
}
