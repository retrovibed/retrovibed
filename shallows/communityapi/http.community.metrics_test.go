package communityapi_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func newMetricsMockClient(communityID string, resp *communityapi.MetricsSyncResponse) *http.Client {
	return httptestx.NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == http.MethodGet && strings.Contains(req.URL.Path, communityID) {
			body, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(body))),
				Header:     make(http.Header),
			}
		}
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader(`{}`)),
			Header:     make(http.Header),
		}
	})
}

func TestMetricsSearchEndpoint(t *testing.T) {
	t.Run("returns empty metrics for community with no data", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var (
			p   meta.Profile
			v   meta.Authz
			q   = sqltestx.Metadatabase(t)
			cid = uuid.Must(uuid.NewV7()).String()
		)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		communityapi.NewHTTPMetrics(
			q,
			communityapi.HTTPMetricsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/c/m").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodGet,
			"/c/m/"+cid,
			nil,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.CommunityMetricsResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotNil(t, result.Summary)
		require.Equal(t, uint32(0), result.Summary.Subscribers)
		require.Empty(t, result.Items)
	})

	t.Run("returns aggregated subscriber count for period", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var (
			p           meta.Profile
			v           meta.Authz
			q           = sqltestx.Metadatabase(t)
			communityID = uuid.Must(uuid.NewV7()).String()
			now         = time.Now()
		)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		m := community.CommunityMetric{
			CommunityID: communityID,
			PeriodStart: now.Add(-24 * time.Hour),
			PeriodEnd:   now,
			Subscribers: 42,
		}
		require.NoError(t, community.CommunityMetricInsertWithDefaults(ctx, q, m).Scan(&m))

		routes := mux.NewRouter()
		communityapi.NewHTTPMetrics(
			q,
			communityapi.HTTPMetricsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/c/m").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodGet,
			"/c/m/"+communityID,
			nil,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.CommunityMetricsResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotNil(t, result.Summary)
		require.Equal(t, uint32(42), result.Summary.Subscribers)
	})

	t.Run("returns content-level archiver metrics", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var (
			p                meta.Profile
			v                meta.Authz
			q                = sqltestx.Metadatabase(t)
			communityID      = uuid.Must(uuid.NewV7()).String()
			publishedContent community.PublishedContent
			now              = time.Now()
		)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		publishedContent = community.NewPublishedContent(community.PublishedContent{
			CommunityID: communityID,
			LibraryID:   uuid.Nil.String(),
		})
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, publishedContent).Scan(&publishedContent))

		pcm := community.PublishedCASMetric{
			PublishedContentID: publishedContent.ID,
			PeriodStart:        now.Add(-24 * time.Hour),
			PeriodEnd:          now,
			Archivers:          5,
			Bytes:              1024,
			Revenue:            100,
		}
		require.NoError(t, community.PublishedCASMetricInsertWithDefaults(ctx, q, pcm).Scan(&pcm))

		routes := mux.NewRouter()
		communityapi.NewHTTPMetrics(
			q,
			communityapi.HTTPMetricsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/c/m").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodGet,
			"/c/m/"+communityID,
			nil,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.CommunityMetricsResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.Len(t, result.Items, 1)
		require.Equal(t, publishedContent.ID, result.Items[0].PublishedContentId)
		require.Equal(t, uint32(5), result.Items[0].Archivers)
		require.Equal(t, int64(1024), result.Items[0].Bytes)
		require.Equal(t, int64(100), result.Items[0].Revenue)
	})
}

func TestMetricsSyncEndpoint(t *testing.T) {
	t.Run("sync succeeds when community has not been upserted locally", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var (
			p           meta.Profile
			v           meta.Authz
			q           = sqltestx.Metadatabase(t)
			communityID = uuid.Must(uuid.NewV7()).String()
			syncedAt    = time.Now().UTC().Truncate(time.Second)
		)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		// no community row inserted — simulates a community only known via deeppool

		syncResp := &communityapi.MetricsSyncResponse{
			SyncedAt: grpcx.EncodeTime(syncedAt),
			CommunityMetrics: []*communityapi.CommunityMetric{
				{
					CommunityId: communityID,
					PeriodStart: grpcx.EncodeTime(syncedAt.Add(-24 * time.Hour)),
					PeriodEnd:   grpcx.EncodeTime(syncedAt),
					Subscribers: 10,
				},
			},
		}

		routes := mux.NewRouter()
		communityapi.NewHTTPMetrics(
			q,
			communityapi.HTTPMetricsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPMetricsOptionHTTPClient(newMetricsMockClient(communityID, syncResp)),
		).Bind(routes.PathPrefix("/c/m").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/m/"+communityID,
			nil,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var updated community.Community
		require.NoError(t, community.CommunityFindByID(ctx, q, communityID).Scan(&updated))
		require.WithinDuration(t, syncedAt, updated.LastSyncAt, time.Second)
	})

	t.Run("returns 503 when no http client configured", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var (
			p           meta.Profile
			v           meta.Authz
			q           = sqltestx.Metadatabase(t)
			communityID = uuid.Must(uuid.NewV7()).String()
		)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		communityapi.NewHTTPMetrics(
			q,
			communityapi.HTTPMetricsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/c/m").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/m/"+communityID,
			nil,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusServiceUnavailable, resp.Code)
	})

	t.Run("stores community and content metrics from deeppool response", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var (
			p                meta.Profile
			v                meta.Authz
			q                = sqltestx.Metadatabase(t)
			communityID      = uuid.Must(uuid.NewV7()).String()
			publishedContent community.PublishedContent
			syncedAt         = time.Now().UTC().Truncate(time.Second)
			periodStart      = syncedAt.Add(-24 * time.Hour)
			periodEnd        = syncedAt
		)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		// community must exist for CommunityUpdateLastSyncAt to succeed
		com := community.Community{ID: communityID}
		require.NoError(t, community.CommunityUpsertAutoDownload(ctx, q, com).Scan(&com))

		publishedContent = community.NewPublishedContent(community.PublishedContent{
			CommunityID: communityID,
			LibraryID:   uuid.Nil.String(),
		})
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, publishedContent).Scan(&publishedContent))

		syncResp := &communityapi.MetricsSyncResponse{
			SyncedAt: grpcx.EncodeTime(syncedAt),
			CommunityMetrics: []*communityapi.CommunityMetric{
				{
					CommunityId: communityID,
					PeriodStart: grpcx.EncodeTime(periodStart),
					PeriodEnd:   grpcx.EncodeTime(periodEnd),
					Subscribers: 99,
				},
			},
			ContentMetrics: []*communityapi.PublishedContentMetric{
				{
					PublishedContentId: publishedContent.ID,
					PeriodStart:        grpcx.EncodeTime(periodStart),
					PeriodEnd:          grpcx.EncodeTime(periodEnd),
					Archivers:          7,
					Bytes:              2048,
					Revenue:            500,
				},
			},
		}

		routes := mux.NewRouter()
		communityapi.NewHTTPMetrics(
			q,
			communityapi.HTTPMetricsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPMetricsOptionHTTPClient(newMetricsMockClient(communityID, syncResp)),
		).Bind(routes.PathPrefix("/c/m").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/m/"+communityID,
			nil,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var progress communityapi.MetricsSyncProgress
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &progress))
		require.Equal(t, "completed", progress.Status)
		require.Equal(t, int32(1), progress.CommunityMetricsCount)
		require.Equal(t, int32(1), progress.ContentMetricsCount)

		// community metrics stored
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM community_metrics WHERE community_id = '"+communityID+"'"))

		// content metrics stored
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM published_cas_metrics WHERE published_content_id = '"+publishedContent.ID+"'"))

		// last_sync_at updated on community
		var updated community.Community
		require.NoError(t, community.CommunityFindByID(ctx, q, communityID).Scan(&updated))
		require.WithinDuration(t, syncedAt, updated.LastSyncAt, time.Second)
	})
}
