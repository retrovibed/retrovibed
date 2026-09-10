package metaapi_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPWireguardCurrent(t *testing.T) {
	t.Run("return current", func(t *testing.T) {
		var (
			result metaapi.WireguardCurrentResponse
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
					meta.WireguardOptionDistribution,
				),
			).Scan(&wg),
		)

		tmpdir := fsx.DirVirtual(t.TempDir())
		path := wg.ID
		require.NoError(t, os.WriteFile(tmpdir.Path(path), testx.IOBytes(testx.Read(testx.Fixture("wireguard", "example.1.conf"))), 0600))

		routes := mux.NewRouter()

		metaapi.NewHTTPWireguard(
			tmpdir.Path(),
			q,
			metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/current?nettype=1", nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		require.Equal(t, path, result.Wireguard.Id)
	})

	t.Run("zerostate", func(t *testing.T) {
		var (
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		tmpdir := fsx.DirVirtual(t.TempDir())

		routes := mux.NewRouter()

		metaapi.NewHTTPWireguard(
			tmpdir.Path(),
			q,
			metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/current?nettype=1", nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Error(t, httpx.ErrorCode(resp.Result()))
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("noncurrent nettype", func(t *testing.T) {
		var (
			claims jwt.RegisteredClaims
			wg     meta.Wireguard
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		wg = meta.NewWireguard(
			testx.Must(uuid.NewV4())(t).String(),
			meta.WireguardOptionDescription("test"),
			meta.WireguardOptionDistribution,
		)

		require.NoError(
			t,
			meta.WireguardInsertWithDefaults(ctx, q, wg).Scan(&wg),
		)

		tmpdir := fsx.DirVirtual(t.TempDir())

		routes := mux.NewRouter()

		metaapi.NewHTTPWireguard(
			tmpdir.Path(),
			q,
			metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/current", nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Error(t, httpx.ErrorCode(resp.Result()))
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("missing configuration", func(t *testing.T) {
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
					meta.WireguardOptionDistribution,
				),
			).Scan(&wg),
		)

		tmpdir := fsx.DirVirtual(t.TempDir())

		routes := mux.NewRouter()

		metaapi.NewHTTPWireguard(
			tmpdir.Path(),
			q,
			metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/current?nettype=1", nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Error(t, httpx.ErrorCode(resp.Result()))
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid configuration", func(t *testing.T) {
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
					meta.WireguardOptionDistribution,
				),
			).Scan(&wg),
		)

		tmpdir := fsx.DirVirtual(t.TempDir())
		require.NoError(t, os.WriteFile(tmpdir.Path(wg.ID), []byte("garbage\n"), 0600))

		routes := mux.NewRouter()

		metaapi.NewHTTPWireguard(
			tmpdir.Path(),
			q,
			metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/current?nettype=1", nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Error(t, httpx.ErrorCode(resp.Result()))
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		tmpdir := fsx.DirVirtual(t.TempDir())

		routes := mux.NewRouter()

		metaapi.NewHTTPWireguard(
			tmpdir.Path(),
			q,
			metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/current?nettype=1", nil)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
