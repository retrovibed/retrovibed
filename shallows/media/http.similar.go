package media

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/acoustics"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

const (
	similarityThreshold = 0.5
	similarLimit        = 20
)

type HTTPSimilarOption func(*HTTPSimilar)

func HTTPSimilarOptionJWTSecret(j jwtx.SecretSource) HTTPSimilarOption {
	return func(t *HTTPSimilar) {
		t.jwtsecret = j
	}
}

func NewHTTPSimilar(q sqlx.Queryer, options ...HTTPSimilarOption) *HTTPSimilar {
	return new(langx.Clone(HTTPSimilar{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...))
}

type HTTPSimilar struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPSimilar) Bind(r *mux.Router) {
	r.Path("/{media_id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.similar))
}

func (t *HTTPSimilar) similar(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		mediaID uuid.UUID
		md      library.Metadata
		result  MediaFindResponse
		req     MediaSearchRequest
	)

	log.Println("acoustic similarity initiated")
	defer log.Println("acoustic similarity completed")
	mediaID, err = uuid.FromString(mux.Vars(r)["media_id"])
	if err != nil {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = t.decoder.Decode(&req, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	seedVec, err := acoustics.FetchFeatures(r.Context(), t.q, mediaID)
	if sqlx.IgnoreNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "acoustics: fetch seed vector"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	} else if err != nil {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	}

	parsed := slicesx.MapTransform(uuid.FromStringOrNil, req.Excluded...)
	exclude := append([]uuid.UUID{mediaID}, parsed...)

	ids, err := acoustics.SimilarMediaIDs(r.Context(), t.q, seedVec, exclude, similarLimit, similarityThreshold)
	if err != nil {
		log.Println(errorsx.Wrap(err, "acoustics: similar"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	for _, id := range ids {
		if err = library.MetadataFindByID(r.Context(), t.q, id.String()).Scan(&md); sqlx.ErrNoRows(err) != nil {
			continue
		} else if err != nil {
			log.Println(errorsx.Wrap(err, "acoustics: metadata"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		result.Media = new(
			langx.Clone(
				Media{},
				MediaOptionFromLibraryMetadata(langx.Clone(md, timex.JSONSafeEncodeOption)),
				MediaOptionImageAuto(r),
			),
		)

		if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &result); err != nil {
			log.Println(errorsx.Wrap(err, "acoustics: write response"))
			return
		}

		return
	}

	errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
}
