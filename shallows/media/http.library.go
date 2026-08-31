package media

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/iox"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type HTTPLibraryOption func(*HTTPLibrary)

func HTTPLibraryOptionJWTSecret(j jwtx.SecretSource) HTTPLibraryOption {
	return func(t *HTTPLibrary) {
		t.jwtsecret = j
	}
}

func HTTPLibraryOptionTorrentStorage(vfs fsx.Virtual) HTTPLibraryOption {
	return func(t *HTTPLibrary) {
		t.torrentstorage = vfs
	}
}

func NewHTTPLibrary(q sqlx.Queryer, archival *asyncx.Wakeup, identification *asyncx.Wakeup, media fsx.Virtual, deeppool *http.Client, options ...HTTPLibraryOption) *HTTPLibrary {
	svc := langx.Clone(HTTPLibrary{
		q:              q,
		jwtsecret:      env.JWTSecret,
		decoder:        formx.NewDecoder(),
		mediastorage:   media,
		fts:            duckdbx.NewLucene(),
		deeppool:       deeppool,
		archival:       archival,
		identification: identification,
	}, options...)

	return &svc
}

type HTTPLibrary struct {
	q              sqlx.Queryer
	jwtsecret      jwtx.SecretSource
	decoder        *form.Decoder
	deeppool       *http.Client
	mediastorage   fsx.Virtual
	torrentstorage fsx.Virtual
	fts            lucenex.Driver
	archival       *asyncx.Wakeup
	identification *asyncx.Wakeup
}

func (t *HTTPLibrary) Bind(r *mux.Router) {
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
		httpx.TimeoutRollingRead(3*time.Second),
	).ThenFunc(t.upload))

	r.Path("/random").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.random))

	r.Path("/{id}/metadatasync").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.patch))

	r.Path("/{id}").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.patch))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout10s(),
	).Then(http.FileServerFS(library.New(t.deeppool, t.mediastorage, func(ctx context.Context, s string) (*library.Metadata, error) {
		var (
			md library.Metadata
		)
		return &md, library.MetadataFindByID(ctx, t.q, strings.TrimPrefix(s, "m/")).Scan(&md)
	}))))

}

func (t *HTTPLibrary) delete(w http.ResponseWriter, r *http.Request) {
	var (
		md library.Metadata
		id = mux.Vars(r)["id"]
	)

	// deleting a folder deletes what it holds. the response describes the row the client
	// named; its descendants went with it.
	tombstoned := sqlx.Scan(library.MetadataTombstoneSubtreeByID(r.Context(), t.q, id))
	for v := range tombstoned.Iter() {
		if v.ID == id {
			md = v
		}
	}

	if err := tombstoned.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to tombstone metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if stringsx.Blank(md.ID) {
		log.Println("unable to tombstone metadata", id, "no such record")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &MediaDeleteResponse{
		Media: new(
			langx.Clone(
				Media{},
				MediaOptionFromLibraryMetadata(langx.Clone(md, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPLibrary) patch(w http.ResponseWriter, r *http.Request) {
	var (
		req MediaUpdateRequest
		md  library.Metadata
		id  = mux.Vars(r)["id"]
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decoded update"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	md = langx.Clone(
		md,
		library.MetadataOptionDescription(req.Media.Description),
		library.MetadataOptionKnownMediaID(req.Media.KnownMediaId),
		library.MetadataOptionArchiveID(req.Media.ArchiveId),
	)

	if err := library.MetadataUpdate(r.Context(), t.q, id, md).Scan(&md); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "update ignored"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to tombstone metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	// the update above has already accounted for a missing id, so the only row the move
	// can fail to match is one whose requested parent sits inside its own subtree.
	if stringsx.Present(req.Media.ParentId) {
		if err := library.MetadataMoveByID(r.Context(), t.q, id, req.Media.ParentId).Scan(&md); sqlx.ErrNoRows(err) != nil {
			log.Println("move rejected", id, "into", req.Media.ParentId, "would cycle")
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		} else if err != nil {
			log.Println(errorsx.Wrap(err, "unable to move metadata"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &MediaUpdateResponse{
		Media: new(
			langx.Clone(
				Media{},
				MediaOptionFromLibraryMetadata(langx.Clone(md, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}

	if uuid.Max == uuid.FromStringOrNil(md.ArchiveID) {
		// signal that archival needs to be done.
		t.archival.Broadcast()
	}
}

func (t *HTTPLibrary) upload(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		f      multipart.File
		fh     *multipart.FileHeader
		buf    [bytesx.MiB]byte
		copied = &iox.Copied{Result: new(uint64)}
		mhash  = md5.New()
		parent = stringsx.FirstNonBlank(r.FormValue("parent_id"), uuid.Nil.String())
	)

	// a folder shares this route because it shares the destination and the response shape,
	// and skips everything downstream that exists to receive bytes.
	if r.FormValue("mimetype") == mimex.Directory {
		t.mkdir(w, r, parent)
		return
	}

	if f, fh, err = r.FormFile("content"); err != nil {
		log.Println(errorsx.Wrap(err, "content parameter required"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	defer f.Close()

	tmp, err := os.MkdirTemp(t.mediastorage.Path(), "retrovibed.upload.*")
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
		errorsx.Log(errorsx.Wrap(os.RemoveAll(tmp), "unable to remove tmp"))
	}()

	bcache, err := blockcache.NewDirectoryCache(tmp)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to create block cache"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if _, err = io.CopyBuffer(io.MultiWriter(io.NewOffsetWriter(bcache, 0), mhash, copied), f, buf[:]); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create temporary file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	lmd := library.NewMetadata(
		md5x.FormatUUID(mhash),
		library.MetadataOptionBytes(*copied.Result),
		library.MetadataOptionDescription(fh.Filename),
		library.MetadataOptionMimetype(fh.Header.Get("Content-Type")),
		library.MetadataOptionKnownMediaID(uuid.Max.String()),
		library.MetadataOptionAutoDescription(library.NormalizedDescription(fh.Filename)),
		library.MetadataOptionParentID(parent),
	)

	if err = library.MetadataInsertWithDefaults(r.Context(), t.q, lmd).Scan(&lmd); err != nil {
		log.Println(errorsx.Wrap(err, "unable to record library metadata record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	// id is the md5 of the content, so re-uploading a file that already exists conflicts
	// onto the filed row, and the conflict clause deliberately leaves an existing parent
	// alone so a torrent completion or a filesystem rescan cannot flatten the library. a
	// client that named a destination meant it, which makes the move this call site's job.
	if stringsx.Present(r.FormValue("parent_id")) {
		if err = library.MetadataMoveByID(r.Context(), t.q, lmd.ID, parent).Scan(&lmd); err != nil {
			log.Println(errorsx.Wrap(err, "unable to file uploaded media"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}
	}

	if err = os.RemoveAll(t.mediastorage.Path(lmd.ID)); err != nil && !os.IsNotExist(err) {
		log.Println(errorsx.Wrap(err, "unable to remove existing media"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = fsx.Rename(t.mediastorage, tmp, lmd.ID); err != nil {
		log.Println(errorsx.Wrap(err, "unable to record library metadata record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &MediaUploadResponse{
		Media: new(
			langx.Clone(
				Media{},
				MediaOptionFromLibraryMetadata(lmd),
				MediaOptionImageAuto(r),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}

	t.identification.Broadcast()
}

// a folder is a library_metadata row with a name, a parent, and nothing on disk.
func (t *HTTPLibrary) mkdir(w http.ResponseWriter, r *http.Request, parent string) {
	var (
		description = r.FormValue("description")
	)

	if stringsx.Blank(description) {
		log.Println("description parameter required")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	// every other row is keyed by the md5 of its content and a folder has none, so the id
	// is generated; nothing in the schema or the scanners re-derives an id from bytes.
	// known_media_id stays at the zero uuid, which is what keeps folders out of the
	// identification daemon and, through it, off the discovery network.
	lmd := library.NewMetadata(
		uuid.Must(uuid.NewV7()).String(),
		library.MetadataOptionDescription(description),
		library.MetadataOptionAutoDescription(library.NormalizedDescription(description)),
		library.MetadataOptionMimetype(mimex.Directory),
		library.MetadataOptionParentID(parent),
	)

	if err := library.MetadataInsertWithDefaults(r.Context(), t.q, lmd).Scan(&lmd); err != nil {
		log.Println(errorsx.Wrap(err, "unable to record library metadata record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &MediaUploadResponse{
		Media: new(
			langx.Clone(
				Media{},
				MediaOptionFromLibraryMetadata(lmd),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPLibrary) random(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req = MediaSearchRequest{
			Limit: 1,
		}
		result MediaFindResponse
	)

	if err = t.decoder.Decode(&req, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	req.Limit = numericx.Min(req.Limit, 1)

	perform := func(id string) (library.Metadata, error) {
		q := library.MetadataSearchBuilder().Where(squirrel.And{
			library.MetadataQueryHidden(req.Hidden),
			library.MetadataQueryNotTombstoned(),
			library.MetadataQueryDirectory(false),
			library.MetadataQueryMimetypes(req.Mimetypes...),
		}).Limit(req.Limit)

		if tmp, err := sqlx.ScanOne(library.MetadataSearch(r.Context(), t.q, q.Where(library.MetadataQueryRandomAfter(id)))); err == nil {
			return tmp, nil
		}

		return sqlx.ScanOne(library.MetadataSearch(r.Context(), t.q, q.Where(library.MetadataQueryRandomBefore(id))))
	}

	tmp, err := perform(uuid.Must(uuid.NewV4()).String())
	if errorsx.Is(err, sql.ErrNoRows) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "retrieval failure"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	result.Media = new(
		langx.Clone(
			Media{},
			MediaOptionFromLibraryMetadata(langx.Clone(tmp, timex.JSONSafeEncodeOption)),
			MediaOptionImageAuto(r),
		),
	)

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &result); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPLibrary) search(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg = MediaSearchResponse{
			Next: &MediaSearchRequest{
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

	ordering := "created_at DESC, description ASC"
	if stringsx.Present(msg.Next.Query) {
		ordering = "description ASC"
	}

	// this handler serves both the flat media grid and the filesystem view. the
	// presence of a parent is what separates them. absent, a folder row is an
	// organizational artifact rather than media and has no place in the grid.
	scope, orderargs := library.MetadataQueryDirectory(false), []any{}
	if stringsx.Present(msg.Next.ParentId) {
		scope = library.MetadataQueryParent(msg.Next.ParentId)
		ordering, orderargs = "mimetype = ? DESC, description ASC", []any{mimex.Directory}
	}

	q := library.MetadataSearchBuilder().Where(squirrel.And{
		library.MetadataQueryNotTombstoned(),
		library.MetadataQueryHidden(msg.Next.Hidden),
		library.MetadataQueryMimetypes(msg.Next.Mimetypes...),
		scope,
		lucenex.Query(t.fts, msg.Next.Query, lucenex.WithDefaultField("auto_description")),
	}).OrderByClause(ordering, orderargs...).Offset(msg.Next.Offset * msg.Next.Limit).Limit(msg.Next.Limit)

	if err = t.breadcrumb(r, &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to resolve breadcrumb"))
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

// fills in the listed folder and its ancestors so the client renders a path without a
// second round trip. the root sentinel is not itself a row and has no ancestors.
func (t *HTTPLibrary) breadcrumb(r *http.Request, msg *MediaSearchResponse) error {
	if uuid.FromStringOrNil(msg.Next.ParentId) == uuid.Nil {
		return nil
	}

	ancestors := sqlx.Scan(library.MetadataAncestorsByID(r.Context(), t.q, msg.Next.ParentId))
	for v := range ancestors.Iter() {
		tmp := langx.Clone(Media{}, MediaOptionFromLibraryMetadata(langx.Clone(v, timex.JSONSafeEncodeOption)))
		msg.Breadcrumb = append(msg.Breadcrumb, &tmp)
	}

	return ancestors.Err()
}
