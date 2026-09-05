package metaapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/audiox"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPAudioSinkActivate(t *testing.T) {
	t.Run("gated", func(t *testing.T) {
		routes := mux.NewRouter()
		metaapi.NewHTTPAudioSink(
			metaapi.HTTPAudioSinkOptionSinker(&fakeAudioSinker{}),
			metaapi.HTTPAudioSinkOptionSupported(false),
			metaapi.HTTPAudioSinkOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		body, err := json.Marshal(&metaapi.AudioSinkTouchRequest{Id: "sink-1"})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(context.Background(), http.MethodPost, "/", body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusServiceUnavailable, resp.Code)
	})

	t.Run("activate", func(t *testing.T) {
		var result metaapi.AudioSinkTouchResponse

		sinker := &fakeAudioSinker{
			currentSink: audiox.Sink{ID: "sink-2", Name: "USB Headset"},
		}

		routes := mux.NewRouter()
		metaapi.NewHTTPAudioSink(
			metaapi.HTTPAudioSinkOptionSinker(sinker),
			metaapi.HTTPAudioSinkOptionSupported(true),
			metaapi.HTTPAudioSinkOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		body, err := json.Marshal(&metaapi.AudioSinkTouchRequest{Id: "sink-2"})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(context.Background(), http.MethodPost, "/", body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, "sink-2", sinker.activatedID)
		require.Equal(t, "sink-2", result.Sink.Id)
		require.Equal(t, "USB Headset", result.Sink.Name)
	})

	t.Run("activate error", func(t *testing.T) {
		routes := mux.NewRouter()
		metaapi.NewHTTPAudioSink(
			metaapi.HTTPAudioSinkOptionSinker(&fakeAudioSinker{
				activateErr: errors.New("no such sink"),
			}),
			metaapi.HTTPAudioSinkOptionSupported(true),
			metaapi.HTTPAudioSinkOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		body, err := json.Marshal(&metaapi.AudioSinkTouchRequest{Id: "missing"})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(context.Background(), http.MethodPost, "/", body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("decode error", func(t *testing.T) {
		routes := mux.NewRouter()
		metaapi.NewHTTPAudioSink(
			metaapi.HTTPAudioSinkOptionSinker(&fakeAudioSinker{}),
			metaapi.HTTPAudioSinkOptionSupported(true),
			metaapi.HTTPAudioSinkOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestContextBytes(context.Background(), http.MethodPost, "/", []byte("not json"), httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
