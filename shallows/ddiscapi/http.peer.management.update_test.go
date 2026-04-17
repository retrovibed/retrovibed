package ddiscapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestHTTPPeerManagementUpdate(t *testing.T) {
	var (
		p      meta.Profile
		v      meta.Authz
		peer   tracking.Peer
		result ddiscapi.PeerUpdateResponse
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

	peer = tracking.NewPeer(int160.Random(), tracking.PeerOptionTestDefaults)
	require.NoError(t, tracking.PeerInsertWithDefaults(ctx, q, peer).Scan(&peer))

	peer.Description = "updated description"
	encoded := testx.Must(ddiscapi.NewPeerFromTrackingPeer(peer))(t)

	routes := mux.NewRouter()
	ddiscapi.NewHTTPPeerManagement(
		q,
		ddiscapi.HTTPPeerManagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

	body := testx.Must(json.Marshal(ddiscapi.PeerUpdateRequest{
		Peer: encoded,
	}))(t)

	resp, req, err := httptestx.BuildRequestBytes(http.MethodPatch, fmt.Sprintf("/%s", peer.ID), body, httptestx.RequestOptionAuthorization(token))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.Equal(t, result.Peer.Id, peer.ID)
	require.Equal(t, result.Peer.Description, "updated description")
}
