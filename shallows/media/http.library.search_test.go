package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestLibrarySearch(t *testing.T) {
	t.Run("basic search", func(t *testing.T) {
		var (
			p      meta.Profile
			authz  meta.Authz
			result media.MediaSearchResponse
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		// Insert test metadata
		var md1, md2, md3 library.Metadata
		require.NoError(t, testx.Fake(&md1, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("ipsum lorem Matrix"), library.MetadataOptionAutoDescription("ipsum lorem Matrix dolor"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md1).Scan(&md1))

		require.NoError(t, testx.Fake(&md2, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("ipsum lorem sit"), library.MetadataOptionAutoDescription("ipsum lorem sit amet"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md2).Scan(&md2))

		require.NoError(t, testx.Fake(&md3, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("ipsum lorem amet"), library.MetadataOptionAutoDescription("ipsum lorem dolor"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md3).Scan(&md3))

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Limit: 10, Offset: 0})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+query.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, 3, len(result.Items))
		require.NotNil(t, result.Next)
		require.Equal(t, uint64(0), result.Next.Offset)
	})

	t.Run("empty database", func(t *testing.T) {
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Limit: 10, Offset: 0})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+query.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, 0, len(result.Items))
		require.NotNil(t, result.Next)
		require.Equal(t, uint64(0), result.Next.Offset)
	})

	t.Run("limit capped at 100", func(t *testing.T) {
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		// Insert 200 items
		for i := range 200 {
			var md library.Metadata
			require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("Title "+string(rune('A'+(i%26)))), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))
		}

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Limit: 1000, Offset: 0})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+query.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, 100, len(result.Items))
		require.NotNil(t, result.Next)
		require.Equal(t, uint64(0), result.Next.Offset)
	})

	t.Run("MimetypeFilter", func(t *testing.T) {
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		var md1, md2, md3 library.Metadata
		require.NoError(t, testx.Fake(&md1, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("Archive 1"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md1).Scan(&md1))

		require.NoError(t, testx.Fake(&md2, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("Binary 1"), library.MetadataOptionMimetype(mimex.Binary)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md2).Scan(&md2))

		require.NoError(t, testx.Fake(&md3, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("Archive 2"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md3).Scan(&md3))

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Mimetypes: []string{mimex.RetrovibedMediaArchive}, Limit: 10})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+query.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, 2, len(result.Items))
		for _, item := range result.Items {
			require.Equal(t, mimex.RetrovibedMediaArchive, item.Mimetype)
		}
	})

	t.Run("filter hidden", func(t *testing.T) {
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		var visible, hidden library.Metadata
		require.NoError(t, testx.Fake(&visible, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("Visible"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive), library.MetadataOptionHidden(false)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, visible).Scan(&visible))

		require.NoError(t, testx.Fake(&hidden, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("Hidden"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive), library.MetadataOptionHidden(true)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, hidden).Scan(&hidden))

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Hidden: false, Limit: 10})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+query.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, 1, len(result.Items))
		require.Equal(t, visible.ID, result.Items[0].Id)
	})

	t.Run("TombstonedExcluded", func(t *testing.T) {
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		var active, tombstoned library.Metadata
		require.NoError(t, testx.Fake(&active, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("Active"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, active).Scan(&active))

		require.NoError(t, testx.Fake(&tombstoned, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("Tombstoned"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, tombstoned).Scan(&tombstoned))

		// Tombstone one item
		require.NoError(t, library.MetadataTombstoneByID(ctx, q, tombstoned.ID).Scan(&tombstoned))

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Limit: 10})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+query.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, 1, len(result.Items))
		require.Equal(t, active.ID, result.Items[0].Id)
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("Test"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?limit=invalid",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
		require.NotNil(t, resp.Body)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("SearchQueryFilter", func(t *testing.T) {
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		var md1, md2, md3 library.Metadata
		require.NoError(t, testx.Fake(&md1, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("ipsum lorem Matrix Reloaded"), library.MetadataOptionAutoDescription("ipsum lorem Matrix"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md1).Scan(&md1))

		require.NoError(t, testx.Fake(&md2, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("ipsum lorem Trek"), library.MetadataOptionAutoDescription("ipsum lorem Star"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md2).Scan(&md2))

		require.NoError(t, testx.Fake(&md3, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("ipsum lorem Matrix Revolutions"), library.MetadataOptionAutoDescription("ipsum lorem Matrix"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md3).Scan(&md3))

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Query: "Matrix", Limit: 10})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+query.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		// Should only return items matching "Matrix"
		require.Equal(t, 2, len(result.Items))
		for _, item := range result.Items {
			require.Contains(t, item.Description, "Matrix")
		}
	})

	t.Run("ordering with query", func(t *testing.T) {
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		for i := range 4 {
			var md library.Metadata
			desc := fmt.Sprintf("Zeta Item %d", i)
			require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription(desc), library.MetadataOptionAutoDescription(desc)))
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))
		}

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Query: "Zeta", Limit: 10})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+query.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		// With query, should order by description ASC
		require.Equal(t, 4, len(result.Items))
		for idx, i := range result.Items {
			require.Contains(t, i.Description, fmt.Sprintf("Zeta Item %d", idx))
		}
	})

	t.Run("OrderingWithoutQuery", func(t *testing.T) {
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Limit: 10})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+query.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		// Without query, should order by created_at DESC, description ASC
		require.NotNil(t, result.Next)
	})

	t.Run("ImageContentHasImageURL", func(t *testing.T) {
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		var imgMD, videoMD library.Metadata
		require.NoError(t, testx.Fake(&imgMD, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("cover.jpg"), library.MetadataOptionMimetype("image/jpeg")))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, imgMD).Scan(&imgMD))

		require.NoError(t, testx.Fake(&videoMD, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("movie.mp4"), library.MetadataOptionMimetype("video/mp4")))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, videoMD).Scan(&videoMD))

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/m").Subrouter())

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Limit: 10})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/m/?"+query.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, 2, len(result.Items))

		byID := make(map[string]*media.Media, len(result.Items))
		for _, item := range result.Items {
			byID[item.Id] = item
		}

		require.Equal(t, fmt.Sprintf("http:///m/%s", imgMD.ID), byID[imgMD.ID].Image)
		require.Empty(t, byID[videoMD.ID].Image)
	})

	t.Run("SearchQueryWithMultipleFiles", func(t *testing.T) {
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
			asyncx.NewWakeup(t.Context()),
			nil,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		)

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("derp.1.mp4"), library.MetadataOptionAutoDescription("derp.1.mp4"), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		for i := range 9 {
			var md library.Metadata
			desc := fmt.Sprintf("example.%d.mp4", i)
			require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription(desc), library.MetadataOptionAutoDescription(desc), library.MetadataOptionMimetype(mimex.RetrovibedMediaArchive)))
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))
		}

		router := mux.NewRouter()
		h.Bind(router.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Query: "example", Limit: 10})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+query.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		var result media.MediaSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.Equal(t, 9, len(result.Items))
		for _, item := range result.Items {
			require.Contains(t, item.Description, "example")
		}
	})

	t.Run("the flat library excludes folders", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testlibrary(t, ctx)

		top := testfolder(t, ctx, q, uuid.Nil.String(), "photos")
		held := testfile(t, ctx, q, top.ID, "held.bin")
		loose := testfile(t, ctx, q, uuid.Nil.String(), "loose.bin")

		// a folder is organization rather than media, and the grid this serves has nothing
		// to do with one.
		result := testsearch(t, routes, token, media.MediaSearchRequest{Limit: 10})
		require.ElementsMatch(t, []string{held.ID, loose.ID}, testresultids(result))
		require.Empty(t, result.Breadcrumb)
	})

	t.Run("the root listing includes folders", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testlibrary(t, ctx)

		top := testfolder(t, ctx, q, uuid.Nil.String(), "photos")
		testfile(t, ctx, q, top.ID, "held.bin")
		loose := testfile(t, ctx, q, uuid.Nil.String(), "loose.bin")

		result := testsearch(t, routes, token, media.MediaSearchRequest{Limit: 10, ParentId: uuid.Nil.String()})
		require.ElementsMatch(t, []string{top.ID, loose.ID}, testresultids(result))
		require.Empty(t, result.Breadcrumb)
	})

	t.Run("a folder listing returns only its contents", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testlibrary(t, ctx)

		top := testfolder(t, ctx, q, uuid.Nil.String(), "photos")
		nested := testfolder(t, ctx, q, top.ID, "2026")
		held := testfile(t, ctx, q, top.ID, "held.bin")
		testfile(t, ctx, q, uuid.Nil.String(), "loose.bin")

		result := testsearch(t, routes, token, media.MediaSearchRequest{Limit: 10, ParentId: top.ID})
		require.ElementsMatch(t, []string{nested.ID, held.ID}, testresultids(result))
	})

	t.Run("a folder listing sorts folders ahead of files", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testlibrary(t, ctx)

		top := testfolder(t, ctx, q, uuid.Nil.String(), "photos")
		zzz := testfolder(t, ctx, q, top.ID, "zzz")
		aaa := testfile(t, ctx, q, top.ID, "aaa.bin")

		result := testsearch(t, routes, token, media.MediaSearchRequest{Limit: 10, ParentId: top.ID})
		require.Equal(t, []string{zzz.ID, aaa.ID}, testresultids(result))
	})

	t.Run("a folder listing carries its breadcrumb root first", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testlibrary(t, ctx)

		top := testfolder(t, ctx, q, uuid.Nil.String(), "photos")
		nested := testfolder(t, ctx, q, top.ID, "2026")
		held := testfile(t, ctx, q, nested.ID, "held.bin")

		result := testsearch(t, routes, token, media.MediaSearchRequest{Limit: 10, ParentId: nested.ID})
		require.Equal(t, []string{held.ID}, testresultids(result))

		require.Equal(t, []string{top.ID, nested.ID}, slicesx.MapTransform(func(m *media.Media) string { return m.Id }, result.Breadcrumb...))
	})
}

func testsearch(t *testing.T, routes *mux.Router, token string, req media.MediaSearchRequest) (result media.MediaSearchResponse) {
	query, err := formx.NewEncoder().Encode(req)
	require.NoError(t, err)

	resp, hreq, err := httptestx.BuildRequestBytes(
		http.MethodGet,
		"/?"+query.Encode(),
		nil,
		httptestx.RequestOptionAuthorization(token),
	)
	require.NoError(t, err)

	routes.ServeHTTP(resp, hreq)
	require.Equal(t, http.StatusOK, resp.Result().StatusCode)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	return result
}

func testresultids(result media.MediaSearchResponse) []string {
	return slicesx.MapTransform(func(m *media.Media) string { return m.Id }, result.Items...)
}
