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
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/stretchr/testify/require"
)

func TestHTTPSocialEnable(t *testing.T) {
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

	token := httpauthtest.UnsafeTokenAuto(t)

	enable := func(t *testing.T) *http.Response {
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/"+com.ID+"/publishers/"+publisher.ID,
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		return resp.Result()
	}

	t.Run("enables a publisher for a community", func(t *testing.T) {
		resp := enable(t)
		require.NoError(t, httpx.ErrorCode(resp))

		var result communityapi.CommunityPublisherEnableResponse
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, com.ID, result.Enabled.CommunityId)
		require.Equal(t, publisher.ID, result.Enabled.PublisherId)

		iter := sqlx.Scan(community.CommunityPublisherFindByCommunityID(ctx, q, com.ID))
		var rows []community.CommunityPublisher
		for cp := range iter.Iter() {
			rows = append(rows, cp)
		}
		require.NoError(t, iter.Err())
		require.Len(t, rows, 1)
	})

	t.Run("enabling twice does not duplicate", func(t *testing.T) {
		require.NoError(t, httpx.ErrorCode(enable(t)))
		require.NoError(t, httpx.ErrorCode(enable(t)))

		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM community_publisher WHERE community_id = '"+com.ID+"' AND publisher_id = '"+publisher.ID+"'"))
	})

	t.Run("requires authentication", func(t *testing.T) {
		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/"+com.ID+"/publishers/"+publisher.ID, nil)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
