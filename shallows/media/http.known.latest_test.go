package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/formx"
	"github.com/retrovibed/retrovibed/internal/grpcx"
	"github.com/retrovibed/retrovibed/internal/httptestx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/library"
	"github.com/retrovibed/retrovibed/media"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/stretchr/testify/require"
)

func TestKnownLatest(t *testing.T) {
	t.Run("handles zero state", func(t *testing.T) {
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

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(&media.KnownLatestRequest{
			Released: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)),
			Limit:    100,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/latest?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.KnownLatestResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Empty(t, result.Items)
	})

	t.Run("returns everything", func(t *testing.T) {
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
		for range 10 {
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))
		}

		routes := mux.NewRouter()
		media.NewHTTPKnown(q, media.HTTPKnownOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(&media.KnownLatestRequest{
			Released: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)),
			Limit:    100,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/latest?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.KnownLatestResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 10)
	})

	t.Run("excludes items without poster or backdrop path", func(t *testing.T) {
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

		// insert items with poster paths - should be returned
		for range 5 {
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))
		}

		// insert items without poster or backdrop paths - should be excluded
		for range 5 {
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionTestNoPoster))
			require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))
		}

		routes := mux.NewRouter()
		media.NewHTTPKnown(q, media.HTTPKnownOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(&media.KnownLatestRequest{
			Released: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)),
			Limit:    100,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/latest?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.KnownLatestResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 5)
	})

	t.Run("order by released date", func(t *testing.T) {
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
		var oldest, middle, newest library.Known
		require.NoError(t, testx.Fake(&oldest, library.KnownOptionTestDefaults, library.KnownOptionReleased(ts.Add(-48*time.Hour))))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, oldest).Scan(&oldest))

		require.NoError(t, testx.Fake(&middle, library.KnownOptionTestDefaults, library.KnownOptionReleased(ts.Add(-24*time.Hour))))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, middle).Scan(&middle))

		require.NoError(t, testx.Fake(&newest, library.KnownOptionTestDefaults, library.KnownOptionReleased(ts)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, newest).Scan(&newest))

		routes := mux.NewRouter()
		media.NewHTTPKnown(q, media.HTTPKnownOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(&media.KnownLatestRequest{
			Released: meta.NewDateRange(timex.NewRangeDuration(72 * time.Hour)),
			Limit:    100,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/latest?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.KnownLatestResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 3)
		require.Equal(t, newest.UID, result.Items[0].Id)
		require.Equal(t, grpcx.EncodeTime(newest.Released), result.Items[0].Released)
		require.Equal(t, middle.UID, result.Items[1].Id)
		require.Equal(t, grpcx.EncodeTime(middle.Released), result.Items[1].Released)
		require.Equal(t, oldest.UID, result.Items[2].Id)
		require.Equal(t, grpcx.EncodeTime(oldest.Released), result.Items[2].Released)
	})
}
