package media

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type HTTPAutoimportOption func(*HTTPAutoimport)

func HTTPAutoimportOptionJWTSecret(j jwtx.SecretSource) HTTPAutoimportOption {
	return func(t *HTTPAutoimport) {
		t.jwtsecret = j
	}
}

func NewHTTPAutoimport(q sqlx.Queryer, autoimport *asyncx.Wakeup, options ...HTTPAutoimportOption) *HTTPAutoimport {
	svc := langx.Clone(HTTPAutoimport{
		q:          q,
		autoimport: autoimport,
		jwtsecret:  env.JWTSecret,
		decoder:    formx.NewDecoder(),
	}, options...)

	return &svc
}

// HTTPAutoimport manages the directories monitored for automatic import into the library.
// creating and editing a directory makes the daemon read, and in move mode remove, files
// at an arbitrary local path so those routes require library modification permissions.
type HTTPAutoimport struct {
	q          sqlx.Queryer
	autoimport *asyncx.Wakeup
	jwtsecret  jwtx.SecretSource
	decoder    *form.Decoder
}

func (t *HTTPAutoimport) Bind(r *mux.Router) {
	r.StrictSlash(false)
	r.Use(httpx.RouteInvoked)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermLibraryModify),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.find))

	r.Path("/{id}").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermLibraryModify),
		httpx.Timeout2s(),
	).ThenFunc(t.update))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		metaapi.AuthzTokenHTTP(t.jwtsecret, metaapi.AuthzPermLibraryModify),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

func (t *HTTPAutoimport) search(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg = AutoimportDirectorySearchResponse{
			Next: &AutoimportDirectorySearchRequest{
				Limit: 100,
			},
		}
	)

	if err = t.decoder.Decode(msg.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	msg.Next.Limit = numericx.Min(msg.Next.Limit, 100)

	b := library.AutoimportDirectorySearchBuilder().Where(squirrel.And{
		squirrel.Expr("1=1"),
		library.AutoimportDirectoryQueryByIDs(msg.Next.Id...),
		library.AutoimportDirectoryQueryText(msg.Next.Query),
	}).
		OrderBy("id ASC").
		Offset(msg.Next.Offset * msg.Next.Limit).
		Limit(msg.Next.Limit)

	q := sqlx.Scan(library.AutoimportDirectorySearch(r.Context(), t.q, b))

	for v := range q.Iter() {
		tmp := langx.Clone(AutoimportDirectory{}, AutoimportDirectoryOptionFromLibrary(langx.Clone(v, timex.JSONSafeEncodeOption)))
		msg.Items = append(msg.Items, &tmp)
	}

	if err = q.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "encoding failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPAutoimport) find(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		d        library.AutoimportDirectory
		pending  uint64
		imported uint64
		id       = mux.Vars(r)["id"]
	)

	if err = library.AutoimportDirectoryFindByID(r.Context(), t.q, id).Scan(&d); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find autoimport directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find autoimport directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if pending, err = library.AutoimportFileCount(r.Context(), t.q, squirrel.And{library.AutoimportFileQueryByDirectoryID(d.ID), library.AutoimportFileQueryPending()}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to count pending files"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if imported, err = library.AutoimportFileCount(r.Context(), t.q, squirrel.And{library.AutoimportFileQueryByDirectoryID(d.ID), library.AutoimportFileQueryImported()}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to count imported files"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &AutoimportDirectoryLookupResponse{
		Directory: new(
			langx.Clone(
				AutoimportDirectory{},
				AutoimportDirectoryOptionFromLibrary(langx.Clone(d, timex.JSONSafeEncodeOption)),
				AutoimportDirectoryOptionCounts(pending, imported),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPAutoimport) create(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg AutoimportDirectoryCreateRequest
	)

	if err = jsonx.UnmarshalRead(r.Body, &msg); err != nil || msg.Directory == nil {
		log.Println(errorsx.Wrap(errorsx.Compact(err, errorsx.String("directory is required")), "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	d := NewLibraryAutoimportDirectory(msg.Directory)
	d.Path = filepath.Clean(d.Path)
	d.LibraryDirectoryID = stringsx.FirstNonBlank(d.LibraryDirectoryID, uuid.Nil.String())

	if !filepath.IsAbs(d.Path) {
		log.Println("autoimport directory path must be absolute", d.Path)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if info, cause := os.Stat(d.Path); cause != nil || !info.IsDir() {
		log.Println(errorsx.Wrapf(errorsx.Compact(cause, errorsx.String("not a directory")), "invalid autoimport directory: %s", d.Path))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	// mode and debounce are enforced by the table's check constraints.
	if err = library.AutoimportDirectoryInsertWithDefaults(r.Context(), t.q, d).Scan(&d); duckdbx.ErrUniqueConstraintViolation(err) != nil {
		log.Println(errorsx.Wrap(err, "invalid autoimport directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to create autoimport directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	t.autoimport.Broadcast()

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &AutoimportDirectoryCreateResponse{
		Directory: new(
			langx.Clone(
				AutoimportDirectory{},
				AutoimportDirectoryOptionFromLibrary(langx.Clone(d, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

// update replaces the editable fields (description, debounce, mode, library directory), the
// path is immutable; remove and recreate the directory to monitor a different path.
func (t *HTTPAutoimport) update(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg AutoimportDirectoryUpdateRequest
		id  = mux.Vars(r)["id"]
	)

	if err = jsonx.UnmarshalRead(r.Body, &msg); err != nil || msg.Directory == nil {
		log.Println(errorsx.Wrap(errorsx.Compact(err, errorsx.String("directory is required")), "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	d := NewLibraryAutoimportDirectory(msg.Directory)
	d.LibraryDirectoryID = stringsx.FirstNonBlank(d.LibraryDirectoryID, uuid.Nil.String())

	// mode and debounce are enforced by the table's check constraints.
	if err = library.AutoimportDirectoryUpdateByID(r.Context(), t.q, id, d).Scan(&d); duckdbx.ErrUniqueConstraintViolation(err) != nil {
		log.Println(errorsx.Wrap(err, "invalid autoimport directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	} else if sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find autoimport directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to update autoimport directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	t.autoimport.Broadcast()

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &AutoimportDirectoryUpdateResponse{
		Directory: new(
			langx.Clone(
				AutoimportDirectory{},
				AutoimportDirectoryOptionFromLibrary(langx.Clone(d, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPAutoimport) delete(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		d   library.AutoimportDirectory
		id  = mux.Vars(r)["id"]
	)

	if err = library.AutoimportDirectoryDeleteByID(r.Context(), t.q, id).Scan(&d); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find autoimport directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to delete autoimport directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = sqlx.Discard(sqlx.Scan(library.AutoimportFileDeleteByDirectory(r.Context(), t.q, d.ID))); err != nil {
		log.Println(errorsx.Wrap(err, "unable to delete autoimport files"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &AutoimportDirectoryDeleteResponse{
		Directory: new(
			langx.Clone(
				AutoimportDirectory{},
				AutoimportDirectoryOptionFromLibrary(langx.Clone(d, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
