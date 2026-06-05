package cmdlibrary

import "github.com/retrovibed/retrovibed/shallows/library"

const exportChunkSize = 16 * 1024 * 1024 // 16 MiB

type exportHeader struct {
	Chunks uint64 `json:"chunks"`
}

type exportChunk struct {
	Data []byte `json:"data"`
}

type exportTrailer struct {
	Metadata library.Metadata `json:"metadata"`
	MD5      string           `json:"md5"`
}
