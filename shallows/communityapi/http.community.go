package communityapi

import (
	"context"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type HTTPOption func(*HTTP)

func HTTPOptionNoop(*HTTP) {}

func HTTPOptionJWTSecret(j jwtx.SecretSource) HTTPOption {
	return func(t *HTTP) {
		t.jwtsecret = j
	}
}

func HTTPOptionHTTPClient(c *http.Client) HTTPOption {
	return func(t *HTTP) {
		t.httpc = c
	}
}

func HTTPOptionCommunitySync(w *asyncx.Wakeup) HTTPOption {
	return func(t *HTTP) {
		t.communitysync = w
	}
}

func NewHTTP(q sqlx.Queryer, options ...HTTPOption) *HTTP {
	svc := langx.Clone(HTTP{
		q:             q,
		jwtsecret:     env.JWTSecret,
		decoder:       formx.NewDecoder(),
		lucene:        duckdbx.NewLucene(),
		communitysync: asyncx.NewWakeup(context.Background()),
	}, options...)

	return &svc
}

type HTTP struct {
	q             sqlx.Queryer
	jwtsecret     jwtx.SecretSource
	decoder       *form.Decoder
	httpc         *http.Client
	lucene        lucenex.Driver
	communitysync *asyncx.Wakeup
}

func (t *HTTP) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/{id}/subscribe").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.subscribe))

	r.Path("/{id}/resync").Methods(http.MethodPost).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout10s(),
	).ThenFunc(t.resync))
}

func (t *HTTP) search(w http.ResponseWriter, r *http.Request) {
	if t.httpc == nil {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusServiceUnavailable))
		return
	}

	var req CommunitySearchRequest
	req.Limit = 100
	if err := t.decoder.Decode(&req, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode search request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	req.Limit = numericx.Min(req.Limit, 100)

	client := NewDeeppoolCommunity(t.httpc)
	resp, err := client.Search(r.Context(), req.Query, req.Offset, req.Limit)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to search communities"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadGateway))
		return
	}

	for _, c := range resp.Items {
		var (
			local community.Community
		)

		if err := community.CommunityInsertWithDefaults(r.Context(), t.q, CommunityFromDeeppool(c)).Scan(&local); err != nil {
			log.Println(errorsx.Wrap(err, "unable to upsert community"))
			continue
		}

		*c = langx.Autoderef(NewCommunity(CommunityOptionFromDB(langx.Clone(local, timex.JSONSafeEncodeOption))))
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTP) subscribe(w http.ResponseWriter, r *http.Request) {
	var (
		existing community.Community
		cid      = mux.Vars(r)["id"]
	)

	client := NewDeeppoolCommunity(t.httpc)
	com, err := client.Find(r.Context(), cid)
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to find community from deeppool - %s", cid))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err := community.CommunityInsertWithDefaults(r.Context(), t.q, CommunityFromDeeppool(com)).Scan(&existing); err != nil {
		log.Println(errorsx.Wrap(err, "unable to upsert community"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	subscribed := existing.SubscribedAt.Before(timex.Inf())

	if subscribed {
		// already subscribed — toggle off
		if err = community.CommunityUnsubscribe(r.Context(), t.q, cid).Scan(&existing); err != nil {
			log.Println(errorsx.Wrap(err, "unable to unsubscribe"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		var (
			feed tracking.RSS
		)
		if err = tracking.RSSDeleteByURL(r.Context(), t.q, existing.URL).Scan(&feed); sqlx.IgnoreNoRows(err) != nil {
			log.Println(errorsx.Wrap(err, "unable to delete rss feed"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}
	} else {
		// not subscribed — subscribe
		if existing, err = SubscribeCommunity(r.Context(), t.q, t.httpc, cid); err != nil {
			log.Println(errorsx.Wrap(err, "unable to subscribe to community"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}
	}

	errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusOK))
}

func (t *HTTP) resync(w http.ResponseWriter, r *http.Request) {
	cid := mux.Vars(r)["id"]

	if t.httpc == nil {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusServiceUnavailable))
		return
	}

	existing, err := ResyncOne(r.Context(), t.q, t.httpc, cid)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to resync community"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadGateway))
		return
	}

	t.communitysync.Broadcast()

	var msg PublishedContentSearchResponse
	msg.Community = new(langx.Clone(Community{}, CommunityOptionFromDB(langx.Clone(*existing, timex.JSONSafeEncodeOption))))
	msg.Next = &PublishedContentSearchRequest{CommunityId: cid, Offset: 1, Limit: 128}

	q := community.PublishedContentSearch(r.Context(), t.q, community.PublishedContentSearchBuilder().Where(
		squirrel.And{
			community.PublishedContentQueryCommunityID(cid),
			community.PublishedContentQueryNotTombstoned(),
		},
	).OrderBy("published_at DESC").Limit(128))

	qi := sqlx.Scan(q)
	for pc := range qi.Iter() {
		tmp := langx.Clone(PublishedContent{}, PublishedContentOptionFromDB(langx.Clone(pc, timex.JSONSafeEncodeOption)))
		msg.Items = append(msg.Items, &tmp)
	}
	if err := qi.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to fetch published content"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
