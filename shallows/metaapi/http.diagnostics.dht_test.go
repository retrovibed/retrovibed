package metaapi_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPDiagnosticsDHT(t *testing.T) {
	t.Run("healthy", func(t *testing.T) {
		var (
			result metaapi.DHTMetricsResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPDHTDiagnostics(
			func() (dht.ServerStats, error) {
				return dht.ServerStats{
					GoodNodes:                             12,
					Nodes:                                 34,
					OutstandingTransactions:               2,
					SuccessfulOutboundAnnouncePeerQueries: 56,
					BadNodes:                              3,
					OutboundQueriesAttempted:              78,
				}, nil
			},
			metaapi.HTTPDHTDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, int32(12), result.Dht.GoodNodes)
		require.Equal(t, int32(34), result.Dht.Nodes)
		require.Equal(t, int32(2), result.Dht.OutstandingTransactions)
		require.Equal(t, int64(56), result.Dht.SuccessfulOutboundAnnouncePeerQueries)
		require.Equal(t, uint32(3), result.Dht.BadNodes)
		require.Equal(t, int64(78), result.Dht.OutboundQueriesAttempted)
	})

	t.Run("snapshot error", func(t *testing.T) {
		var claims jwt.RegisteredClaims

		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPDHTDiagnostics(
			func() (dht.ServerStats, error) { return dht.ServerStats{}, errors.New("dht unavailable") },
			metaapi.HTTPDHTDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("unauthenticated", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPDHTDiagnostics(
			func() (dht.ServerStats, error) { return dht.ServerStats{}, nil },
			metaapi.HTTPDHTDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
