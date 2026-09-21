package library_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

// neural output format: {input.literal}\n{media.episode}\n{input.subtitle}\n{media.release}
// empty fields are omitted entirely (no blank lines).
func TestParseReleaseEpisode(t *testing.T) {
	for _, tc := range []struct {
		name                               string
		input                              string
		remaining, subtitle, date, episode string
	}{
		{name: "empty input", input: "", remaining: "", subtitle: "", date: "", episode: ""},
		{name: "literal only", input: "God of Japan", remaining: "God of Japan"},
		{name: "literal episode", input: "God of Japan\nS01E01", remaining: "God of Japan", episode: "S01E01"},
		{name: "literal release", input: "God of Japan\n2021-09-09", remaining: "God of Japan", date: "2021-09-09"},
		{name: "literal episode subtitle", input: "Breaking Bad\nS01E01\nCat's in the Bag", remaining: "Breaking Bad", episode: "S01E01", subtitle: "Cat's in the Bag"},
		{name: "literal subtitle release", input: "Mission Impossible\nFallout\n2018-07-27", remaining: "Mission Impossible", subtitle: "Fallout", date: "2018-07-27"},
		{name: "literal episode release", input: "God of Japan\nS01E01\n2021-09-09", remaining: "God of Japan", episode: "S01E01", date: "2021-09-09"},
		{name: "literal episode subtitle release", input: "Breaking Bad\nS01E01\nCat's in the Bag\n2008-01-20", remaining: "Breaking Bad", episode: "S01E01", subtitle: "Cat's in the Bag", date: "2008-01-20"},
		{name: "episode lowercase", input: "God of Japan\ns01e01", remaining: "God of Japan", episode: "s01e01"},
		{name: "episode single digits", input: "God of Japan\nS1E1", remaining: "God of Japan", episode: "S1E1"},
		{name: "episode extra digit", input: "God of Japan\nS01E0049", remaining: "God of Japan", episode: "S01E0049"},
		{name: "episode e49 is not an episode", input: "God of Japan\ne49", remaining: "God of Japan", subtitle: "e49"},
		{name: "episode EP01 is not an episode", input: "God of Japan\nEP01", remaining: "God of Japan", subtitle: "EP01"},
		{name: "episode 1x01 is not an episode", input: "God of Japan\n1x01", remaining: "God of Japan", subtitle: "1x01"},
		{name: "date bare year is not a date", input: "God of Japan\n2006", remaining: "God of Japan", subtitle: "2006"},
		{name: "date YYYY-MM is not a date", input: "God of Japan\n2006-01", remaining: "God of Japan", subtitle: "2006-01"},
		{name: "date MM/YYYY is not a date", input: "God of Japan\n01/2006", remaining: "God of Japan", subtitle: "01/2006"},
		{name: "date YYYY/MM/DD is not a date", input: "God of Japan\n2006/01/02", remaining: "God of Japan", subtitle: "2006/01/02"},
		{name: "date MM-DD-YYYY is not a date", input: "God of Japan\n01-02-2006", remaining: "God of Japan", subtitle: "01-02-2006"},
		{name: "space separated literal is left intact", input: "God of Japan 2021-09-09 S01E01", remaining: "God of Japan 2021-09-09 S01E01"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			remaining, subtitle, datish, episodish := library.ParseReleaseEpisode(tc.input)
			require.Equal(t, tc.remaining, remaining)
			require.Equal(t, tc.subtitle, subtitle)
			require.Equal(t, tc.date, datish)
			require.Equal(t, tc.episode, episodish)
		})
	}
}
