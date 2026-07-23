package library

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/neurals"
)

const (
	NeuralMediaIDCached = "mediaid.onnx"
)

type QueryerCleanerV0 struct {
	text *neurals.Text
}

func NewQueryerCleanerAuto() QueryCleaner {
	return langx.FirstNonNil[QueryCleaner](
		NewQueryerCleanerV0(
			envx.String(
				userx.DefaultCacheDirectory(userx.DefaultRelRoot(), NeuralMediaIDCached),
				env.NeuralMediaID,
			),
		),
		QueryCleanerNoop(),
	)
}

func NewQueryerCleanerV0(path string, options ...func(*neurals.Text)) *QueryerCleanerV0 {
	if stringsx.Blank(path) {
		return nil
	}

	if !fsx.Exists(path) {
		log.Println("unable to locate", path, "not attempting to load")
		return nil
	}

	log.Println("neural located at", path, "loading...")

	return &QueryerCleanerV0{text: neurals.NewText(path, options...)}
}

func (t QueryerCleanerV0) Clean(_ context.Context, input string) (r string, err error) {
	if stringsx.Blank(input) {
		return "", nil
	}

	if r, err = t.text.Predict(input); err != nil {
		return "", errorsx.Wrapf(err, "failed cleaning input: %s", input)
	}

	if len(r) > len(input) {
		return "", errorsx.Errorf("result long than input: %s v %s", input, r)
	}

	return r, nil
}

var (
	reEpisode = regexp.MustCompile(`(?i)^(?:s\d{1,5}e\d{1,5}|ep\d{1,5}|e\d{1,5}|\d{1,5}x\d{1,5})$`)

	reDate = regexp.MustCompile(`^(?:` +
		`\d{3,5}[-/]\d{1,3}[-/]\d{1,3}` + // YYYY-MM-DD / YYYY/MM/DD (loose digit counts)
		`|\d{1,3}[-/]\d{1,3}[-/]\d{3,5}` + // MM-DD-YYYY
		`|\d{3,5}[-/]\d{1,3}` + // YYYY-MM
		`|\d{1,3}/\d{3,5}` + // MM/YYYY
		`|\d{3,5}` + // bare year
		`)$`)
)

// ParseReleaseEpisode splits a cleaned/predicted media string into its
// title remainder plus any trailing release-date and episode markers.
// Matching is intentionally loose on digit counts: the source is neural
// model output, not the synthetic generator, so a token like a year or
// episode number may carry a spurious extra digit from a generation error.
// Release always precedes episode, so the episode marker is only looked
// for as the final token, and the release marker only immediately before it.
func ParseReleaseEpisode(input string) (remaining, datish, episodish string) {
	tokens := strings.Fields(input)

	end := len(tokens)
	if end > 0 && reEpisode.MatchString(tokens[end-1]) {
		episodish = tokens[end-1]
		end--
	}

	if end > 0 && reDate.MatchString(tokens[end-1]) {
		datish = tokens[end-1]
		end--
	}

	remaining = strings.Join(tokens[:end], " ")
	return remaining, datish, episodish
}
