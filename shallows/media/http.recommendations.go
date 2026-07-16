package media

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type HTTPRecommendationsOption func(*HTTPRecommendations)

func HTTPRecommendationsOptionJWTSecret(j jwtx.SecretSource) HTTPRecommendationsOption {
	return func(t *HTTPRecommendations) {
		t.jwtsecret = j
	}
}

func NewHTTPRecommendations(q sqlx.Queryer, options ...HTTPRecommendationsOption) *HTTPRecommendations {
	return new(langx.Clone(HTTPRecommendations{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...))
}

type HTTPRecommendations struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPRecommendations) Bind(r *mux.Router) {
	r.StrictSlash(false)
	r.Use(httpx.RouteInvoked)
	// r.Use(httpx.DebugRequest)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.latest))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.refresh))

	r.Path("/random").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.random))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.find))

	r.Path("/content/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.content))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

func (t *HTTPRecommendations) latest(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg = RecommendationSearchResponse{
			Next: &RecommendationSearchRequest{
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

	query := library.RecommendationSearchBuilder().
		Where(library.RecommendationQueryNotTombstoned()).
		Where(library.RecommendationQueryMimetype(msg.Next.Mimetype)).
		OrderBy("library_recommendations.updated_at DESC").
		Offset(msg.Next.Offset * msg.Next.Limit).
		Limit(msg.Next.Limit)

	q := sqlx.Scan(library.RecommendationSearch(r.Context(), t.q, query))
	for rec := range q.Iter() {
		tmp := langx.Clone(Known{}, KnownOptionFromRecommendation(langx.Clone(rec, timex.JSONSafeEncodeOption)))
		msg.Items = append(msg.Items, &tmp)
	}

	if err = q.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "recommendations query failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRecommendations) find(w http.ResponseWriter, r *http.Request) {
	var (
		rec library.Recommendation
		id  = mux.Vars(r)["id"]
	)

	if err := library.RecommendationFindByID(r.Context(), t.q, id).Scan(&rec); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find recommendation"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find recommendation"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &RecommendationFindResponse{
		Recommendation: new(langx.Clone(Known{}, KnownOptionFromRecommendation(langx.Clone(rec, timex.JSONSafeEncodeOption)))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRecommendations) content(w http.ResponseWriter, r *http.Request) {
	var (
		rec library.Recommendation
		id  = mux.Vars(r)["id"]
	)

	if err := library.RecommendationFindByContentID(r.Context(), t.q, id).Scan(&rec); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find recommendation"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find recommendation"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &RecommendationFindResponse{
		Recommendation: new(langx.Clone(Known{}, KnownOptionFromRecommendation(langx.Clone(rec, timex.JSONSafeEncodeOption)))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRecommendations) delete(w http.ResponseWriter, r *http.Request) {
	var (
		rec library.Recommendation
		id  = mux.Vars(r)["id"]
	)

	if err := library.RecommendationDeleteByID(r.Context(), t.q, id).Scan(&rec); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to delete recommendation"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to delete recommendation"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &RecommendationDeleteResponse{
		Recommendation: new(langx.Clone(Known{}, KnownOptionFromRecommendation(langx.Clone(rec, timex.JSONSafeEncodeOption)))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRecommendations) random(w http.ResponseWriter, r *http.Request) {
	var req RecommendationSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	rec, err := library.RecommendationFromRandomKnown(r.Context(), t.q, req.Mimetype, req.Language, req.Adult)
	if sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "no recommendation available"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to select random known media"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	tmp := langx.Clone(Known{}, KnownOptionFromRecommendation(langx.Clone(rec, timex.JSONSafeEncodeOption)))
	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &RecommendationSearchResponse{
		Next:  &RecommendationSearchRequest{},
		Items: []*Known{&tmp},
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRecommendations) refresh(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req RecommendationSearchRequest
		rec library.Recommendation
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if rec, err = library.RecommendationFromRandomKnown(r.Context(), t.q, req.Mimetype, req.Language, req.Adult); sqlx.ErrNoRows(err) != nil {
		// no known media available - nothing to do
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "refresh failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &RecommendationRefreshResponse{
		Recommendation: new(langx.Clone(Known{}, KnownOptionFromRecommendation(langx.Clone(rec, timex.JSONSafeEncodeOption)))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
