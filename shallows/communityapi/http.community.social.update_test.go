package communityapi_test

import (
	"fmt"
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
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPSocialUpdate(t *testing.T) {
	t.Run("updates the community publisher", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var (
			p   meta.Profile
			aid = uuidx.WithSuffix(1)
		)
		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		var v meta.Authz
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		var owned community.Community
		require.NoError(t, community.CommunityInsertWithDefaults(ctx, q, community.Community{
			ID: uuid.Must(uuid.NewV7()).String(), AccountID: aid, Description: "owned",
		}).Scan(&owned))

		var publisher community.PluginPublisher
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: uuid.Must(uuid.NewV7()).String(), Path: "/plugins/youtube", Description: "YouTube", Mimetype: "application/vnd.retrovibe.publisher.youtube",
		}).Scan(&publisher))

		var compub community.CommunityPublisher
		require.NoError(t, community.CommunityPublisherInsertWithDefaults(ctx, q, community.CommunityPublisher{
			ID: uuid.Must(uuid.NewV7()).String(), CommunityID: owned.ID, PublisherID: publisher.ID,
		}).Scan(&compub))

		routes := mux.NewRouter()
		communityapi.NewHTTPSocial(
			q,
			communityapi.HTTPSocialOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration(), jwtx.ClaimsOptionIssuer(aid)), metaapi.TokenOptionFromAuthz(v)))

		encoded, err := jsonx.Marshal(&communityapi.CommunityPublisherUpdateRequest{
			Compub: communityapi.NewCommunityPublisher(
				communityapi.CommunityPublisherOptionFromDB(langx.Clone(compub, timex.JSONSafeEncodeOption)),
				communityapi.CommunityPublisherOptionTemplates("Derp 0", "Derp 1"),
			),
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			fmt.Sprintf("/%s", compub.ID),
			encoded,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result communityapi.CommunityPublisherUpdateResponse
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		assert.EqualValues(t, "Derp 0", result.Compub.TemplateTitle)
		assert.EqualValues(t, "Derp 1", result.Compub.TemplateDescription)
	})

	t.Run("requires authentication", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var (
			p   meta.Profile
			aid = uuidx.WithSuffix(1)
		)
		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		var v meta.Authz
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		var owned community.Community
		require.NoError(t, community.CommunityInsertWithDefaults(ctx, q, community.Community{
			ID: uuid.Must(uuid.NewV7()).String(), AccountID: aid, Description: "owned",
		}).Scan(&owned))

		var publisher community.PluginPublisher
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: uuid.Must(uuid.NewV7()).String(), Path: "/plugins/youtube", Description: "YouTube", Mimetype: "application/vnd.retrovibe.publisher.youtube",
		}).Scan(&publisher))

		var compub community.CommunityPublisher
		require.NoError(t, community.CommunityPublisherInsertWithDefaults(ctx, q, community.CommunityPublisher{
			ID: uuid.Must(uuid.NewV7()).String(), CommunityID: owned.ID, PublisherID: publisher.ID,
		}).Scan(&compub))

		routes := mux.NewRouter()
		communityapi.NewHTTPSocial(
			q,
			communityapi.HTTPSocialOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		encoded, err := jsonx.Marshal(&communityapi.CommunityPublisherUpdateRequest{
			Compub: communityapi.NewCommunityPublisher(
				communityapi.CommunityPublisherOptionFromDB(langx.Clone(compub, timex.JSONSafeEncodeOption)),
				communityapi.CommunityPublisherOptionTemplates("Derp 0", "Derp 1"),
			),
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			fmt.Sprintf("/%s", compub.ID),
			encoded,
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
