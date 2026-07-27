package neurals_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/neurals"
	"github.com/stretchr/testify/require"
)

func TestLimitVocab(t *testing.T) {
	for _, tc := range []struct {
		name      string
		input     string
		numTokens int64
		expected  string
	}{
		{name: "ascii passthrough", input: "Nirvana In Utero", numTokens: 4096, expected: "Nirvana In Utero"},
		{name: "emoji stripped", input: "In Utero ⭐️ FLAC", numTokens: 4096, expected: "In Utero FLAC"},
		{name: "leading and trailing out of vocab", input: "⭐ Utero ⭐", numTokens: 4096, expected: "Utero"},
		{name: "all out of vocab collapses to empty", input: "⭐⭐⭐", numTokens: 4096, expected: ""},
		{name: "zero num tokens strips everything", input: "abc", numTokens: 0, expected: ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, neurals.LimitVocab(tc.input, tc.numTokens))
		})
	}
}
