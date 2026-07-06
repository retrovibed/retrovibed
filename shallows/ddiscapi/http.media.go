package ddiscapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type HTTPMediaOption func(*HTTPMedia)

func HTTPMediaOptionJWTSecret(j jwtx.SecretSource) HTTPMediaOption {
	return func(t *HTTPMedia) {
		t.jwtsecret = j
	}
}

func NewHTTPMedia(q sqlx.Queryer, options ...HTTPMediaOption) *HTTPMedia {
	svc := langx.Clone(HTTPMedia{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPMedia struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPMedia) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, AuthzPermPeer),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, AuthzPermPeer),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

func (t *HTTPMedia) search(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp = MediaSearchResponse{
			Next: &MediaSearchRequest{
				Offset: 0,
				Limit:  100,
			},
		}
	)

	if err = t.decoder.Decode(resp.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	resp.Next.Limit = numericx.Min(resp.Next.Limit, 100)

	query := ddisc.DiscoveredSearchBuilder().
		Offset(resp.Next.Offset * resp.Next.Limit).Limit(resp.Next.Limit).
		Where(squirrel.And{
			squirrel.Expr("1=1"),
			ddisc.DiscoveredQueryKnownMediaID(resp.Next.KnownMediaId),
			ddisc.DiscoveredQueryByIDs(resp.Next.Id...),
			ddisc.DiscoveredQueryNeedsCheck(resp.Next.NeedsCheck),
			ddisc.DiscoveredQueryText(resp.Next.Query),
		})

	q := sqlx.Scan(ddisc.DiscoveredSearch(r.Context(), t.q, query))
	for d := range q.Iter() {
		resp.Items = append(resp.Items, NewMediaFromDiscovered(d))
	}

	if err = q.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "media generation failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPMedia) create(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg MediaCreateRequest
		d   ddisc.Discovered
	)

	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	options := []ddisc.DiscoveredOption{
		func(d *ddisc.Discovered) {
			d.Title = msg.Media.GetTitle()
			d.Description = msg.Media.GetDescription()
		},
		ddisc.DiscoveredOptionMimetype(msg.Media.GetMimetype()),
	}

	if kid := msg.Media.GetKnownMediaId(); kid != "" {
		options = append(options, ddisc.DiscoveredOptionKnownMedia(kid))
	}

	if partition := msg.Media.GetPartition(); partition != "" {
		options = append(options, ddisc.DiscoveredOptionPartition(partition))
	}

	id := int160.FromBytesOrZero(msg.Media.GetInfohash())
	d = ddisc.NewDiscovered(&id, options...)

	if err = ddisc.DiscoveredInsertWithDefaults(r.Context(), t.q, d).Scan(&d); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create media"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &MediaCreateResponse{
		Media: NewMediaFromDiscovered(d),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPMedia) delete(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		vars = mux.Vars(r)
		d    ddisc.Discovered
	)

	if err = ddisc.DiscoveredDeleteByID(r.Context(), t.q, vars["id"]).Scan(&d); err != nil {
		log.Println(errorsx.Wrap(err, "unable to delete media"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &MediaDeleteResponse{
		Media: NewMediaFromDiscovered(d),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
