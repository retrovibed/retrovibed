package metaapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/httpauth"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/formx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/numericx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/sqlxx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/meta"
)

type HTTPDaemonsOption func(*HTTPDaemons)

func HTTPDaemonsOptionJWTSecret(j jwtx.JWTSecretSource) HTTPDaemonsOption {
	return func(t *HTTPDaemons) {
		t.jwtsecret = j
	}
}

func NewHTTPDaemons(q sqlx.Queryer, options ...HTTPDaemonsOption) *HTTPDaemons {
	svc := langx.Clone(HTTPDaemons{
		q:         q,
		jwtsecret: authn.JWTSecretFromEnv,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPDaemons struct {
	q         sqlx.Queryer
	jwtsecret jwtx.JWTSecretSource
	decoder   *form.Decoder
}

func (t *HTTPDaemons) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.create))
}

func (t *HTTPDaemons) search(w http.ResponseWriter, r *http.Request) {
	const resplimit = 128
	var (
		err  error
		resp = DaemonSearchResponse{
			Next: &DaemonSearchRequest{
				Offset: 0,
				Limit:  resplimit,
			},
		}
	)

	if err = t.decoder.Decode(resp.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	resp.Next.Limit = numericx.Min(resp.Next.Limit, resplimit)

	query := meta.DaemonSearchBuilder().
		Offset(resp.Next.Offset * resp.Next.Limit).Limit(resp.Next.Limit)

	err = sqlxx.ScanEach(meta.DaemonSearch(r.Context(), t.q, query), func(v *meta.Daemon) (failed error) {
		var (
			encoded *Daemon
		)

		if encoded, failed = NewDaemonFromMetaDaemon(*v); failed != nil {
			return failed
		}

		resp.Items = append(resp.Items, encoded)
		return nil
	})

	if err != nil {
		log.Println(errorsx.Wrap(err, "response generation failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	resp.Next.Offset += 1

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDaemons) create(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg DaemonCreateRequest
		v   meta.Daemon
	)

	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if v, err = NewMetadaemonFromDaemon(msg.Daemon, meta.DaemonOptionMaybeID, meta.DaemonOptionEnsureDescription, timex.JSONSafeDecodeOption); err != nil {
		log.Println(errorsx.Wrap(err, "converting data failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = meta.DaemonInsertWithDefaults(r.Context(), t.q, v).Scan(&v); err != nil {
		log.Println(errorsx.Wrap(err, "unable to insert record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DaemonCreateResponse{
		Daemon: errorsx.Must(NewDaemonFromMetaDaemon(v)),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
