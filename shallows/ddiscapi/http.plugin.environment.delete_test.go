package ddiscapi_test

import (
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

func TestHTTPPluginEnvironmentDelete(t *testing.T) {
	configDir := t.TempDir()

	routes := mux.NewRouter()
	ddiscapi.NewHTTPPluginEnvironment(
		ddiscapi.HTTPPluginEnvironmentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		ddiscapi.HTTPPluginEnvironmentOptionDir(searchplugin.SearchPluginDir(configDir)),
	).Bind(routes.PathPrefix("/").Subrouter())

	var v meta.Authz
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
	token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)

	t.Run("successfully deletes and returns prior content", func(t *testing.T) {
		require.NoError(t, os.MkdirAll(searchplugin.SearchPluginDir(configDir), 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "foo.wasm"), []byte("foocontent"), 0o600))

		const content = "FOO=bar\n"
		path := filepath.Join(searchplugin.SearchPluginDir(configDir), "foo.env")
		require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+md5x.String("foo"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.Equal(t, content, resp.Body.String())
		require.NoFileExists(t, path)
	})

	t.Run("missing environment is a no-op success", func(t *testing.T) {
		require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "missing.wasm"), []byte("missingcontent"), 0o600))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+md5x.String("missing"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.Empty(t, resp.Body.String())
	})

	t.Run("unknown id not found", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+md5x.String("nonexistent"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("requires privileged token", func(t *testing.T) {
		require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "bar.wasm"), []byte("barcontent"), 0o600))

		const content = "FOO=bar\n"
		path := filepath.Join(searchplugin.SearchPluginDir(configDir), "bar.env")
		require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

		claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		unprivileged := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+md5x.String("bar"), nil, httptestx.RequestOptionAuthorization(unprivileged))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.FileExists(t, path)
	})
}
