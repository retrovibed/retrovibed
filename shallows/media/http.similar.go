package media

import (
	"log"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/acoustics"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

const similarityThreshold = 0.5

type HTTPSimilarOption func(*HTTPSimilar)

func HTTPSimilarOptionJWTSecret(j jwtx.SecretSource) HTTPSimilarOption {
	return func(t *HTTPSimilar) {
		t.jwtsecret = j
	}
}

func NewHTTPSimilar(q sqlx.Queryer, idx *acoustics.Index, stats *acoustics.RunningStats, options ...HTTPSimilarOption) *HTTPSimilar {
	return langx.Autoptr(langx.Clone(HTTPSimilar{
		q:         q,
		idx:       idx,
		stats:     stats,
		jwtsecret: env.JWTSecret,
	}, options...))
}

type HTTPSimilar struct {
	q         sqlx.Queryer
	idx       *acoustics.Index
	stats     *acoustics.RunningStats
	jwtsecret jwtx.SecretSource
}

func (t *HTTPSimilar) Bind(r *mux.Router) {
	r.Path("/{media_id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.similar))
}

type SimilarResponse struct {
	Items []*Media `json:"items"`
}

func (t *HTTPSimilar) similar(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		mediaID uuid.UUID
		md      library.Metadata
	)

	mediaID, err = uuid.FromString(mux.Vars(r)["media_id"])
	if err != nil {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	count, err := acoustics.IndexedCount(r.Context(), t.q, acoustics.StatsVersion)
	if err != nil {
		log.Println(errorsx.Wrap(err, "acoustics: count indexed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}
	if count < acoustics.ColdStartThreshold {
		w.Header().Set("Retry-After", "60")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusServiceUnavailable))
		return
	}

	seedVecs, _, err := acoustics.FetchCandidateVectors(r.Context(), t.q, []uuid.UUID{mediaID})
	if err != nil {
		log.Println(errorsx.Wrap(err, "acoustics: fetch seed vector"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}
	if len(seedVecs) == 0 {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	}
	seedVec := t.stats.Normalize(seedVecs[0])

	exclude := make(map[uuid.UUID]struct{})
	exclude[mediaID] = struct{}{}
	if raw := r.FormValue("exclude"); raw != "" {
		for _, s := range strings.Split(raw, ",") {
			uid, parseErr := uuid.FromString(strings.TrimSpace(s))
			if parseErr == nil {
				exclude[uid] = struct{}{}
			}
		}
	}

	candidateIDs := t.idx.Candidates(seedVec, exclude)
	candidateVecs, candidateMediaIDs, err := acoustics.FetchCandidateVectors(r.Context(), t.q, candidateIDs)
	if err != nil {
		log.Println(errorsx.Wrap(err, "acoustics: fetch candidate vectors"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	normalized := make([]acoustics.FeatureVector, len(candidateVecs))
	for i, v := range candidateVecs {
		normalized[i] = t.stats.Normalize(v)
	}

	results := acoustics.RankCandidates(seedVec, normalized, candidateMediaIDs, 20, similarityThreshold)

	resp := SimilarResponse{Items: make([]*Media, 0, len(results))}
	for _, res := range results {
		err = library.MetadataFindByID(r.Context(), t.q, res.MediaID.String()).Scan(&md)
		if err != nil {
			continue
		}

		resp.Items = append(resp.Items, langx.Autoptr(langx.Clone(
			Media{},
			MediaOptionFromLibraryMetadata(langx.Clone(md, timex.JSONSafeEncodeOption)),
			MediaOptionImageAuto(r),
		)))
	}

	err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp)
	if err != nil {
		log.Println(errorsx.Wrap(err, "acoustics: write response"))
	}
}
