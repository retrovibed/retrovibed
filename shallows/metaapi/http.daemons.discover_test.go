package metaapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/hashicorp/mdns"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPDaemonDiscover(t *testing.T) {
	t.Run("gated", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		routes := mux.NewRouter()
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			metaapi.HTTPDaemonsOptionMDNSDiscovery(false),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestContextBytes(context.Background(), http.MethodGet, "/discover", nil, httptestx.RequestOptionAuthorization(token), httptestx.RequestOptionHeader("Upgrade", "websocket"))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusServiceUnavailable, resp.Code)
	})

	t.Run("discover", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		routes := mux.NewRouter()
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			metaapi.HTTPDaemonsOptionMDNSDiscovery(true),
			metaapi.HTTPDaemonsOptionMDNSLookup(func(ctx context.Context, service string, timeout time.Duration) ([]*mdns.ServiceEntry, error) {
				return []*mdns.ServiceEntry{
					{Host: "peer1.local."},
					{Host: "peer2.local."},
				}, nil
			}),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		server := httptest.NewServer(routes)
		defer server.Close()

		wsURL := fmt.Sprintf("ws://%s/discover", server.Listener.Addr().String())
		c, wsResp, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		})
		require.NoError(t, err)
		defer c.Close(websocket.StatusNormalClosure, "") //nolint: errcheck
		require.Equal(t, http.StatusSwitchingProtocols, wsResp.StatusCode)

		var hostnames []string
		for i := 0; i < 2; i++ {
			var d metaapi.Daemon
			messageType, data, err := c.Read(ctx)
			require.NoError(t, err)
			require.Equal(t, websocket.MessageText, messageType)
			require.NoError(t, json.Unmarshal(data, &d))
			hostnames = append(hostnames, d.Hostname)
		}

		require.ElementsMatch(t, []string{"peer1.local.", "peer2.local."}, hostnames)
		require.Equal(t, 2, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_daemons"))(t))
	})

	t.Run("already known peers are not re-streamed", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		existing := meta.Daemon{Hostname: "peer1.local."}
		meta.DaemonOptionMaybeID(&existing)
		meta.DaemonOptionEnsureDescription(&existing)
		require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, existing).Scan(&existing))

		routes := mux.NewRouter()
		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			metaapi.HTTPDaemonsOptionMDNSDiscovery(true),
			metaapi.HTTPDaemonsOptionMDNSLookup(func(ctx context.Context, service string, timeout time.Duration) ([]*mdns.ServiceEntry, error) {
				return []*mdns.ServiceEntry{
					{Host: "peer1.local."},
					{Host: "peer2.local."},
				}, nil
			}),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		server := httptest.NewServer(routes)
		defer server.Close()

		wsURL := fmt.Sprintf("ws://%s/discover", server.Listener.Addr().String())
		c, wsResp, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		})
		require.NoError(t, err)
		defer c.Close(websocket.StatusNormalClosure, "") //nolint: errcheck
		require.Equal(t, http.StatusSwitchingProtocols, wsResp.StatusCode)

		var d metaapi.Daemon
		messageType, data, err := c.Read(ctx)
		require.NoError(t, err)
		require.Equal(t, websocket.MessageText, messageType)
		require.NoError(t, json.Unmarshal(data, &d))
		require.Equal(t, "peer2.local.", d.Hostname)

		_, _, err = c.Read(ctx)
		require.Error(t, err) // socket closes after the single newly discovered peer.

		require.Equal(t, 2, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_daemons"))(t))
	})
}
