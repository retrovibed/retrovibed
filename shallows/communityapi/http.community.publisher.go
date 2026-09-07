package communityapi

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
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

func NewHTTPCommunityPublisher(q sqlx.Queryer, reg *publishplugin.Registry, options ...HTTPCommunityPublisherOption) *HTTPCommunityPublisher {
	svc := langx.Clone(HTTPCommunityPublisher{
		q:         q,
		reg:       reg,
		dir:       fsx.DirVirtual(userx.DefaultConfigDir(userx.DefaultRelRoot(), "publish.d")),
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
		lucene:    duckdbx.NewLucene(),
	}, options...)

	return &svc
}

type HTTPCommunityPublisher struct {
	q         sqlx.Queryer
	reg       *publishplugin.Registry
	dir       fsx.Virtual
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
	lucene    lucenex.Driver
}

func (t *HTTPCommunityPublisher) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout10s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout10s(),
	).ThenFunc(t.find))

	r.Path("/{id}").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.update))

	// touches the filesystem and compiles the cloned module, hence the longer
	// budget - the same reasoning the environment service's GET uses.
	r.Path("/{id}/clone").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout10s(),
	).ThenFunc(t.clone))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

func (t *HTTPCommunityPublisher) search(w http.ResponseWriter, r *http.Request) {
	var resp = PluginPublisherSearchResponse{
		Next: &PluginPublisherSearchRequest{Limit: 100},
	}

	if err := t.decoder.Decode(resp.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode search request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	resp.Next.Limit = numericx.Min(resp.Next.Limit, 100)

	q := community.PluginPublisherSearch(r.Context(), t.q, community.PluginPublisherSearchBuilder().
		Where(
			squirrel.And{
				squirrel.Expr("1=1"),
				lucenex.Query(t.lucene, resp.Next.Query, lucenex.WithDefaultField("description")),
				squirrelx.NotIn("plugin_publishers.id", resp.Next.Excluded...),
			},
		).OrderBy("created_at DESC").Offset(resp.Next.Offset*resp.Next.Limit).Limit(resp.Next.Limit))
	qi := sqlx.Scan(q)
	for p := range qi.Iter() {
		resp.Items = append(resp.Items, NewPluginPublisher(PluginPublisherOptionFromDB(langx.Clone(p, timex.JSONSafeEncodeOption))))
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

	if err = publishplugin.VerifyWasmMagicPath(tmp.Name()); err != nil {
		log.Println(errorsx.Wrap(err, "invalid wasm module"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	id := md5x.FormatUUID(digest)
	dst := t.dir.Path(id + ".wasm")
	if err = os.Rename(tmp.Name(), dst); err != nil {
		log.Println(errorsx.Wrap(err, "unable to install publisher plugin"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = t.reg.Load(r.Context(), dst); err != nil {
		log.Println(errorsx.Wrap(err, "unable to load publisher plugin"))
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

func (t *HTTPCommunityPublisher) find(w http.ResponseWriter, r *http.Request) {
	var (
		pub community.PluginPublisher
	)

	id := mux.Vars(r)["id"]

	if err := community.PluginPublisherFindByID(r.Context(), t.q, id).Scan(&pub); errors.Is(err, sql.ErrNoRows) {
		log.Println(errorsx.Wrap(err, "unable to locate plugin publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to locate plugin publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &PluginPublisherFindResponse{
		Publisher: NewPluginPublisher(PluginPublisherOptionFromDB(langx.Clone(pub, timex.JSONSafeEncodeOption))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

// update records the fields an operator owns on an installed publisher. The
// request carries an entire publisher, but only description and mimetype are
// taken from it: id and path describe what is installed on disk, and a client
// cannot move or rename a module by editing a field.
func (t *HTTPCommunityPublisher) update(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg PluginPublisherUpdateRequest
		pub community.PluginPublisher
	)

	id := mux.Vars(r)["id"]

	if err = jsonx.UnmarshalRead(r.Body, &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if msg.Publisher == nil {
		log.Println("publisher is required")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = community.PluginPublisherFindByID(r.Context(), t.q, id).Scan(&pub); errors.Is(err, sql.ErrNoRows) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to locate plugin publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	pub.Description = msg.Publisher.Description
	pub.Mimetype = msg.Publisher.Mimetype

	if err = community.PluginPublisherUpdateByID(r.Context(), t.q, id, pub).Scan(&pub); err != nil {
		log.Println(errorsx.Wrap(err, "unable to update plugin publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &PluginPublisherUpdateResponse{
		Publisher: NewPluginPublisher(PluginPublisherOptionFromDB(langx.Clone(pub, timex.JSONSafeEncodeOption))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

// clone installs a second identity for an already installed module: a symlink
// in publish.d under a stem of its own. publishplugin.Identity folds that stem
// into the id, and every per-installation resource - the .env sidecar, the
// config and cache directories - is derived from the path, so the clone is
// configured entirely independently of what it points at.
//
// It starts as a copy of the original's configuration rather than an empty
// one: the reason to clone a publisher is a second account on the same
// service, which shares everything but the credentials.
func (t *HTTPCommunityPublisher) clone(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		source   community.PluginPublisher
		inserted community.PluginPublisher
	)

	if err = community.PluginPublisherFindByID(r.Context(), t.q, mux.Vars(r)["id"]).Scan(&source); errors.Is(err, sql.ErrNoRows) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to locate plugin publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	// point at the module itself, not at another link: a chain never forms,
	// and deleting the original - which removes only its own link - cannot
	// strand the clone.
	target, err := filepath.EvalSymlinks(source.Path)
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to resolve publisher plugin: %s", source.Path))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	dst := t.dir.Path(uuid.Must(uuid.NewV7()).String() + ".wasm")
	if err = os.Symlink(target, dst); err != nil {
		log.Println(errorsx.Wrapf(err, "unable to install publisher plugin clone: %s", dst))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}
	defer func() {
		if err == nil {
			return
		}

		log.Println("failure cloning publisher plugin, removing attempt", err)
		errorsx.Log(errorsx.Wrap(fsx.IgnoreIsNotExist(os.Remove(publishplugin.EnvPath(dst))), "unable to remove clone environment"))
		errorsx.Log(errorsx.Wrap(fsx.IgnoreIsNotExist(os.Remove(dst)), "unable to remove clone"))
	}()

	// a publisher that was never configured has no sidecar at all; the clone
	// simply starts without one too.
	if env, cause := os.ReadFile(publishplugin.EnvPath(source.Path)); fsx.IgnoreIsNotExist(cause) != nil {
		err = cause
		log.Println(errorsx.Wrap(err, "unable to read publisher environment"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	} else if cause == nil {
		if err = os.WriteFile(publishplugin.EnvPath(dst), env, 0o600); err != nil {
			log.Println(errorsx.Wrap(err, "unable to write publisher environment"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}
	}

	id, err := publishplugin.Identity(dst)
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to identify publisher plugin clone: %s", dst))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = t.reg.Load(r.Context(), dst); err != nil {
		log.Println(errorsx.Wrap(err, "unable to load publisher plugin clone"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	// deliberately unlabelled: the console falls back to the id, and the
	// details card is where a clone gets named.
	if err = community.PluginPublisherInsertWithDefaults(r.Context(), t.q, community.PluginPublisher{
		ID:       id,
		Path:     dst,
		Mimetype: source.Mimetype,
	}).Scan(&inserted); err != nil {
		log.Println(errorsx.Wrap(err, "unable to insert plugin publisher clone"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &PluginPublisherCloneResponse{
		Publisher: NewPluginPublisher(PluginPublisherOptionFromDB(langx.Clone(inserted, timex.JSONSafeEncodeOption))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPCommunityPublisher) delete(w http.ResponseWriter, r *http.Request) {
	var existing community.PluginPublisher

	id := mux.Vars(r)["id"]

	err := community.PluginPublisherDeleteByID(r.Context(), t.q, id).Scan(&existing)
	if errors.Is(err, sql.ErrNoRows) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to delete plugin publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	// a selection outlives the module it points at as a lookup that fails on
	// every subsequent publish, so uninstalling detaches it everywhere.
	errorsx.Log(errorsx.Wrap(
		sqlx.Discard(sqlx.Scan(community.CommunityPublisherDeleteByPublisherID(r.Context(), t.q, existing.ID))),
		"unable to detach publisher from communities",
	))

	t.reg.Unload(existing.Path)

	errorsx.Log(errorsx.Wrap(fsx.IgnoreIsNotExist(os.Remove(existing.Path)), "unable to remove publisher plugin from disk"))

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &PluginPublisherDeleteResponse{
		Publisher: NewPluginPublisher(PluginPublisherOptionFromDB(langx.Clone(existing, timex.JSONSafeEncodeOption))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
