package community_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func newCommunityMockClient(communityID string) *http.Client {
	return httptestx.NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == http.MethodGet && strings.Contains(req.URL.Path, communityID) {
			body, _ := json.Marshal(&meta.CommunityFindResponse{
				Community: &meta.Community{
					Id:          communityID,
					Domain:      "community",
					Description: communityID,
					Entropy:     uuidx.WithSuffix(1),
				},
			})
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(string(body))),
				Header:     make(http.Header),
			}
		}
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader(`{}`)),
			Header:     make(http.Header),
		}
	})
}

func TestSubscribeEndpoint(t *testing.T) {
	t.Run("subscribes to a community", func(t *testing.T) {
		var (
			p   meta.Profile
			v   meta.Authz
			sub community.CommunitySubscription
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		communityID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		community.NewHTTP(
			q,
			community.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			community.HTTPOptionHTTPClient(newCommunityMockClient(communityID)),
			community.HTTPOptionMediaStorage(fsx.DirVirtual(t.TempDir())),
			community.HTTPOptionTorrentStorage(fsx.DirVirtual(t.TempDir())),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/c/"+communityID+"/subscribe",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		require.NoError(t, community.CommunitySubscriptionFindByCommunityID(ctx, q, communityID).Scan(&sub))
		require.Equal(t, communityID, sub.CommunityID)
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))

		feedURL := "https://community.community.retrovibe.space"
		var feed tracking.RSS
		require.NoError(t, tracking.RSSFindByURL(ctx, q, feedURL).Scan(&feed))
		require.Equal(t, fmt.Sprintf("community - %s", communityID), feed.Description)
		require.Equal(t, uuidx.WithSuffix(1), feed.EncryptionSeed)
		require.WithinDuration(t, time.Now(), feed.NextCheck, time.Second)
		require.True(t, feed.Autodownload)
	})

	t.Run("toggle unsubscribes", func(t *testing.T) {
		var (
			p   meta.Profile
			v   meta.Authz
			sub community.CommunitySubscription
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		communityID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		community.NewHTTP(
			q,
			community.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			community.HTTPOptionHTTPClient(newCommunityMockClient(communityID)),
			community.HTTPOptionMediaStorage(fsx.DirVirtual(t.TempDir())),
			community.HTTPOptionTorrentStorage(fsx.DirVirtual(t.TempDir())),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)
		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))
		// subscribe
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/c/"+communityID+"/subscribe",
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.NoError(t, community.CommunitySubscriptionFindByCommunityID(ctx, q, communityID).Scan(&sub))
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))

		// unsubscribe
		resp, req, err = httptestx.BuildRequestBytes(
			http.MethodPost,
			"/c/"+communityID+"/subscribe",
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.ErrorIs(t, community.CommunitySubscriptionFindByCommunityID(ctx, q, communityID).Scan(&sub), sql.ErrNoRows)
		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))
	})

	t.Run("resubscribe after unsubscribe", func(t *testing.T) {
		var (
			p   meta.Profile
			v   meta.Authz
			sub community.CommunitySubscription
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		communityID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		community.NewHTTP(
			q,
			community.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			community.HTTPOptionHTTPClient(newCommunityMockClient(communityID)),
			community.HTTPOptionMediaStorage(fsx.DirVirtual(t.TempDir())),
			community.HTTPOptionTorrentStorage(fsx.DirVirtual(t.TempDir())),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)

		// subscribe
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/c/"+communityID+"/subscribe",
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.NoError(t, community.CommunitySubscriptionFindByCommunityID(ctx, q, communityID).Scan(&sub))
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))

		// unsubscribe
		resp, req, err = httptestx.BuildRequestBytes(
			http.MethodPost,
			"/c/"+communityID+"/subscribe",
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.ErrorIs(t, community.CommunitySubscriptionFindByCommunityID(ctx, q, communityID).Scan(&sub), sql.ErrNoRows)
		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))

		// resubscribe
		resp, req, err = httptestx.BuildRequestBytes(
			http.MethodPost,
			"/c/"+communityID+"/subscribe",
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.NoError(t, community.CommunitySubscriptionFindByCommunityID(ctx, q, communityID).Scan(&sub))
		require.Equal(t, communityID, sub.CommunityID)
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))
	})
}
