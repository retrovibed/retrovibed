package metaapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestHTTPDiagnosticsTorrent(t *testing.T) {
	t.Run("aggregates metadata totals", func(t *testing.T) {
		var (
			result   metaapi.TorrentMetricsResponse
			claims   jwt.RegisteredClaims
			seeding  tracking.Metadata
			leeching tracking.Metadata
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		seeding = tracking.NewMetadata(
			new(int160.Random()),
			tracking.MetadataOptionBytes(bytesx.MiB),
			tracking.MetadataOptionDownloaded(bytesx.MiB),
			tracking.MetadataOptionUploaded(bytesx.MiB),
			tracking.MetadataOptionAutoSeeding,
		)
		seeding.Peers = 4
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, seeding).Scan(&seeding))

		leeching = tracking.NewMetadata(
			new(int160.Random()),
			tracking.MetadataOptionBytes(2*bytesx.MiB),
			tracking.MetadataOptionDownloaded(bytesx.MiB),
		)
		leeching.Peers = 2
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, leeching).Scan(&leeching))

		routes := mux.NewRouter()
		metaapi.NewHTTPTorrentDiagnostics(
			q,
			metaapi.HTTPTorrentDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.EqualValues(t, 2, result.Torrent.Total)
		require.EqualValues(t, 1, result.Torrent.Seeding)
		require.EqualValues(t, 3*bytesx.MiB, result.Torrent.Bytes)
		require.EqualValues(t, 2*bytesx.MiB, result.Torrent.Downloaded)
		require.EqualValues(t, bytesx.MiB, result.Torrent.Uploaded)
		require.EqualValues(t, 6, result.Torrent.Peers)
	})

	t.Run("no metadata", func(t *testing.T) {
		var (
			result metaapi.TorrentMetricsResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		routes := mux.NewRouter()
		metaapi.NewHTTPTorrentDiagnostics(
			q,
			metaapi.HTTPTorrentDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.EqualValues(t, 0, result.Torrent.Total)
		require.EqualValues(t, 0, result.Torrent.Bytes)
	})

	t.Run("query error", func(t *testing.T) {
		var claims jwt.RegisteredClaims

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		require.NoError(t, q.Close())

		routes := mux.NewRouter()
		metaapi.NewHTTPTorrentDiagnostics(
			q,
			metaapi.HTTPTorrentDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
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
		metaapi.NewHTTPTorrentDiagnostics(
			q,
			metaapi.HTTPTorrentDiagnosticsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
