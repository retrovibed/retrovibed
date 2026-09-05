package metaapi_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
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
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

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
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

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
	require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

	require.Equal(t, result.Profile.Id, p.ID)
}
