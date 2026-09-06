package communityapi_test

import (
	"net/http"
	"os"
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
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPCommunityPublisherUpdate(t *testing.T) {
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

	var v meta.Authz
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
	token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)

	seed := func(t *testing.T) community.PluginPublisher {
		var row community.PluginPublisher
		id := uuid.Must(uuid.NewV7()).String()
		path := filepath.Join(dir, id+".wasm")
		require.NoError(t, os.WriteFile(path, []byte("installed"), 0o600))
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: id, Path: path, Description: "installed", Mimetype: "application/vnd.retrovibe.publisher.test",
		}).Scan(&row))
		return row
	}

	t.Run("records the submitted description and mimetype", func(t *testing.T) {
		row := seed(t)

		body := errorsx.Must(jsonx.Marshal(&communityapi.PluginPublisherUpdateRequest{
			Publisher: &communityapi.PluginPublisher{
				Id:          row.ID,
				Description: "mastodon - alternate account",
				Mimetype:    "application/vnd.retrovibe.publisher.mastodon",
			},
		}))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/"+row.ID, body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result communityapi.PluginPublisherUpdateResponse
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, "mastodon - alternate account", result.Publisher.Description)
		require.Equal(t, "application/vnd.retrovibe.publisher.mastodon", result.Publisher.Mimetype)

		var stored community.PluginPublisher
		require.NoError(t, community.PluginPublisherFindByID(ctx, q, row.ID).Scan(&stored))
		require.Equal(t, "mastodon - alternate account", stored.Description)
		require.Equal(t, "application/vnd.retrovibe.publisher.mastodon", stored.Mimetype)
	})

	t.Run("ignores an id and path submitted by the client", func(t *testing.T) {
		row := seed(t)

		// id and path describe what is installed on disk; the catalog follows
		// the filesystem, so neither is the client's to move.
		body := errorsx.Must(jsonx.Marshal(&communityapi.PluginPublisherUpdateRequest{
			Publisher: &communityapi.PluginPublisher{
				Id:          uuid.Must(uuid.NewV7()).String(),
				Path:        filepath.Join(dir, "elsewhere.wasm"),
				Description: "renamed",
				Mimetype:    row.Mimetype,
			},
		}))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/"+row.ID, body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var stored community.PluginPublisher
		require.NoError(t, community.PluginPublisherFindByID(ctx, q, row.ID).Scan(&stored))
		require.Equal(t, row.Path, stored.Path)
		require.Equal(t, "renamed", stored.Description)
	})

	t.Run("a blank name is what unnames a publisher", func(t *testing.T) {
		row := seed(t)

		// the console falls back to the id for a blank description, which is
		// how a clone arrives - so clearing the name has to be expressible.
		body := errorsx.Must(jsonx.Marshal(&communityapi.PluginPublisherUpdateRequest{
			Publisher: &communityapi.PluginPublisher{Id: row.ID, Mimetype: row.Mimetype},
		}))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/"+row.ID, body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var stored community.PluginPublisher
		require.NoError(t, community.PluginPublisherFindByID(ctx, q, row.ID).Scan(&stored))
		require.Empty(t, stored.Description)
	})

	t.Run("missing publisher returns 404", func(t *testing.T) {
		body := errorsx.Must(jsonx.Marshal(&communityapi.PluginPublisherUpdateRequest{
			Publisher: &communityapi.PluginPublisher{Description: "renamed"},
		}))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/"+uuid.Must(uuid.NewV7()).String(), body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("a request without a publisher returns 400", func(t *testing.T) {
		row := seed(t)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/"+row.ID, []byte(`{}`), httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("requires a privileged token", func(t *testing.T) {
		row := seed(t)

		claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		unprivileged := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		body := errorsx.Must(jsonx.Marshal(&communityapi.PluginPublisherUpdateRequest{
			Publisher: &communityapi.PluginPublisher{Id: row.ID, Description: "renamed"},
		}))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/"+row.ID, body, httptestx.RequestOptionAuthorization(unprivileged))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)

		var stored community.PluginPublisher
		require.NoError(t, community.PluginPublisherFindByID(ctx, q, row.ID).Scan(&stored))
		require.Equal(t, "installed", stored.Description)
	})
}
