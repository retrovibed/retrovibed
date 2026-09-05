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

func TestHTTPCommunityPublisherDelete(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	dir := t.TempDir()

	routes := mux.NewRouter()
	communityapi.NewHTTPCommunityPublisher(
		q,
		communityapi.HTTPCommunityPublisherOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		communityapi.HTTPCommunityPublisherOptionDir(dir),
	).Bind(routes.PathPrefix("/").Subrouter())

	var v meta.Authz
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
	token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)

	seed := func(t *testing.T, id string, content []byte) string {
		var row community.PluginPublisher
		path := filepath.Join(dir, id)
		require.NoError(t, os.WriteFile(path, content, 0o600))
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: id, Path: path, Description: "installed", Mimetype: "application/vnd.retrovibe.publisher.test",
		}).Scan(&row))
		return path
	}

	t.Run("successfully deletes", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7()).String()
		path := seed(t, id, []byte("removeme"))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+id, nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result communityapi.PluginPublisherDeleteResponse
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, id, result.Publisher.Id)

		require.NoFileExists(t, path)

		var row community.PluginPublisher
		require.Error(t, community.PluginPublisherFindByID(ctx, q, id).Scan(&row))
	})

	t.Run("missing publisher returns 404", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+uuid.Must(uuid.NewV7()).String(), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("requires a privileged token", func(t *testing.T) {
		claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		unprivileged := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		id := uuid.Must(uuid.NewV7()).String()
		path := seed(t, id, []byte("keepme"))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+id, nil, httptestx.RequestOptionAuthorization(unprivileged))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.FileExists(t, path)
	})
}
