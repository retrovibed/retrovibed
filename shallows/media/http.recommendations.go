package media

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/shallows/httpauth"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jwtx"
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
	return langx.Autoptr(langx.Clone(HTTPRecommendations{
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

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

func (t *HTTPRecommendations) latest(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg = RecommendationsResponse{
			Next: &RecommendationsRequest{
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

	query := library.RecommendationKnownSearchBuilder().Where(
		library.RecommendationQueryNotTombstoned(),
	).OrderBy("library_recommendations.updated_at DESC").
		Offset(msg.Next.Offset * msg.Next.Limit).
		Limit(msg.Next.Limit)

	q := sqlx.Scan2(library.RecommendationKnownSearch(r.Context(), t.q, query))
	for _, k := range q.Iter() {
		tmp := langx.Clone(Known{}, KnownOptionFromLibraryKnown(langx.Clone(k, timex.JSONSafeEncodeOption)))
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

	if err := httpx.WriteEmptyJSON(w, http.StatusOK); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRecommendations) random(w http.ResponseWriter, r *http.Request) {
	rec, err := library.RecommendationFromRandomKnown(r.Context(), t.q)
	if sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "no known media available"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to select random known media"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	var known library.Known
	if err = library.KnownFindByID(r.Context(), t.q, rec.KnownMediaID).Scan(&known); err != nil {
		log.Println(errorsx.Wrap(err, "unable to find known media"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	tmp := langx.Clone(Known{}, KnownOptionFromLibraryKnown(langx.Clone(known, timex.JSONSafeEncodeOption)))
	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &RecommendationsResponse{
		Next:  &RecommendationsRequest{},
		Items: []*Known{&tmp},
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRecommendations) refresh(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg RecommendationsRequest
	)

	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if _, err = library.RecommendationFromRandomKnown(r.Context(), t.q); sqlx.ErrNoRows(err) != nil {
		// no known media available - nothing to do
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "refresh failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &RecommendationsResponse{}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
