package cmdmedia

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"iter"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"

	tmdb "github.com/cyruzin/golang-tmdb"
)

type tmdbimport struct {
	APIKey      string    `flag:"" name:"apikey" help:"api key for requests" require:"true"  env:"RETROVIBED_TMDB_APIKEY"`
	URL         string    `flag:"" name:"baseurl" help:"url base for image assets" default:"https://image.tmdb.org/t/p/original"`
	StartAt     time.Time `flag:"" name:"start" help:"date to start retrieving data from, format: 2006-01-02" format:"2006-01-02" default:"1900-01-01"`
	EndAt       time.Time `flag:"" name:"end" help:"date to end retrieving data from, format: 2006-01-02" format:"2006-01-02" default:"${vars_date_started}"`
	Source      string    `flag:"" name:"source" help:"short id for the data source" hidden:"true" default:"tmdb"`
	Attempts    uint      `flag:"" name:"attempts" help:"set maximum number of attempts per request" default:"5"`
	Concurrency uint16    `flag:"" name:"concurrency" help:"maximum number of shows to fetch episode details for concurrently" default:"8"`
	cause       error
}

// imgpath takes a pointer receiver, not a value one, so calling it doesn't
// copy the whole tmdbimport struct - series() calls it concurrently from
// multiple pool workers while its own dispatcher goroutine writes t.cause,
// and a value-receiver copy would read that field as part of copying every
// field, racing with those writes even though imgpath itself only uses URL.
func (t *tmdbimport) imgpath(s string) string {
	if stringsx.Blank(s) {
		return ""
	}

	return fmt.Sprintf("%s%s", t.URL, s)
}

// tmdbLanguageOptions returns urlOptions requesting lang from TMDB, or nil
// when lang is blank (TMDB then falls back to its default, en-US). Season
// and episode detail requests pass the show's own OriginalLanguage here so
// TMDB returns names in that language instead - it sometimes has an episode
// title when the English translation is missing and falls back to a
// generic "Episode N" placeholder.
func tmdbLanguageOptions(lang string) map[string]string {
	if stringsx.Blank(lang) {
		return nil
	}

	return map[string]string{"language": lang}
}

// tmdbNotFound reports whether err is TMDB's permanent "resource not found"
// response (status_code 34), e.g. a season TVDetails lists that has no
// actual season-details page.
func tmdbNotFound(err error) bool {
	var terr tmdb.Error
	return errors.As(err, &terr) && terr.StatusCode == 34
}

func (t *tmdbimport) movies(ctx context.Context, c *tmdb.Client) iter.Seq[library.Known] {
	return func(yield func(library.Known) bool) {
		cctx, cancel := context.WithCancel(ctx)
		defer cancel()

		rmedia := make(chan library.Known, 128)
		media := asynccompute.New(func(ctx context.Context, results []tmdb.MovieResult) error {
			for _, mr := range results {
				_md5 := md5x.JSON(mr)
				uidmd5 := uuid.FromBytesOrNil(_md5.Sum(nil))
				v := library.Known{
					Source:           t.Source,
					UID:              ddiscapi.ImportedMediaUintID(t.Source, uint64(mr.ID)),
					Md5:              uidmd5.String(),
					Md5Lower:         binary.LittleEndian.Uint64(uuidx.LowN(uidmd5, 64)),
					ID:               strconv.FormatInt(mr.ID, 10),
					Adult:            mr.Adult,
					BackdropPath:     t.imgpath(mr.BackdropPath),
					OriginalLanguage: mr.OriginalLanguage,
					OriginalTitle:    mr.OriginalTitle,
					Overview:         mr.Overview,
					Popularity:       float64(mr.Popularity),
					PosterPath:       t.imgpath(mr.PosterPath),
					Title:            mr.Title,
					Released:         errorsx.ZeroSilent(time.Parse(time.DateOnly, mr.ReleaseDate)),
					Mimetype:         mimex.Video,
					ParentUID:        uuid.Nil.String(),
				}

				if !sendResult(ctx, rmedia, v) {
					return context.Cause(ctx)
				}
			}
			return nil
		}, asynccompute.Workers[[]tmdb.MovieResult](t.Concurrency), asynccompute.Backlog[[]tmdb.MovieResult](t.Concurrency))
		dates := asynccompute.New(func(ctx context.Context, year time.Time) (err error) {
			var (
				resp = &tmdb.DiscoverMovie{
					PaginatedResultsMeta: tmdb.PaginatedResultsMeta{
						TotalResults: math.MaxInt64,
						TotalPages:   math.MaxInt64,
					},
				}
			)

			bs := backoffx.New(backoffx.Exponential(time.Second), backoffx.Maximum(time.Minute))
			for page := int64(1); page <= resp.TotalPages; page = resp.Page + 1 {
				resp, err = backoffx.AttemptV(ctx, bs, func(ctx context.Context, attempts uint) (*tmdb.DiscoverMovie, error) {
					log.Println("retrieving movies", attempts, year, page)
					if attempts > t.Attempts {
						return nil, backoffx.ErrStopAttempts
					}

					return c.GetDiscoverMovie(map[string]string{
						"include_adult":            "true",
						"page":                     strconv.FormatInt(page, 10),
						"primary_release_date.gte": year.Format(time.DateOnly),
						"primary_release_date.lte": year.Format(time.DateOnly),
						"sort_by":                  "primary_release_date.asc",
					})
				})

				if err != nil {
					return errorsx.Wrapf(err, "failed to discover movies %v - %d", year, page)
				}

				if err = media.Run(ctx, resp.Results); err != nil {
					return errorsx.Wrap(err, "failed to push results")
				}
			}

			return nil
		}, asynccompute.Workers[time.Time](t.Concurrency), asynccompute.Backlog[time.Time](t.Concurrency))
		go func() {
			defer func() {
				if err := asynccompute.Shutdown(context.WithoutCancel(ctx), dates, media); err != nil && t.cause == nil {
					t.cause = err
				}
				close(rmedia)
			}()

			for date, edate := t.StartAt, t.EndAt.Add(24*time.Hour); date.Before(edate); date = date.Add(24 * time.Hour) {
				if err := dates.Run(cctx, date); err != nil {
					t.cause = err
					return
				}
			}
		}()

		stopped := false
		for v := range rmedia {
			if stopped {
				continue
			}

			if !yield(v) {
				cancel()
				stopped = true
			}
		}
	}
}

// sendResult delivers v onto results, or reports false if ctx is done first -
// used so a blocked send unblocks as soon as the consumer stops pulling
// results (see series()'s cancel on early yield-stop).
func sendResult(ctx context.Context, results chan<- library.Known, v library.Known) bool {
	select {
	case results <- v:
		return true
	case <-ctx.Done():
		return false
	}
}

// retrieveSeriesJob is one discovered show whose TVDetails (needed for its
// season list) is fetched by the "retrieve series" pool.
type retrieveSeriesJob struct {
	mr tmdb.TVShowResult
}

// processSeriesJob converts a show's TVDetails into its library.Known record
// and dispatches one retrieveEpisodesJob per season - work done by the
// "process a series" pool.
type processSeriesJob struct {
	mr      tmdb.TVShowResult
	details *tmdb.TVDetails
}

// retrieveEpisodesJob is one show/season pair whose episode details are
// fetched and emitted by the "retrieve episodes" pool.
type retrieveEpisodesJob struct {
	showID int64
	season tmdb.Season
	parent library.Known
}

func (t *tmdbimport) seriesKnown(mr tmdb.TVShowResult) library.Known {
	_md5 := md5x.JSON(mr)
	uidmd5 := uuid.FromBytesOrNil(_md5.Sum(nil))

	return library.Known{
		UID:              ddiscapi.ImportedMediaUintID(t.Source, uint64(mr.ID)),
		Source:           t.Source,
		Md5:              uidmd5.String(),
		Md5Lower:         binary.LittleEndian.Uint64(uuidx.LowN(uidmd5, 64)),
		ID:               strconv.FormatInt(mr.ID, 10),
		BackdropPath:     t.imgpath(mr.BackdropPath),
		OriginalLanguage: mr.OriginalLanguage,
		OriginalTitle:    mr.OriginalName,
		Overview:         mr.Overview,
		Popularity:       float64(mr.Popularity),
		PosterPath:       t.imgpath(mr.PosterPath),
		Title:            mr.Name,
		Adult:            mr.Adult,
		Released:         errorsx.ZeroSilent(time.Parse(time.DateOnly, mr.FirstAirDate)),
		Mimetype:         mimex.Video,
		ParentUID:        uuid.Nil.String(),
	}
}

// series discovers TV shows page by page (necessarily sequential - it's a
// paged cursor walk) and fans the slow per-show work out across three
// chained asynccompute pools so many shows' details are fetched
// concurrently instead of one at a time:
//
//	discovery -> poolSeries (GetTVDetails) -> poolProcess (emit the show's
//	own record, dispatch one job per season) -> poolEpisodes
//	(GetTVSeasonDetails, emit episode records)
//
// Shutting the pools down (see asynccompute.Shutdown) closes poolSeries
// first, which - because Close drains that pool's own workers to
// completion before returning - guarantees every job it might still submit
// to poolProcess has already been submitted by the time poolProcess is
// closed next, and the same holds one level further down for poolEpisodes.
// That ordering is what makes it safe to close results once Shutdown
// returns: nothing can still be sending to it. Results are fanned back in
// over one shared channel that the returned iterator ranges over; because a
// show and its episodes are produced by different pools running
// concurrently, cross-show and show-to-episode ordering are both no longer
// guaranteed (only within a single GetTVSeasonDetails response are an
// episode's siblings emitted in order).
func (t *tmdbimport) series(ctx context.Context, c *tmdb.Client) iter.Seq[library.Known] {
	return func(yield func(library.Known) bool) {
		cctx, cancel := context.WithCancel(ctx)
		defer cancel()

		results := make(chan library.Known, 128)

		var (
			poolEpisodes *asynccompute.Pool[retrieveEpisodesJob]
			poolProcess  *asynccompute.Pool[processSeriesJob]
			poolSeries   *asynccompute.Pool[retrieveSeriesJob]
		)

		poolEpisodes = asynccompute.New(func(ctx context.Context, job retrieveEpisodesJob) error {
			bs := backoffx.New(backoffx.Exponential(time.Second), backoffx.Maximum(time.Minute))

			sdetails, err := backoffx.AttemptV(ctx, bs, func(ctx context.Context, attempts uint) (*tmdb.TVSeasonDetails, error) {
				log.Println("retrieving season details", attempts, job.showID, job.season.SeasonNumber)
				if attempts > t.Attempts {
					return nil, backoffx.ErrStopAttempts
				}

				d, err := c.GetTVSeasonDetails(int(job.showID), job.season.SeasonNumber, tmdbLanguageOptions(job.parent.OriginalLanguage))
				if tmdbNotFound(err) {
					return nil, errors.Join(backoffx.ErrStopAttempts, err)
				}

				return d, err
			})

			if tmdbNotFound(err) {
				log.Println("season not found, skipping", job.showID, job.season.SeasonNumber)
				return nil
			}

			if err != nil {
				errorsx.Debug(err)
				return errorsx.Wrapf(err, "failed to retrieve season details %d - %d", job.showID, job.season.SeasonNumber)
			}

			for _, ep := range sdetails.Episodes {
				_md5 := md5x.JSON(ep)
				uidmd5 := uuid.FromBytesOrNil(_md5.Sum(nil))
				imgpath := t.imgpath(ep.StillPath)
				v := library.Known{
					Source:           t.Source,
					UID:              ddiscapi.ImportedMediaUintID(t.Source, uint64(ep.ID)),
					Md5:              uidmd5.String(),
					Md5Lower:         binary.LittleEndian.Uint64(uuidx.LowN(uidmd5, 64)),
					ID:               strconv.FormatInt(ep.ID, 10),
					OriginalLanguage: job.parent.OriginalLanguage,
					OriginalTitle:    job.parent.OriginalTitle,
					Title:            job.parent.Title,
					Overview:         ep.Overview,
					Subtitle:         ep.Name,
					ParentUID:        job.parent.UID,
					PosterPath:       stringsx.FirstNonBlank(imgpath, job.parent.PosterPath),
					BackdropPath:     stringsx.FirstNonBlank(imgpath, job.parent.BackdropPath),
					Collation:        library.KnownCollationEpisode(uint16(ep.SeasonNumber), uint16(ep.EpisodeNumber)),
					Released:         errorsx.ZeroSilent(time.Parse(time.DateOnly, ep.AirDate)),
					Mimetype:         mimex.Video,
				}

				if !sendResult(ctx, results, v) {
					return context.Cause(ctx)
				}
			}

			return nil
		}, asynccompute.Workers[retrieveEpisodesJob](t.Concurrency), asynccompute.Backlog[retrieveEpisodesJob](t.Concurrency))

		poolProcess = asynccompute.New(func(ctx context.Context, job processSeriesJob) error {
			v := t.seriesKnown(job.mr)
			if !sendResult(ctx, results, v) {
				return context.Cause(ctx)
			}

			for _, season := range job.details.Seasons {
				if err := poolEpisodes.Run(ctx, retrieveEpisodesJob{showID: job.mr.ID, season: season, parent: v}); err != nil {
					return err
				}
			}

			return nil
		}, asynccompute.Workers[processSeriesJob](t.Concurrency), asynccompute.Backlog[processSeriesJob](t.Concurrency))

		poolSeries = asynccompute.New(func(ctx context.Context, job retrieveSeriesJob) error {
			bs := backoffx.New(backoffx.Exponential(time.Second), backoffx.Maximum(time.Minute))

			details, err := backoffx.AttemptV(ctx, bs, func(ctx context.Context, attempts uint) (*tmdb.TVDetails, error) {
				log.Println("retrieving tv details", attempts, job.mr.ID)
				if attempts > t.Attempts {
					return nil, backoffx.ErrStopAttempts
				}

				return c.GetTVDetails(int(job.mr.ID), nil)
			})

			if err != nil {
				errorsx.Debug(err)
				return errorsx.Wrapf(err, "failed to retrieve tv details %d", job.mr.ID)
			}

			return poolProcess.Run(ctx, processSeriesJob{mr: job.mr, details: details})
		}, asynccompute.Workers[retrieveSeriesJob](t.Concurrency), asynccompute.Backlog[retrieveSeriesJob](t.Concurrency))

		go func() {
			defer func() {
				// poolSeries is the only pool the dispatch loop below
				// submits to; Shutdown closes it first, and - because
				// Close drains a pool's own workers to completion before
				// returning - every job it might still submit downstream
				// has been submitted by the time poolProcess is closed
				// next, and the same holds one level further down for
				// poolEpisodes. So by the time Shutdown returns, nothing
				// can still be sending to results.
				if err := asynccompute.Shutdown(context.WithoutCancel(ctx), poolSeries, poolProcess, poolEpisodes); err != nil && t.cause == nil {
					t.cause = err
				}
				close(results)
			}()

			var (
				resp *tmdb.DiscoverTV
			)

			bs := backoffx.New(backoffx.Exponential(time.Second), backoffx.Maximum(time.Minute))
			year := t.StartAt
			cyear := t.EndAt.Add(24 * time.Hour)

			for page := int64(1); cyear.After(year); {
				var (
					err error
				)

				resp, err = backoffx.AttemptV(cctx, bs, func(ctx context.Context, attempts uint) (*tmdb.DiscoverTV, error) {
					log.Println("retrieving series", attempts, year, page)
					if attempts > t.Attempts {
						return nil, backoffx.ErrStopAttempts
					}

					return c.GetDiscoverTV(map[string]string{
						"include_adult":      "true",
						"page":               strconv.FormatInt(page, 10),
						"first_air_date.gte": year.Format(time.DateOnly),
						"first_air_date.lte": year.Format(time.DateOnly),
						"sort_by":            "first_air_date.asc",
					})
				})

				if err != nil {
					errorsx.Debug(err)
					t.cause = errorsx.Wrapf(err, "failed to discover series %v - %d", year, page)
					return
				}

				for _, mr := range resp.Results {
					if err := poolSeries.Run(cctx, retrieveSeriesJob{mr: mr}); err != nil {
						t.cause = err
						return
					}
				}

				year = slicesx.LastOrDefault(year, slicesx.MapTransform(func(mr tmdb.TVShowResult) time.Time {
					return timex.Max(errorsx.ZeroSilent(time.Parse(time.DateOnly, mr.FirstAirDate)), year)
				}, resp.Results...)...)

				if page >= resp.TotalPages {
					year = year.Add(24 * time.Hour)
					page = 1
				} else {
					page = resp.Page + 1
				}
			}
		}()

		stopped := false
		for v := range results {
			if stopped {
				// drain without yielding: cancel() only asks the pools to
				// stop, it doesn't happen instantly, and results only
				// closes once the dispatcher goroutine's Shutdown call
				// above has fully drained them. Returning before that
				// would leave their goroutines running detached from this
				// call - e.g. still holding and using c after the caller
				// has moved on and possibly reused or discarded it.
				continue
			}

			if !yield(v) {
				cancel()
				stopped = true
			}
		}
	}
}

func (t tmdbimport) Run(gctx *cmdopts.Global) (err error) {
	c, err := tmdb.Init(t.APIKey)
	if err != nil {
		return errorsx.Wrap(err, "unable to initialize tmdb client")
	}
	c.SetClientAutoRetry()

	encoder := jsonl.NewEncoder(os.Stdout)

	for v := range t.movies(gctx.Context, c) {
		if err := encoder.Encode(v); err != nil {
			return errorsx.Wrap(err, "unable to encode media")
		}
	}

	if t.cause != nil {
		return t.cause
	}

	for v := range t.series(gctx.Context, c) {
		if err := encoder.Encode(v); err != nil {
			return errorsx.Wrap(err, "unable to encode media")
		}
	}

	return t.cause
}
