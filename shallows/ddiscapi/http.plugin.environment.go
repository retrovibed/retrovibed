package ddiscapi

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/bytesx"
	"github.com/retrovibed/retrovibed/retroapi/envfile"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

// HTTPPluginEnvironment proxies a search plugin's <name>.env file to/from
// HTTP as raw bytes - the same "no server-side parsing" convention as
// metaapi.HTTPFileConfig, except the target path is resolved per-request
// from {name} rather than fixed at construction, since there's one .env
// per plugin rather than one file for the whole feature. Comment-derived
// hints, quoting, etc. are a frontend concern; the server never interprets
// the content.
//
// The one exception: a GET asks the plugin which variables it understands
// (its "env" command) and fills whatever has been configured in over the
// top, so the console can render a settings form for a plugin it knows nothing
// about.
type HTTPPluginEnvironmentOption func(*HTTPPluginEnvironment)

func HTTPPluginEnvironmentOptionJWTSecret(j jwtx.SecretSource) HTTPPluginEnvironmentOption {
	return func(t *HTTPPluginEnvironment) {
		t.jwtsecret = j
	}
}

// HTTPPluginEnvironmentOptionDir overrides the search plugin directory
// (default: searchplugin.SearchPluginDir(userx.DefaultConfigDir(userx.DefaultRelRoot()))),
// letting callers (e.g. tests) point it at a directory of their choosing
// instead of the real, hardcoded userx config directory.
func HTTPPluginEnvironmentOptionDir(dir string) HTTPPluginEnvironmentOption {
	return func(t *HTTPPluginEnvironment) {
		t.dir = fsx.DirVirtual(dir)
	}
}

func NewHTTPPluginEnvironment(reg searchplugin.E, options ...HTTPPluginEnvironmentOption) *HTTPPluginEnvironment {
	svc := langx.Clone(HTTPPluginEnvironment{
		dir:       fsx.DirVirtual(searchplugin.SearchPluginDir(userx.DefaultConfigDir(userx.DefaultRelRoot()))),
		reg:       reg,
		jwtsecret: env.JWTSecret,
	}, options...)

	return &svc
}

type HTTPPluginEnvironment struct {
	dir       fsx.Virtual
	reg       searchplugin.E
	jwtsecret jwtx.SecretSource
}

func (t *HTTPPluginEnvironment) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
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

// path resolves {id} to the plugin's .env file, keying off the same
// installed *.wasm plugin the id was originally derived from - the
// environment for a plugin is only addressable once the plugin exists.
func (t *HTTPPluginEnvironment) path(r *http.Request) (string, error) {
	name, err := resolvePluginName(t.dir, mux.Vars(r)["id"])
	if err != nil {
		return "", err
	}

	return t.dir.Path(name + ".env"), nil
}

// declared asks the plugin which variables it understands. A plugin that
// predates the env command, or one that isn't loaded, simply contributes
// nothing - the saved file is still served, so an operator never loses
// access to configuration they already wrote.
func (t *HTTPPluginEnvironment) declared(r *http.Request, envpath string) string {
	wasmpath := strings.TrimSuffix(envpath, ".env") + ".wasm"
	content, err := t.reg.Environment(r.Context(), wasmpath)
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to read declared plugin environment, serving the saved configuration alone: %s", wasmpath))
		return ""
	}

	return string(content)
}

func (t *HTTPPluginEnvironment) get(w http.ResponseWriter, r *http.Request) {
	path, err := t.path(r)
	if os.IsNotExist(err) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	saved, err := os.ReadFile(path)
	if fsx.IgnoreIsNotExist(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to load plugin environment"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	// the declaration supplies the keys and their descriptions, the saved
	// file supplies the values; Apply keeps each declared line's comment
	// while swapping in what was configured, and appends anything
	// configured that the plugin never declared.
	content := envfile.Apply(t.declared(r, path), envfile.Parse(string(saved)))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(content)); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPPluginEnvironment) update(w http.ResponseWriter, r *http.Request) {
	path, err := t.path(r)
	if os.IsNotExist(err) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	content, err := io.ReadAll(io.LimitReader(r.Body, bytesx.MiB))
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to read request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = t.dir.MkDirAll("", 0o700); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create search plugin directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = os.WriteFile(path, content, 0o600); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write plugin environment"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write(content); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPPluginEnvironment) delete(w http.ResponseWriter, r *http.Request) {
	path, err := t.path(r)
	if os.IsNotExist(err) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	content, err := os.ReadFile(path)
	if fsx.IgnoreIsNotExist(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to load plugin environment"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = os.Remove(path); fsx.IgnoreIsNotExist(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to remove plugin environment"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write(content); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
