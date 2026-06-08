//go:build darwin || android || ffmpeg_disabled

package ddisc

import (
	"io"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

func Extract(src io.ReadSeeker) (_zero Extracted, err error) {
	mime, err := Mimetype(src)
	if err != nil {
		return _zero, errorsx.Wrap(err, "unable to determine mimetype")
	}

	if strings.HasPrefix(mime.String(), "audio/") {
		v, err := Audio(mime)
		return Extracted{Music: v}, errorsx.Wrap(err, "unable to extract audio metadata")
	}

	if strings.HasPrefix(mime.String(), "video/") {
		v, err := Video(mime)
		return Extracted{Video: v}, errorsx.Wrap(err, "unable to extract video metadata")
	}

	return _zero, errorsx.String("unknown content")
}

func Audio(mime *mimetype.MIME) (m *MetadataAudio, err error) {
	return &MetadataAudio{
		Metadata: Metadata{
			Mimetype: mime.String(),
		},
	}, nil
}

func Video(mime *mimetype.MIME) (m *MetadataVideo, err error) {
	return &MetadataVideo{
		Metadata: Metadata{
			Mimetype: mime.String(),
		},
	}, nil
}
