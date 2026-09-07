package metaapi_test

import (
	"errors"
	"net/http"
	"net/netip"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestHTTPDiagnosticsDiscovery(t *testing.T) {
	t.Run("aggregates peer and unknown hash totals", func(t *testing.T) {
		var (
			result   metaapi.DiscoveryMetricsResponse
			claims   jwt.RegisteredClaims
			ddiscp   tracking.Peer
			bep51p   tracking.Peer
			expiredp tracking.Peer
			unknown  tracking.UnknownHash
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		ddiscp = tracking.NewPeer(int160.Random(), tracking.PeerOptionDdisc(uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())))
		require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, ddiscp).Scan(&ddiscp))

		bep51p = tracking.NewPeer(int160.Random(), tracking.PeerOptionBEP51(5, 60))
		require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, bep51p).Scan(&bep51p))

		expiredp = tracking.NewPeer(int160.Random(), tracking.PeerOptionTombstone(time.Now().Add(-time.Hour)))
		require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, expiredp).Scan(&expiredp))

		unknown = tracking.NewUnknownHash(int160.Random(), tracking.OptionUnknownHashPeer(int160.Random(), netip.MustParseAddrPort("1.2.3.4:6881")))
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, unknown).Scan(&unknown))

		routes := mux.NewRouter()
		metaapi.NewHTTPDiscoveryDiagnostics(
			q,
			func() (ddisc.Snapshot, error) {
				return ddisc.Snapshot{
					Enabled:        true,
					Ratio:          50,
					Partitions:     4,
					Workloads:      2,
					LocalPartition: "partition-xyz",
				}, nil
			},
			metaapi.HTTPDiscoveryDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.True(t, result.Discovery.Enabled)
		require.EqualValues(t, 50, result.Discovery.Ratio)
		require.EqualValues(t, 4, result.Discovery.Partitions)
		require.EqualValues(t, 2, result.Discovery.Workloads)
		require.Equal(t, "partition-xyz", result.Discovery.LocalPartition)
		require.EqualValues(t, 2, result.Discovery.Peers)
		require.EqualValues(t, 1, result.Discovery.PeersDdisc)
		require.EqualValues(t, 1, result.Discovery.PeersBep51)
		require.EqualValues(t, 1, result.Discovery.Unidentified)
	})

	t.Run("snapshot error", func(t *testing.T) {
		var claims jwt.RegisteredClaims

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		routes := mux.NewRouter()
		metaapi.NewHTTPDiscoveryDiagnostics(
			q,
			func() (ddisc.Snapshot, error) { return ddisc.Snapshot{}, errors.New("not yet initialized") },
			metaapi.HTTPDiscoveryDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
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

		q := sqltestx.Metadatabase(t)

		routes := mux.NewRouter()
		metaapi.NewHTTPDiscoveryDiagnostics(
			q,
			func() (ddisc.Snapshot, error) { return ddisc.Snapshot{}, nil },
			metaapi.HTTPDiscoveryDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
