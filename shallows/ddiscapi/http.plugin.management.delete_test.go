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

func TestHTTPPluginManagementDelete(t *testing.T) {
	configDir := t.TempDir()

	ctx, done := testx.Context(t)
	defer done()

	reg := testx.Must(searchplugin.NewRegistry(ctx, searchplugin.OptionConfigDir(configDir), searchplugin.OptionCacheDir(t.TempDir())))(t)

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

	t.Run("successfully delete", func(t *testing.T) {
		path := filepath.Join(searchplugin.SearchPluginDir(configDir), "foo.wasm")
		require.NoError(t, os.WriteFile(path, []byte("foocontent"), 0o600))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+md5x.String("foo"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result ddiscapi.PluginDeleteResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "foo", result.Plugin.Name)
		require.EqualValues(t, len("foocontent"), result.Plugin.Size)
		require.NoFileExists(t, path)
	})

	t.Run("missing plugin", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+md5x.String("missing"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("requires privileged token", func(t *testing.T) {
		path := filepath.Join(searchplugin.SearchPluginDir(configDir), "bar.wasm")
		require.NoError(t, os.WriteFile(path, []byte("barcontent"), 0o600))

		claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		unprivileged := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+md5x.String("bar"), nil, httptestx.RequestOptionAuthorization(unprivileged))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.FileExists(t, path)
	})
}
