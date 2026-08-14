package metaapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coder/websocket"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/audiox"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPAudioSinkListen(t *testing.T) {
	t.Run("gated", func(t *testing.T) {
		routes := mux.NewRouter()
		metaapi.NewHTTPAudioSink(
			metaapi.HTTPAudioSinkOptionSinker(&fakeAudioSinker{}),
			metaapi.HTTPAudioSinkOptionSupported(false),
			metaapi.HTTPAudioSinkOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestContextBytes(context.Background(), http.MethodGet, "/", nil, httptestx.RequestOptionAuthorization(token), httptestx.RequestOptionHeader("Upgrade", "websocket"))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusServiceUnavailable, resp.Code)
	})

	t.Run("listen", func(t *testing.T) {
		var result metaapi.AudioSinkSearchResponse

		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPAudioSink(
			metaapi.HTTPAudioSinkOptionSinker(&fakeAudioSinker{
				sinks: []audiox.Sink{
					{ID: "sink-1", Name: "Built-in Audio"},
					{ID: "sink-2", Name: "USB Headset"},
				},
			}),
			metaapi.HTTPAudioSinkOptionSupported(true),
			metaapi.HTTPAudioSinkOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		server := httptest.NewServer(routes)
		defer server.Close()

		wsURL := fmt.Sprintf("ws://%s/", server.Listener.Addr().String())
		c, wsResp, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		})
		require.NoError(t, err)
		defer c.Close(websocket.StatusNormalClosure, "") //nolint: errcheck
		require.Equal(t, http.StatusSwitchingProtocols, wsResp.StatusCode)

		messageType, data, err := c.Read(ctx)
		require.NoError(t, err)
		require.Equal(t, websocket.MessageBinary, messageType)
		require.NoError(t, json.Unmarshal(data, &result))
		require.Len(t, result.Items, 2)
		require.Equal(t, "sink-1", result.Items[0].Id)
		require.Equal(t, "sink-2", result.Items[1].Id)
	})

	t.Run("list error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes := mux.NewRouter()
		metaapi.NewHTTPAudioSink(
			metaapi.HTTPAudioSinkOptionSinker(&fakeAudioSinker{
				sinksErr: errors.New("pulseaudio unreachable"),
			}),
			metaapi.HTTPAudioSinkOptionSupported(true),
			metaapi.HTTPAudioSinkOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		server := httptest.NewServer(routes)
		defer server.Close()

		wsURL := fmt.Sprintf("ws://%s/", server.Listener.Addr().String())
		c, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		})
		require.NoError(t, err)
		defer c.Close(websocket.StatusNormalClosure, "") //nolint: errcheck

		_, _, err = c.Read(ctx)
		require.Error(t, err)
	})
}
