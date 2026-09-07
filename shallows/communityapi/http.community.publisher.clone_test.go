package communityapi_test

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPCommunityPublisherClone(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	dir := t.TempDir()

	reg := testx.Must(publishplugin.NewRegistry(ctx, publishplugin.OptionConfigDir(t.TempDir()), publishplugin.OptionCacheDir(t.TempDir())))(t)

	routes := mux.NewRouter()
	communityapi.NewHTTPCommunityPublisher(
		q,
		reg,
		communityapi.HTTPCommunityPublisherOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		communityapi.HTTPCommunityPublisherOptionDir(dir),
	).Bind(routes.PathPrefix("/").Subrouter())

	// the clone is loaded into the registry like any other install, so the
	// module has to actually compile.
	buildPath := filepath.Join(t.TempDir(), "noop.wasm")
	build := exec.Command("go", "build", "-o", buildPath, "./.fixtures/noopplugin")
	build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))
	wasm := errorsx.Must(os.ReadFile(buildPath))

	var v meta.Authz
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
	token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)

	seed := func(t *testing.T, environment string) community.PluginPublisher {
		var row community.PluginPublisher
		path := filepath.Join(dir, uuid.Must(uuid.NewV7()).String()+".wasm")
		require.NoError(t, os.WriteFile(path, wasm, 0o600))
		if environment != "" {
			require.NoError(t, os.WriteFile(publishplugin.EnvPath(path), []byte(environment), 0o600))
		}

		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID:          errorsx.Must(publishplugin.Identity(path)),
			Path:        path,
			Description: "mastodon",
			Mimetype:    "application/vnd.retrovibe.publisher.test",
		}).Scan(&row))

		return row
	}

	t.Run("installs a second identity for the same module", func(t *testing.T) {
		source := seed(t, "MASTODON_TOKEN=original\n")

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodPost, "/"+source.ID+"/clone", nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result communityapi.PluginPublisherCloneResponse
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		require.NotEqual(t, source.ID, result.Publisher.Id)
		require.NotEqual(t, source.Path, result.Publisher.Path)
		require.Equal(t, source.Mimetype, result.Publisher.Mimetype)
		// unlabelled on purpose: the console falls back to the id and the
		// details card is where a clone gets named.
		require.Empty(t, result.Publisher.Description)

		// the id is the path's identity, which is what keeps the clone from
		// colliding with the module it points at.
		require.Equal(t, result.Publisher.Id, errorsx.Must(publishplugin.Identity(result.Publisher.Path)))

		var stored community.PluginPublisher
		require.NoError(t, community.PluginPublisherFindByID(ctx, q, result.Publisher.Id).Scan(&stored))
		require.Equal(t, result.Publisher.Path, stored.Path)

		// a link, not a copy: replacing the module replaces both identities.
		info, err := os.Lstat(stored.Path)
		require.NoError(t, err)
		require.NotEqual(t, 0, info.Mode()&os.ModeSymlink)
		require.Equal(t, wasm, errorsx.Must(os.ReadFile(stored.Path)))
	})

	t.Run("starts from the source's configuration and diverges from it", func(t *testing.T) {
		source := seed(t, "MASTODON_TOKEN=original\n")

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodPost, "/"+source.ID+"/clone", nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result communityapi.PluginPublisherCloneResponse
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		clonedenv := publishplugin.EnvPath(result.Publisher.Path)
		require.Equal(t, "MASTODON_TOKEN=original\n", string(errorsx.Must(os.ReadFile(clonedenv))))

		// EnvPath is derived from the path, so the clone's configuration is
		// its own from here on - the point of cloning at all.
		require.NoError(t, os.WriteFile(clonedenv, []byte("MASTODON_TOKEN=second\n"), 0o600))
		require.Equal(t, "MASTODON_TOKEN=original\n", string(errorsx.Must(os.ReadFile(publishplugin.EnvPath(source.Path)))))
	})

	t.Run("an unconfigured publisher clones without a sidecar", func(t *testing.T) {
		source := seed(t, "")

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodPost, "/"+source.ID+"/clone", nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result communityapi.PluginPublisherCloneResponse
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.NoFileExists(t, publishplugin.EnvPath(result.Publisher.Path))
	})

	t.Run("a clone of a clone points at the module itself", func(t *testing.T) {
		source := seed(t, "")

		clone := func(t *testing.T, id string) *communityapi.PluginPublisherCloneResponse {
			resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodPost, "/"+id+"/clone", nil, httptestx.RequestOptionAuthorization(token))
			require.NoError(t, err)

			routes.ServeHTTP(resp, req)
			require.NoError(t, httpx.ErrorCode(resp.Result()))

			var result communityapi.PluginPublisherCloneResponse
			require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
			return &result
		}

		first := clone(t, source.ID)
		second := clone(t, first.Publisher.Id)

		// no chain forms, so deleting the first clone - which removes only its
		// own link - cannot strand the second.
		require.Equal(t, errorsx.Must(filepath.EvalSymlinks(source.Path)), errorsx.Must(os.Readlink(second.Publisher.Path)))
	})

	t.Run("missing publisher returns 404", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodPost, "/"+uuid.Must(uuid.NewV7()).String()+"/clone", nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("requires a privileged token", func(t *testing.T) {
		source := seed(t, "")

		claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		unprivileged := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		before := errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM plugin_publishers"))

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodPost, "/"+source.ID+"/clone", nil, httptestx.RequestOptionAuthorization(unprivileged))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)

		require.Equal(t, before, errorsx.Zero(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM plugin_publishers")))
	})
}
