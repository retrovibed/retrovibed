package metaapi_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPWireguardSearch(t *testing.T) {
	t.Run("search", func(t *testing.T) {
		var (
			result metaapi.WireguardSearchResponse
			claims jwt.RegisteredClaims
			wg     meta.Wireguard
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(
			t,
			meta.WireguardInsertWithDefaults(
				ctx,
				q,
				meta.NewWireguard(testx.Must(uuid.NewV4())(t).String(), meta.WireguardOptionDescription("test")),
			).Scan(&wg),
		)

		tmpdir := fsx.DirVirtual(t.TempDir())
		path := wg.ID
		d, err := os.Create(tmpdir.Path(path))
		require.NoError(t, err)
		require.NoError(t, d.Close())

		routes := mux.NewRouter()

		metaapi.NewHTTPWireguard(
			tmpdir.Path(),
			q,
			metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(formx.NewEncoder().Encode(&metaapi.WireguardSearchRequest{
			Offset: 0,
		}))(t)

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.Equal(t, len(result.Items), 1)
		require.Equal(t, result.Items[0].Id, wg.ID)
		require.False(t, result.Items[0].Default)
		require.Equal(t, result.Items[0].Description, wg.Description)
	})

	t.Run("zero state", func(t *testing.T) {
		var (
			result metaapi.WireguardSearchResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		tmpdir := fsx.DirVirtual(t.TempDir())

		routes := mux.NewRouter()

		metaapi.NewHTTPWireguard(
			tmpdir.Path("does.not.exist"),
			q,
			metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(formx.NewEncoder().Encode(&metaapi.WireguardSearchRequest{
			Offset: 0,
		}))(t)

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.Equal(t, len(result.Items), 0)
	})
}
