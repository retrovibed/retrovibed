package communityapi_test

import (
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestHTTPSocialDisable(t *testing.T) {
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

	var enabled community.CommunityPublisher
	require.NoError(t, community.CommunityPublisherInsertWithDefaults(ctx, q, community.CommunityPublisher{
		ID: uuid.Must(uuid.NewV7()).String(), CommunityID: com.ID, PublisherID: publisher.ID,
	}).Scan(&enabled))

	routes := mux.NewRouter()
	communityapi.NewHTTPSocial(
		q,
		communityapi.HTTPSocialOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	token := httpauthtest.UnsafeTokenAuto(t)

	t.Run("disables a publisher for a community", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodDelete,
			"/"+com.ID+"/publishers/"+publisher.ID,
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result communityapi.CommunityPublisherDisableResponse
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, com.ID, result.Disabled.CommunityId)
		require.Equal(t, publisher.ID, result.Disabled.PublisherId)

		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM community_publisher WHERE community_id = '"+com.ID+"' AND publisher_id = '"+publisher.ID+"'"))
	})

	t.Run("disabling an already-disabled publisher returns 404", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodDelete,
			"/"+com.ID+"/publishers/"+publisher.ID,
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("requires authentication", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, "/"+com.ID+"/publishers/"+publisher.ID, nil)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
