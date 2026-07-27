package neurals

import (
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

type Text struct {
	model     string
	seqLen    int
	numTokens int64
	pad       int64
	bos       int64
	eos       int64
	outputLen int
}

// TextOption is a slice of mutators enabling chaining:
//
//	NewText(TextOptions().Model("x").NumTokens(4096)...)
type TextOption []func(*Text)

func TextOptions() TextOption {
	return TextOption(nil)
}

func (t TextOption) SeqLen(n int) TextOption {
	return append(t, func(tt *Text) { tt.seqLen = n })
}

func (t TextOption) NumTokens(n int64) TextOption {
	return append(t, func(tt *Text) { tt.numTokens = n })
}

func (t TextOption) PAD(n int64) TextOption {
	return append(t, func(tt *Text) { tt.pad = n })
}

func (t TextOption) BOS(n int64) TextOption {
	return append(t, func(tt *Text) { tt.bos = n })
}

func (t TextOption) EOS(n int64) TextOption {
	return append(t, func(tt *Text) { tt.eos = n })
}

func (t TextOption) OutputLen(n int) TextOption {
	return append(t, func(tt *Text) { tt.outputLen = n })
}

func NewText(path string, options ...func(*Text)) *Text {
	return new(langx.Clone(Text{
		model:     path,
		seqLen:    256,
		numTokens: 4096,
		pad:       0,
		bos:       1,
		eos:       2,
		outputLen: 4096,
	}, options...))
}

func (t *Text) Predict(input string) (string, error) {
	return predict(t, LimitVocab(input, t.numTokens))
}
