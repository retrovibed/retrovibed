package library_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestStripHallucinations(t *testing.T) {
	for _, tc := range []struct {
		name      string
		input     string
		generated string
		expected  string
	}{
		{name: "all words present", input: "God of Japan 2021", generated: "God of Japan", expected: "God of Japan"},
		{name: "hallucinated word removed", input: "God of Japan", generated: "God of Japan Ultra", expected: "God of Japan"},
		{name: "all hallucinated", input: "God of Japan", generated: "Totally Fake Title", expected: ""},
		{name: "empty generated", input: "God of Japan", generated: "", expected: ""},
		{name: "empty input", input: "", generated: "God of Japan", expected: ""},
		{name: "field with no exact substring match is dropped", input: "God of Japan", generated: "God of Japans", expected: "God of"},
		{name: "duplicate words kept if present", input: "God God of Japan", generated: "God God", expected: "God God"},
		{name: "extra whitespace normalized", input: "God of Japan", generated: "  God   of   Japan  ", expected: "God of Japan"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actual := library.StripHallucinations(tc.input, tc.generated)
			require.Equal(t, tc.expected, actual)
		})
	}
}
