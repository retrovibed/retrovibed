package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestRecommendationsLatest(t *testing.T) {
	t.Run("zero state", func(t *testing.T) {
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

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Empty(t, result.Items)
	})

	t.Run("returns results", func(t *testing.T) {
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

		for idx := range 3 {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

			var rec library.Recommendation
			require.NoError(t, testx.Fake(&rec, library.RecommendationOptionTestDefaults, library.RecommendationOptionID(uuidx.WithSuffix(idx)), library.RecommendationOptionContentID(known.UID)))
			require.NoError(t, library.RecommendationInsertWithDefaults(ctx, q, rec).Scan(&rec))
		}

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 3)
	})

	t.Run("returns recommendation with future non-infinity tombstone", func(t *testing.T) {
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

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		var rec library.Recommendation
		require.NoError(t, testx.Fake(&rec, library.RecommendationOptionTestDefaults, library.RecommendationOptionContentID(known.UID), func(r *library.Recommendation) {
			r.TombstoneAt = time.Now().Add(30 * 24 * time.Hour)
		}))
		require.NoError(t, library.RecommendationInsertWithDefaults(ctx, q, rec).Scan(&rec))

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
	})

	t.Run("excludes tombstoned recommendation", func(t *testing.T) {
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

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		var rec library.Recommendation
		require.NoError(t, testx.Fake(&rec, library.RecommendationOptionTestDefaults, library.RecommendationOptionContentID(known.UID), func(r *library.Recommendation) {
			r.TombstoneAt = time.Now().Add(-time.Hour)
		}))
		require.NoError(t, library.RecommendationInsertWithDefaults(ctx, q, rec).Scan(&rec))

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Empty(t, result.Items)
	})

	t.Run("filters by mimetype", func(t *testing.T) {
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

		for idx := range 3 {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

			var rec library.Recommendation
			require.NoError(t, testx.Fake(&rec, library.RecommendationOptionTestDefaults, library.RecommendationOptionID(uuidx.WithSuffix(idx)), library.RecommendationOptionContentID(known.UID), func(r *library.Recommendation) {
				r.Mimetype = known.Mimetype
			}))
			require.NoError(t, library.RecommendationInsertWithDefaults(ctx, q, rec).Scan(&rec))
		}

		for idx := range 2 {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Audio)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

			var rec library.Recommendation
			require.NoError(t, testx.Fake(&rec, library.RecommendationOptionTestDefaults, library.RecommendationOptionID(uuidx.WithSuffix(idx+10)), library.RecommendationOptionContentID(known.UID), func(r *library.Recommendation) {
				r.Mimetype = known.Mimetype
			}))
			require.NoError(t, library.RecommendationInsertWithDefaults(ctx, q, rec).Scan(&rec))
		}

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(&media.RecommendationsSearchRequest{Mimetype: mimex.Video})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 3)
		for _, item := range result.Items {
			require.Equal(t, mimex.Video, item.Mimetype)
		}
	})

	t.Run("blank mimetype returns all", func(t *testing.T) {
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

		for idx := range 2 {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

			var rec library.Recommendation
			require.NoError(t, testx.Fake(&rec, library.RecommendationOptionTestDefaults, library.RecommendationOptionID(uuidx.WithSuffix(idx)), library.RecommendationOptionContentID(known.UID), func(r *library.Recommendation) {
				r.Mimetype = known.Mimetype
			}))
			require.NoError(t, library.RecommendationInsertWithDefaults(ctx, q, rec).Scan(&rec))
		}

		for idx := range 2 {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Audio)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

			var rec library.Recommendation
			require.NoError(t, testx.Fake(&rec, library.RecommendationOptionTestDefaults, library.RecommendationOptionID(uuidx.WithSuffix(idx+10)), library.RecommendationOptionContentID(known.UID), func(r *library.Recommendation) {
				r.Mimetype = known.Mimetype
			}))
			require.NoError(t, library.RecommendationInsertWithDefaults(ctx, q, rec).Scan(&rec))
		}

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 4)
	})
}
