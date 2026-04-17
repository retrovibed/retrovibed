package ddiscapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestHTTPPeerManagementSearch(t *testing.T) {
	t.Run("search", func(t *testing.T) {
		var (
			p      tracking.Peer
			result ddiscapi.PeerSearchResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		p = tracking.NewPeerFromInfo(krpc.RandomNodeInfo(16), tracking.PeerOptionTestDefaults)
		require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, p).Scan(&p))

		routes := mux.NewRouter()

		ddiscapi.NewHTTPPeerManagement(
			q,
			ddiscapi.HTTPPeerManagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(formx.NewEncoder().Encode(&ddiscapi.PeerSearchRequest{
			Offset: 0,
		}))(t)

		claims = jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		encoded := testx.Must(ddiscapi.NewPeerFromTrackingPeer(p))(t)
		require.Equal(t, p.Description, encoded.Description)
		require.Equal(t, result.Next.Offset, uint64(0))
		require.Contains(t, result.Items, encoded)
	})
}
