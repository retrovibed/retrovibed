package ddiscapi_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPPluginManagementFind(t *testing.T) {
	configDir := t.TempDir()

	ctx, done := testx.Context(t)
	defer done()

	reg := testx.Must(searchplugin.NewRegistry(ctx, searchplugin.OptionConfigDir(configDir), searchplugin.OptionCacheDir(t.TempDir())))(t)

	require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "foo.wasm"), []byte("foocontent"), 0o600))

	routes := mux.NewRouter()
	ddiscapi.NewHTTPPluginManagement(
		reg,
		ddiscapi.HTTPPluginManagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		ddiscapi.HTTPPluginManagementOptionDir(searchplugin.SearchPluginDir(configDir)),
	).Bind(routes.PathPrefix("/").Subrouter())

	var v meta.Authz
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
	token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)

	t.Run("existing plugin", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+md5x.String("foo"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result ddiscapi.PluginFindResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "foo", result.Plugin.Name)
		require.Equal(t, md5x.String("foo"), result.Plugin.Id)
		require.EqualValues(t, len("foocontent"), result.Plugin.Size)
	})

	t.Run("missing plugin", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+md5x.String("missing"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("requires privileged token", func(t *testing.T) {
		unprivileged := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+md5x.String("foo"), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&unprivileged, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
