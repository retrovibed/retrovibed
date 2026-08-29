package cmdmedia

import (
	"context"
	"encoding/binary"
	"fmt"
	"iter"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/michiwend/gomusicbrainz"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"golang.org/x/text/language"
	"golang.org/x/time/rate"
)

type mbimport struct {
	StartAt  time.Time `flag:"" name:"start" help:"date to start retrieving data from" format:"2006-01-02" default:"1900-01-01"`
	EndAt    time.Time `flag:"" name:"end" help:"date to end retrieving data from" format:"2006-01-02" default:"${vars_date_started}"`
	Source   string    `flag:"" name:"source" help:"short id for the data source" hidden:"true" default:"musicbrainz"`
	Attempts uint      `flag:"" name:"attempts" help:"set maximum number of attempts per request" default:"5"`
	Contact  string    `flag:"" name:"contact" help:"required email for user agent" required:"true"`
	cause    error
}

func (t *mbimport) releases(ctx context.Context, c *gomusicbrainz.WS2Client, l *rate.Limiter, bs backoffx.Strategy) iter.Seq[library.Known] {
	return func(yield func(library.Known) bool) {
		const (
			limit = 100
		)

		currentDate := t.StartAt
		endDate := t.EndAt.Add(24 * time.Hour)

		for currentDate.Before(endDate) {
			offset := 0

			nextDate := currentDate.Add(24 * time.Hour)
			// Query a specific 24-hour slice of modifications
			// Format: [YYYY-MM-DD TO YYYY-MM-DD]
			query := fmt.Sprintf("lastmodified:[%s TO %s]",
				currentDate.Format(time.DateOnly),
				nextDate.Format(time.DateOnly),
			)

			for {
				resp, err := backoffx.AttemptV(ctx, bs, func(ctx context.Context, attempts uint) (*gomusicbrainz.ReleaseSearchResponse, error) {
					if err := l.Wait(ctx); err != nil {
						return nil, errorsx.Wrap(err, "rate limit failure")
					}

					if attempts > t.Attempts {
						return nil, backoffx.ErrStopAttempts
					}

					return c.SearchRelease(query, limit, offset)
				})

				if err != nil {
					errorsx.Debug(err)
					t.cause = errorsx.Wrapf(err, "failed to search releases %v - %d", currentDate, offset)
					return
				}

				log.Println("retrieved musicbrainz releases", query, offset, len(resp.Releases))

				if len(resp.Releases) == 0 {
					break
				}

				for _, rel := range resp.Releases {
					uidmd5 := uuid.FromBytesOrNil(md5x.JSON(rel).Sum(nil))

					lang := langx.FirstNonZero(
						errorsx.ZeroSilent(language.Parse(rel.TextRepresentation.Language)),
						language.Und,
					)

					credit := slicesx.FirstOrZero(rel.ArtistCredit.NameCredits...)

					// We use the RG ID so that all versions of this album collapse into one record.
					id := uuid.FromStringOrNil(string(rel.ReleaseGroup.ID))

					// Note: Search results don't include full tracklists (recordings).
					// To get songs, a separate Lookup call per release is required.
					v := library.Known{
						Source:           t.Source,
						UID:              ddiscapi.ImportedMediaUintID(t.Source, uint64(binary.BigEndian.Uint64(id.Bytes()[:8]))),
						Md5:              uidmd5.String(),
						Md5Lower:         binary.LittleEndian.Uint64(uuidx.LowN(uidmd5, 64)),
						ID:               id.String(),
						OriginalTitle:    rel.ReleaseGroup.Title,
						Title:            strings.TrimSpace(fmt.Sprintf("%s %s", credit.Artist.Name, rel.ReleaseGroup.Title)),
						Released:         rel.ReleaseGroup.FirstReleaseDate.Time,
						PosterPath:       fmt.Sprintf("https://coverartarchive.org/release-group/%s/front-500", id),
						OriginalLanguage: lang.String(),
						Mimetype:         mimex.Audio,
					}

					if !yield(v) {
						return
					}
				}

				offset += len(resp.Releases)
				if offset >= resp.Count || len(resp.Releases) < limit {
					break
				}
			}
			currentDate = currentDate.Add(24 * time.Hour)
		}
	}
}

func (t mbimport) run(ctx context.Context, c *gomusicbrainz.WS2Client, encoder *jsonl.Encoder, l *rate.Limiter, bs backoffx.Strategy) (err error) {
	for v := range t.releases(ctx, c, l, bs) {
		if err := encoder.Encode(v); err != nil {
			return errorsx.Wrap(err, "unable to encode music metadata")
		}
	}

	return t.cause
}

func (t mbimport) Run(gctx *cmdopts.Global) (err error) {
	// MusicBrainz requires a meaningful User-Agent
	c, err := gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"Retrovibed",
		"1.0.0",
		t.Contact,
	)

	if err != nil {
		return errorsx.Wrap(err, "failed to setup client")
	}

	encoder := jsonl.NewEncoder(os.Stdout)

	return t.run(
		gctx.Context,
		c,
		encoder,
		rate.NewLimiter(rate.Every(time.Second), 1),
		backoffx.New(backoffx.Exponential(time.Second), backoffx.Maximum(time.Minute)),
	)
}
