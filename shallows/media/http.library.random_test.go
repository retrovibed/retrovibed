package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestRandom(t *testing.T) {
	t.Run("BasicRetrieval", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		h := media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestID(uuid.Must(uuid.NewV4()).String())))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/random",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaFindResponse

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.NotNil(t, result.Media)
		require.Equal(t, md.ID, result.Media.Id)
	})

	t.Run("EmptyDatabase", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		h := media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/random",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
		require.NotNil(t, resp.Body)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("LimitCappedToOne", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		h := media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		var md1, md2 library.Metadata
		require.NoError(t, testx.Fake(&md1, library.MetadataOptionTestDefaults))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md1).Scan(&md1))
		require.NoError(t, testx.Fake(&md2, library.MetadataOptionTestDefaults))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md2).Scan(&md2))

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Limit: 5})
		require.NoError(t, err)
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/random?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaFindResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.NotNil(t, result.Media)
	})

	t.Run("InvalidRequest", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		h := media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/random?limit=invalid_number",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
		require.NotNil(t, resp.Body)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("RespectMimetypes", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		h := media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		var md1, md2 library.Metadata
		require.NoError(t, testx.Fake(&md1, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md1).Scan(&md1))

		require.NoError(t, testx.Fake(&md2, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionMimetype(mimex.Binary)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md2).Scan(&md2))

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Mimetypes: []string{mimex.Binary}})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/random?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaFindResponse

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.NotNil(t, result.Media)
		require.Equal(t, md2.ID, result.Media.Id)
		require.Equal(t, md2.Description, result.Media.Description)
	})
}
