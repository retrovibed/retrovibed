package communityapi_test

import (
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPSocialDisable(t *testing.T) {
	t.Run("disables a publisher for a community", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var (
			aid = uuidx.WithSuffix(1)
			p   meta.Profile
			com community.Community
			v   meta.Authz
		)
		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, community.CommunityInsertWithDefaults(ctx, q, community.Community{
			ID: uuid.Must(uuid.NewV7()).String(), AccountID: uuid.Nil.String(),
		}).Scan(&com))

		var publisher community.PluginPublisher
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: uuid.Must(uuid.NewV7()).String(), Path: "/plugins/x", Description: "X", Mimetype: "application/vnd.retrovibe.publisher.x",
		}).Scan(&publisher))

		var enabled community.CommunityPublisher
		require.NoError(t, community.CommunityPublisherInsertWithDefaults(ctx, q, community.CommunityPublisher{
			ID: uuid.Must(uuid.NewV7()).String(), CommunityID: com.ID, PublisherID: publisher.ID,
		}).Scan(&enabled))

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration(), jwtx.ClaimsOptionIssuer(aid)), metaapi.TokenOptionFromAuthz(v)))

		routes := mux.NewRouter()
		communityapi.NewHTTPSocial(
			q,
			communityapi.HTTPSocialOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodDelete,
			"/"+com.ID+"/publishers/"+publisher.ID,
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result communityapi.CommunityPublisherDisableResponse
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, com.ID, result.Compub.CommunityId)
		require.Equal(t, publisher.ID, result.Compub.PublisherId)

		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM community_publisher WHERE community_id = '"+com.ID+"' AND publisher_id = '"+publisher.ID+"'"))
	})

	t.Run("disabling an already-disabled publisher returns 404", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var (
			aid = uuidx.WithSuffix(1)
			p   meta.Profile
			com community.Community
			v   meta.Authz
		)
		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, community.CommunityInsertWithDefaults(ctx, q, community.Community{
			ID: uuid.Must(uuid.NewV7()).String(), AccountID: uuid.Nil.String(),
		}).Scan(&com))

		var publisher community.PluginPublisher
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: uuid.Must(uuid.NewV7()).String(), Path: "/plugins/x", Description: "X", Mimetype: "application/vnd.retrovibe.publisher.x",
		}).Scan(&publisher))

		routes := mux.NewRouter()
		communityapi.NewHTTPSocial(
			q,
			communityapi.HTTPSocialOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration(), jwtx.ClaimsOptionIssuer(aid)), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodDelete,
			"/"+com.ID+"/publishers/"+publisher.ID,
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("requires authentication", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var (
			p   meta.Profile
			com community.Community
			v   meta.Authz
		)
		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, community.CommunityInsertWithDefaults(ctx, q, community.Community{
			ID: uuid.Must(uuid.NewV7()).String(), AccountID: uuid.Nil.String(),
		}).Scan(&com))

		var publisher community.PluginPublisher
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: uuid.Must(uuid.NewV7()).String(), Path: "/plugins/x", Description: "X", Mimetype: "application/vnd.retrovibe.publisher.x",
		}).Scan(&publisher))

		var enabled community.CommunityPublisher
		require.NoError(t, community.CommunityPublisherInsertWithDefaults(ctx, q, community.CommunityPublisher{
			ID: uuid.Must(uuid.NewV7()).String(), CommunityID: com.ID, PublisherID: publisher.ID,
		}).Scan(&enabled))
		routes := mux.NewRouter()
		communityapi.NewHTTPSocial(
			q,
			communityapi.HTTPSocialOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+com.ID+"/publishers/"+publisher.ID, nil)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
