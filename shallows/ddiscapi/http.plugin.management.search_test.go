package ddiscapi_test

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
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

func TestHTTPPluginManagementSearch(t *testing.T) {
	configDir := t.TempDir()

	ctx, done := testx.Context(t)
	defer done()

	reg := testx.Must(searchplugin.NewRegistry(ctx, searchplugin.OptionConfigDir(configDir), searchplugin.OptionCacheDir(t.TempDir())))(t)

	require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "bar.wasm"), []byte("barcontent!"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "foo.wasm"), []byte("foocontent"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(searchplugin.SearchPluginDir(configDir), "foo.env"), []byte("IGNORED=true"), 0o600))

	routes := mux.NewRouter()
	ddiscapi.NewHTTPPluginManagement(
		reg,
		ddiscapi.HTTPPluginManagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		ddiscapi.HTTPPluginManagementOptionDir(searchplugin.SearchPluginDir(configDir)),
	).Bind(routes.PathPrefix("/").Subrouter())

	var v meta.Authz
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

	resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/", nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))

	var result ddiscapi.PluginSearchResponse
	require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

	require.Len(t, result.Items, 2, "the .env sidecar file must not be listed as a plugin")

	byname := map[string]*ddiscapi.Plugin{}
	for _, p := range result.Items {
		byname[p.Name] = p
	}

	require.Contains(t, byname, "bar")
	require.EqualValues(t, len("barcontent!"), byname["bar"].Size)
	require.Contains(t, byname, "foo")
	require.EqualValues(t, len("foocontent"), byname["foo"].Size)
}

func TestHTTPPluginManagementSearchRequiresPrivilegedToken(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	reg := testx.Must(searchplugin.NewRegistry(ctx, searchplugin.OptionConfigDir(t.TempDir()), searchplugin.OptionCacheDir(t.TempDir())))(t)

	routes := mux.NewRouter()
	ddiscapi.NewHTTPPluginManagement(
		reg,
		ddiscapi.HTTPPluginManagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

	resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/", nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHTTPPluginManagementSearchUnauthenticated(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	reg := testx.Must(searchplugin.NewRegistry(ctx, searchplugin.OptionConfigDir(t.TempDir()), searchplugin.OptionCacheDir(t.TempDir())))(t)

	routes := mux.NewRouter()
	ddiscapi.NewHTTPPluginManagement(
		reg,
		ddiscapi.HTTPPluginManagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/", nil)
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.Error(t, httpx.ErrorCode(resp.Result()))
	require.Equal(t, http.StatusUnauthorized, resp.Code)
}
