package communityapi_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
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

func newResyncMockClient(communityID string) *http.Client {
	return httptestx.NewTestClient(func(req *http.Request) *http.Response {
		switch {
		case req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "/c/") && strings.Contains(req.URL.Path, communityID):
			body, _ := json.Marshal(&communityapi.CommunityFindResponse{
				Community: &communityapi.Community{
					Id:          communityID,
					Description: communityID,
					Entropy:     uuidx.WithSuffix(1),
					Url:         "https://resynced.community.retrovibe.space",
				},
			})
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		case req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "/p/"):
			body, _ := json.Marshal(&communityapi.PublishedContentSearchResponse{})
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		default:
			return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(strings.NewReader(`{}`)), Header: make(http.Header)}
		}
	})
}

func TestResyncEndpoint(t *testing.T) {
	t.Run("refreshes community metadata and returns local published content", func(t *testing.T) {
		var (
			p   meta.Profile
			v   meta.Authz
			sub community.Community
			pc  community.PublishedContent
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		communityID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		before := time.Now().Add(-time.Hour)
		sub = community.Community{ID: communityID, AccountID: uuid.Nil.String(), LastSyncAt: before}
		require.NoError(t, community.CommunityInsertWithDefaults(ctx, q, sub).Scan(&sub))

		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		routes := mux.NewRouter()
		communityapi.NewHTTP(
			q,
			communityapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPOptionHTTPClient(newResyncMockClient(communityID)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/c/"+communityID+"/resync",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishedContentSearchResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.Equal(t, communityID, result.Community.Id)
		require.Equal(t, "https://resynced.community.retrovibe.space", result.Community.Url)
		require.Len(t, result.Items, 1)
		require.Equal(t, pc.ID, result.Items[0].Id)

		var updated community.Community
		require.NoError(t, community.CommunityFindByID(ctx, q, communityID).Scan(&updated))
		require.Equal(t, "https://resynced.community.retrovibe.space", updated.URL)
		require.True(t, updated.LastSyncAt.After(before))
	})

	t.Run("returns 503 when the deeppool http client is not configured", func(t *testing.T) {
		var (
			p meta.Profile
			v meta.Authz
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
		communityapi.NewHTTP(
			q,
			communityapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/c/"+communityID+"/resync",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusServiceUnavailable, resp.Code)
	})
}
