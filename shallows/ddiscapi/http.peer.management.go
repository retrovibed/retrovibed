package ddiscapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/shallows/authn"
	"github.com/retrovibed/retrovibed/shallows/httpauth"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type HTTPPeerManagementOption func(*HTTPPeerManagement)

func HTTPPeerManagementOptionJWTSecret(j jwtx.SecretSource) HTTPPeerManagementOption {
	return func(t *HTTPPeerManagement) {
		t.jwtsecret = j
	}
}

func NewHTTPPeerManagement(q sqlx.Queryer, options ...HTTPPeerManagementOption) *HTTPPeerManagement {
	svc := langx.Clone(HTTPPeerManagement{
		q:         q,
		jwtsecret: authn.JWTSecretFromEnv,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPPeerManagement struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPPeerManagement) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, AuthzPermPeerManagement),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, AuthzPermPeerManagement),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))

	r.Path("/{id}").Methods(http.MethodPatch).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, AuthzPermPeerManagement),
		httpx.Timeout2s(),
	).ThenFunc(t.update))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.find))

	// r.Path("/{id}").Methods(http.MethodPut).Handler(alice.New(
	// 	httpx.ContextBufferPool512(),
	// 	httpx.ParseForm,
	// 	httpauth.AuthenticateWithToken(t.jwtsecret),
	// 	httpx.Timeout2s(),
	// ).ThenFunc(t.sync))
}

func (t *HTTPPeerManagement) search(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp = PeerSearchResponse{
			Next: &PeerSearchRequest{
				Offset: 0,
				Limit:  100,
			},
		}
	)

	if err = t.decoder.Decode(resp.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	resp.Next.Limit = numericx.Min(resp.Next.Limit, 100)

	query := tracking.PeerSearchBuilder().
		Where(
			squirrel.And{
				tracking.PeerQueryDdiscEnabled(),
				lucenex.Query(duckdbx.NewLucene(), resp.Next.Query, lucenex.WithDefaultField("description")),
			},
		).
		Offset(resp.Next.Offset * resp.Next.Limit).Limit(resp.Next.Limit)

	q := sqlx.Scan(tracking.PeerSearch(r.Context(), t.q, query))
	for p := range q.Iter() {
		var (
			encoded *Peer
		)

		if encoded, err = NewPeerFromTrackingPeer(p); err != nil {
			log.Println(errorsx.Wrap(err, "peer encoding failed"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		resp.Items = append(resp.Items, encoded)
	}

	if err = q.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "peer generation failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPPeerManagement) find(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp = PeerFindResponse{
			Peer: &Peer{},
		}
		vars = mux.Vars(r)
		p    tracking.Peer
	)

	if err = tracking.PeerFindByID(r.Context(), t.q, vars["id"]).Scan(&p); err != nil {
		log.Println(errorsx.Wrap(err, "unable to find peer"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	}

	if resp.Peer, err = NewPeerFromTrackingPeer(p); err != nil {
		log.Println(errorsx.Wrap(err, "failed to encode peer"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPPeerManagement) delete(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		vars = mux.Vars(r)
		p    tracking.Peer
		pp   *Peer
	)

	if err := tracking.PeerDeleteByID(r.Context(), t.q, vars["id"]).Scan(&p); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if pp, err = NewPeerFromTrackingPeer(p); err != nil {
		log.Println(errorsx.Wrap(err, "conversion failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &PeerDeleteResponse{
		Peer: pp,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPPeerManagement) update(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		msg  PeerUpdateRequest
		p    tracking.Peer
		vars = mux.Vars(r)
	)

	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if p, err = NewTrackingPeerFromPeer(msg.Peer); err != nil {
		log.Println(errorsx.Wrap(err, "converting peer failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = tracking.PeerUpdateDdisc(r.Context(), t.q, vars["id"], p).Scan(&p); err != nil {
		log.Println(errorsx.Wrap(err, "unable to insert peer"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	pp, err := NewPeerFromTrackingPeer(p)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to encode peer"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &PeerUpdateResponse{
		Peer: pp,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPPeerManagement) create(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg PeerCreateRequest
		p   tracking.Peer
	)

	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if p, err = NewTrackingPeerFromPeer(msg.Peer, tracking.PeerOptionTombstone(timex.Inf())); err != nil {
		log.Println(errorsx.Wrap(err, "unable to transform peer"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = tracking.PeerInsertWithDefaults(r.Context(), t.q, p).Scan(&p); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create peer"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	pp, err := NewPeerFromTrackingPeer(langx.Clone(p, timex.JSONSafeEncodeOption, timex.UTCEncodeOption))
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to encode peer"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &PeerCreateResponse{
		Peer: pp,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

// TODO: sync is to manually kick off a synchronization task for the peer.
// func (t *HTTPPeerManagement) sync(w http.ResponseWriter, r *http.Request) {
// 	// var (
// 	// 	err  error
// 	// 	resp = PeerSyncResponse{
// 	// 		Peer: &Peer{},
// 	// 	}
// 	// 	vars = mux.Vars(r)
// 	// 	p    tracking.Peer
// 	// )

// 	// if err = tracking.PeerFindByID(r.Context(), t.q, vars["id"]).Scan(&p); err != nil {
// 	// 	log.Println(errorsx.Wrap(err, "unable to find peer"))
// 	// 	errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
// 	// 	return
// 	// }

// 	// if resp.Peer, err = NewPeerFromTrackingPeer(p); err != nil {
// 	// 	log.Println(errorsx.Wrap(err, "failed to encode peer"))
// 	// 	errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
// 	// 	return
// 	// }

// 	// if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
// 	// 	log.Println(errorsx.Wrap(err, "unable to write response"))
// 	// 	return
// 	// }
// }
