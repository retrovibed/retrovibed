package communityapi

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/envfile"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type HTTPPublisherEnvironmentOption func(*HTTPPublisherEnvironment)

func HTTPPublisherEnvironmentOptionJWTSecret(j jwtx.SecretSource) HTTPPublisherEnvironmentOption {
	return func(t *HTTPPublisherEnvironment) {
		t.jwtsecret = j
	}
}

// NewHTTPPublisherEnvironment proxies a publisher plugin's .env sidecar
// to and from HTTP as raw bytes - the same "no server-side parsing"
// convention ddiscapi.HTTPPluginEnvironment uses for search plugins, since
// quoting and comment-derived hints are a frontend concern.
//
// The one thing it does interpret is the plugin's own declaration: a GET
// returns the variables the plugin says it understands, with whatever has
// actually been configured filled in over the top. That is what lets the
// console render a settings form for a plugin it knows nothing about.
func NewHTTPPublisherEnvironment(q sqlx.Queryer, reg publishplugin.E, options ...HTTPPublisherEnvironmentOption) *HTTPPublisherEnvironment {
	svc := langx.Clone(HTTPPublisherEnvironment{
		q:         q,
		reg:       reg,
		jwtsecret: env.JWTSecret,
	}, options...)

	return &svc
}

type HTTPPublisherEnvironment struct {
	q         sqlx.Queryer
	reg       publishplugin.E
	jwtsecret jwtx.SecretSource
}

func (t *HTTPPublisherEnvironment) Bind(r *mux.Router) {
	r.StrictSlash(false)

	// a GET runs the plugin to read its declaration, unlike the write
	// paths which only touch the filesystem - hence the longer budget.
	// The module is already resident and compiled, so this is headroom
	// rather than an expected cost.
	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout10s(),
	).ThenFunc(t.get))

	r.Path("/{id}").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.update))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

// publisher resolves {id} to its catalog row. Going through the database
// rather than scanning publish.d for a matching filename is what makes a
// symlinked second copy of a plugin addressable: its row records the path
// it was installed at, and every one of a plugin's per-installation
// resources is derived from that path.
func (t *HTTPPublisherEnvironment) publisher(r *http.Request) (pub community.PluginPublisher, err error) {
	err = community.PluginPublisherFindByID(r.Context(), t.q, mux.Vars(r)["id"]).Scan(&pub)
	return pub, err
}

// declared asks the plugin which variables it understands. A plugin that
// predates the env command, or one whose module is missing, simply
// contributes nothing - the saved file is still served, so an operator
// never loses access to configuration they already wrote.
func (t *HTTPPublisherEnvironment) declared(r *http.Request, pub community.PluginPublisher) string {
	content, err := t.reg.Environment(r.Context(), pub.Path)
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to read declared plugin environment, serving the saved configuration alone: %s", pub.Path))
		return ""
	}

	return string(content)
}

func (t *HTTPPublisherEnvironment) get(w http.ResponseWriter, r *http.Request) {
	pub, err := t.publisher(r)
	if errors.Is(err, sql.ErrNoRows) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find plugin publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	path := publishplugin.EnvPath(pub.Path)
	saved, err := os.ReadFile(path)
	if fsx.IgnoreIsNotExist(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to load publisher environment"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	// the declaration supplies the keys and their descriptions, the saved
	// file supplies the values; Apply keeps each declared line's comment
	// while swapping in what was configured, and appends anything
	// configured that the plugin never declared.
	content := envfile.Apply(t.declared(r, pub), envfile.Parse(string(saved)))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(content)); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPPublisherEnvironment) update(w http.ResponseWriter, r *http.Request) {
	pub, err := t.publisher(r)
	if errors.Is(err, sql.ErrNoRows) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find plugin publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	content, err := io.ReadAll(io.LimitReader(r.Body, bytesx.MiB))
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to read request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = os.WriteFile(publishplugin.EnvPath(pub.Path), content, 0o600); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write publisher environment"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write(content); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPPublisherEnvironment) delete(w http.ResponseWriter, r *http.Request) {
	pub, err := t.publisher(r)
	if errors.Is(err, sql.ErrNoRows) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find plugin publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	path := publishplugin.EnvPath(pub.Path)
	content, err := os.ReadFile(path)
	if fsx.IgnoreIsNotExist(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to load publisher environment"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = os.Remove(path); fsx.IgnoreIsNotExist(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to remove publisher environment"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write(content); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
