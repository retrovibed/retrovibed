package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestLocateFind(t *testing.T) {
	var (
		p      meta.Profile
		authz  meta.Authz
		d      library.Locate
		result media.LocateLookupResponse
	)
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
	require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))
	require.NoError(t, testx.Fake(&d, library.LocateOptionTestDefaults))
	require.NoError(t, library.LocateInsertWithDefaults(ctx, q, d).Scan(&d))

	routes := mux.NewRouter()

	media.NewHTTPLocate(
		q,
		media.HTTPLocateOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

	resp, req, err := httptestx.BuildRequestBytes(
		http.MethodGet,
		fmt.Sprintf("/%s", d.ID),
		nil,
		httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
	)
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Result().StatusCode)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.Equal(t, d.ID, result.Locate.Id)
}
