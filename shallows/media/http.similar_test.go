package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/acoustics"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestSimilar(t *testing.T) {
	t.Run("bad media id", func(t *testing.T) {
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

		router := mux.NewRouter()
		media.NewHTTPSimilar(q, media.HTTPSimilarOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(router.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/not-a-uuid",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})

	t.Run("no seed vector indexed", func(t *testing.T) {
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

		router := mux.NewRouter()
		media.NewHTTPSimilar(q, media.HTTPSimilarOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(router.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		seedID := uuid.Must(uuid.NewV4())
		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/%s", seedID),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("match found", func(t *testing.T) {
		var (
			p         meta.Profile
			authz     meta.Authz
			seed      library.Metadata
			candidate library.Metadata
			vec       acoustics.FeatureVector
			result    media.MediaFindResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		for i := range vec {
			vec[i] = float32(i%7) + 1
		}

		require.NoError(t, testx.Fake(&seed, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, seed).Scan(&seed))
		require.NoError(t, acoustics.StoreFeatures(ctx, q, seed.ID, vec, acoustics.StatsVersion))

		require.NoError(t, testx.Fake(&candidate, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, candidate).Scan(&candidate))
		require.NoError(t, acoustics.StoreFeatures(ctx, q, candidate.ID, vec, acoustics.StatsVersion))

		router := mux.NewRouter()
		media.NewHTTPSimilar(q, media.HTTPSimilarOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(router.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/%s", seed.ID),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.NotNil(t, result.Media)
		require.Equal(t, candidate.ID, result.Media.Id)
	})

	t.Run("candidate has no surviving metadata", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
			seed  library.Metadata
			vec   acoustics.FeatureVector
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		for i := range vec {
			vec[i] = float32(i%7) + 1
		}

		require.NoError(t, testx.Fake(&seed, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, seed).Scan(&seed))
		require.NoError(t, acoustics.StoreFeatures(ctx, q, seed.ID, vec, acoustics.StatsVersion))

		// audio_features row with no corresponding (surviving) library_metadata row.
		orphan := uuid.Must(uuid.NewV4())
		require.NoError(t, acoustics.StoreFeatures(ctx, q, orphan.String(), vec, acoustics.StatsVersion))

		router := mux.NewRouter()
		media.NewHTTPSimilar(q, media.HTTPSimilarOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(router.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/%s", seed.ID),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("exclusion is honored", func(t *testing.T) {
		var (
			p        meta.Profile
			authz    meta.Authz
			seed     library.Metadata
			excluded library.Metadata
			eligible library.Metadata
			vec      acoustics.FeatureVector
			result   media.MediaFindResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		for i := range vec {
			vec[i] = float32(i%7) + 1
		}

		require.NoError(t, testx.Fake(&seed, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, seed).Scan(&seed))
		require.NoError(t, acoustics.StoreFeatures(ctx, q, seed.ID, vec, acoustics.StatsVersion))

		require.NoError(t, testx.Fake(&excluded, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, excluded).Scan(&excluded))
		require.NoError(t, acoustics.StoreFeatures(ctx, q, excluded.ID, vec, acoustics.StatsVersion))

		require.NoError(t, testx.Fake(&eligible, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, eligible).Scan(&eligible))
		require.NoError(t, acoustics.StoreFeatures(ctx, q, eligible.ID, vec, acoustics.StatsVersion))

		router := mux.NewRouter()
		media.NewHTTPSimilar(q, media.HTTPSimilarOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(router.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{Excluded: []string{excluded.ID}})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/%s?%s", seed.ID, query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.NotNil(t, result.Media)
		require.Equal(t, eligible.ID, result.Media.Id)
	})
}
