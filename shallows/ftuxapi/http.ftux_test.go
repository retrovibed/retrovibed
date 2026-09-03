package ftuxapi_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/ftuxapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestDefaultsEndpoint(t *testing.T) {
	t.Run("returns the curated default communities", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())

		var (
			p meta.Profile
			v meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		ftuxapi.NewHTTP(
			q,
			ftuxapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ftux").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/ftux/communities",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var out ftuxapi.CommunitySuggestions
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &out))
		require.Len(t, out.Community, 2)
	})
}

func TestSubscribeEndpoint(t *testing.T) {
	t.Run("subscribes to resolvable communities and skips ones that fail to resolve", func(t *testing.T) {
		var (
			p meta.Profile
			v meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		goodID := uuid.Must(uuid.NewV7()).String()
		badID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		client := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			if req.Method == http.MethodGet && strings.Contains(req.URL.Path, goodID) {
				body, _ := json.Marshal(&communityapi.CommunityFindResponse{
					Community: &communityapi.Community{
						Id:          goodID,
						Description: goodID,
						Entropy:     uuidx.WithSuffix(1),
						Url:         "https://community.community.retrovibe.space",
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

		routes := mux.NewRouter()
		ftuxapi.NewHTTP(
			q,
			ftuxapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			ftuxapi.HTTPOptionHTTPClient(client),
		).Bind(routes.PathPrefix("/ftux").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		body, err := json.Marshal(&ftuxapi.SubscribeCommunitiesRequest{CommunityId: []string{goodID, badID}})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/ftux/communities",
			body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var sub community.Community
		require.NoError(t, community.CommunityFindByID(ctx, q, goodID).Scan(&sub))
		require.Equal(t, goodID, sub.ID)
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM community"))
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))
	})

	t.Run("responds service unavailable when no deeppool client is configured", func(t *testing.T) {
		var (
			p meta.Profile
			v meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		ftuxapi.NewHTTP(
			q,
			ftuxapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/ftux").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		body, err := json.Marshal(&ftuxapi.SubscribeCommunitiesRequest{CommunityId: []string{uuid.Must(uuid.NewV7()).String()}})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/ftux/communities",
			body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusServiceUnavailable, resp.Code)
	})
}
