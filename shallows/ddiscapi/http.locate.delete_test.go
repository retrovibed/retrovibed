package ddiscapi_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPLocateDelete(t *testing.T) {
	t.Run("successfully delete", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			l      ddisc.Locate
			result ddiscapi.LocateDeleteResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, ddisc.LocateInsertWithDefaults(ctx, q, ddisc.NewLocate("ubuntu", "video")).Scan(&l))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPLocate(
			q,
			asyncx.NewWakeup(t.Context()),
			ddiscapi.HTTPLocateOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, fmt.Sprintf("/%s", l.ID), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		require.Equal(t, l.ID, result.Locate.Id)

		var target = sql.ErrNoRows
		require.ErrorAs(t, ddisc.LocateFindByID(ctx, q, l.ID).Scan(&l), &target)
	})

	t.Run("unknown id returns not found", func(t *testing.T) {
		var (
			p meta.Profile
			v meta.Authz
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPLocate(
			q,
			asyncx.NewWakeup(t.Context()),
			ddiscapi.HTTPLocateOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, fmt.Sprintf("/%s", uuid.Must(uuid.NewV4()).String()), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})
}
