//go:build !(darwin || ffmpeg_disabled)

package ddisc

import (
	"io"
	"log"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/mutablelogic/go-media/pkg/ffmpeg"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

func Extract(src io.ReadSeeker) (_zero Extracted, err error) {
	mime, err := Mimetype(src)
	if err != nil {
		return _zero, errorsx.Wrap(err, "unable to determine mimetype")
	}

	if strings.HasPrefix(mime.String(), "audio/") {
		c, err := ffmpeg.NewReader(src)
		if err != nil {
			return _zero, errorsx.Wrap(err, "unable to create ffmpeg reader")
		}

		v, err := Audio(mime, c)
		return Extracted{Music: v}, errorsx.Wrap(err, "unable to extract audio metadata")
	}

	if strings.HasPrefix(mime.String(), "video/") {
		c, err := ffmpeg.NewReader(src)
		if err != nil {
			return _zero, errorsx.Wrap(err, "unable to create ffmpeg reader")
		}

		v, err := Video(mime, c)
		return Extracted{Video: v}, errorsx.Wrap(err, "unable to extract video metadata")
	}

	return _zero, errorsx.String("unknown content")
}

func Audio(mime *mimetype.MIME, src *ffmpeg.Reader) (m *MetadataAudio, err error) {
	const (
		KeyTitle  = "TITLE"
		KeyAlbum  = "ALBUM"
		KeyArtist = "ARTIST"
		KeyDate   = "DATE"
		KeyTrack  = "TRACK"
	)

	var (
		album  string
		artist string
		title  string
		track  string
		date   = timex.NegInf()
	)

	for _, md := range src.Metadata() {
		switch strings.ToUpper(md.Key()) {
		case KeyAlbum:
			album = md.Value()
		case KeyArtist:
			artist = md.Value()
		case KeyTitle:
			title = md.Value()
		case KeyDate:
			date = errorsx.Zero(parseDate(md.Value()))
		case KeyTrack:
			track = md.Value()
		default:
			debugx.Println("unknown metadata key", md.Key(), "->", md.Value())
		}
	}

	return &MetadataAudio{
		Metadata: Metadata{
			Mimetype:  mime.String(),
			Title:     stringsx.Join(" - ", slicesx.Filter(stringsx.Present, artist, album)...),
			Subtitle:  title,
			Collation: track,
			Date:      date,
		},
	}, nil
}

func Video(mime *mimetype.MIME, src *ffmpeg.Reader) (m *MetadataVideo, err error) {
	const (
		KeyTitle = "TITLE"
	)

	var (
		title string
		date  = timex.NegInf()
	)

	for _, md := range src.Metadata() {
		switch strings.ToUpper(md.Key()) {
		case KeyTitle:
			title = md.Value()
		default:
			log.Println("unknown metadata key", md.Key(), "->", md.Value())
		}
	}

	return &MetadataVideo{
		Metadata: Metadata{
			Mimetype: mime.String(),
			Title:    title,
			Date:     date,
		},
	}, nil
}

func parseDate(encoded string) (_ time.Time, err error) {
	const (
		YearOnly = "2006"
	)

	formats := []string{
		time.RFC822,
		time.RFC850,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		YearOnly,
	}

	for _, format := range formats {
		if ts, failed := time.Parse(format, encoded); failed == nil {
			return ts, nil
		} else {
			err = errorsx.Compact(err, failed)
		}
	}

	return timex.NegInf(), err
}
