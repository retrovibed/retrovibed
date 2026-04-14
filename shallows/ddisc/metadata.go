package ddisc

import (
	"io"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/iox"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
)

type Metadata struct {
	Mimetype  string
	Title     string
	Subtitle  string
	Collation string
	Date      time.Time
}

func (t Metadata) String() string {
	return stringsx.Join(" ", slicesx.Filter(stringsx.Present, t.Title, t.Collation, t.Subtitle)...)
}

type MetadataVideo struct {
	Metadata
}

type MetadataAudio struct {
	Metadata
}

type Extracted struct {
	Video *MetadataVideo
	Music *MetadataAudio
}

func (t Extracted) Metadata() Metadata {
	return langx.FirstNonZero(langx.Autoderef(t.Video).Metadata, langx.Autoderef(t.Music).Metadata)
}

func Mimetype(src io.ReadSeeker) (*mimetype.MIME, error) {
	mime, err := mimetype.DetectReader(src)
	return mime, errorsx.Compact(errorsx.Wrap(err, "unable to determine mimetype"), iox.Rewind(src))
}
