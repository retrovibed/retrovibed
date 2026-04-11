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
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/meta/identityssh"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPUserManagementCreate(t *testing.T) {
	const testPubKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKQUl9ov/ly5SVfZAoqkzZsfACWTAxoZMCF9TXa22Rm8 test@test"

	t.Run("creates profile with permissions", func(t *testing.T) {
		var (
			admin  meta.Profile
			authz  meta.Authz
			result metaapi.ProfileCreateResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&admin, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, admin).Scan(&admin))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(admin.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/u12t").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(admin.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))
		body, err := json.Marshal(&metaapi.ProfileCreateRequest{Profile: &metaapi.Profile{Display: "Test User"}, PublicKey: testPubKey})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/u12t/", body, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.NotEmpty(t, result.Profile.Id)

		var created meta.Profile
		require.NoError(t, meta.ProfileFindByID(ctx, q, result.Profile.Id).Scan(&created))
		require.Equal(t, metaapi.ProfileStatus_ENABLED, metaapi.ProfileStatusOf(created))
		require.Equal(t, "Test User", created.Display)

		var iden identityssh.Identity
		require.NoError(t, identityssh.IdentityFindByProfileID(ctx, q, sqlx.NewNullString(created.ID)).Scan(&iden))
		require.Equal(t, created.ID, iden.ProfileID)
	})

	t.Run("rejects unauthorized", func(t *testing.T) {
		var (
			user  meta.Profile
			authz meta.Authz
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&user, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, user).Scan(&user))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(user.ID), meta.AuthzOptionNoPrivileges))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/u12t").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(user.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))
		body, err := json.Marshal(&metaapi.ProfileCreateRequest{Profile: &metaapi.Profile{Display: "Unauthorized User"}, PublicKey: testPubKey})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/u12t/", body, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode)
	})

	t.Run("rejects invalid public key", func(t *testing.T) {
		var (
			admin meta.Profile
			authz meta.Authz
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&admin, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, admin).Scan(&admin))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(admin.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/u12t").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(admin.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))
		body, err := json.Marshal(&metaapi.ProfileCreateRequest{Profile: &metaapi.Profile{Display: "Invalid Key User"}, PublicKey: "invalid-key"})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/u12t/", body, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})

	t.Run("uses ssh comment as display when display is empty", func(t *testing.T) {
		var (
			admin  meta.Profile
			authz  meta.Authz
			result metaapi.ProfileCreateResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&admin, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, admin).Scan(&admin))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(admin.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/u12t").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(admin.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))
		// Display is empty — should fall back to the SSH comment "test@test"
		body, err := json.Marshal(&metaapi.ProfileCreateRequest{Profile: &metaapi.Profile{}, PublicKey: testPubKey})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/u12t/", body, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var created meta.Profile
		require.NoError(t, meta.ProfileFindByID(ctx, q, result.Profile.Id).Scan(&created))
		require.Equal(t, "test@test", created.Display)
	})

	t.Run("idempotent - same key does not create duplicates", func(t *testing.T) {
		var (
			admin  meta.Profile
			authz  meta.Authz
			result metaapi.ProfileCreateResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&admin, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, admin).Scan(&admin))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(admin.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/u12t").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(admin.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))
		body, err := json.Marshal(&metaapi.ProfileCreateRequest{Profile: &metaapi.Profile{Display: "Test User"}, PublicKey: testPubKey})
		require.NoError(t, err)

		for range 2 {
			resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/u12t/", body, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
			require.NoError(t, err)
			routes.ServeHTTP(resp, req)
			require.NoError(t, httpx.ErrorCode(resp.Result()))
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		}

		require.Equal(t, 2, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
	})

	t.Run("explicit display takes precedence over ssh comment", func(t *testing.T) {
		var (
			admin  meta.Profile
			authz  meta.Authz
			result metaapi.ProfileCreateResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&admin, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, admin).Scan(&admin))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(admin.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/u12t").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(admin.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))
		// Display is set — SSH comment "test@test" should not override it
		body, err := json.Marshal(&metaapi.ProfileCreateRequest{Profile: &metaapi.Profile{Display: "Explicit Name"}, PublicKey: testPubKey})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/u12t/", body, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var created meta.Profile
		require.NoError(t, meta.ProfileFindByID(ctx, q, result.Profile.Id).Scan(&created))
		require.Equal(t, "Explicit Name", created.Display)
	})
}
