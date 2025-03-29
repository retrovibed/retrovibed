package metaapi_test

import (
	"encoding/json"
	"fmt"
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

func TestHTTPAuthz(t *testing.T) {
	// ensure that the authz endpoint returns the correct permissions based on the provided bearer token
	var (
		p      meta.Profile
		v      meta.Authz
		result metaapi.AuthzResponse
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	defer q.Close()

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

	routes := mux.NewRouter()

	metaapi.NewHTTPAuthz(
		q,
		metaapi.HTTPAuthzOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
	resp, req, err := httptestx.BuildRequest(http.MethodGet, "/", nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.True(t, result.Token.Usermanagement)
}

func TestHTTPProfileUnauthorized(t *testing.T) {
	// ensure that the profile specific endpoint rejects requests from users who do not have usermanagement permissions
	var (
		p1 meta.Profile
		v  meta.Authz
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	defer q.Close()

	require.NoError(t, testx.Fake(&p1, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p1).Scan(&p1))

	routes := mux.NewRouter()

	metaapi.NewHTTPAuthz(
		q,
		metaapi.HTTPAuthzOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p1.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
	resp, req, err := httptestx.BuildRequest(http.MethodGet, fmt.Sprintf("/%s", p1.ID), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode)
}

func TestHTTPProfile(t *testing.T) {
	// ensure that the profile specific endpoint returns the correct permissions based on the provided profile id
	var (
		p1     meta.Profile
		p2     meta.Profile
		v      meta.Authz
		result metaapi.AuthzResponse
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	defer q.Close()

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
	resp, req, err := httptestx.BuildRequest(http.MethodGet, fmt.Sprintf("/%s", p1.ID), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.True(t, result.Token.Usermanagement)

	resp, req, err = httptestx.BuildRequest(http.MethodGet, fmt.Sprintf("/%s", p2.ID), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.False(t, result.Token.Usermanagement)
}

func TestHTTPGrantUnauthorized(t *testing.T) {
	// ensure that the grant endpoint rejects requests from users who do not have usermanagement permissions
	var (
		p1 meta.Profile
		v  meta.Authz
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	defer q.Close()

	require.NoError(t, testx.Fake(&p1, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p1).Scan(&p1))

	routes := mux.NewRouter()

	metaapi.NewHTTPAuthz(
		q,
		metaapi.HTTPAuthzOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p1.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

	resp, req, err := httptestx.BuildRequest(http.MethodPost, fmt.Sprintf("/%s", p1.ID), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode)
}

func TestHTTPGrant(t *testing.T) {
	// ensure that the grant endpoint updates permissions
	var (
		p1     meta.Profile
		p2     meta.Profile
		v      meta.Authz
		result metaapi.AuthzGrantResponse
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	defer q.Close()

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

	resp, req, err := httptestx.BuildRequest(http.MethodGet, fmt.Sprintf("/%s", p2.ID), encoded, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.False(t, result.Token.Usermanagement)

	resp, req, err = httptestx.BuildRequest(http.MethodPost, fmt.Sprintf("/%s", p2.ID), encoded, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.True(t, result.Token.Usermanagement)
}

func TestHTTPRevokeUnauthorized(t *testing.T) {
	// ensure that the revoke endpoint rejects requests from users who do not have usermanagement permissions
	var (
		p1 meta.Profile
		v  meta.Authz
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	defer q.Close()

	require.NoError(t, testx.Fake(&p1, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p1).Scan(&p1))

	routes := mux.NewRouter()

	metaapi.NewHTTPAuthz(
		q,
		metaapi.HTTPAuthzOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p1.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

	resp, req, err := httptestx.BuildRequest(http.MethodDelete, fmt.Sprintf("/%s", p1.ID), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode)
}
