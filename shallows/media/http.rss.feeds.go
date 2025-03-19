package media

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/x/errorsx"
	"github.com/retrovibed/retrovibed/internal/x/formx"
	"github.com/retrovibed/retrovibed/internal/x/httpx"
	"github.com/retrovibed/retrovibed/internal/x/jwtx"
	"github.com/retrovibed/retrovibed/internal/x/langx"
	"github.com/retrovibed/retrovibed/internal/x/md5x"
	"github.com/retrovibed/retrovibed/internal/x/numericx"
	"github.com/retrovibed/retrovibed/internal/x/sqlx"
	"github.com/retrovibed/retrovibed/internal/x/sqlxx"
	"github.com/retrovibed/retrovibed/internal/x/stringsx"
	"github.com/retrovibed/retrovibed/rss"
	"github.com/retrovibed/retrovibed/tracking"
)

type HTTPRSSFeedOption func(*HTTPRSSFeed)

func HTTPRSSFeedOptionJWTSecret(j jwtx.JWTSecretSource) HTTPRSSFeedOption {
	return func(t *HTTPRSSFeed) {
		t.jwtsecret = j
	}
}

func NewHTTPRSSFeed(q sqlx.Queryer, options ...HTTPRSSFeedOption) *HTTPRSSFeed {
	svc := langx.Clone(HTTPRSSFeed{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPRSSFeed struct {
	q         sqlx.Queryer
	jwtsecret jwtx.JWTSecretSource
	decoder   *form.Decoder
}

func (t *HTTPRSSFeed) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

func (t *HTTPRSSFeed) search(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg rss.FeedSearchResponse = rss.FeedSearchResponse{
			Next: &rss.FeedSearchRequest{
				Limit: 100,
			},
		}
	)

	if err = t.decoder.Decode(msg.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	msg.Next.Limit = numericx.Min(msg.Next.Limit, 100)

	q := tracking.RSSSearchBuilder().Where(squirrel.And{
		tracking.RSSQuerySearch(msg.Next.Query),
		squirrel.Expr("1=1"),
	}).OrderBy("created_at DESC").Offset(msg.Next.Offset * msg.Next.Limit).Limit(msg.Next.Limit)

	err = sqlxx.ScanEach(tracking.RSSSearch(r.Context(), t.q, q), func(p *tracking.RSS) error {
		tmp := langx.Clone(rss.Feed{}, rss.FeedOptionFromTorrentRSS(langx.Clone(*p, tracking.RSSOptionJSONSafeEncode)))
		msg.Items = append(msg.Items, &tmp)
		return nil
	})

	if err != nil {
		log.Println(errorsx.Wrap(err, "encoding failed"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
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
		req rss.FeedUpdateRequest = rss.FeedUpdateRequest{}
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	feed := tracking.RSS{
		ID:           stringsx.DefaultIfBlank(req.Feed.Id, md5x.FormatString(md5x.Digest(req.Feed.Url))),
		Description:  req.Feed.Description,
		URL:          req.Feed.Url,
		Autodownload: req.Feed.Autodownload,
		Autoarchive:  req.Feed.Autoarchive,
		Contributing: req.Feed.Contributing,
	}

	if err = tracking.RSSInsertWithDefaults(r.Context(), t.q, feed).Scan(&feed); err != nil {
		log.Println(errorsx.Wrap(err, "encoding failed"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	tmp := langx.Clone(rss.Feed{}, rss.FeedOptionFromTorrentRSS(langx.Clone(feed, tracking.RSSOptionJSONSafeEncode)))

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &rss.FeedCreateResponse{Feed: &tmp}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRSSFeed) delete(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		req  rss.FeedUpdateRequest = rss.FeedUpdateRequest{}
		feed tracking.RSS
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = tracking.RSSDeleteByID(r.Context(), t.q, req.Feed.Id).Scan(&feed); err != nil {
		log.Println(errorsx.Wrap(err, "encoding failed"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	tmp := langx.Clone(rss.Feed{}, rss.FeedOptionFromTorrentRSS(langx.Clone(feed, tracking.RSSOptionJSONSafeEncode)))

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &rss.FeedDeleteResponse{Feed: &tmp}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
