package mediaapi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/websocketx"
	"github.com/retrovibed/retrovibed/shallows/mediaapi"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

func TestHTTPRemoteControl(t *testing.T) {
	t.Run("connect relays to listen and listen broadcasts to all connects", func(t *testing.T) {
		routes := mux.NewRouter()
		mediaapi.NewHTTPRemoteControl(
			true,
			mediaapi.HTTPRemoteControlOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/rc").Subrouter())
		server := httptest.NewServer(routes)
		defer server.Close()

		listentoken, err := mediaapi.RemoteControlListenToken()
		require.NoError(t, err)

		listenconn, _, err := websocket.Dial(t.Context(), fmt.Sprintf("ws://%s/rc/listen", server.Listener.Addr().String()), &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", listentoken)},
			},
		})
		require.NoError(t, err)
		defer listenconn.Close(websocket.StatusNormalClosure, "") //nolint: errcheck

		connecttoken := httpauthtest.UnsafeClaimsToken(
			metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(
				jwtx.NewJWTClaims("profile-1", jwtx.ClaimsOptionAuthnExpiration()),
				func(t *metaapi.Token) { t.RemoteControl = true },
			)),
			httpauthtest.UnsafeJWTSecretSource,
		)

		connect1, _, err := websocket.Dial(t.Context(), fmt.Sprintf("ws://%s/rc/connect", server.Listener.Addr().String()), &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", connecttoken)},
			},
		})
		require.NoError(t, err)
		defer connect1.Close(websocket.StatusNormalClosure, "") //nolint: errcheck

		connect2, _, err := websocket.Dial(t.Context(), fmt.Sprintf("ws://%s/rc/connect", server.Listener.Addr().String()), &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", connecttoken)},
			},
		})
		require.NoError(t, err)
		defer connect2.Close(websocket.StatusNormalClosure, "") //nolint: errcheck

		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		cmd, err := jsonx.Marshal(&mediaapi.Stream{
			Sid:     "cmd-1",
			Command: &mediaapi.Stream_Queue{Queue: &mediaapi.Queue{}},
		})
		require.NoError(t, err)
		require.NoError(t, connect1.Write(ctx, websocket.MessageBinary, cmd))

		_, received, err := listenconn.Read(ctx)
		require.NoError(t, err)

		var relayed mediaapi.Stream
		require.NoError(t, jsonx.Unmarshal(received, &relayed))
		require.Equal(t, "cmd-1", relayed.Sid)
		require.NotNil(t, relayed.GetQueue())

		reply, err := jsonx.Marshal(&mediaapi.Stream{
			Sid:     "reply-1",
			Command: &mediaapi.Stream_Pause{},
		})
		require.NoError(t, err)
		require.NoError(t, listenconn.Write(ctx, websocket.MessageBinary, reply))

		for _, conn := range []*websocket.Conn{connect1, connect2} {
			_, broadcast, err := conn.Read(ctx)
			require.NoError(t, err)

			var got mediaapi.Stream
			require.NoError(t, jsonx.Unmarshal(broadcast, &got))
			require.Equal(t, "reply-1", got.Sid)
			require.True(t, proto.Equal(&mediaapi.Pause{}, got.GetPause()))
		}
	})

	t.Run("listen rejects a bearer signed with the wrong secret", func(t *testing.T) {
		routes := mux.NewRouter()
		mediaapi.NewHTTPRemoteControl(true).Bind(routes.PathPrefix("/rc").Subrouter())
		server := httptest.NewServer(routes)
		defer server.Close()

		var claims jwt.RegisteredClaims
		wrongtoken := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource) // signed with a secret that isn't the listen secret

		_, resp, err := websocket.Dial(t.Context(), fmt.Sprintf("ws://%s/rc/listen", server.Listener.Addr().String()), &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", wrongtoken)},
			},
		})
		require.Error(t, err)
		require.NotNil(t, resp)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("connect closes with try again later when no listener is attached", func(t *testing.T) {
		routes := mux.NewRouter()
		mediaapi.NewHTTPRemoteControl(
			true,
			mediaapi.HTTPRemoteControlOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/rc").Subrouter())
		server := httptest.NewServer(routes)
		defer server.Close()

		token := httpauthtest.UnsafeClaimsToken(
			metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(
				jwtx.NewJWTClaims("profile-1", jwtx.ClaimsOptionAuthnExpiration()),
				func(t *metaapi.Token) { t.RemoteControl = true },
			)),
			httpauthtest.UnsafeJWTSecretSource,
		)

		conn, _, err := websocket.Dial(t.Context(), fmt.Sprintf("ws://%s/rc/connect", server.Listener.Addr().String()), &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		})
		require.NoError(t, err)
		defer conn.Close(websocket.StatusNormalClosure, "") //nolint: errcheck

		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		cmd, err := jsonx.Marshal(&mediaapi.Stream{
			Sid:     "cmd-no-listener",
			Command: &mediaapi.Stream_Queue{Queue: &mediaapi.Queue{}},
		})
		require.NoError(t, err)
		require.NoError(t, conn.Write(ctx, websocket.MessageBinary, cmd))

		_, _, err = conn.Read(ctx)
		require.Error(t, err)
		require.Equal(t, websocket.StatusTryAgainLater, websocket.CloseStatus(err))
	})

	t.Run("connect rejects a token without the remote control permission", func(t *testing.T) {
		routes := mux.NewRouter()
		mediaapi.NewHTTPRemoteControl(
			true,
			mediaapi.HTTPRemoteControlOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/rc").Subrouter())
		server := httptest.NewServer(routes)
		defer server.Close()

		token := httpauthtest.UnsafeClaimsToken(
			metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(
				jwtx.NewJWTClaims("profile-1", jwtx.ClaimsOptionAuthnExpiration()),
			)),
			httpauthtest.UnsafeJWTSecretSource,
		)

		_, resp, err := websocket.Dial(t.Context(), fmt.Sprintf("ws://%s/rc/connect", server.Listener.Addr().String()), &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		})
		require.Error(t, err)
		require.NotNil(t, resp)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("second listen connection evicts the first", func(t *testing.T) {
		routes := mux.NewRouter()
		mediaapi.NewHTTPRemoteControl(true).Bind(routes.PathPrefix("/rc").Subrouter())
		server := httptest.NewServer(routes)
		defer server.Close()

		listentoken, err := mediaapi.RemoteControlListenToken()
		require.NoError(t, err)

		first, _, err := websocket.Dial(t.Context(), fmt.Sprintf("ws://%s/rc/listen", server.Listener.Addr().String()), &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", listentoken)},
			},
		})
		require.NoError(t, err)
		defer first.Close(websocket.StatusNormalClosure, "") //nolint: errcheck

		second, _, err := websocket.Dial(t.Context(), fmt.Sprintf("ws://%s/rc/listen", server.Listener.Addr().String()), &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", listentoken)},
			},
		})
		require.NoError(t, err)
		defer second.Close(websocket.StatusNormalClosure, "") //nolint: errcheck

		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		_, _, err = first.Read(ctx)
		require.Error(t, err)
		require.Equal(t, websocketx.PrivateStatus(http.StatusConflict), websocket.CloseStatus(err))
	})

	t.Run("disabled rejects listen and connect with forbidden", func(t *testing.T) {
		routes := mux.NewRouter()
		mediaapi.NewHTTPRemoteControl(
			false,
			mediaapi.HTTPRemoteControlOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/rc").Subrouter())
		server := httptest.NewServer(routes)
		defer server.Close()

		listentoken, err := mediaapi.RemoteControlListenToken()
		require.NoError(t, err)

		_, resp, err := websocket.Dial(t.Context(), fmt.Sprintf("ws://%s/rc/listen", server.Listener.Addr().String()), &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", listentoken)},
			},
		})
		require.Error(t, err)
		require.NotNil(t, resp)
		require.Equal(t, http.StatusForbidden, resp.StatusCode)

		var claims jwt.RegisteredClaims
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		_, resp, err = websocket.Dial(t.Context(), fmt.Sprintf("ws://%s/rc/connect", server.Listener.Addr().String()), &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		})
		require.Error(t, err)
		require.NotNil(t, resp)
		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
