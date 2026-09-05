package communityapi_test

import (
	"net/http"
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

func TestHTTPCommunityPublisherSearch(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	var (
		youtube community.PluginPublisher
		spotify community.PluginPublisher
	)
	require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
		ID: uuid.Must(uuid.NewV7()).String(), Path: "/plugins/youtube", Description: "YouTube", Mimetype: "application/vnd.retrovibe.publisher.youtube",
	}).Scan(&youtube))
	require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
		ID: uuid.Must(uuid.NewV7()).String(), Path: "/plugins/spotify", Description: "Spotify", Mimetype: "application/vnd.retrovibe.publisher.spotify",
	}).Scan(&spotify))

	routes := mux.NewRouter()
	communityapi.NewHTTPCommunityPublisher(
		q,
		communityapi.HTTPCommunityPublisherOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	var v meta.Authz
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

	t.Run("lists the catalog", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result communityapi.SocialsSearchResponse
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		bymimetype := map[string]*communityapi.PluginPublisher{}
		for _, p := range result.Catalog {
			bymimetype[p.Mimetype] = p
		}
		require.Contains(t, bymimetype, youtube.Mimetype)
		require.Equal(t, youtube.ID, bymimetype[youtube.Mimetype].Id)
		require.Contains(t, bymimetype, spotify.Mimetype)
	})

	t.Run("requires a privileged token", func(t *testing.T) {
		claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		unprivileged := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/", nil, httptestx.RequestOptionAuthorization(unprivileged))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("requires authentication", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/", nil)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
