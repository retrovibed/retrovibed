package communityapi_test

import (
	"fmt"
	"net/http"
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
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPCommunityPublisherFind(t *testing.T) {
	t.Run("locate a publisher by id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var (
			youtube community.PluginPublisher
			spotify community.PluginPublisher
			v       meta.Authz
			result  communityapi.PluginPublisherFindResponse
		)
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: uuid.Must(uuid.NewV7()).String(), Path: "/plugins/youtube", Description: "YouTube", Mimetype: "application/vnd.retrovibe.publisher.youtube",
		}).Scan(&youtube))
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: uuid.Must(uuid.NewV7()).String(), Path: "/plugins/spotify", Description: "Spotify", Mimetype: "application/vnd.retrovibe.publisher.spotify",
		}).Scan(&spotify))

		reg := testx.Must(publishplugin.NewRegistry(ctx, publishplugin.OptionConfigDir(t.TempDir()), publishplugin.OptionCacheDir(t.TempDir())))(t)

		routes := mux.NewRouter()
		communityapi.NewHTTPCommunityPublisher(
			q,
			reg,
			communityapi.HTTPCommunityPublisherOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/%s", spotify.ID),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		assert.Equal(t, result.Publisher.Id, spotify.ID)
		assert.Equal(t, result.Publisher.Description, spotify.Description)
		assert.Equal(t, result.Publisher.Mimetype, spotify.Mimetype)
		assert.Equal(t, result.Publisher.Path, spotify.Path)
	})

	t.Run("locate a publisher by id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var (
			v meta.Authz
		)

		reg := testx.Must(publishplugin.NewRegistry(ctx, publishplugin.OptionConfigDir(t.TempDir()), publishplugin.OptionCacheDir(t.TempDir())))(t)

		routes := mux.NewRouter()
		communityapi.NewHTTPCommunityPublisher(
			q,
			reg,
			communityapi.HTTPCommunityPublisherOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/%s", uuid.Must(uuid.NewV7()).String()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.EqualValues(t, http.StatusNotFound, resp.Result().StatusCode)
	})
}
