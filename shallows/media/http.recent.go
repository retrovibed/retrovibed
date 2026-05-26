package media

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
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
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"google.golang.org/protobuf/proto"
)

type HTTPRecentOption func(*HTTPRecent)

func HTTPRecentOptionJWTSecret(j jwtx.SecretSource) HTTPRecentOption {
	return func(t *HTTPRecent) {
		t.jwtsecret = j
	}
}

func NewHTTPRecent(q sqlx.Queryer, options ...HTTPRecentOption) *HTTPRecent {
	svc := langx.Clone(HTTPRecent{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...)
	return &svc
}

type HTTPRecent struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPRecent) Bind(r *mux.Router) {
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
	).ThenFunc(t.record))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.tombstone))
}

func (t *HTTPRecent) latest(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		created timex.Range
		msg     = RecentSearchResponse{
			Next: &RecentSearchRequest{
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

	if created, err = timex.NewRangeISO8601(msg.Next.Created.Oldest, msg.Next.Created.Newest); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode created"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	query := library.RecentSessionLibrarySearchBuilder().Where(
		squirrel.And{
			library.MetadataQueryNotTombstoned(),
			library.RecentSessionQueryCreated(created),
			library.RecentSessionQueryMimetype(msg.Next.Mimetype),
		},
	).OrderBy("library_recent_sessions.last_played_at DESC").Limit(msg.Next.Limit)

	q := sqlx.Scan2(library.RecentSessionLibrarySearch(r.Context(), t.q, query))
	for s, md := range q.Iter() {
		var (
			req MediaSearchRequest
		)

		m := langx.Clone(Media{}, MediaOptionFromLibraryMetadata(langx.Clone(md, timex.JSONSafeEncodeOption)))
		if err := proto.Unmarshal(s.Query, &req); err != nil {
			log.Println(errorsx.Wrap(err, "unable to decode query"))
			continue
		}

		msg.Items = append(msg.Items, &RecentRecordRequest{
			Id:       s.ID,
			Media:    &m,
			Query:    &req,
			Duration: uint64(s.Duration / time.Millisecond),
			Position: uint64(s.Position / time.Millisecond),
			Mimetype: s.Mimetype,
		})
	}

	if err = q.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "recent session query failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRecent) tombstone(w http.ResponseWriter, r *http.Request) {
	var (
		rs library.RecentSession
		id = mux.Vars(r)["id"]
	)

	if err := library.RecentSessionDeleteByID(r.Context(), t.q, id).Scan(&rs); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to tombstone recent session"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to tombstone recent session"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &RecentDeleteResponse{}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPRecent) record(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		rs  library.RecentSession
		msg RecentRecordRequest
	)

	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	encoded, err := proto.Marshal(msg.Query)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to encode query"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = library.RecentSessionInsertWithDefaults(r.Context(), t.q, library.RecentSession{
		ID:       md5x.FormatUUID(md5x.Digest(encoded)),
		Mimetype: msg.Mimetype,
		MediaID:  msg.Media.Id,
		Duration: time.Duration(msg.Duration) * time.Millisecond,
		Position: time.Duration(msg.Position) * time.Millisecond,
		Query:    encoded,
	}).Scan(&rs); err != nil {
		log.Println(errorsx.Wrap(err, "upsert failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &RecentRecordResponse{}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
