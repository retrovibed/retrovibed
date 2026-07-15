package cmdmedia

import (
	"encoding/binary"
	"fmt"
	"iter"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/dashotv/tvdb"
	"github.com/dashotv/tvdb/openapi/models/shared"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type tvdbimport struct {
	APIKey string  `flag:"" name:"apikey" help:"api key for requests" require:"true"`
	URL    string  `flag:"" name:"baseurl" help:"url base for image assets" default:"https://thetvdb.com"`
	Source string  `flag:"" name:"source" help:"short id for the data source" hidden:"true" default:"tvdb"`
	Lang   string  `flag:"" name:"lang" help:"language to retrieve from aliases" default:"eng"`
	Limit  float64 `flag:"" name:"maxpage" help:"maximum page to retrieve"`
	cause  error
}

func (t tvdbimport) imgpath(s string) string {
	if stringsx.Blank(s) {
		return ""
	}

	return fmt.Sprintf("%s%s", t.URL, s)
}

func (t *tvdbimport) records(c *tvdb.Client) iter.Seq[library.Known] {
	return func(yield func(library.Known) bool) {
		epochts := time.Unix(0, 0).Format(time.DateOnly)

		for i := 0.; i < langx.FirstNonZero(t.Limit, math.MaxFloat64); i++ {
			page, err := c.GetAllSeries(new(i))
			if err != nil {
				t.cause = errorsx.Wrap(err, "failed to discover records")
				return
			}

			for _, mr := range page.Data {
				if stringsx.Blank(langx.Autoderef(mr.Image)) {
					debugx.Println("skipping", langx.Autoderef(mr.Name), "missing poster")
					continue
				}

				_md5 := md5x.JSON(mr)
				uidmd5 := uuid.FromBytesOrNil(_md5.Sum(nil))

				alias := slicesx.FindOrZero(func(v shared.Alias) bool {
					return langx.Autoderef(v.Language) == t.Lang
				}, mr.Aliases...)

				var translation = &shared.Translation{
					Name:     new(langx.FirstNonZero(langx.Autoderef(alias.Name), langx.Autoderef(mr.Name))),
					Overview: new(langx.Autoderef(mr.Overview)),
				}

				if langx.Autoderef(mr.OriginalLanguage) != t.Lang && slicesx.FindOrZero(func(lang string) bool { return t.Lang == lang }, mr.OverviewTranslations...) != "" {
					r, err := c.GetSeriesTranslation(float64(*mr.ID), t.Lang)
					if err != nil {
						log.Println("failed to retrieve translation", *mr.ID, err)
					} else {
						translation = r.Data
					}
				}

				v := library.Known{
					Source:           t.Source,
					UID:              library.KnownImportedUintID(t.Source, uint64(langx.Autoderef(mr.ID))),
					Md5:              uidmd5.String(),
					Md5Lower:         binary.LittleEndian.Uint64(uuidx.LowN(uidmd5, 64)),
					ID:               strconv.FormatInt(int64(langx.Autoderef(mr.ID)), 10),
					OriginalLanguage: langx.Autoderef(mr.OriginalLanguage),
					OriginalTitle:    langx.Autoderef(mr.Name),
					Title:            langx.Autoderef(translation.Name),
					Overview:         langx.Autoderef(translation.Overview),
					PosterPath:       t.imgpath(langx.Autoderef(mr.Image)),
					Released:         errorsx.Zero(time.Parse(time.DateOnly, langx.FirstNonZero(langx.Autoderef(mr.FirstAired), epochts))),
					Mimetype:         mimex.Video,
					// Popularity:       max(langx.Autoderef(mr.Score)/100000, 1.0), // essentially a useless field, they only include it for sorting.
				}

				if !yield(v) {
					return
				}
			}

			if stringsx.Blank(langx.Autoderef(page.Links.Next)) {
				return
			}

			log.Println("retrieving", langx.Autoderef(page.Links.Next))
		}
	}
}

func (t tvdbimport) Run(gctx *cmdopts.Global) (err error) {
	c, err := tvdb.Login(t.APIKey)
	if err != nil {
		return errorsx.Wrap(err, "unable to initialize tvdb client")
	}

	encoder := jsonl.NewEncoder(os.Stdout)

	for v := range t.records(c) {
		if err := encoder.Encode(v); err != nil {
			return errorsx.Wrap(err, "unable to encode media")
		}
	}

	return t.cause
}
