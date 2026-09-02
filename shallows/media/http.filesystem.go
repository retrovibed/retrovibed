package media

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type HTTPFilesystemOption func(*HTTPFilesystem)

func HTTPFilesystemOptionJWTSecret(j jwtx.SecretSource) HTTPFilesystemOption {
	return func(t *HTTPFilesystem) {
		t.jwtsecret = j
	}
}

func NewHTTPFilesystem(q sqlx.Queryer, options ...HTTPFilesystemOption) *HTTPFilesystem {
	svc := langx.Clone(HTTPFilesystem{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
		fts:       duckdbx.NewLucene(),
	}, options...)

	return &svc
}

// the directory tree over library_metadata. it shares the Media message with the library
// because a directory and a file are the same kind of row, and shares nothing else: the
// library endpoint serves a flat grid that never wants to see a directory.
type HTTPFilesystem struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
	fts       lucenex.Driver
}

func (t *HTTPFilesystem) Bind(r *mux.Router) {
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
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.move))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

// lists one directory: its contents, and the path that reaches it.
func (t *HTTPFilesystem) search(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg = FilesystemSearchResponse{
			Next: &FilesystemSearchRequest{
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
	msg.Next.DirectoryId = stringsx.FirstNonBlank(msg.Next.DirectoryId, uuid.Nil.String())

	// directories float to the top of their own listing, which is the ordering every file
	// manager uses. the mimetype is bound rather than interpolated.
	q := library.MetadataSearchBuilder().Where(squirrel.And{
		library.MetadataQueryNotTombstoned(),
		library.MetadataQueryHidden(msg.Next.Hidden),
		library.MetadataQueryMimetypes(msg.Next.Mimetypes...),
		library.MetadataQueryDirectoryID(msg.Next.DirectoryId),
		lucenex.Query(t.fts, msg.Next.Query, lucenex.WithDefaultField("auto_description")),
	}).OrderByClause("mimetype = ? DESC, description ASC", mimex.Directory).
		Offset(msg.Next.Offset * msg.Next.Limit).
		Limit(msg.Next.Limit)

	if err = t.path(r, &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to resolve directory path"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	qi := sqlx.Scan(library.MetadataSearch(r.Context(), t.q, q))
	for v := range qi.Iter() {
		tmp := langx.Clone(Media{}, MediaOptionFromLibraryMetadata(langx.Clone(v, timex.JSONSafeEncodeOption)), MediaOptionImageAuto(r))
		msg.Items = append(msg.Items, &tmp)
	}

	if err = qi.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "encoding failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

// the listed directory and its ancestors. the root is not itself a row and has no path.
func (t *HTTPFilesystem) path(r *http.Request, msg *FilesystemSearchResponse) error {
	if uuid.FromStringOrNil(msg.Next.DirectoryId) == uuid.Nil {
		return nil
	}

	ancestors := sqlx.Scan(library.MetadataAncestorsByID(r.Context(), t.q, msg.Next.DirectoryId))
	for v := range ancestors.Iter() {
		tmp := langx.Clone(Media{}, MediaOptionFromLibraryMetadata(langx.Clone(v, timex.JSONSafeEncodeOption)))
		msg.Breadcrumb = append(msg.Breadcrumb, &tmp)
	}

	return ancestors.Err()
}

// a directory is a library_metadata row with a name, a directory, and nothing on disk.
func (t *HTTPFilesystem) create(w http.ResponseWriter, r *http.Request) {
	var (
		req FilesystemCreateRequest
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if stringsx.Blank(req.Name) {
		log.Println("name is required")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	// every other row is keyed by the md5 of its content and a directory has none, so the
	// id is generated; nothing in the schema or the scanners re-derives an id from bytes.
	// known_media_id stays at the zero uuid, which keeps directories out of the
	// identification daemon and, through it, off the discovery network.
	md := library.NewMetadata(
		uuid.Must(uuid.NewV7()).String(),
		library.MetadataOptionDescription(req.Name),
		library.MetadataOptionAutoDescription(library.NormalizedDescription(req.Name)),
		library.MetadataOptionMimetype(mimex.Directory),
		library.MetadataOptionDirectoryID(stringsx.FirstNonBlank(req.DirectoryId, uuid.Nil.String())),
	)

	if err := library.DirectoryUpsert(r.Context(), t.q, md).Scan(&md); err != nil {
		log.Println(errorsx.Wrap(err, "unable to record directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &FilesystemCreateResponse{
		Media: new(langx.Clone(Media{}, MediaOptionFromLibraryMetadata(langx.Clone(md, timex.JSONSafeEncodeOption)))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPFilesystem) move(w http.ResponseWriter, r *http.Request) {
	var (
		req FilesystemMoveRequest
		md  library.Metadata
		id  = mux.Vars(r)["id"]
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	// a destination inside the row's own subtree builds a directory_id cycle, and every
	// recursive descent here would then run until the process is killed. such a move
	// matches no row, which is indistinguishable from a missing id and equally a bad
	// request either way.
	if err := library.MetadataMoveByID(r.Context(), t.q, id, stringsx.FirstNonBlank(req.DirectoryId, uuid.Nil.String())).Scan(&md); sqlx.ErrNoRows(err) != nil {
		log.Println("move rejected", id, "into", req.DirectoryId)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to move"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &FilesystemMoveResponse{
		Media: new(langx.Clone(Media{}, MediaOptionFromLibraryMetadata(langx.Clone(md, timex.JSONSafeEncodeOption)))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

// deleting a directory deletes what it holds. directory_id carries no cascade, so
// tombstoning the directory alone would leave its contents pointing at a row
// NewTombstonedCleanup is about to hard delete, listed by nothing and still counted
// against disk usage.
func (t *HTTPFilesystem) delete(w http.ResponseWriter, r *http.Request) {
	var (
		md      library.Metadata
		removed uint64
		id      = mux.Vars(r)["id"]
	)

	tombstoned := sqlx.Scan(library.MetadataTombstoneSubtreeByID(r.Context(), t.q, id))
	for v := range tombstoned.Iter() {
		if v.ID == id {
			md = v
			continue
		}

		removed++
	}

	if err := tombstoned.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to tombstone directory"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if stringsx.Blank(md.ID) {
		log.Println("unable to tombstone directory", id, "no such record")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &FilesystemDeleteResponse{
		Media:   new(langx.Clone(Media{}, MediaOptionFromLibraryMetadata(langx.Clone(md, timex.JSONSafeEncodeOption)))),
		Removed: removed,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
