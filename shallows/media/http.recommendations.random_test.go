package media_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestRecommendationsRandomFiltering(t *testing.T) {
	t.Run("excludes adult content", func(t *testing.T) {
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
		known.Adult = true
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/random", func() []byte { b, _ := json.Marshal(&media.RecommendationsRequest{}); return b }(),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("excludes known without poster and backdrop path", func(t *testing.T) {
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
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionTestNoPoster))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/random", func() []byte { b, _ := json.Marshal(&media.RecommendationsRequest{}); return b }(),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("returns known media with only a poster path", func(t *testing.T) {
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
		known.BackdropPath = ""
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/random", func() []byte { b, _ := json.Marshal(&media.RecommendationsRequest{}); return b }(),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
	})

	t.Run("returns known media with only a backdrop path", func(t *testing.T) {
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
		known.PosterPath = ""
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/random", func() []byte { b, _ := json.Marshal(&media.RecommendationsRequest{}); return b }(),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
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

		var video, audio library.Known
		require.NoError(t, testx.Fake(&video, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, video).Scan(&video))
		require.NoError(t, testx.Fake(&audio, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Audio)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, audio).Scan(&audio))

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		body, err := json.Marshal(&media.RecommendationsRequest{Mimetype: mimex.Video})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/random", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
		require.Equal(t, mimex.Video, result.Items[0].Mimetype)
	})

	t.Run("filters by language", func(t *testing.T) {
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

		var en, fr library.Known
		require.NoError(t, testx.Fake(&en, library.KnownOptionTestDefaults))
		en.OriginalLanguage = "en"
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, en).Scan(&en))
		require.NoError(t, testx.Fake(&fr, library.KnownOptionTestDefaults))
		fr.OriginalLanguage = "fr"
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, fr).Scan(&fr))

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		body, err := json.Marshal(&media.RecommendationsRequest{Language: "en"})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/random", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
		require.Equal(t, en.UID, result.Items[0].Id)
	})

	t.Run("returns adult content when requested", func(t *testing.T) {
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
		known.Adult = true
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		body, err := json.Marshal(&media.RecommendationsRequest{Adult: true})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/random", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
		require.Equal(t, known.UID, result.Items[0].Id)
		require.True(t, result.Items[0].Adult)
	})

	t.Run("includes non-adult content when adult requested", func(t *testing.T) {
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

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		body, err := json.Marshal(&media.RecommendationsRequest{Adult: true})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/random", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
		require.Equal(t, known.UID, result.Items[0].Id)
	})

	t.Run("returns known media with description and image", func(t *testing.T) {
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

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/random", func() []byte { b, _ := json.Marshal(&media.RecommendationsRequest{}); return b }(),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
		require.False(t, result.Items[0].Adult)
		require.NotEmpty(t, result.Items[0].Description)
		require.NotEmpty(t, result.Items[0].Image)
	})
}

func TestRecommendationsRandom(t *testing.T) {
	t.Run("creates recommendation from random known media", func(t *testing.T) {
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

		routes := mux.NewRouter()
		media.NewHTTPRecommendations(q, media.HTTPRecommendationsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		body, err := json.Marshal(&media.RecommendationsRequest{})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/random", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)

		var result media.RecommendationsResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, known.UID, result.Items[0].Id)
	})

	t.Run("empty database returns 404", func(t *testing.T) {
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

		body, err := json.Marshal(&media.RecommendationsRequest{})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/random", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})
}
