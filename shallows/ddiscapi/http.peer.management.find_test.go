package ddiscapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/ddiscapi"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/httptestx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/tracking"
	"github.com/stretchr/testify/require"
)

func TestHTTPPeerManagementFind(t *testing.T) {
	var (
		p      tracking.Peer
		result ddiscapi.PeerFindResponse
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

	token := httpauthtest.UnsafeClaimsToken(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), httpauthtest.UnsafeJWTSecretSource)
	resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, fmt.Sprintf("/%s", p.ID), nil, httptestx.RequestOptionAuthorization(token))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.Equal(t, result.Peer.Id, p.ID)
}
