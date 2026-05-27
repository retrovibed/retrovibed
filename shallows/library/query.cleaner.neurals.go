package library

import (
	"context"
	"log"

	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/userx"
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

func (t QueryerCleanerV0) Clean(_ context.Context, text string) (r string, err error) {
	return t.text.Predict(text)
}
