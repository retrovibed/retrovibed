package metaapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

type HTTPFileConfigOption func(*HTTPFileConfig)

func HTTPFileConfigOptionJWTSecret(j jwtx.SecretSource) HTTPFileConfigOption {
	return func(t *HTTPFileConfig) {
		t.jwtsecret = j
	}
}

func NewHTTPFileConfig(path string, options ...HTTPFileConfigOption) *HTTPFileConfig {
	svc := langx.Clone(HTTPFileConfig{
		path:      path,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPFileConfig struct {
	path      string
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPFileConfig) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.get))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.update))

	r.Path("/").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

func (t *HTTPFileConfig) get(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		encoded json.RawMessage
	)

	if encoded, err = fsx.AutoCached(t.path, func() ([]byte, error) { return nil, nil }); err != nil {
		log.Println(errorsx.Wrap(err, "unable to load file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if len(encoded) == 0 && encoded != nil {
		encoded = nil
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), json.RawMessage(encoded)); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPFileConfig) update(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		encoded json.RawMessage
	)

	if encoded, err = io.ReadAll(io.LimitReader(r.Body, bytesx.MiB)); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = os.WriteFile(t.path, encoded, 0600); err != nil {
		log.Println(errorsx.Wrap(err, "unable to insert record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), encoded); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPFileConfig) delete(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		encoded json.RawMessage
	)

	if encoded, err = fsx.AutoCached(t.path, func() ([]byte, error) { return nil, nil }); err != nil {
		log.Println(errorsx.Wrap(err, "unable to load file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = os.Remove(t.path); err != nil {
		log.Println(errorsx.Wrap(err, "unable to delete file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), json.RawMessage(encoded)); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
