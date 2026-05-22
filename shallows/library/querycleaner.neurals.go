package library

import (
	"context"

	"github.com/retrovibed/retrovibed/shallows/neurals"
)

type QueryerCleanerV0 struct {
	text *neurals.Text
}

func NewQueryerCleanerV0(path string, options ...func(*neurals.Text)) QueryerCleanerV0 {
	return QueryerCleanerV0{text: neurals.NewText(path, options...)}
}

func (t QueryerCleanerV0) Clean(_ context.Context, text string) (string, error) {
	return t.text.Predict(text)
}
