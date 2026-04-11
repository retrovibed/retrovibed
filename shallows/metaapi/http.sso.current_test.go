package metaapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/httptestx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPSSOCurrent(t *testing.T) {
	t.Run("returns profile for authenticated user", func(t *testing.T) {
		var (
			p      meta.Profile
			result metaapi.Authn
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		routes := mux.NewRouter()
		metaapi.NewHTTP(
			q,
			metaapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes)

		claims := jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, p.ID, result.Profile.Id)
		require.NotEmpty(t, result.Token)
	})

	t.Run("rejects unauthenticated request", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		routes := mux.NewRouter()
		metaapi.NewHTTP(
			q,
			metaapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes)

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("not found when profile does not exist", func(t *testing.T) {
		var (
			p meta.Profile
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		// intentionally not inserted

		routes := mux.NewRouter()
		metaapi.NewHTTP(
			q,
			metaapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes)

		claims := jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/", nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Code)
	})
}
