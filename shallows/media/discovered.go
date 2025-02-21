package media

import (
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/deeppool/internal/env"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/formx"
	"github.com/james-lawrence/deeppool/internal/x/httpx"
	"github.com/james-lawrence/deeppool/internal/x/jwtx"
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/numericx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/internal/x/sqlxx"
	"github.com/james-lawrence/deeppool/tracking"
	"github.com/justinas/alice"
)

type HTTPDiscoveredOption func(*HTTPDiscovered)

func HTTPDiscoveredOptionJWTSecret(j jwtx.JWTSecretSource) HTTPDiscoveredOption {
	return func(t *HTTPDiscovered) {
		t.jwtsecret = j
	}
}

func NewHTTPDiscovered(q sqlx.Queryer, options ...HTTPDiscoveredOption) *HTTPDiscovered {
	svc := langx.Clone(HTTPDiscovered{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPDiscovered struct {
	q         sqlx.Queryer
	jwtsecret jwtx.JWTSecretSource
	decoder   *form.Decoder
}

func (t *HTTPDiscovered) Bind(r *mux.Router) {
	r.StrictSlash(false)
	r.Use(httpx.RouteInvoked)

	r.Path("/available").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/downloading").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.downloading))
}

func (t *HTTPDiscovered) downloading(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg DownloadSearchResponse = DownloadSearchResponse{
			Next: &DownloadSearchRequest{
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

	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryInitiated(),
			tracking.MetadataQueryIncomplete(),
		},
	).OrderBy("created_at DESC").Offset(msg.Next.Offset * msg.Next.Limit).Limit(msg.Next.Limit)

	err = sqlxx.ScanEach(tracking.MetadataSearch(r.Context(), t.q, q), func(p *tracking.Metadata) error {
		tmp := langx.Clone(Download{}, DownloadOptionFromTorrentMetadata(langx.Clone(*p, tracking.MetadataOptionJSONSafeEncode)))
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

func (t *HTTPDiscovered) search(w http.ResponseWriter, r *http.Request) {
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

	q := tracking.MetadataSearchBuilder().OrderBy("created_at DESC").Offset(msg.Next.Offset * msg.Next.Limit).Limit(msg.Next.Limit)

	err = sqlxx.ScanEach(tracking.MetadataSearch(r.Context(), t.q, q), func(p *tracking.Metadata) error {
		tmp := langx.Clone(Media{}, MediaOptionFromTorrentMetadata(langx.Clone(*p, tracking.MetadataOptionJSONSafeEncode)))
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
