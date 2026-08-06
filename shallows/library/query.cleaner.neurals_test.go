package library_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestParseReleaseEpisode(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		input                    string
		remaining, date, episode string
	}{
		{name: "title date episode", input: "God of Japan 2021-09-09 e49", remaining: "God of Japan", date: "2021-09-09", episode: "e49"},
		{name: "episode before date is not unpacked", input: "God of Japan e49 2021-09-09", remaining: "God of Japan e49", date: "2021-09-09", episode: ""},
		{name: "title date only", input: "God of Japan 2021-09-09", remaining: "God of Japan", date: "2021-09-09", episode: ""},
		{name: "title episode only", input: "God of Japan e49", remaining: "God of Japan", date: "", episode: "e49"},
		{name: "title only", input: "God of Japan", remaining: "God of Japan", date: "", episode: ""},
		{name: "empty input", input: "", remaining: "", date: "", episode: ""},
		{name: "extra digit year", input: "God of Japan 20211", remaining: "God of Japan", date: "20211", episode: ""},
		{name: "extra digit episode", input: "God of Japan e0049", remaining: "God of Japan", date: "", episode: "e0049"},
		{name: "episode s01e01", input: "God of Japan s01e01", remaining: "God of Japan", date: "", episode: "s01e01"},
		{name: "episode S1E1", input: "God of Japan S1E1", remaining: "God of Japan", date: "", episode: "S1E1"},
		{name: "episode 1x01", input: "God of Japan 1x01", remaining: "God of Japan", date: "", episode: "1x01"},
		{name: "episode EP01", input: "God of Japan EP01", remaining: "God of Japan", date: "", episode: "EP01"},
		{name: "episode e1", input: "God of Japan e1", remaining: "God of Japan", date: "", episode: "e1"},
		{name: "date bare year", input: "God of Japan 2006", remaining: "God of Japan", date: "2006", episode: ""},
		{name: "date YYYY-MM-DD", input: "God of Japan 2006-01-02", remaining: "God of Japan", date: "2006-01-02", episode: ""},
		{name: "date YYYY-MM", input: "God of Japan 2006-01", remaining: "God of Japan", date: "2006-01", episode: ""},
		{name: "date MM/YYYY", input: "God of Japan 01/2006", remaining: "God of Japan", date: "01/2006", episode: ""},
		{name: "date YYYY/MM/DD", input: "God of Japan 2006/01/02", remaining: "God of Japan", date: "2006/01/02", episode: ""},
		{name: "date MM-DD-YYYY", input: "God of Japan 01-02-2006", remaining: "God of Japan", date: "01-02-2006", episode: ""},
		{name: "v2 title date episode", input: "God of Japan\n2021-09-09\ne49", remaining: "God of Japan", date: "2021-09-09", episode: "e49"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			remaining, datish, episodish := library.ParseReleaseEpisode(tc.input)
			require.Equal(t, tc.remaining, remaining)
			require.Equal(t, tc.date, datish)
			require.Equal(t, tc.episode, episodish)
		})
	}
}
