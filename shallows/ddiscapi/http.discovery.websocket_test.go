package ddiscapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coder/websocket"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestHTTPDiscoveryLocateSocket(t *testing.T) {
	t.Run("streams every ranked candidate then closes with the best last", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		kid := uuid.Must(uuid.NewV4()).String()

		worseID := int160.Random()
		worse := ddisc.NewDiscovered(&worseID,
			ddisc.DiscoveredOptionKnownMedia(kid),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionTitle("Some.Movie.2024.480p.SDTV.x264"),
			ddisc.DiscoveredOptionAutoMagnet,
		)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, worse).Scan(&worse))

		bestID := int160.Random()
		best := ddisc.NewDiscovered(&bestID,
			ddisc.DiscoveredOptionKnownMedia(kid),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionTitle("Some.Movie.2024.1080p.BluRay.x264"),
			ddisc.DiscoveredOptionAutoMagnet,
		)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, best).Scan(&best))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			nil,
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(kid, jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		server := httptest.NewServer(routes)
		defer server.Close()

		wsURL := fmt.Sprintf("ws://%s/locate?known_media_id=%s", server.Listener.Addr().String(), kid)

		c, wsResp, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		})
		require.NoError(t, err)
		defer c.Close(websocket.StatusNormalClosure, "") //nolint: errcheck
		require.Equal(t, http.StatusSwitchingProtocols, wsResp.StatusCode)

		var received []*ddiscapi.Discovery
		for {
			messageType, data, err := c.Read(ctx)
			if err != nil {
				break
			}
			require.Equal(t, websocket.MessageBinary, messageType)

			var result ddiscapi.Discovery
			require.NoError(t, json.Unmarshal(data, &result))
			received = append(received, &result)
		}

		require.NotEmpty(t, received)
		require.Equal(t, best.ID, received[len(received)-1].Id)
	})
}
