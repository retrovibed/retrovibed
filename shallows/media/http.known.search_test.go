package media_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestKnownSearch(t *testing.T) {
	t.Run("returns all items", func(t *testing.T) {
		var (
			p      meta.Profile
			authz  meta.Authz
			known  library.Known
			result media.KnownSearchResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))
		for range 10 {
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))
		}

		routes := mux.NewRouter()

		media.NewHTTPKnown(
			q,
			media.HTTPKnownOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, 10, len(result.Items))
	})

	t.Run("filters by mimetype", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
			known library.Known
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		for range 5 {
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))
		}
		for range 5 {
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Audio)))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))
		}

		routes := mux.NewRouter()
		media.NewHTTPKnown(q, media.HTTPKnownOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/?mimetype=%s", mimex.Video),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.KnownSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Len(t, result.Items, 5)
		for _, item := range result.Items {
			require.Equal(t, mimex.Video, item.Mimetype)
		}
	})

	t.Run("filters by source", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
			known library.Known
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		for range 5 {
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionSource("source-a")))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))
		}
		for range 5 {
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionSource("source-b")))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))
		}

		routes := mux.NewRouter()
		media.NewHTTPKnown(q, media.HTTPKnownOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(&media.KnownSearchRequest{
			Source: []string{"source-a"},
			Limit:  100,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.KnownSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Len(t, result.Items, 5)
	})

	t.Run("filters by id", func(t *testing.T) {
		var (
			p       meta.Profile
			authz   meta.Authz
			match   library.Known
			nomatch library.Known
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		require.NoError(t, testx.Fake(&match, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, match).Scan(&match))

		require.NoError(t, testx.Fake(&nomatch, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, nomatch).Scan(&nomatch))

		routes := mux.NewRouter()
		media.NewHTTPKnown(q, media.HTTPKnownOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(&media.KnownSearchRequest{
			Id:    []string{match.UID},
			Limit: 100,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.KnownSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Len(t, result.Items, 1)
		require.Equal(t, match.UID, result.Items[0].Uid)
	})

	t.Run("filters by released range", func(t *testing.T) {
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

		ts := time.Now()
		var inrange, outofrange library.Known
		require.NoError(t, testx.Fake(&inrange, library.KnownOptionTestDefaults, library.KnownOptionReleased(ts.Add(-12*time.Hour))))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, inrange).Scan(&inrange))

		require.NoError(t, testx.Fake(&outofrange, library.KnownOptionTestDefaults, library.KnownOptionReleased(ts.Add(-48*time.Hour))))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, outofrange).Scan(&outofrange))

		routes := mux.NewRouter()
		media.NewHTTPKnown(q, media.HTTPKnownOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(&media.KnownSearchRequest{
			Released: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)),
			Limit:    100,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.KnownSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Len(t, result.Items, 1)
		require.Equal(t, inrange.UID, result.Items[0].Uid)
		require.Equal(t, grpcx.EncodeTime(inrange.Released), result.Items[0].Released)
	})

	t.Run("rejects malformed released range", func(t *testing.T) {
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
		media.NewHTTPKnown(q, media.HTTPKnownOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?released[oldest]=not-a-date&released[newest]=not-a-date",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})
}
