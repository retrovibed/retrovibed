package ddiscapi_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

// stubEnvironment stands in for a loaded plugin, returning whatever a test
// wants that plugin to have declared.
type stubEnvironment struct {
	declared []byte
	err      error
}

func (t stubEnvironment) Environment(ctx context.Context, path string) ([]byte, error) {
	return t.declared, t.err
}

func TestHTTPPluginEnvironmentGet(t *testing.T) {
	configDir := t.TempDir()

	routes := mux.NewRouter()
	ddiscapi.NewHTTPPluginEnvironment(
		searchplugin.Unimplemented{},
		ddiscapi.HTTPPluginEnvironmentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		ddiscapi.HTTPPluginEnvironmentOptionDir(searchplugin.SearchPluginDir(configDir)),
	).Bind(routes.PathPrefix("/").Subrouter())

	var v meta.Authz
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
	token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)

	t.Run("without a declaration serves the configured values", func(t *testing.T) {
		const content = "FOO=\"bar\" # derp 0\n# derp 1\nBAR=\"baz\"\nBIZ=\"BAN\"\n# derp 2\n"

		require.NoError(t, os.MkdirAll(searchplugin.SearchPluginDir(configDir), 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "foo.wasm"), []byte("foocontent"), 0o600))
		require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "foo.env"), []byte(content), 0o600))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+md5x.String("foo"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.Equal(t, "FOO=bar\nBAR=baz\nBIZ=BAN\n", resp.Body.String())
	})

	t.Run("declaration is served when nothing is configured yet", func(t *testing.T) {
		const declaration = "# api key for requests\nUNIT3D_APIKEY=\"\"\n# base url for the unit3d api\nUNIT3D_DOMAIN=\"\"\n"

		require.NoError(t, os.MkdirAll(searchplugin.SearchPluginDir(configDir), 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "declared.wasm"), []byte("declaredcontent"), 0o600))

		declared := mux.NewRouter()
		ddiscapi.NewHTTPPluginEnvironment(
			stubEnvironment{declared: []byte(declaration)},
			ddiscapi.HTTPPluginEnvironmentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			ddiscapi.HTTPPluginEnvironmentOptionDir(searchplugin.SearchPluginDir(configDir)),
		).Bind(declared.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+md5x.String("declared"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		declared.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.Equal(t, declaration, resp.Body.String())
	})

	t.Run("configured values are merged over the declaration, hints intact", func(t *testing.T) {
		const declaration = "# api key for requests\nUNIT3D_APIKEY=\"\"\n# base url for the unit3d api\nUNIT3D_DOMAIN=\"\"\n"

		require.NoError(t, os.MkdirAll(searchplugin.SearchPluginDir(configDir), 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "merged.wasm"), []byte("mergedcontent"), 0o600))
		require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "merged.env"), []byte("UNIT3D_APIKEY=secret\nEXTRA=kept\n"), 0o600))

		merged := mux.NewRouter()
		ddiscapi.NewHTTPPluginEnvironment(
			stubEnvironment{declared: []byte(declaration)},
			ddiscapi.HTTPPluginEnvironmentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			ddiscapi.HTTPPluginEnvironmentOptionDir(searchplugin.SearchPluginDir(configDir)),
		).Bind(merged.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+md5x.String("merged"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		merged.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.Equal(t, "# api key for requests\nUNIT3D_APIKEY=secret\n# base url for the unit3d api\nUNIT3D_DOMAIN=\"\"\nEXTRA=kept\n", resp.Body.String())
	})

	t.Run("missing environment returns empty body", func(t *testing.T) {
		require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "missing.wasm"), []byte("missingcontent"), 0o600))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+md5x.String("missing"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.Empty(t, resp.Body.String())
	})

	t.Run("unknown id not found", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+md5x.String("nonexistent"), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("unauthenticated rejected", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+md5x.String("foo"), nil)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("requires privileged token", func(t *testing.T) {
		unprivileged := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+md5x.String("foo"), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&unprivileged, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
