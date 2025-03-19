package media

import (
	"crypto/md5"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/x/bytesx"
	"github.com/retrovibed/retrovibed/internal/x/errorsx"
	"github.com/retrovibed/retrovibed/internal/x/formx"
	"github.com/retrovibed/retrovibed/internal/x/fsx"
	"github.com/retrovibed/retrovibed/internal/x/httpx"
	"github.com/retrovibed/retrovibed/internal/x/iox"
	"github.com/retrovibed/retrovibed/internal/x/jwtx"
	"github.com/retrovibed/retrovibed/internal/x/langx"
	"github.com/retrovibed/retrovibed/internal/x/md5x"
	"github.com/retrovibed/retrovibed/internal/x/numericx"
	"github.com/retrovibed/retrovibed/internal/x/sqlx"
	"github.com/retrovibed/retrovibed/internal/x/sqlxx"
	"github.com/retrovibed/retrovibed/library"
)

type HTTPLibraryOption func(*HTTPLibrary)

func HTTPLibraryOptionJWTSecret(j jwtx.JWTSecretSource) HTTPLibraryOption {
	return func(t *HTTPLibrary) {
		t.jwtsecret = j
	}
}

func NewHTTPLibrary(q sqlx.Queryer, s fsx.Virtual, options ...HTTPLibraryOption) *HTTPLibrary {
	svc := langx.Clone(HTTPLibrary{
		q:            q,
		jwtsecret:    env.JWTSecret,
		decoder:      formx.NewDecoder(),
		mediastorage: s,
	}, options...)

	return &svc
}

type HTTPLibrary struct {
	q            sqlx.Queryer
	jwtsecret    jwtx.JWTSecretSource
	decoder      *form.Decoder
	mediastorage fsx.Virtual
}

func (t *HTTPLibrary) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.TimeoutRollingRead(3*time.Second),
	).ThenFunc(t.upload))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.DebugRequest,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.TimeoutRollingWrite(3*time.Second),
	).Then(http.FileServerFS(fsx.VirtualAsFSWithRewrite(t.mediastorage, func(s string) string {
		return strings.TrimPrefix(s, "m/")
	}))))
}

func (t *HTTPLibrary) delete(w http.ResponseWriter, r *http.Request) {
	var (
		md library.Metadata
		id = mux.Vars(r)["id"]
	)

	if err := library.MetadataTombstoneByID(r.Context(), t.q, id).Scan(&md); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to tombstone metadata"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to tombstone metadata"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &MediaDeleteResponse{
		Media: langx.Autoptr(
			langx.Clone(
				Media{},
				MediaOptionFromLibraryMetadata(langx.Clone(md, library.MetadataOptionJSONSafeEncode))),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
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
	)

	if f, fh, err = r.FormFile("content"); err != nil {
		log.Println(errorsx.Wrap(err, "content parameter required"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	defer f.Close()

	tmp, err := fsx.CreateTemp(t.mediastorage, "retrovibed.upload.*")
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to create temporary file"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}
	defer func() {
		if err == nil {
			return
		}

		log.Println("failure receiving upload, removing attempt", err)
		errorsx.Log(errorsx.Wrap(os.Remove(tmp.Name()), "unable to remove tmp"))
	}()
	defer tmp.Close()

	if _, err = io.CopyBuffer(io.MultiWriter(tmp, mhash, copied), f, buf[:]); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create temporary file"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	lmd := library.Metadata{
		ID:          md5x.FormatString(mhash),
		Description: fh.Filename,
		Bytes:       *copied.Result,
		Mimetype:    fh.Header.Get("Content-Type"),
	}

	if err = library.MetadataInsertWithDefaults(r.Context(), t.q, lmd).Scan(&lmd); err != nil {
		log.Println(errorsx.Wrap(err, "unable to record library metadata record"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = fsx.Rename(t.mediastorage, tmp.Name(), lmd.ID); err != nil {
		log.Println(errorsx.Wrap(err, "unable to record library metadata record"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &MediaUploadResponse{
		Media: langx.Autoptr(
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

func (t *HTTPLibrary) search(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg MediaSearchResponse = MediaSearchResponse{
			Next: &MediaSearchRequest{
				Limit: 100,
			},
		}
	)

	if err = t.decoder.Decode(msg.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	msg.Next.Limit = numericx.Min(msg.Next.Limit, 100)

	q := library.MetadataSearchBuilder().Where(squirrel.And{
		library.MetadataQueryVisible(),
		library.MetadataQuerySearch(msg.Next.Query, "description"),
	}).OrderBy("created_at DESC").Offset(msg.Next.Offset * msg.Next.Limit).Limit(msg.Next.Limit)

	err = sqlxx.ScanEach(library.MetadataSearch(r.Context(), t.q, q), func(p *library.Metadata) error {
		tmp := langx.Clone(Media{}, MediaOptionFromLibraryMetadata(langx.Clone(*p, library.MetadataOptionJSONSafeEncode)))
		msg.Items = append(msg.Items, &tmp)
		return nil
	})

	if err != nil {
		log.Println(errorsx.Wrap(err, "encoding failed"))
		errorsx.MaybeLog(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
