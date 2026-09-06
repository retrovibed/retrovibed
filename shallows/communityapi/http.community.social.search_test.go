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

func TestHTTPSocialSearch(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	var p meta.Profile
	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

	var v meta.Authz
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

	var owned, other community.Community
	require.NoError(t, community.CommunityInsertWithDefaults(ctx, q, community.Community{
		ID: uuid.Must(uuid.NewV7()).String(), AccountID: p.ID, Description: "owned by profile",
	}).Scan(&owned))
	require.NoError(t, community.CommunityInsertWithDefaults(ctx, q, community.Community{
		ID: uuid.Must(uuid.NewV7()).String(), AccountID: uuid.Must(uuid.NewV7()).String(), Description: "owned by someone else",
	}).Scan(&other))

	var publisher community.PluginPublisher
	require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
		ID: uuid.Must(uuid.NewV7()).String(), Path: "/plugins/youtube", Description: "YouTube", Mimetype: "application/vnd.retrovibe.publisher.youtube",
	}).Scan(&publisher))

	var enabled community.CommunityPublisher
	require.NoError(t, community.CommunityPublisherInsertWithDefaults(ctx, q, community.CommunityPublisher{
		ID: uuid.Must(uuid.NewV7()).String(), CommunityID: owned.ID, PublisherID: publisher.ID,
	}).Scan(&enabled))

	routes := mux.NewRouter()
	communityapi.NewHTTPSocial(
		q,
		communityapi.HTTPSocialOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

	t.Run("returns only the account's communities with their enabled publishers and the full catalog", func(t *testing.T) {
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

		require.Len(t, result.Items, 1)
		require.Equal(t, owned.ID, result.Items[0].Community.Id)
		require.Len(t, result.Items[0].Publishers, 1)
		require.Equal(t, publisher.ID, result.Items[0].Publishers[0].PublisherId)
	})

	t.Run("requires authentication", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/", nil)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
