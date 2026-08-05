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
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPPluginEnvironmentUpdate(t *testing.T) {
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

	t.Run("writes raw content, comments included", func(t *testing.T) {
		const content = "FOO=\"bar\" # derp 0\n# derp 1\nBAR=\"baz\"\nBIZ=\"BAN\"\n# derp 2\n"

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/foo", []byte(content), httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.Equal(t, content, resp.Body.String())

		written, err := os.ReadFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "foo.env"))
		require.NoError(t, err)
		require.Equal(t, content, string(written))
	})

	t.Run("overwrites existing content", func(t *testing.T) {
		path := filepath.Join(searchplugin.SearchPluginDir(configDir), "bar.env")
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
		require.NoError(t, os.WriteFile(path, []byte("OLD=value\n"), 0o600))

		const content = "NEW=value\n"
		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/bar", []byte(content), httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))

		written, err := os.ReadFile(path)
		require.NoError(t, err)
		require.Equal(t, content, string(written))
	})

	t.Run("requires privileged token", func(t *testing.T) {
		claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		unprivileged := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/baz", []byte("FOO=bar\n"), httptestx.RequestOptionAuthorization(unprivileged))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.NoFileExists(t, filepath.Join(searchplugin.SearchPluginDir(configDir), "baz.env"))
	})
}
