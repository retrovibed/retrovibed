package ddiscapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

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

func TestHTTPLocateRetry(t *testing.T) {
	t.Run("clears tombstone and cooldown, keeps attempts", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			l      ddisc.Locate
			result ddiscapi.LocateRetryResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, ddisc.LocateInsertWithDefaults(ctx, q, ddisc.NewLocate("ubuntu", "video")).Scan(&l))
		require.NoError(t, ddisc.LocateCooldown(ctx, q, l.ID, time.Now().Add(time.Hour)).Scan(&l))
		require.NoError(t, ddisc.LocateCompleted(ctx, q, l.ID).Scan(&l))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPLocate(
			q,
			asyncx.NewWakeup(t.Context()),
			ddiscapi.HTTPLocateOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

		body := testx.Must(json.Marshal(ddiscapi.LocateRetryRequest{}))(t)
		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, fmt.Sprintf("/%s/retry", l.ID), body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, l.ID, result.Locate.Id)

		require.NoError(t, ddisc.LocateFindByID(ctx, q, l.ID).Scan(&l))
		require.EqualValues(t, 1, l.Attempts)
		require.True(t, l.TombstonedAt.After(time.Now()))
		require.False(t, l.NextCheckAt.After(time.Now()))
	})

	t.Run("reset attempts", func(t *testing.T) {
		var (
			p meta.Profile
			v meta.Authz
			l ddisc.Locate
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, ddisc.LocateInsertWithDefaults(ctx, q, ddisc.NewLocate("ubuntu", "video")).Scan(&l))
		require.NoError(t, ddisc.LocateCooldown(ctx, q, l.ID, time.Now().Add(time.Hour)).Scan(&l))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPLocate(
			q,
			asyncx.NewWakeup(t.Context()),
			ddiscapi.HTTPLocateOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

		body := testx.Must(json.Marshal(ddiscapi.LocateRetryRequest{ResetAttempts: true}))(t)
		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, fmt.Sprintf("/%s/retry", l.ID), body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		require.NoError(t, ddisc.LocateFindByID(ctx, q, l.ID).Scan(&l))
		require.EqualValues(t, 0, l.Attempts)
		require.False(t, l.NextCheckAt.After(time.Now()))
	})
}
