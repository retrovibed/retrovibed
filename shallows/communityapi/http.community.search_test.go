package communityapi_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func newSearchMockClient(items ...*communityapi.Community) *http.Client {
	return httptestx.NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == http.MethodGet && req.URL.Path == "/c/" {
			body, _ := json.Marshal(&communityapi.CommunitySearchResponse{Items: items})
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

func TestSearchEndpoint(t *testing.T) {
	t.Run("upserts returned communities into local DB", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			p           meta.Profile
			v           meta.Authz
			communityID = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		communityapi.NewHTTP(
			q,
			communityapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPOptionHTTPClient(newSearchMockClient(&communityapi.Community{
				Id:          communityID,
				Domain:      "testcommunity",
				Description: "a test community",
			})),
			communityapi.HTTPOptionMediaStorage(fsx.DirVirtual(t.TempDir())),
			communityapi.HTTPOptionTorrentStorage(fsx.DirVirtual(t.TempDir())),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/c/",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var found community.Community
		require.NoError(t, community.CommunityFindByID(ctx, q, communityID).Scan(&found))
		require.Equal(t, communityID, found.ID)
		require.Equal(t, "testcommunity", found.Domain)
		require.Equal(t, "a test community", found.Description)
	})

	t.Run("response items reflect local DB state", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			p           meta.Profile
			v           meta.Authz
			communityID = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		communityapi.NewHTTP(
			q,
			communityapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPOptionHTTPClient(newSearchMockClient(&communityapi.Community{
				Id:     communityID,
				Domain: "testcommunity",
			})),
			communityapi.HTTPOptionMediaStorage(fsx.DirVirtual(t.TempDir())),
			communityapi.HTTPOptionTorrentStorage(fsx.DirVirtual(t.TempDir())),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/c/",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var searchResp communityapi.CommunitySearchResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&searchResp))
		require.Len(t, searchResp.Items, 1)
		require.Equal(t, communityID, searchResp.Items[0].Id)
		require.Equal(t, "testcommunity", searchResp.Items[0].Domain)
		require.NotEmpty(t, searchResp.Items[0].SubscribedAt)
	})
}
