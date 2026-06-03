package communityapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type HTTPOption func(*HTTP)

func HTTPOptionNoop(*HTTP) {}

func HTTPOptionJWTSecret(j jwtx.SecretSource) HTTPOption {
	return func(t *HTTP) {
		t.jwtsecret = j
	}
}

func HTTPOptionHTTPClient(c *http.Client) HTTPOption {
	return func(t *HTTP) {
		t.httpc = c
	}
}

func HTTPOptionPublishing(p *asyncx.Wakeup) HTTPOption {
	return func(t *HTTP) {
		t.publishing = p
	}
}

func HTTPOptionMediaStorage(vfs fsx.Virtual) HTTPOption {
	return func(t *HTTP) {
		t.mediastorage = vfs
	}
}

func HTTPOptionTorrentStorage(vfs fsx.Virtual) HTTPOption {
	return func(t *HTTP) {
		t.torrentstorage = vfs
	}
}

func NewHTTP(q sqlx.Queryer, options ...HTTPOption) *HTTP {
	svc := langx.Clone(HTTP{
		q:          q,
		jwtsecret:  env.JWTSecret,
		decoder:    formx.NewDecoder(),
		publishing: asyncx.NewWakeup(context.Background()),
		lucene:     duckdbx.NewLucene(),
	}, options...)

	return &svc
}

type HTTP struct {
	q              sqlx.Queryer
	jwtsecret      jwtx.SecretSource
	decoder        *form.Decoder
	httpc          *http.Client
	publishing     *asyncx.Wakeup
	mediastorage   fsx.Virtual
	torrentstorage fsx.Virtual
	lucene         lucenex.Driver
}

func (t *HTTP) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/published/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.tombstoned))

	r.Path("/{id}/publish").Methods(http.MethodPost).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.publish))

	r.Path("/{id}/published").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.published))

	r.Path("/{id}/metrics").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.metrics))

	r.Path("/{id}/metrics/sync").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.metricsSync))

	r.Path("/{id}/subscribe").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.subscribe))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.resync))
}

func (t *HTTP) tombstoned(w http.ResponseWriter, r *http.Request) {
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
		PublishedContent: langx.Autoptr(
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

func (t *HTTP) publish(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		lmd         library.Metadata
		communityID = mux.Vars(r)["id"]
		req         PublishContentRequest
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
		CommunityID:   communityID,
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
		PublishedContent: langx.Autoptr(
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

func (t *HTTP) resync(w http.ResponseWriter, r *http.Request) {
	var (
		cid   = mux.Vars(r)["id"]
		msg   CommunityFindRequest
		resp  CommunityFindResponse
		syncd community.Community
	)
	_ = &msg

	if err := community.CommunityFindByID(r.Context(), t.q, cid).Scan(&syncd); err != nil {
		log.Println(errorsx.Wrap(err, "unable lookup sync state"))
		return
	}

	published := NewDeeppoolPublished(t.httpc)
	published.List(r.Context(), cid, nil)

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTP) published(w http.ResponseWriter, r *http.Request) {
	communityID := mux.Vars(r)["id"]

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

func (t *HTTP) metrics(w http.ResponseWriter, r *http.Request) {
	communityID := mux.Vars(r)["id"]

	var options []community.MetricPeriodOption

	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		if ts, err := time.Parse(time.RFC3339, startDate); err == nil {
			options = append(options, community.MetricPeriodOptionStartDate(ts))
		}
	}
	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		if ts, err := time.Parse(time.RFC3339, endDate); err == nil {
			options = append(options, community.MetricPeriodOptionEndDate(ts))
		}
	}

	periodStart, periodEnd := community.ResolvePeriod(options...)

	summary, err := community.CommunityMetricAggregateSearch(r.Context(), t.q, communityID, periodStart, periodEnd)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to aggregate community metrics"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	totalArchivers, err := community.PublishedCASMetricAggregateSearch(r.Context(), t.q, communityID, periodStart, periodEnd)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to aggregate archiver metrics"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	msg := CommunityMetricsResponse{
		Summary: langx.Autoptr(langx.Clone(
			CommunityMetric{},
			CommunityMetricOptionFromDB(langx.Clone(summary, timex.JSONSafeEncodeOption)),
		)),
		TotalArchivers: totalArchivers,
	}

	items, err := community.PublishedCASMetricPerContentSearch(r.Context(), t.q, communityID, periodStart, periodEnd)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to fetch content metrics"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}
	for _, m := range items {
		tmp := langx.Clone(PublishedContentMetric{}, PublishedCASMetricOptionFromDB(langx.Clone(m, timex.JSONSafeEncodeOption)))
		msg.Items = append(msg.Items, &tmp)
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

// PublishedCASMetricOptionFromDB converts a database model to proto options.
func PublishedCASMetricOptionFromDB(pcm community.PublishedCASMetric) func(*PublishedContentMetric) {
	return func(m *PublishedContentMetric) {
		m.Id = pcm.ID
		m.PublishedContentId = pcm.PublishedContentID
		m.PeriodStart = grpcx.EncodeTime(pcm.PeriodStart)
		m.PeriodEnd = grpcx.EncodeTime(pcm.PeriodEnd)
		m.Archivers = pcm.Archivers
		m.Bytes = pcm.Bytes
		m.Revenue = pcm.Revenue
	}
}

func (t *HTTP) search(w http.ResponseWriter, r *http.Request) {
	if t.httpc == nil {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusServiceUnavailable))
		return
	}

	var req CommunitySearchRequest
	req.Limit = 100
	if err := t.decoder.Decode(&req, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode search request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	req.Limit = numericx.Min(req.Limit, 100)

	client := NewDeeppoolCommunity(t.httpc)
	resp, err := client.Search(r.Context(), req.Query, req.Offset, req.Limit)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to search communities"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadGateway))
		return
	}

	for _, c := range resp.Items {
		var (
			local community.Community
		)

		if err := community.CommunityInsertWithDefaults(r.Context(), t.q, CommunityFromDeeppool(c)).Scan(&local); err != nil {
			log.Println(errorsx.Wrap(err, "unable to upsert community"))
			continue
		}

		*c = langx.Autoderef(NewCommunity(CommunityOptionFromDB(langx.Clone(local, timex.JSONSafeEncodeOption))))
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTP) subscribe(w http.ResponseWriter, r *http.Request) {
	var (
		existing community.Community
		cid      = mux.Vars(r)["id"]
	)

	client := NewDeeppoolCommunity(t.httpc)
	com, err := client.Find(r.Context(), cid)
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to find community from deeppool - %s", cid))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	err = community.CommunityFindByID(r.Context(), t.q, cid).Scan(&existing)
	if sqlx.IgnoreNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to look up subscription"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err == nil {
		// already subscribed — toggle off
		if err = community.CommunityDeleteByID(r.Context(), t.q, cid).Scan(&existing); err != nil {
			log.Println(errorsx.Wrap(err, "unable to delete subscription"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		var (
			feed    tracking.RSS
			feedURL = fmt.Sprintf("https://%s.community.retrovibe.space", com.Domain)
		)
		if err = tracking.RSSDeleteByURL(r.Context(), t.q, feedURL).Scan(&feed); sqlx.IgnoreNoRows(err) != nil {
			log.Println(errorsx.Wrap(err, "unable to delete rss feed"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}
	} else {
		// not subscribed — subscribe
		var sub = community.Community{
			ID: cid,
		}

		if err = community.CommunityUpsertAutoDownload(r.Context(), t.q, sub).Scan(&sub); err != nil {
			log.Println(errorsx.Wrap(err, "unable to insert subscription"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		var (
			feedURL = fmt.Sprintf("https://%s.community.retrovibe.space", com.Domain)
		)

		feed := tracking.NewFeedRSS(
			"",
			tracking.RSSOptionURL(feedURL),
			tracking.RSSOptionDescription(stringsx.Join(" - ", slicesx.Filter(stringsx.Present, com.Domain, com.Description)...)),
			tracking.RSSOptionEncryptionSeed(com.Entropy),
			tracking.RSSOptionAutodownload(true),
			tracking.RSSOptionAutoarchive(true),
			tracking.RSSOptionAutoID,
		)

		if err = tracking.RSSInsertWithDefaults(r.Context(), t.q, feed).Scan(&feed); err != nil {
			log.Println(errorsx.Wrap(err, "unable to register rss feed"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}
	}

	errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusOK))
}

func (t *HTTP) metricsSync(w http.ResponseWriter, r *http.Request) {
	if t.httpc == nil {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusServiceUnavailable))
		return
	}

	var (
		err         error
		resp        *MetricsSyncResponse
		communityID = mux.Vars(r)["id"]
		metrics     = NewMetrics(t.httpc)
	)

	if resp, err = metrics.Sync(r.Context(), communityID); err != nil {
		log.Println(errorsx.Wrap(err, "failed to fetch metrics from deeppool"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = t.storeMetrics(r.Context(), communityID, resp); err != nil {
		log.Println(errorsx.Wrap(err, "failed to store metrics"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &MetricsSyncProgress{
		Status:                "completed",
		CommunityMetricsCount: int32(len(resp.CommunityMetrics)),
		ContentMetricsCount:   int32(len(resp.ContentMetrics)),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTP) storeMetrics(ctx context.Context, communityID string, resp *MetricsSyncResponse) error {
	var (
		err         error
		periodStart time.Time
		periodEnd   time.Time
		syncedAt    time.Time
		syncState   community.Community
	)

	for _, cm := range resp.CommunityMetrics {
		if periodStart, err = grpcx.DecodeTime(cm.PeriodStart); err != nil {
			return errorsx.Wrap(err, "failed to decode period start")
		}
		if periodEnd, err = grpcx.DecodeTime(cm.PeriodEnd); err != nil {
			return errorsx.Wrap(err, "failed to decode period end")
		}
		local := community.CommunityMetric{
			CommunityID: cm.CommunityId,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
			Subscribers: cm.Subscribers,
		}
		if err = community.CommunityMetricInsertWithDefaults(ctx, t.q, local).Scan(&local); err != nil {
			return errorsx.Wrap(err, "failed to insert community metric")
		}
	}

	for _, pcm := range resp.ContentMetrics {
		if periodStart, err = grpcx.DecodeTime(pcm.PeriodStart); err != nil {
			return errorsx.Wrap(err, "failed to decode period start")
		}
		if periodEnd, err = grpcx.DecodeTime(pcm.PeriodEnd); err != nil {
			return errorsx.Wrap(err, "failed to decode period end")
		}
		local := community.PublishedCASMetric{
			PublishedContentID: pcm.PublishedContentId,
			PeriodStart:        periodStart,
			PeriodEnd:          periodEnd,
			Archivers:          pcm.Archivers,
			Bytes:              pcm.Bytes,
			Revenue:            pcm.Revenue,
		}
		if err = community.PublishedCASMetricInsertWithDefaults(ctx, t.q, local).Scan(&local); err != nil {
			return errorsx.Wrap(err, "failed to insert content metric")
		}
	}

	if syncedAt, err = grpcx.DecodeTime(resp.SyncedAt); err != nil {
		return errorsx.Wrap(err, "failed to decode synced at")
	}
	syncState = community.Community{
		ID:         communityID,
		LastSyncAt: syncedAt,
	}
	if err = community.CommunityInsertWithDefaults(ctx, t.q, syncState).Scan(&syncState); err != nil {
		return errorsx.Wrap(err, "failed to update sync state")
	}

	return nil
}
