package ddiscapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestHTTPDiscoveringCreate(t *testing.T) {
	var (
		p        meta.Profile
		v        meta.Authz
		result   ddiscapi.DiscoveryCreateResponse
		infohash = int160.Random()
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

	routes := mux.NewRouter()
	ddiscapi.NewHTTPDiscovery(
		q,
		searchplugin.Unimplemented{},
		ddisc.UnimplementedStrategy{},
		tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir())),
		ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

	body := testx.Must(json.Marshal(ddiscapi.DiscoveryCreateRequest{
		Discovery: &ddiscapi.Discovery{
			Infohash: infohash.Bytes(),
		},
	}))(t)

	resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/", body, httptestx.RequestOptionAuthorization(token))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

	require.Equal(t, result.Discovery.Infohash, infohash.Bytes())

	found := sqlx.Scan(tracking.UnknownSearch(ctx, q, tracking.UnknownSearchBuilder().Where(tracking.UnknownHashQueryByIDs(result.Discovery.Id))))
	rows := 0
	for range found.Iter() {
		rows++
	}
	require.NoError(t, found.Err())
	require.Equal(t, 1, rows)
}
