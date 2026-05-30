//go:build retrovibed && neural

package neurals

// -Bstatic/-Bdynamic wrapping forces predicttext to link statically without affecting other libs.
// Done here rather than CGO_LDFLAGS env to avoid Go repeating the flags once per CGo module.

// #cgo LDFLAGS: -Wl,-Bstatic -lpredicttext -Wl,-Bdynamic
// #include <stdlib.h>
// extern int predict(const char* model_path, const char* input, size_t seq_len, long long num_tokens, long long pad, long long bos, long long eos, char* output, size_t output_len);
import "C"

import (
	"fmt"
	"log"
	"unsafe"
)

func predict(t *Text, input string) (res string, err error) {
	log.Println("-----------------------------------------------------------------------")
	defer func() {
		if err != nil {
			return
		}
		log.Println("cleaned", input, "->", res)
	}()
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
