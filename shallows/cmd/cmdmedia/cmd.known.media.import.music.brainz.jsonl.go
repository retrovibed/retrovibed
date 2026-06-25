package cmdmedia

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"iter"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"golang.org/x/text/language"
)

// mbJSONRelease matches the MusicBrainz JSON API release format.
// This is distinct from gomusicbrainz.Release, which uses XML tags and Go
// PascalCase field names.
type mbJSONRelease struct {
	ID                 string               `json:"id"`
	Title              string               `json:"title"`
	TextRepresentation mbJSONTextRepr       `json:"text-representation"`
	ReleaseGroup       mbJSONReleaseGroup   `json:"release-group"`
	ArtistCredit       []mbJSONArtistCredit `json:"artist-credit"`
	Asin               string               `json:"asin"`
}

type mbJSONTextRepr struct {
	Language string `json:"language"`
	Script   string `json:"script"`
}

type mbJSONReleaseGroup struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	FirstReleaseDate string `json:"first-release-date"`
}

type mbJSONArtistCredit struct {
	JoinPhrase string       `json:"joinphrase"`
	Artist     mbJSONArtist `json:"artist"`
}

type mbJSONArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type mbjsonlimport struct {
	Source string        `flag:"" name:"source" help:"short id for the data source" hidden:"true" default:"musicbrainz"`
	Output cmdopts.IOOut `flag:"" name:"output" default:"-" help:"output destination; '-' for stdout"`
	cause  error
}

func (t *mbjsonlimport) releases(ctx context.Context, r io.Reader) iter.Seq[library.Known] {
	return func(yield func(library.Known) bool) {
		seq := jsonl.Iter[mbJSONRelease](jsonl.NewDecoder(r))
		for rel := range seq.Each(ctx) {
			rawID := langx.FirstNonZero(rel.ReleaseGroup.ID, rel.ID)
			id := uuid.FromStringOrNil(rawID)

			uidmd5 := uuid.FromBytesOrNil(md5x.JSON(rel).Sum(nil))

			lang := langx.FirstNonZero(
				errorsx.ZeroSilent(language.Parse(rel.TextRepresentation.Language)),
				language.Und,
			)

			title := langx.FirstNonZero(rel.ReleaseGroup.Title, rel.Title)
			credit := slicesx.FirstOrZero(rel.ArtistCredit...)

			v := library.Known{
				Source:           t.Source,
				UID:              library.KnownImportedUintID(t.Source, uint64(binary.BigEndian.Uint64(id.Bytes()[:8]))),
				Md5:              uidmd5.String(),
				Md5Lower:         binary.LittleEndian.Uint64(uuidx.LowN(uidmd5, 64)),
				ID:               id.String(),
				OriginalTitle:    title,
				Title:            strings.TrimSpace(fmt.Sprintf("%s %s", credit.Artist.Name, title)),
				Released:         mbjsonlParseDate(rel.ReleaseGroup.FirstReleaseDate),
				PosterPath:       fmt.Sprintf("https://coverartarchive.org/release-group/%s/front-500", id),
				OriginalLanguage: lang.String(),
				Mimetype:         mimex.Audio,
			}

			if !yield(v) {
				return
			}
		}
		t.cause = seq.Err()
	}
}

// mbjsonlParseDate parses MusicBrainz date strings: YYYY-MM-DD, YYYY-MM, YYYY.
func mbjsonlParseDate(s string) time.Time {
	for _, format := range []string{"2006-01-02", "2006-01", "2006"} {
		if t, err := time.Parse(format, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func (t mbjsonlimport) run(ctx context.Context, r io.Reader, encoder *jsonl.Encoder) (err error) {
	for v := range t.releases(ctx, r) {
		if err := encoder.Encode(v); err != nil {
			return errorsx.Wrap(err, "unable to encode music metadata")
		}
	}
	return t.cause
}

func (t mbjsonlimport) Run(gctx *cmdopts.Global) (err error) {
	out, err := t.Output.Open(os.Stdout)
	if err != nil {
		return errorsx.Wrap(err, "unable to open output")
	}

	var in io.Reader = bytes.NewReader(nil)
	if cmdopts.Readable(os.Stdin) {
		in = os.Stdin
	}

	return t.run(gctx.Context, in, jsonl.NewEncoder(out))
}
