package metaapi_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/httptestx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPWireguardUpdate(t *testing.T) {
	t.Run("update payload", func(t *testing.T) {
		var (
			result metaapi.WireguardUpdateResponse
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
		d, err := os.Create(tmpdir.Path(path))
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

		request := metaapi.WireguardUpdateRequest{
			Wireguard: &metaapi.Wireguard{
				Id:          wg.ID,
				Description: "updated description",
			},
		}

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodPatch, fmt.Sprintf("/%s", path), testx.Must(json.Marshal(&request))(t), httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.Equal(t, wg.ID, result.Wireguard.Id)
		require.Equal(t, "updated description", result.Wireguard.Description)
	})
}
