package communityapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
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

func NewHTTP(q sqlx.Queryer, options ...HTTPOption) *HTTP {
	svc := langx.Clone(HTTP{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
		lucene:    duckdbx.NewLucene(),
	}, options...)

	return &svc
}

type HTTP struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
	httpc     *http.Client
	lucene    lucenex.Driver
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

	err = community.CommunityFindByID(r.Context(), t.q, cid).Scan(&existing)
	if sqlx.IgnoreNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to look up subscription"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err == nil {
		// already subscribed — toggle off
		if err = community.CommunityDeleteByID(r.Context(), t.q, cid).Scan(&existing); err != nil {
			log.Println(errorsx.Wrap(err, "unable to delete subscription"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		var (
			feed    tracking.RSS
			feedURL = fmt.Sprintf("https://%s.community.retrovibe.space", com.Domain)
		)
		if err = tracking.RSSDeleteByURL(r.Context(), t.q, feedURL).Scan(&feed); sqlx.IgnoreNoRows(err) != nil {
			log.Println(errorsx.Wrap(err, "unable to delete rss feed"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}
	} else {
		// not subscribed — subscribe
		var sub = community.Community{
			ID: cid,
		}

		if err = community.CommunityUpsertAutoDownload(r.Context(), t.q, sub).Scan(&sub); err != nil {
			log.Println(errorsx.Wrap(err, "unable to insert subscription"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		var (
			feedURL = fmt.Sprintf("https://%s.community.retrovibe.space", com.Domain)
		)

		feed := tracking.NewFeedRSS(
			"",
			tracking.RSSOptionURL(feedURL),
			tracking.RSSOptionDescription(stringsx.Join(" - ", slicesx.Filter(stringsx.Present, com.Domain, com.Description)...)),
			tracking.RSSOptionEncryptionSeed(com.Entropy),
			tracking.RSSOptionAutodownload(true),
			tracking.RSSOptionAutoarchive(true),
			tracking.RSSOptionAutoID,
		)

		if err = tracking.RSSInsertWithDefaults(r.Context(), t.q, feed).Scan(&feed); err != nil {
			log.Println(errorsx.Wrap(err, "unable to register rss feed"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}
	}

	errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusOK))
}
