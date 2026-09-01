package ftuxapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/ftux"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

type HTTPOption func(*HTTP)

func HTTPOptionNoop(*HTTP) {}

func HTTPOptionHTTPClient(c *http.Client) HTTPOption {
	return func(t *HTTP) { t.httpc = c }
}

func HTTPOptionJWTSecret(j jwtx.SecretSource) HTTPOption {
	return func(t *HTTP) { t.jwtsecret = j }
}

func NewHTTP(q sqlx.Queryer, options ...HTTPOption) *HTTP {
	svc := langx.Clone(HTTP{
		q:         q,
		jwtsecret: env.JWTSecret,
	}, options...)

	return &svc
}

type HTTP struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	httpc     *http.Client
}

func (t *HTTP) Bind(r *mux.Router) {
	r.Path("/communities").Methods(http.MethodGet).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.defaults))

	r.Path("/communities").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout10s(),
	).ThenFunc(t.subscribe))
}

func (t *HTTP) defaults(w http.ResponseWriter, r *http.Request) {
	communities, err := ftux.PrepareDefaultCommunities()
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to load default communities"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	resp := &CommunitySuggestions{Community: communities}
	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
	}
}

func (t *HTTP) subscribe(w http.ResponseWriter, r *http.Request) {
	var req SubscribeCommunitiesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode subscribe request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if t.httpc == nil {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusServiceUnavailable))
		return
	}

	for _, cid := range req.CommunityId {
		if _, err := communityapi.SubscribeCommunity(r.Context(), t.q, t.httpc, cid); err != nil {
			log.Println(errorsx.Wrapf(err, "unable to subscribe to community - %s", cid))
			continue
		}
	}

	errorsx.Log(httpx.WriteJSON(w, httpx.GetBuffer(r), &SubscribeCommunitiesResponse{}))
}
