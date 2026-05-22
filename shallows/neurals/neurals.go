package neurals

// #cgo LDFLAGS: -lpredicttext
// #include <stdlib.h>
// extern int predict(const char* model_path, const char* input, size_t seq_len, long long num_tokens, long long pad, long long bos, long long eos, char* output, size_t output_len);
import "C"

import (
	"fmt"
	"unsafe"
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
	t := Text{
		model:     path,
		seqLen:    256,
		numTokens: 4096,
		pad:       0,
		bos:       1,
		eos:       2,
		outputLen: 4096,
	}
	for _, o := range options {
		o(&t)
	}
	return &t
}

func (t *Text) Predict(input string) (string, error) {
	cModel := C.CString(t.model)
	defer C.free(unsafe.Pointer(cModel))
	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))

	buf := make([]byte, t.outputLen)
	ret := C.predict(
		cModel,
		cInput,
		C.size_t(t.seqLen),
		C.longlong(t.numTokens),
		C.longlong(t.pad),
		C.longlong(t.bos),
		C.longlong(t.eos),
		(*C.char)(unsafe.Pointer(&buf[0])),
		C.size_t(t.outputLen),
	)
	if ret != 0 {
		return "", fmt.Errorf("predict failed for %q", input)
	}
	return C.GoString((*C.char)(unsafe.Pointer(&buf[0]))), nil
}
