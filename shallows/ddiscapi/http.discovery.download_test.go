package ddiscapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestHTTPDiscoveryDownload(t *testing.T) {
	t.Run("downloads a discovered candidate", func(t *testing.T) {
		var (
			disc   ddisc.Discovered
			md     tracking.Metadata
			result ddiscapi.DiscoveryDownloadResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		disc = ddisc.NewDiscovered(&id)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, disc).Scan(&disc))

		md = tracking.NewMetadata(&id, tracking.MetadataOptionCompleted, tracking.MetadataOptionBytes(1), tracking.MetadataOptionDownloaded(1), tracking.MetadataOptionUploaded(1))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(disc.ID, jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, fmt.Sprintf("/%s", disc.ID), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.Equal(t, ddiscapi.NewDiscoveryFromDiscovered(disc), result.Discovery)
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_metadata WHERE initiated_at <= NOW()"))(t))
	})

	t.Run("unknown id returns not found", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims("missing", jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/missing", nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("missing tracking metadata returns internal server error", func(t *testing.T) {
		var (
			disc ddisc.Discovered
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		disc = ddisc.NewDiscovered(&id)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, disc).Scan(&disc))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(disc.ID, jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, fmt.Sprintf("/%s", disc.ID), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)
	})
}
