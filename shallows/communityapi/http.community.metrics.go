package communityapi

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

type HTTPMetricsOption func(*HTTPMetrics)

func HTTPMetricsOptionNoop(*HTTPMetrics) {}

func HTTPMetricsOptionHTTPClient(c *http.Client) HTTPMetricsOption {
	return func(t *HTTPMetrics) {
		t.httpc = c
	}
}

func HTTPMetricsOptionJWTSecret(j jwtx.SecretSource) HTTPMetricsOption {
	return func(t *HTTPMetrics) {
		t.jwtsecret = j
	}
}

func NewHTTPMetrics(q sqlx.Queryer, options ...HTTPMetricsOption) *HTTPMetrics {
	svc := langx.Clone(HTTPMetrics{
		q:         q,
		jwtsecret: env.JWTSecret,
	}, options...)
	return &svc
}

type HTTPMetrics struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	httpc     *http.Client
}

func (t *HTTPMetrics) Bind(r *mux.Router) {
	r.StrictSlash(false)

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
}

func (t *HTTPMetrics) metrics(w http.ResponseWriter, r *http.Request) {
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

func (t *HTTPMetrics) metricsSync(w http.ResponseWriter, r *http.Request) {
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

func (t *HTTPMetrics) storeMetrics(ctx context.Context, communityID string, resp *MetricsSyncResponse) error {
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
