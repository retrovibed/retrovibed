package metaapi_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPWireguardDelete(t *testing.T) {
	var (
		result metaapi.WireguardDeleteResponse
		claims jwt.RegisteredClaims
		wg     meta.Wireguard
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	require.NoError(
		t,
		meta.WireguardInsertWithDefaults(
			ctx,
			q,
			meta.NewWireguard(testx.Must(uuid.NewV4())(t).String(), meta.WireguardOptionDescription("test")),
		).Scan(&wg),
	)

	tmpdir := fsx.DirVirtual(t.TempDir())
	path := wg.ID
	d, err := os.Create(tmpdir.Path(wg.ID))
	require.NoError(t, err)
	_, err = io.Copy(d, testx.Read(testx.Fixture("wireguard", "example.1.conf")))
	require.NoError(t, err)
	require.NoError(t, d.Close())

	routes := mux.NewRouter()

	metaapi.NewHTTPWireguard(
		tmpdir.Path(),
		q,
		metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

	resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodDelete, fmt.Sprintf("/%s", path), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_wireguard"))
	routes.ServeHTTP(resp, req)

	require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_wireguard"))
	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

	require.Equal(t, wg.ID, result.Wireguard.Id)
}
