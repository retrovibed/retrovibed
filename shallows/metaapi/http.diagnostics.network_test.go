package metaapi_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/wireguardx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPDiagnosticsNetwork(t *testing.T) {
	t.Run("healthy", func(t *testing.T) {
		var (
			result metaapi.NetworkMetricsResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPDiagnostics(
			func() (wireguardx.Statistics, error) {
				return wireguardx.Statistics{
					PeerKey:           "abc123",
					KeepaliveInterval: 25,
					TXBytes:           1024,
					RXBytes:           512,
					LastHandshakeSec:  time.Now().Unix() - 30,
				}, nil
			},
			metaapi.HTTPDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "Healthy", result.Wireguard.Status)
		require.Equal(t, "abc123", result.Wireguard.PeerKey)
		require.NotNil(t, result.Network)
	})

	t.Run("inactive", func(t *testing.T) {
		var (
			result metaapi.NetworkMetricsResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPDiagnostics(
			func() (wireguardx.Statistics, error) { return wireguardx.Statistics{}, nil },
			metaapi.HTTPDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "Inactive", result.Wireguard.Status)
	})

	t.Run("snapshot error", func(t *testing.T) {
		var (
			result metaapi.NetworkMetricsResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPDiagnostics(
			func() (wireguardx.Statistics, error) { return wireguardx.Statistics{}, errors.New("device closed") },
			metaapi.HTTPDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "Inactive", result.Wireguard.Status)
	})

	t.Run("no handshake", func(t *testing.T) {
		var (
			result metaapi.NetworkMetricsResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPDiagnostics(
			func() (wireguardx.Statistics, error) {
				return wireguardx.Statistics{PeerKey: "abc123", LastHandshakeSec: 0}, nil
			},
			metaapi.HTTPDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "No Handshake", result.Wireguard.Status)
	})

	t.Run("stale handshake", func(t *testing.T) {
		var (
			result metaapi.NetworkMetricsResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPDiagnostics(
			func() (wireguardx.Statistics, error) {
				return wireguardx.Statistics{
					PeerKey:          "abc123",
					LastHandshakeSec: time.Now().Unix() - 300,
					TXBytes:          100,
					RXBytes:          100,
				}, nil
			},
			metaapi.HTTPDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "Stale Handshake", result.Wireguard.Status)
	})

	t.Run("unbalanced pipe", func(t *testing.T) {
		var (
			result metaapi.NetworkMetricsResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPDiagnostics(
			func() (wireguardx.Statistics, error) {
				return wireguardx.Statistics{
					PeerKey:          "abc123",
					LastHandshakeSec: time.Now().Unix() - 30,
					TXBytes:          1024,
					RXBytes:          0,
				}, nil
			},
			metaapi.HTTPDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "Unbalanced Pipe", result.Wireguard.Status)
	})

	t.Run("unauthenticated", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPDiagnostics(
			func() (wireguardx.Statistics, error) { return wireguardx.Statistics{}, nil },
			metaapi.HTTPDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
