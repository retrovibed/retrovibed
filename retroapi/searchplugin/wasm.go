package searchplugin

import (
	"bytes"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
)

// wasmMagic is the 4-byte binary header every valid wasm module begins with
// ("\0asm").
var wasmMagic = []byte{0x00, 0x61, 0x73, 0x6D}

// VerifyWasmMagic confirms path begins with the wasm binary header, so an
// unexpected build output or upload never gets installed into search.d.
func VerifyWasmMagic(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errorsx.Wrapf(err, "unable to open compiled plugin: %s", path)
	}
	defer f.Close()

	magic := make([]byte, len(wasmMagic))
	if _, err := f.Read(magic); err != nil {
		return errorsx.Wrapf(err, "unable to read compiled plugin: %s", path)
	}

	if !bytes.Equal(magic, wasmMagic) {
		return errorsx.Errorf("compiled output is not a valid wasm module: %s", path)
	}

	return nil
}
