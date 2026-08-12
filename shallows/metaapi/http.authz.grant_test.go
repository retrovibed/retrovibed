package metaapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPAuthzGrant(t *testing.T) {
	t.Run("rejects requests from users who do not have usermanagement permissions", func(t *testing.T) {
		var (
			p1 meta.Profile
			v  meta.Authz
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p1, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p1).Scan(&p1))

		routes := mux.NewRouter()

		metaapi.NewHTTPAuthz(
			q,
			metaapi.HTTPAuthzOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p1.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, fmt.Sprintf("/%s", p1.ID), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode)
	})

	t.Run("updates permissions", func(t *testing.T) {
		var (
			p1     meta.Profile
			p2     meta.Profile
			v      meta.Authz
			result metaapi.AuthzGrantResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p1, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p1).Scan(&p1))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p1.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, testx.Fake(&p2, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p2).Scan(&p2))

		routes := mux.NewRouter()

		metaapi.NewHTTPAuthz(
			q,
			metaapi.HTTPAuthzOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p1.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		encoded, err := json.Marshal(&metaapi.AuthzGrantRequest{
			Token: metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p1.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)),
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, fmt.Sprintf("/%s", p2.ID), encoded, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.False(t, result.Token.Usermanagement)
		require.False(t, result.Token.RemoteControl)

		resp, req, err = httptestx.BuildRequestBytes(http.MethodPost, fmt.Sprintf("/%s", p2.ID), encoded, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.True(t, result.Token.Usermanagement)
		require.True(t, result.Token.RemoteControl)
	})

	t.Run("persists a second grant against an already granted profile", func(t *testing.T) {
		var (
			p1     meta.Profile
			p2     meta.Profile
			v      meta.Authz
			result metaapi.AuthzGrantResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p1, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p1).Scan(&p1))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p1.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, testx.Fake(&p2, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p2).Scan(&p2))

		routes := mux.NewRouter()

		metaapi.NewHTTPAuthz(
			q,
			metaapi.HTTPAuthzOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p1.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		authz := httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource))

		// first grant: takes the INSERT path since p2 has no existing authz row yet.
		var granted meta.Authz
		require.NoError(t, testx.Fake(&granted, meta.AuthzOptionProfileID(p2.ID), meta.AuthzOptionAdmin))
		encoded, err := json.Marshal(&metaapi.AuthzGrantRequest{
			Token: metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p2.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(granted)),
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, fmt.Sprintf("/%s", p2.ID), encoded, authz)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.True(t, result.Token.Usermanagement)
		require.True(t, result.Token.RemoteControl)

		// second grant against the same profile: takes the ON CONFLICT DO UPDATE
		// path. Every field granted above must actually flip, not silently keep
		// its prior value.
		var revoked meta.Authz
		require.NoError(t, testx.Fake(&revoked, meta.AuthzOptionProfileID(p2.ID), meta.AuthzOptionNoPrivileges))
		encoded, err = json.Marshal(&metaapi.AuthzGrantRequest{
			Token: metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p2.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(revoked)),
		})
		require.NoError(t, err)

		resp, req, err = httptestx.BuildRequestBytes(http.MethodPost, fmt.Sprintf("/%s", p2.ID), encoded, authz)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.False(t, result.Token.Usermanagement)
		require.False(t, result.Token.RemoteControl)
	})
}
