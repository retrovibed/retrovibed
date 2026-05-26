package media

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type HTTPKnownOption func(*HTTPKnown)

func HTTPKnownOptionJWTSecret(j jwtx.SecretSource) HTTPKnownOption {
	return func(t *HTTPKnown) {
		t.jwtsecret = j
	}
}

func NewHTTPKnown(q sqlx.Queryer, options ...HTTPKnownOption) *HTTPKnown {
	svc := langx.Clone(HTTPKnown{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
		fts:       duckdbx.NewLucene(),
	}, options...)

	return &svc
}

type HTTPKnown struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
	fts       lucenex.Driver
}

func (t *HTTPKnown) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/latest").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.latest))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.find))

}

func (t *HTTPKnown) search(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg = KnownSearchResponse{
			Next: &KnownSearchRequest{
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

	q := sqlx.Scan(library.KnownSearch(r.Context(), t.q, library.KnownSearchBuilder().Where(squirrel.And{
		library.KnownQueryMimetype(msg.Next.Mimetype),
		library.KnownQueryLanguage(msg.Next.Language),
		library.KnownQueryDetectLanguage(msg.Next.Language),
		library.KnownQueryExplicit(msg.Next.Adult),
		library.KnownQueryWithPoster(),
		lucenex.Query(t.fts, msg.Next.Query, lucenex.WithDefaultField("auto_description")),
	}).OrderBy("released DESC").Offset(msg.Next.Offset*msg.Next.Limit).Limit(msg.Next.Limit)))

	for v := range q.Iter() {
		tmp := langx.Clone(Known{}, KnownOptionFromLibraryKnown(langx.Clone(v, timex.JSONSafeEncodeOption)))
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

func (t *HTTPKnown) find(w http.ResponseWriter, r *http.Request) {
	var (
		meta library.Known
		id   = mux.Vars(r)["id"]
	)

	if err := library.KnownFindByID(r.Context(), t.q, id).Scan(&meta); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &KnownLookupResponse{
		Known: langx.Autoptr(
			langx.Clone(
				Known{},
				KnownOptionFromLibraryKnown(langx.Clone(meta, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPKnown) latest(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		released timex.Range
		msg      = KnownLatestResponse{
			Next: &KnownLatestRequest{
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

	if released, err = timex.NewRangeISO8601(msg.Next.Released.Oldest, msg.Next.Released.Newest); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode created"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	q := sqlx.Scan(library.KnownSearch(r.Context(), t.q, library.KnownSearchBuilder().Where(squirrel.And{
		library.KnownQueryExplicit(false),
		library.KnownQueryReleased(released),
		library.KnownQueryMimetype(msg.Next.Mimetype),
		library.KnownQueryWithPoster(),
	}).OrderBy("released DESC", "title DESC").Offset(msg.Next.Offset*msg.Next.Limit).Limit(msg.Next.Limit)))

	for v := range q.Iter() {
		tmp := langx.Clone(Known{}, KnownOptionFromLibraryKnown(langx.Clone(v, timex.JSONSafeEncodeOption)))
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

func (t *HTTPKnown) create(w http.ResponseWriter, r *http.Request) {
	var (
		req  KnownCreateRequest
		meta library.Known
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	known := req.GetKnown()
	if known == nil || known.Description == "" {
		log.Println("known.description is required")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	released, err := time.Parse(time.RFC3339Nano, known.Released)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to parse time for released media"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	contentmd5 := md5x.Digest(known.Description + known.Summary)
	uidmd5 := uuid.FromBytesOrNil(contentmd5.Sum(nil))

	if err := library.KnownFindByMd5(r.Context(), t.q, uidmd5.String()).Scan(&meta); err == nil {
		if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &KnownCreateResponse{
			Known: langx.Autoptr(
				langx.Clone(
					Known{},
					KnownOptionFromLibraryKnown(langx.Clone(meta, timex.JSONSafeEncodeOption)),
				),
			),
		}); err != nil {
			log.Println(errorsx.Wrap(err, "unable to write response"))
		}
		return
	}

	localid := errorsx.Must(uuid.NewV7())
	uid := library.KnownImportedUUID("local", localid)

	meta = library.Known{
		ID:              localid.String(),
		UID:             uid.String(),
		Md5:             uidmd5.String(),
		Md5Lower:        binary.LittleEndian.Uint64(uuidx.LowN(uidmd5, 64)),
		Title:           known.Description,
		Overview:        known.Summary,
		PosterPath:      known.Image,
		Released:        released,
		Adult:           known.Adult,
		Source:          "retrovibed",
		AutoDescription: known.Description,
	}

	if err := library.KnownInsertWithDefaults(r.Context(), t.q, meta).Scan(&meta); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create known media record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &KnownCreateResponse{
		Known: langx.Autoptr(
			langx.Clone(
				Known{},
				KnownOptionFromLibraryKnown(langx.Clone(meta, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
