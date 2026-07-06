package ddiscapi

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"net/netip"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type HTTPDiscoveryOption func(*HTTPDiscovery)

func HTTPDiscoveryOptionJWTSecret(j jwtx.SecretSource) HTTPDiscoveryOption {
	return func(t *HTTPDiscovery) {
		t.jwtsecret = j
	}
}

func NewHTTPDiscovery(q sqlx.Queryer, options ...HTTPDiscoveryOption) *HTTPDiscovery {
	svc := langx.Clone(HTTPDiscovery{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPDiscovery struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPDiscovery) Bind(r *mux.Router) {
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
		metaapi.AuthzTokenHTTP(t.jwtsecret, AuthzPermPeer),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, AuthzPermPeer),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

func (t *HTTPDiscovery) search(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp = DiscoverySearchResponse{
			Next: &DiscoverySearchRequest{
				Offset:      0,
				Limit:       100,
				NextCheck:   meta.NewDateRange(timex.NewRangeEverything()),
				AttemptsMax: math.MaxInt64,
			},
		}
	)

	if err = t.decoder.Decode(resp.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	resp.Next.Limit = numericx.Min(resp.Next.Limit, 1024)

	query := tracking.UnknownSearchBuilder().
		Offset(resp.Next.Offset * resp.Next.Limit).Limit(resp.Next.Limit).
		Where(squirrel.And{
			tracking.UnknownHashQueryAttemptsRange(resp.Next.AttemptsMin, resp.Next.AttemptsMax),
			tracking.UnknownHashQueryByIDs(resp.Next.Id...),
			tracking.UnknownHashQueryNextCheck(meta.TimexRange(resp.Next.NextCheck, timex.NewRangeEverything())),
		})

	q := sqlx.Scan(tracking.UnknownSearch(r.Context(), t.q, query))
	for uh := range q.Iter() {
		var (
			encoded *Discovery
		)

		if encoded, err = NewDiscoveryFromTrackingUnknownHash(uh); err != nil {
			log.Println(errorsx.Wrap(err, "discovery encoding failed"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		resp.Items = append(resp.Items, encoded)
	}

	if err = q.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "discovery generation failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovery) create(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg DiscoveryCreateRequest
		uh  tracking.UnknownHash
	)

	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	uh = tracking.NewUnknownHash(
		int160.FromBytesOrZero(msg.Discovery.GetInfohash()),
		func(uh *tracking.UnknownHash) { uh.IP = netip.IPv6Unspecified() },
	)

	if err = tracking.UnknownHashInsertWithDefaults(r.Context(), t.q, uh).Scan(&uh); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create discovery entry"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	dd, err := NewDiscoveryFromTrackingUnknownHash(uh)
	if err != nil {
		log.Println(errorsx.Wrap(err, "conversion failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DiscoveryCreateResponse{
		Discovery: dd,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovery) delete(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		vars = mux.Vars(r)
		uh   tracking.UnknownHash
		dd   *Discovery
	)

	if err := tracking.UnknownHashDeleteByID(r.Context(), t.q, vars["id"]).Scan(&uh); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if dd, err = NewDiscoveryFromTrackingUnknownHash(uh); err != nil {
		log.Println(errorsx.Wrap(err, "conversion failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DiscoveryDeleteResponse{
		Discovery: dd,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
