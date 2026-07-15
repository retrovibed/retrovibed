package ddiscapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

type HTTPLocateOption func(*HTTPLocate)

func HTTPLocateOptionJWTSecret(j jwtx.SecretSource) HTTPLocateOption {
	return func(t *HTTPLocate) {
		t.jwtsecret = j
	}
}

func NewHTTPLocate(q sqlx.Queryer, pub *asyncx.Wakeup, options ...HTTPLocateOption) *HTTPLocate {
	svc := langx.Clone(HTTPLocate{
		q:         q,
		locate:    pub,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
		lucene:    duckdbx.NewLucene(),
	}, options...)

	return &svc
}

type HTTPLocate struct {
	q         sqlx.Queryer
	locate    *asyncx.Wakeup
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
	lucene    lucenex.Driver
}

func (t *HTTPLocate) Bind(r *mux.Router) {
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
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.find))
}

func (t *HTTPLocate) search(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg = LocateSearchResponse{
			Next: &LocateSearchRequest{
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

	q := sqlx.Scan(ddisc.LocateSearch(r.Context(), t.q, ddisc.LocateSearchBuilder().OrderBy("id ASC").Offset(msg.Next.Offset*msg.Next.Limit).Limit(msg.Next.Limit)))

	for v := range q.Iter() {
		tmp := langx.Clone(Locate{}, LocateOptionFromDdiscLocate(langx.Clone(v, timex.JSONSafeEncodeOption)))
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

func (t *HTTPLocate) find(w http.ResponseWriter, r *http.Request) {
	var (
		meta ddisc.Locate
		id   = mux.Vars(r)["id"]
	)

	if err := ddisc.LocateFindByID(r.Context(), t.q, id).Scan(&meta); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &LocateLookupResponse{
		Locate: new(
			langx.Clone(
				Locate{},
				LocateOptionFromDdiscLocate(langx.Clone(meta, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPLocate) create(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		d   ddisc.Locate
		msg LocateCreateRequest
	)

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if d, err = NewDdiscLocateFromLocate(msg.Locate); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	d = ddisc.NewLocate(d.Query, d.Mimetype, ddisc.LocateOptionKnownMedia(langx.FirstNonZero(d.KnownMediaID, uuid.Nil.String())))

	if d.KnownMediaID == uuid.Nil.String() && d.Query == "" {
		log.Println("locate requires a known media id or a query")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err := ddisc.LocateInsertWithDefaults(r.Context(), t.q, d).Scan(&d); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create locate record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	t.locate.Broadcast()

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &LocateCreateResponse{
		Locate: new(
			langx.Clone(
				Locate{},
				LocateOptionFromDdiscLocate(langx.Clone(d, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
