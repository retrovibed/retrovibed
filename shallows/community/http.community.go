package community

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/meta"
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

func HTTPOptionArchival(a *asyncx.Wakeup) HTTPOption {
	return func(t *HTTP) {
		t.archival = a
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
		archival:   asyncx.NewWakeup(context.Background()),
		publishing: asyncx.NewWakeup(context.Background()),
	}, options...)

	return &svc
}

type HTTP struct {
	q              sqlx.Queryer
	jwtsecret      jwtx.SecretSource
	decoder        *form.Decoder
	httpc          *http.Client
	archival       *asyncx.Wakeup
	publishing     *asyncx.Wakeup
	mediastorage   fsx.Virtual
	torrentstorage fsx.Virtual
}

func (t *HTTP) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/{id}/publish").Methods(http.MethodPost).Handler(alice.New(
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
}

func (t *HTTP) publish(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		lmd         library.Metadata
		communityID = mux.Vars(r)["id"]
		req         meta.PublishContentRequest
	)

	if t.httpc == nil {
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

	pc := PublishedContent{
		CommunityID:   communityID,
		KnownMediaID:  stringsx.FirstNonBlank(req.PublishedContent.KnownMediaId, lmd.KnownMediaID),
		LibraryID:     lmd.ID,
		PublishMode:   int32(req.PublishMode),
		OAuthGoogleID: stringsx.FirstNonBlank(req.PublishedContent.OauthGoogleId, uuid.Nil.String()),
		Bytes:         lmd.Bytes,
	}

	if err = PublishedContentInsertWithDefaults(r.Context(), t.q, pc).Scan(&pc); err != nil {
		log.Println(errorsx.Wrap(err, "unable to insert published content"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &meta.PublishContentResponse{
		PublishedContent: langx.Autoptr(
			langx.Clone(
				meta.PublishedContent{},
				PublishedContentOptionFromDB(langx.Clone(pc, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}

	if req.PublishMode > meta.PublishMode_UNLISTED {
		t.publishing.Broadcast()
	}
}

func (t *HTTP) published(w http.ResponseWriter, r *http.Request) {
	communityID := mux.Vars(r)["id"]

	var msg meta.PublishedContentListResponse
	msg.Next = &meta.PublishedContentListRequest{
		CommunityId: communityID,
		Limit:       100,
	}

	if err := t.decoder.Decode(msg.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	msg.Next.Limit = numericx.Min(msg.Next.Limit, 100)

	q := sqlx.Scan(PublishedContentFindByCommunityID(r.Context(), t.q, communityID))
	for pc := range q.Iter() {
		tmp := langx.Clone(meta.PublishedContent{}, PublishedContentOptionFromDB(langx.Clone(pc, timex.JSONSafeEncodeOption)))
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

	var options []MetricPeriodOption

	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		if ts, err := time.Parse(time.RFC3339, startDate); err == nil {
			options = append(options, MetricPeriodOptionStartDate(ts))
		}
	}
	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		if ts, err := time.Parse(time.RFC3339, endDate); err == nil {
			options = append(options, MetricPeriodOptionEndDate(ts))
		}
	}

	periodStart, periodEnd := ResolvePeriod(options...)

	summary, err := CommunityMetricAggregateSearch(r.Context(), t.q, communityID, periodStart, periodEnd)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to aggregate community metrics"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	totalArchivers, err := PublishedCASMetricAggregateSearch(r.Context(), t.q, communityID, periodStart, periodEnd)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to aggregate archiver metrics"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	msg := meta.CommunityMetricsResponse{
		Summary: langx.Autoptr(langx.Clone(
			meta.CommunityMetric{},
			CommunityMetricOptionFromDB(langx.Clone(summary, timex.JSONSafeEncodeOption)),
		)),
		TotalArchivers: totalArchivers,
	}

	items, err := PublishedCASMetricPerContentSearch(r.Context(), t.q, communityID, periodStart, periodEnd)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to fetch content metrics"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}
	for _, m := range items {
		tmp := langx.Clone(meta.PublishedContentMetric{}, PublishedCASMetricOptionFromDB(langx.Clone(m, timex.JSONSafeEncodeOption)))
		msg.Items = append(msg.Items, &tmp)
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

// PublishedContentOptionFromDB converts a database model to proto options.
func PublishedContentOptionFromDB(pc PublishedContent) func(*meta.PublishedContent) {
	return func(p *meta.PublishedContent) {
		p.Id = pc.ID
		p.CommunityId = pc.CommunityID
		p.KnownMediaId = pc.KnownMediaID
		p.MagnetUri = pc.MagnetURI
		p.LibraryId = pc.LibraryID
		p.OauthGoogleId = pc.OAuthGoogleID
		p.PublishedAt = grpcx.EncodeTime(pc.PublishedAt)
		p.CreatedAt = grpcx.EncodeTime(pc.CreatedAt)
		p.UpdatedAt = grpcx.EncodeTime(pc.UpdatedAt)
		p.Bytes = pc.Bytes
	}
}

// CommunityMetricOptionFromDB converts a database model to proto options.
func CommunityMetricOptionFromDB(cm CommunityMetric) func(*meta.CommunityMetric) {
	return func(m *meta.CommunityMetric) {
		m.Id = cm.ID
		m.CommunityId = cm.CommunityID
		m.PeriodStart = grpcx.EncodeTime(cm.PeriodStart)
		m.PeriodEnd = grpcx.EncodeTime(cm.PeriodEnd)
		m.Subscribers = cm.Subscribers
	}
}

// PublishedCASMetricOptionFromDB converts a database model to proto options.
func PublishedCASMetricOptionFromDB(pcm PublishedCASMetric) func(*meta.PublishedContentMetric) {
	return func(m *meta.PublishedContentMetric) {
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

	var req meta.CommunitySearchRequest
	req.Limit = 100
	if err := t.decoder.Decode(&req, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode search request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	req.Limit = numericx.Min(req.Limit, 100)

	client := deeppool.NewPublished(t.httpc)
	resp, err := client.Search(r.Context(), req.Query, req.Offset, req.Limit)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to search communities"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadGateway))
		return
	}

	subs := make(map[string]time.Time)
	q := sqlx.Scan(CommunitySubscriptionFindAll(r.Context(), t.q))
	for sub := range q.Iter() {
		subs[sub.CommunityID] = sub.CreatedAt
	}

	if err := q.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to fetch subscriptions"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	for _, c := range resp.Items {
		if subscribedAt, ok := subs[c.Id]; ok {
			c.SubscribedAt = grpcx.EncodeTime(subscribedAt)
		}
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTP) subscribe(w http.ResponseWriter, r *http.Request) {
	var (
		existing CommunitySubscription
		cid      = mux.Vars(r)["id"]
	)

	client := deeppool.NewPublished(t.httpc)
	com, err := client.Find(r.Context(), cid)
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to find community from deeppool - %s", cid))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	err = CommunitySubscriptionFindByCommunityID(r.Context(), t.q, cid).Scan(&existing)
	if sqlx.IgnoreNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to look up subscription"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err == nil {
		// already subscribed — toggle off
		if err = CommunitySubscriptionDeleteByCommunityID(r.Context(), t.q, cid).Scan(&existing); err != nil {
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
		var sub = CommunitySubscription{
			CommunityID: cid,
		}

		if err = CommunitySubscriptionInsertWithDefaults(r.Context(), t.q, sub).Scan(&sub); err != nil {
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
		resp        *meta.MetricsSyncResponse
		communityID = mux.Vars(r)["id"]
		metrics     = deeppool.NewMetrics(t.httpc)
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

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &meta.MetricsSyncProgress{
		Status:                "completed",
		CommunityMetricsCount: int32(len(resp.CommunityMetrics)),
		ContentMetricsCount:   int32(len(resp.ContentMetrics)),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTP) storeMetrics(ctx context.Context, communityID string, resp *meta.MetricsSyncResponse) error {
	var (
		err         error
		periodStart time.Time
		periodEnd   time.Time
		syncedAt    time.Time
		syncState   CommunitySyncState
	)

	for _, cm := range resp.CommunityMetrics {
		if periodStart, err = grpcx.DecodeTime(cm.PeriodStart); err != nil {
			return errorsx.Wrap(err, "failed to decode period start")
		}
		if periodEnd, err = grpcx.DecodeTime(cm.PeriodEnd); err != nil {
			return errorsx.Wrap(err, "failed to decode period end")
		}
		local := CommunityMetric{
			CommunityID: cm.CommunityId,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
			Subscribers: cm.Subscribers,
		}
		if err = CommunityMetricInsertWithDefaults(ctx, t.q, local).Scan(&local); err != nil {
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
		local := PublishedCASMetric{
			PublishedContentID: pcm.PublishedContentId,
			PeriodStart:        periodStart,
			PeriodEnd:          periodEnd,
			Archivers:          pcm.Archivers,
			Bytes:              pcm.Bytes,
			Revenue:            pcm.Revenue,
		}
		if err = PublishedCASMetricInsertWithDefaults(ctx, t.q, local).Scan(&local); err != nil {
			return errorsx.Wrap(err, "failed to insert content metric")
		}
	}

	if syncedAt, err = grpcx.DecodeTime(resp.SyncedAt); err != nil {
		return errorsx.Wrap(err, "failed to decode synced at")
	}
	syncState = CommunitySyncState{
		CommunityID: communityID,
		LastSyncAt:  syncedAt,
	}
	if err = CommunitySyncStateInsertWithDefaults(ctx, t.q, syncState).Scan(&syncState); err != nil {
		return errorsx.Wrap(err, "failed to update sync state")
	}

	return nil
}
