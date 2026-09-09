package communityapi_test

import (
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPSocialEnable(t *testing.T) {
	t.Run("enables a publisher for a community", func(t *testing.T) {
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

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration(), jwtx.ClaimsOptionIssuer(aid)), metaapi.TokenOptionFromAuthz(v)))

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

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/"+com.ID+"/publishers/"+publisher.ID,
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result communityapi.CommunityPublisherEnableResponse
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, com.ID, result.Compub.CommunityId)
		require.Equal(t, publisher.ID, result.Compub.PublisherId)

		iter := sqlx.Scan(community.CommunityPublisherFindByCommunityID(ctx, q, com.ID))
		var rows []community.CommunityPublisher
		for cp := range iter.Iter() {
			rows = append(rows, cp)
		}
		require.NoError(t, iter.Err())
		require.Len(t, rows, 1)
	})

	t.Run("enabling twice does not duplicate", func(t *testing.T) {
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

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration(), jwtx.ClaimsOptionIssuer(aid)), metaapi.TokenOptionFromAuthz(v)))

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

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/"+com.ID+"/publishers/"+publisher.ID,
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		resp, req, err = httptestx.BuildRequestBytes(
			http.MethodPost,
			"/"+com.ID+"/publishers/"+publisher.ID,
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM community_publisher WHERE community_id = '"+com.ID+"' AND publisher_id = '"+publisher.ID+"'"))
	})

	t.Run("requires authentication", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var com community.Community
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

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/"+com.ID+"/publishers/"+publisher.ID, nil)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
