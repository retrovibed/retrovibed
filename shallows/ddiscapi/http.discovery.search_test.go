package ddiscapi_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestHTTPDiscoverySearch(t *testing.T) {
	t.Run("search", func(t *testing.T) {
		var (
			uh     tracking.UnknownHash
			result ddiscapi.DiscoverySearchResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&uh, tracking.UnknownHashOptionTestDefaults))
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, uh).Scan(&uh))

		routes := mux.NewRouter()

		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(formx.NewEncoder().Encode(&ddiscapi.DiscoverySearchRequest{
			Offset: 0,
		}))(t)

		claims = jwtx.NewJWTClaims(uh.ID, jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.Equal(t, result.Next.Offset, uint64(0))
		require.Contains(t, result.Items, ddiscapi.NewDiscoveryFromTrackingUnknownHash(uh))
	})

	t.Run("next check filter", func(t *testing.T) {
		var (
			due    tracking.UnknownHash
			future tracking.UnknownHash
			result ddiscapi.DiscoverySearchResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&due, tracking.UnknownHashOptionTestDefaults))
		due.NextCheck = time.Now().Add(-time.Minute)
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, due).Scan(&due))

		require.NoError(t, testx.Fake(&future, tracking.UnknownHashOptionTestDefaults))
		future.NextCheck = time.Now().Add(time.Hour)
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, future).Scan(&future))

		routes := mux.NewRouter()

		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(formx.NewEncoder().Encode(&ddiscapi.DiscoverySearchRequest{
			NextCheck: meta.NewDateRange(timex.NewRangeDuration(2 * time.Minute)),
		}))(t)

		claims = jwtx.NewJWTClaims(due.ID, jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.Contains(t, result.Items, ddiscapi.NewDiscoveryFromTrackingUnknownHash(due))
		require.Len(t, result.Items, 1)
	})

	t.Run("needs check disabled returns entries regardless of next_check", func(t *testing.T) {
		var (
			due    tracking.UnknownHash
			future tracking.UnknownHash
			result ddiscapi.DiscoverySearchResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&due, tracking.UnknownHashOptionTestDefaults))
		due.NextCheck = time.Now().Add(-time.Minute)
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, due).Scan(&due))

		require.NoError(t, testx.Fake(&future, tracking.UnknownHashOptionTestDefaults))
		future.NextCheck = time.Now().Add(time.Hour)
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, future).Scan(&future))

		routes := mux.NewRouter()

		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(formx.NewEncoder().Encode(&ddiscapi.DiscoverySearchRequest{}))(t)

		claims = jwtx.NewJWTClaims(due.ID, jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		encodeddue := ddiscapi.NewDiscoveryFromTrackingUnknownHash(due)
		encodedfuture := ddiscapi.NewDiscoveryFromTrackingUnknownHash(future)
		require.Contains(t, result.Items, encodeddue)
		require.Contains(t, result.Items, encodedfuture)
		require.Len(t, result.Items, 2)
	})
}
