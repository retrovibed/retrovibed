package ddiscapi

import (
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type HTTPPluginManagementOption func(*HTTPPluginManagement)

func HTTPPluginManagementOptionJWTSecret(j jwtx.SecretSource) HTTPPluginManagementOption {
	return func(t *HTTPPluginManagement) {
		t.jwtsecret = j
	}
}

// HTTPPluginManagementOptionDir overrides the search plugin directory
// (default: searchplugin.SearchPluginDir(userx.DefaultConfigDir(userx.DefaultRelRoot()))),
// letting callers (e.g. tests) point it at a directory of their choosing
// instead of the real, hardcoded userx config directory.
func HTTPPluginManagementOptionDir(dir string) HTTPPluginManagementOption {
	return func(t *HTTPPluginManagement) {
		t.dir = fsx.DirVirtual(dir)
	}
}

func NewHTTPPluginManagement(reg *searchplugin.Registry, options ...HTTPPluginManagementOption) *HTTPPluginManagement {
	svc := langx.Clone(HTTPPluginManagement{
		reg:       reg,
		dir:       fsx.DirVirtual(searchplugin.SearchPluginDir(userx.DefaultConfigDir(userx.DefaultRelRoot()))),
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPPluginManagement struct {
	reg       *searchplugin.Registry
	dir       fsx.Virtual
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPPluginManagement) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.find))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

// pluginID derives the stable, opaque id used in URLs from a plugin's
// on-disk name, so clients never need to pass the raw name back to
// find/delete/environment endpoints.
func pluginID(name string) string {
	return md5x.String(name)
}

// resolvePluginName scans dir for the *.wasm entry whose name hashes to id,
// returning its bare name (no extension). Returns an fs.ErrNotExist-wrapping
// error if no entry matches, or the underlying error if the directory
// itself couldn't be read.
func resolvePluginName(dir fsx.Virtual, id string) (string, error) {
	entries, err := os.ReadDir(dir.Path())
	if err != nil {
		return "", err
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".wasm") {
			continue
		}

		name := strings.TrimSuffix(e.Name(), ".wasm")
		if pluginID(name) == id {
			return name, nil
		}
	}

	return "", os.ErrNotExist
}

func pluginFromFileInfo(name string, info os.FileInfo) *Plugin {
	return &Plugin{
		Id:          pluginID(name),
		Name:        name,
		Size:        uint64(info.Size()),
		InstalledAt: info.ModTime().UTC().Format(time.RFC3339),
	}
}

func (t *HTTPPluginManagement) search(w http.ResponseWriter, r *http.Request) {
	const resplimit = 1024
	var (
		err  error
		resp = PluginSearchResponse{
			Next: &PluginSearchRequest{
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
	resp.Next.Limit = numericx.Min(resp.Next.Limit, resplimit)

	entries, err := os.ReadDir(t.dir.Path())
	if err != nil && !os.IsNotExist(err) {
		log.Println(errorsx.Wrap(err, "unable to list search plugins"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".wasm") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)

	offset := resp.Next.Offset * resp.Next.Limit
	for i, name := range names {
		if uint64(i) < offset {
			continue
		}
		if uint64(len(resp.Items)) >= resp.Next.Limit {
			break
		}

		info, err := os.Stat(t.dir.Path(name))
		if err != nil {
			continue
		}

		resp.Items = append(resp.Items, pluginFromFileInfo(strings.TrimSuffix(name, ".wasm"), info))
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPPluginManagement) find(w http.ResponseWriter, r *http.Request) {
	name, err := resolvePluginName(t.dir, mux.Vars(r)["id"])
	if os.IsNotExist(err) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to resolve search plugin"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	info, err := os.Stat(t.dir.Path(name + ".wasm"))
	if os.IsNotExist(err) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find search plugin"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &PluginFindResponse{
		Plugin: pluginFromFileInfo(name, info),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPPluginManagement) create(w http.ResponseWriter, r *http.Request) {
	var err error

	name := searchplugin.SanitizeName(r.FormValue("name"))
	if name == "" {
		log.Println("invalid plugin name", name)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	f, _, err := r.FormFile("content")
	if err != nil {
		log.Println(errorsx.Wrap(err, "content parameter required"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	defer f.Close()

	if err = t.dir.MkDirAll("", 0o700); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create search plugin directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	tmp, err := fsx.CreateTemp(t.dir, "retrovibed.upload.*")
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to create temporary file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}
	defer func() {
		if err == nil {
			return
		}

		log.Println("failure receiving upload, removing attempt", err)
		errorsx.Log(errorsx.Wrap(fsx.IgnoreIsNotExist(os.Remove(tmp.Name())), "unable to remove tmp"))
	}()
	defer tmp.Close()

	if _, err = io.Copy(tmp, f); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write temporary file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = searchplugin.VerifyWasmMagic(tmp.Name()); err != nil {
		log.Println(errorsx.Wrap(err, "invalid wasm module"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	dst := t.dir.Path(name + ".wasm")
	if err = os.Rename(tmp.Name(), dst); err != nil {
		log.Println(errorsx.Wrap(err, "unable to install search plugin"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = t.reg.Load(r.Context(), dst); err != nil {
		log.Println(errorsx.Wrap(err, "unable to load search plugin"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	info, err := os.Stat(dst)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to stat installed search plugin"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &PluginCreateResponse{
		Plugin: pluginFromFileInfo(name, info),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPPluginManagement) delete(w http.ResponseWriter, r *http.Request) {
	name, err := resolvePluginName(t.dir, mux.Vars(r)["id"])
	if os.IsNotExist(err) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to resolve search plugin"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	path := t.dir.Path(name + ".wasm")
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find search plugin"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	t.reg.Unload(path)

	if err = os.Remove(path); err != nil {
		log.Println(errorsx.Wrap(err, "unable to remove search plugin"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &PluginDeleteResponse{
		Plugin: pluginFromFileInfo(name, info),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
