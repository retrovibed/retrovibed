package media

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/rss"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type HTTPRSSFeedOption func(*HTTPRSSFeed)

func HTTPRSSFeedOptionJWTSecret(j jwtx.SecretSource) HTTPRSSFeedOption {
	return func(t *HTTPRSSFeed) {
		t.jwtsecret = j
	}
}

func NewHTTPRSSFeed(q sqlx.Queryer, options ...HTTPRSSFeedOption) *HTTPRSSFeed {
	svc := langx.Clone(HTTPRSSFeed{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
		fts:       duckdbx.NewLucene(),
	}, options...)

	return &svc
}

type HTTPRSSFeed struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
	fts       lucenex.Driver
}

func (t *HTTPRSSFeed) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

func (t *HTTPRSSFeed) search(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg = rss.FeedSearchResponse{
			Next: &rss.FeedSearchRequest{
				Limit: 100,
			},
		}
	)

	if err = t.decoder.Decode(msg.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	msg.Next.Limit = numericx.Min(msg.Next.Limit, 100)

	q := tracking.RSSSearchBuilder().Where(squirrel.And{
		lucenex.Query(t.fts, msg.Next.Query, lucenex.WithDefaultField("description")),
		squirrel.Expr("1=1"),
	}).OrderBy("created_at DESC").Offset(msg.Next.Offset * msg.Next.Limit).Limit(msg.Next.Limit)

	qs := sqlx.Scan(tracking.RSSSearch(r.Context(), t.q, q))
	for p := range qs.Iter() {
		tmp := langx.Clone(rss.Feed{}, rss.FeedOptionFromTorrentRSS(langx.Clone(p, timex.JSONSafeEncodeOption)))
		msg.Items = append(msg.Items, &tmp)
	}

	if err = qs.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "encoding failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRSSFeed) create(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req = rss.FeedCreateRequest{}
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	feed := tracking.NewFeedRSS(
		langx.FirstNonZero(req.Feed.Id, md5x.FormatUUID(md5x.Digest(req.Feed.Url))),
		NewTrackingFeedRSSFromFeedRSS(&req),
		timex.JSONSafeDecodeOption,
	)

	if err = tracking.RSSInsertWithDefaults(r.Context(), t.q, feed).Scan(&feed); err != nil {
		log.Println(errorsx.Wrap(err, "unable to record failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}
	feed = langx.Clone(feed, timex.JSONSafeEncodeOption)

	tmp := langx.Clone(rss.Feed{}, rss.FeedOptionFromTorrentRSS(feed))
	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &rss.FeedCreateResponse{Feed: &tmp}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRSSFeed) delete(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		feed tracking.RSS
		id   = mux.Vars(r)["id"]
	)

	if err = tracking.RSSDeleteByID(r.Context(), t.q, id).Scan(&feed); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "delete failed: not found"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "delete failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	tmp := langx.Clone(rss.Feed{}, rss.FeedOptionFromTorrentRSS(langx.Clone(feed, timex.JSONSafeEncodeOption)))

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &rss.FeedDeleteResponse{Feed: &tmp}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
