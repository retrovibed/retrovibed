package metaapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPWireguardTouch(t *testing.T) {
	t.Run("touch to activate", func(t *testing.T) {
		var (
			claims jwt.RegisteredClaims
			wg     meta.Wireguard
			result metaapi.WireguardTouchResponse
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

		path1 := wg.ID
		d, err := os.Create(tmpdir.Path(path1))
		require.NoError(t, err)
		require.NoError(t, d.Close())

		routes := mux.NewRouter()

		metaapi.NewHTTPWireguard(
			tmpdir.Path(),
			q,
			metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodPut, fmt.Sprintf("/%s", path1), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_wireguard WHERE \"default\""))

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))

		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_wireguard WHERE \"default\""))

		resp, req, err = httptestx.BuildRequestContextBytes(ctx, http.MethodPut, fmt.Sprintf("/%s", path1), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_wireguard WHERE \"default\""))
		require.Equal(t, wg.ID, result.Wireguard.Id)
	})

	t.Run("touch to deactivate", func(t *testing.T) {
		var (
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
				meta.NewWireguard(
					testx.Must(uuid.NewV4())(t).String(),
					meta.WireguardOptionDescription("test"),
					meta.WireguardOptionDefault,
				),
			).Scan(&wg),
		)

		tmpdir := fsx.DirVirtual(t.TempDir())

		path1 := wg.ID
		d, err := os.Create(tmpdir.Path(path1))
		require.NoError(t, err)
		require.NoError(t, d.Close())

		routes := mux.NewRouter()

		metaapi.NewHTTPWireguard(
			tmpdir.Path(),
			q,
			metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_wireguard WHERE \"default\""))

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodPut, fmt.Sprintf("/%s", uuid.Max.String()), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.EqualValues(t, http.StatusNotFound, resp.Code)

		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_wireguard WHERE \"default\""))
	})
}
