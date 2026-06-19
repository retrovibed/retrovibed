package communityapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type HTTPPublishedOption func(*HTTPPublished)

func HTTPPublishedOptionNoop(*HTTPPublished) {}

func HTTPPublishedOptionJWTSecret(j jwtx.SecretSource) HTTPPublishedOption {
	return func(t *HTTPPublished) {
		t.jwtsecret = j
	}
}

func HTTPPublishedOptionHTTPClient(c *http.Client) HTTPPublishedOption {
	return func(t *HTTPPublished) {
		t.httpc = c
	}
}

func HTTPPublishedOptionPublishing(p *asyncx.Wakeup) HTTPPublishedOption {
	return func(t *HTTPPublished) {
		t.publishing = p
	}
}

func HTTPPublishedOptionMediaStorage(vfs fsx.Virtual) HTTPPublishedOption {
	return func(t *HTTPPublished) {
		t.mediastorage = vfs
	}
}

func HTTPPublishedOptionTorrentStorage(vfs fsx.Virtual) HTTPPublishedOption {
	return func(t *HTTPPublished) {
		t.torrentstorage = vfs
	}
}

func NewHTTPPublished(q sqlx.Queryer, options ...HTTPPublishedOption) *HTTPPublished {
	svc := langx.Clone(HTTPPublished{
		q:          q,
		jwtsecret:  env.JWTSecret,
		decoder:    formx.NewDecoder(),
		publishing: asyncx.NewWakeup(context.Background()),
		lucene:     duckdbx.NewLucene(),
	}, options...)
	return &svc
}

type HTTPPublished struct {
	q              sqlx.Queryer
	jwtsecret      jwtx.SecretSource
	httpc          *http.Client
	publishing     *asyncx.Wakeup
	decoder        *form.Decoder
	lucene         lucenex.Driver
	mediastorage   fsx.Virtual
	torrentstorage fsx.Virtual
}

func (t *HTTPPublished) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/{cid}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/{cid}").Methods(http.MethodPost).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.publish))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.tombstoned))

	// r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
	// 	httpx.RouteInvoked,
	// 	httpx.ContextBufferPool512(),
	// 	httpauth.AuthenticateWithToken(t.jwtsecret),
	// 	httpx.Timeout2s(),
	// ).ThenFunc(t.resync))
}

func (t *HTTPPublished) tombstoned(w http.ResponseWriter, r *http.Request) {
	pid := mux.Vars(r)["id"]

	var (
		cs community.Community
		pc community.PublishedContent
	)

	if err := community.PublishedContentDeleteByID(r.Context(), t.q, pid).Scan(&pc); errors.Is(err, sql.ErrNoRows) {
		log.Println(errorsx.Wrap(err, "unable to tombstone missing record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	}

	if err := community.CommunityRequestFeedSync(r.Context(), t.q, community.Community{
		ID:         pc.CommunityID,
		SyncFeedAt: time.Now(),
	}).Scan(&cs); err != nil {
		log.Println(errorsx.Wrap(err, "unable to request feed sync for community"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &PublishContentDeleteResponse{
		PublishedContent: new(
			langx.Clone(
				PublishedContent{},
				PublishedContentOptionFromDB(langx.Clone(pc, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}

	t.publishing.Broadcast()
}

func (t *HTTPPublished) publish(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		lmd library.Metadata
		cid = mux.Vars(r)["cid"]
		req PublishContentRequest
	)

	if t.httpc == nil {
		log.Println("http client is missing - unable to publish content")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusServiceUnavailable))
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode publish request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if req.PublishedContent.LibraryId == "" {
		log.Println("library_id is required")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = library.MetadataFindByID(r.Context(), t.q, req.PublishedContent.LibraryId).Scan(&lmd); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "library item not found"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find library metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	pc := community.NewPublishedContent(community.PublishedContent{
		Title:         stringsx.FirstNonBlank(req.PublishedContent.Title, lmd.Description),
		Description:   req.PublishedContent.Description,
		CommunityID:   cid,
		KnownMediaID:  stringsx.FirstNonBlank(req.PublishedContent.KnownMediaId, lmd.KnownMediaID),
		LibraryID:     lmd.ID,
		PublishMode:   int32(req.PublishMode),
		OAuthGoogleID: req.PublishedContent.OauthGoogleId,
		Bytes:         lmd.Bytes,
		Mimetype:      lmd.Mimetype,
		PublishedAt:   errorsx.Zero(grpcx.DecodeTime(langx.FirstNonZero(req.PublishedContent.PublishedAt, grpcx.EncodeTime(timex.Inf())))),
	})

	if err = community.PublishedContentInsertWithDefaults(r.Context(), t.q, pc).Scan(&pc); err != nil {
		log.Println(errorsx.Wrap(err, "unable to insert published content"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &PublishContentResponse{
		PublishedContent: new(
			langx.Clone(
				PublishedContent{},
				PublishedContentOptionFromDB(langx.Clone(pc, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}

	if req.PublishMode > PublishMode_UNLISTED {
		t.publishing.Broadcast()
	}
}

// func (t *HTTPPublished) resync(w http.ResponseWriter, r *http.Request) {
// 	var (
// 		cid   = mux.Vars(r)["id"]
// 		syncd community.Community
// 	)

// 	if err := community.CommunityFindByID(r.Context(), t.q, cid).Scan(&syncd); err != nil {
// 		log.Println(errorsx.Wrap(err, "unable lookup community"))
// 		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
// 		return
// 	}

// 	published := NewDeeppoolPublished(t.httpc)
// 	// TODO JAL 2026-06-02: kickoff background async if # of results equals the limit.
// 	pubed, err := published.List(r.Context(), cid, &PublishedContentSearchRequest{
// 		Sync:   syncd.SyncCursorPublishedContent,
// 		Limit:  1024,
// 		Offset: 0,
// 	})
// 	if err != nil {
// 		log.Println(errorsx.Wrap(err, "unable lookup published content to sync"))
// 		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
// 		return
// 	}

// 	for _, pc := range pubed.Items {
// 		dbpc := community.NewPublishedContent(community.PublishedContent{
// 			ID:            pc.Id,
// 			CommunityID:   pc.CommunityId,
// 			MagnetURI:     pc.MagnetUri,
// 			LibraryID:     stringsx.FirstNonBlank(pc.LibraryId, uuid.Nil.String()),
// 			OAuthGoogleID: pc.OauthGoogleId,
// 			KnownMediaID:  pc.KnownMediaId,
// 		})
// 		if err := community.PublishedContentInsertWithDefaults(r.Context(), t.q, dbpc).Scan(&dbpc); err != nil {
// 			log.Println(errorsx.Wrap(err, "failed to sync published content item"))
// 		}
// 	}

// 	syncd.SyncCursorPublishedContent = langx.FirstNonZero(pubed.GetNext().GetSync(), syncd.SyncCursorPublishedContent)
// 	if err := community.CommunityInsertWithDefaults(r.Context(), t.q, syncd).Scan(&syncd); err != nil {
// 		log.Println(errorsx.Wrap(err, "failed to update sync cursor"))
// 		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
// 		return
// 	}

// 	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), pubed); err != nil {
// 		log.Println(errorsx.Wrap(err, "unable to write response"))
// 		return
// 	}
// }

func (t *HTTPPublished) search(w http.ResponseWriter, r *http.Request) {
	communityID := mux.Vars(r)["cid"]

	var msg PublishedContentSearchResponse
	msg.Next = &PublishedContentSearchRequest{
		CommunityId: communityID,
		Limit:       100,
	}

	if err := t.decoder.Decode(msg.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	msg.Next.Limit = numericx.Min(msg.Next.Limit, 100)

	q := sqlx.Scan(community.PublishedContentSearch(r.Context(), t.q, community.PublishedContentSearchBuilder().Where(
		squirrel.And{
			community.PublishedContentQueryCommunityID(communityID),
			community.PublishedContentQueryNotTombstoned(),
			lucenex.Query(t.lucene, msg.Next.Query, lucenex.WithDefaultField("title")),
		},
	).OrderBy("published_at DESC")))

	for pc := range q.Iter() {
		tmp := langx.Clone(PublishedContent{}, PublishedContentOptionFromDB(langx.Clone(pc, timex.JSONSafeEncodeOption)))
		msg.Items = append(msg.Items, &tmp)
	}

	if err := q.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to fetch published content"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
