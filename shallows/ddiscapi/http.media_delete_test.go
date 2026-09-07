package ddiscapi_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPMediaDelete(t *testing.T) {
	t.Run("successfully delete", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			d      ddisc.Discovered
			result ddiscapi.MediaDeleteResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		id := int160.Random()
		d = ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPMedia(
			q,
			ddiscapi.HTTPMediaOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, fmt.Sprintf("/%s", d.ID), nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		require.Equal(t, result.Media.GetId(), d.ID)

		var (
			target = sql.ErrNoRows
		)
		require.ErrorAs(t, ddisc.DiscoveredFindByID(ctx, q, d.ID).Scan(&d), &target)
	})
}
