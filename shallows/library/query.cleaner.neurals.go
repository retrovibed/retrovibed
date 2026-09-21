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
	reEpisode = regexp.MustCompile(`(?i)^s\d+e\d+$`)

	reDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

// ParseReleaseEpisode splits a cleaned/predicted media string, formatted as
// {input.literal}\n{media.episode}\n{input.subtitle}\n{media.release}, into
// its parts. Empty fields are omitted from the output entirely (no blank
// lines), so only the literal is guaranteed: it is always the first line,
// the episode is only recognized as the line directly after it, and the
// release only as the last line. Whatever remains between is the subtitle.
func ParseReleaseEpisode(input string) (remaining, subtitle, datish, episodish string) {
	lines := strings.FieldsFunc(input, func(r rune) bool { return r == '\n' })
	if len(lines) == 0 {
		return "", "", "", ""
	}

	remaining, lines = lines[0], lines[1:]

	if len(lines) > 0 && reEpisode.MatchString(lines[0]) {
		episodish, lines = lines[0], lines[1:]
	}

	if end := len(lines) - 1; end >= 0 && reDate.MatchString(lines[end]) {
		datish, lines = lines[end], lines[:end]
	}

	return remaining, strings.Join(lines, " "), datish, episodish
}

// removes any hallucinated strings from the title.
func StripHallucinations(input string, generated string) string {
	var b strings.Builder
	for s := range strings.FieldsSeq(generated) {
		if !strings.Contains(input, s) {
			continue
		}

		if _, err := b.WriteString(" "); err != nil {
			log.Println("failed to write fragment", err)
			return ""
		}
		if _, err := b.WriteString(s); err != nil {
			log.Println("failed to write fragment", err)
			return ""
		}
	}

	return strings.TrimSpace(b.String())
}
