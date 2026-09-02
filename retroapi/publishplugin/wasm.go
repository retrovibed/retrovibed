package publishplugin

import (
	"bytes"
	"io"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
)

const (
	ErrNotWasi = errorsx.String("not a wasm file")
)

// wasmMagic is the 4-byte binary header every valid wasm module begins with
// ("\0asm").
var wasmMagic = []byte{0x00, 0x61, 0x73, 0x6D}

func VerifyWasmMagic(r io.Reader) error {
	magic := make([]byte, len(wasmMagic))
	if _, err := io.ReadFull(r, magic); errorsx.Ignore(err, io.EOF, io.ErrUnexpectedEOF) != nil {
		return errorsx.Wrap(err, "unable to sniff wasm magic bytes")
	}

	if !bytes.Equal(magic, wasmMagic) {
		return ErrNotWasi
	}

	return nil
}

// VerifyWasmMagicPath confirms path begins with the wasm binary header, so an
// unexpected build output or upload never gets installed into publish.d.
func VerifyWasmMagicPath(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errorsx.Wrapf(err, "unable to open compiled plugin: %s", path)
	}
	defer f.Close()

	if err := VerifyWasmMagic(f); err != nil {
		return errorsx.Wrapf(err, "file: %s", path)
	}

	return nil
}
