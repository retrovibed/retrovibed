package acoustics

import (
	"math"
	"testing"
)

func TestFormatDoubleArray(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := formatDoubleArray([]float64{})
		if got != "[]" {
			t.Fatalf("expected [], got %s", got)
		}
	})

	t.Run("single value", func(t *testing.T) {
		got := formatDoubleArray([]float64{3.14})
		if got != "[3.14]" {
			t.Fatalf("expected [3.14], got %s", got)
		}
	})

	t.Run("multiple values", func(t *testing.T) {
		got := formatDoubleArray([]float64{1, 2.5, -3})
		if got != "[1,2.5,-3]" {
			t.Fatalf("expected [1,2.5,-3], got %s", got)
		}
	})

	t.Run("roundtrip", func(t *testing.T) {
		var original [VectorDim]float64
		for i := range original {
			original[i] = float64(i)*0.1 - 5.0
		}

		s := formatDoubleArray(original[:])
		parsed, err := parseDoubleArray(s, VectorDim)
		if err != nil {
			t.Fatal(err)
		}

		for i, v := range parsed {
			if math.Abs(v-original[i]) > 1e-10 {
				t.Fatalf("dim %d: expected %f, got %f", i, original[i], v)
			}
		}
	})
}

func TestParseDoubleArray(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		vals, err := parseDoubleArray("[1,2,3]", 3)
		if err != nil {
			t.Fatal(err)
		}
		if vals[0] != 1 || vals[1] != 2 || vals[2] != 3 {
			t.Fatalf("unexpected values: %v", vals)
		}
	})

	t.Run("wrong count", func(t *testing.T) {
		_, err := parseDoubleArray("[1,2]", 3)
		if err == nil {
			t.Fatal("expected error for wrong element count")
		}
	})

	t.Run("empty brackets", func(t *testing.T) {
		vals, err := parseDoubleArray("[]", 0)
		if err != nil {
			t.Fatal(err)
		}
		if vals != nil {
			t.Fatalf("expected nil, got %v", vals)
		}
	})

	t.Run("negative values", func(t *testing.T) {
		vals, err := parseDoubleArray("[-1.5,0,1.5]", 3)
		if err != nil {
			t.Fatal(err)
		}
		if vals[0] != -1.5 || vals[2] != 1.5 {
			t.Fatalf("unexpected values: %v", vals)
		}
	})
}
