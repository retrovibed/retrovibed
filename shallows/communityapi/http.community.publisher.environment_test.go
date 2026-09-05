package communityapi_test

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
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

func TestHTTPPublisherEnvironment(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	dir := t.TempDir()

	declaration := "# base url of the lemmy instance to post to\nLEMMY_INSTANCE=\"\"\n# lemmy community to post into\nLEMMY_COMMUNITY=\"\"\n"

	var v meta.Authz
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
	token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)

	t.Run("declaration is served when nothing is configured yet", func(t *testing.T) {
		var row community.PluginPublisher
		id := uuid.Must(uuid.NewV7()).String()
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: id, Path: filepath.Join(dir, id+".wasm"), Description: "lemmy",
		}).Scan(&row))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublisherEnvironment(
			q,
			stubEnvironment{declared: []byte(declaration)},
			communityapi.HTTPPublisherEnvironmentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+id, nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.Equal(t, declaration, resp.Body.String())
	})

	t.Run("configured values are merged over the declaration, hints intact", func(t *testing.T) {
		var row community.PluginPublisher
		id := uuid.Must(uuid.NewV7()).String()
		path := filepath.Join(dir, id+".wasm")
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: id, Path: path, Description: "lemmy",
		}).Scan(&row))
		require.NoError(t, os.WriteFile(publishplugin.EnvPath(path), []byte("LEMMY_INSTANCE=https://lemmy.ml\nEXTRA=kept\n"), 0o600))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublisherEnvironment(
			q,
			stubEnvironment{declared: []byte(declaration)},
			communityapi.HTTPPublisherEnvironmentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+id, nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		body := resp.Body.String()
		require.Contains(t, body, "# base url of the lemmy instance to post to")
		require.Contains(t, body, "LEMMY_INSTANCE=https://lemmy.ml")
		require.Contains(t, body, "LEMMY_COMMUNITY=")
		require.Contains(t, body, "EXTRA=kept")
	})

	t.Run("an unreadable plugin still serves what was configured", func(t *testing.T) {
		var row community.PluginPublisher
		id := uuid.Must(uuid.NewV7()).String()
		path := filepath.Join(dir, id+".wasm")
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: id, Path: path, Description: "lemmy",
		}).Scan(&row))
		require.NoError(t, os.WriteFile(publishplugin.EnvPath(path), []byte("LEMMY_INSTANCE=https://lemmy.ml\n"), 0o600))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublisherEnvironment(
			q,
			stubEnvironment{err: errors.New("not loaded")},
			communityapi.HTTPPublisherEnvironmentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+id, nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.Contains(t, resp.Body.String(), "LEMMY_INSTANCE=https://lemmy.ml")
	})

	t.Run("update writes the sidecar and delete removes it", func(t *testing.T) {
		var row community.PluginPublisher
		id := uuid.Must(uuid.NewV7()).String()
		path := filepath.Join(dir, id+".wasm")
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: id, Path: path, Description: "lemmy",
		}).Scan(&row))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublisherEnvironment(
			q,
			stubEnvironment{declared: []byte(declaration)},
			communityapi.HTTPPublisherEnvironmentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		updated := []byte("LEMMY_INSTANCE=https://lemmy.world\nLEMMY_COMMUNITY=movies\n")
		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/"+id, updated, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		written, err := os.ReadFile(publishplugin.EnvPath(path))
		require.NoError(t, err)
		require.Equal(t, updated, written)

		resp, req, err = httptestx.BuildRequestBytes(http.MethodDelete, "/"+id, nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoFileExists(t, publishplugin.EnvPath(path))
	})

	t.Run("missing publisher returns 404", func(t *testing.T) {
		routes := mux.NewRouter()
		communityapi.NewHTTPPublisherEnvironment(
			q,
			stubEnvironment{declared: []byte(declaration)},
			communityapi.HTTPPublisherEnvironmentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+uuid.Must(uuid.NewV7()).String(), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("requires a privileged token", func(t *testing.T) {
		unprivileged := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

		routes := mux.NewRouter()
		communityapi.NewHTTPPublisherEnvironment(
			q,
			stubEnvironment{declared: []byte(declaration)},
			communityapi.HTTPPublisherEnvironmentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/"+uuid.Must(uuid.NewV7()).String(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&unprivileged, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
