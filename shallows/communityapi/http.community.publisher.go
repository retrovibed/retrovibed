package communityapi

import (
	"crypto/md5"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type HTTPCommunityPublisherOption func(*HTTPCommunityPublisher)

func HTTPCommunityPublisherOptionJWTSecret(j jwtx.SecretSource) HTTPCommunityPublisherOption {
	return func(t *HTTPCommunityPublisher) {
		t.jwtsecret = j
	}
}

// HTTPCommunityPublisherOptionDir overrides the directory installed publisher
// plugin binaries are written to (default: {config}/publish.d), letting
// callers (e.g. tests) point it at a directory of their choosing.
func HTTPCommunityPublisherOptionDir(dir string) HTTPCommunityPublisherOption {
	return func(t *HTTPCommunityPublisher) {
		t.dir = fsx.DirVirtual(dir)
	}
}

func NewHTTPCommunityPublisher(q sqlx.Queryer, options ...HTTPCommunityPublisherOption) *HTTPCommunityPublisher {
	svc := langx.Clone(HTTPCommunityPublisher{
		q:         q,
		dir:       fsx.DirVirtual(userx.DefaultConfigDir(userx.DefaultRelRoot(), "publish.d")),
		jwtsecret: env.JWTSecret,
	}, options...)

	return &svc
}

type HTTPCommunityPublisher struct {
	q         sqlx.Queryer
	dir       fsx.Virtual
	jwtsecret jwtx.SecretSource
}

func (t *HTTPCommunityPublisher) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout10s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

func (t *HTTPCommunityPublisher) search(w http.ResponseWriter, r *http.Request) {
	var resp SocialsSearchResponse

	q := community.PluginPublisherFindAll(r.Context(), t.q)
	qi := sqlx.Scan(q)
	for p := range qi.Iter() {
		resp.Catalog = append(resp.Catalog, NewPluginPublisher(PluginPublisherOptionFromDB(langx.Clone(p, timex.JSONSafeEncodeOption))))
	}
	if err := qi.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to list plugin publishers"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPCommunityPublisher) create(w http.ResponseWriter, r *http.Request) {
	var (
		existing community.PluginPublisher
	)

	description := r.FormValue("description")
	mimetype := r.FormValue("mimetype")
	if mimetype == "" {
		log.Println("mimetype is required")
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
		log.Println(errorsx.Wrap(err, "unable to create publisher plugin directory"))
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

	digest := md5.New()
	if _, err = io.Copy(io.MultiWriter(tmp, digest), f); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write temporary file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	id := md5x.FormatUUID(digest)
	dst := t.dir.Path(id)
	if err = os.Rename(tmp.Name(), dst); err != nil {
		log.Println(errorsx.Wrap(err, "unable to install publisher plugin"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	publisher := community.PluginPublisher{
		ID:          id,
		Path:        dst,
		Description: description,
		Mimetype:    mimetype,
	}

	if err = community.PluginPublisherInsertWithDefaults(r.Context(), t.q, publisher).Scan(&existing); err != nil {
		log.Println(errorsx.Wrap(err, "unable to insert plugin publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &PluginPublisherCreateResponse{
		Publisher: NewPluginPublisher(PluginPublisherOptionFromDB(langx.Clone(existing, timex.JSONSafeEncodeOption))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPCommunityPublisher) delete(w http.ResponseWriter, r *http.Request) {
	var existing community.PluginPublisher

	id := mux.Vars(r)["id"]

	if err := community.PluginPublisherDeleteByID(r.Context(), t.q, id).Scan(&existing); err != nil {
		log.Println(errorsx.Wrap(err, "unable to delete plugin publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	errorsx.Log(errorsx.Wrap(fsx.IgnoreIsNotExist(os.Remove(existing.Path)), "unable to remove publisher plugin from disk"))

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &PluginPublisherDeleteResponse{
		Publisher: NewPluginPublisher(PluginPublisherOptionFromDB(langx.Clone(existing, timex.JSONSafeEncodeOption))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
