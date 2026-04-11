package metaapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/formx"
	"github.com/retrovibed/retrovibed/internal/httptestx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPUserManagementSearch(t *testing.T) {
	t.Run("search", func(t *testing.T) {
		var (
			p      meta.Profile
			result metaapi.ProfileSearchResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, meta.ProfileEnable(ctx, q, p.ID).Scan(&p))

		routes := mux.NewRouter()

		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(formx.NewEncoder().Encode(&metaapi.ProfileSearchRequest{
			Offset: 0,
		}))(t)

		claims = jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		encoded := testx.Must(metaapi.NewProfileFromMetaProfile(p))(t)
		require.Equal(t, p.Display, encoded.Display)
		require.Equal(t, result.Next.Offset, uint64(0))
		require.Contains(t, result.Items, encoded)
	})

	t.Run("search pending", func(t *testing.T) {
		var (
			p      meta.Profile
			result metaapi.ProfileSearchResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		routes := mux.NewRouter()

		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(formx.NewEncoder().Encode(&metaapi.ProfileSearchRequest{
			Status: 2,
			Offset: 0,
		}))(t)

		claims = jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		encoded := testx.Must(metaapi.NewProfileFromMetaProfile(p))(t)
		require.Equal(t, result.Next.Offset, uint64(0))
		require.Contains(t, result.Items, encoded)
	})
}

func TestHTTPUserManagementFind(t *testing.T) {
	var (
		p      meta.Profile
		result metaapi.ProfileLookupResponse
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

	routes := mux.NewRouter()
	metaapi.NewHTTPUsermanagement(
		q,
		metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	mut := p

	token := httpauthtest.UnsafeClaimsToken(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), httpauthtest.UnsafeJWTSecretSource)
	resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, fmt.Sprintf("/%s", mut.ID), nil, httptestx.RequestOptionAuthorization(token))
	require.NoError(t, err)
	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.Equal(t, result.Profile.Id, p.ID)
}
