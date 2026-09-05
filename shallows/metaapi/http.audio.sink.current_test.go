package metaapi_test

import (
	"context"
	"errors"
	"iter"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
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

// fakeAudioSinker is a metaapi.Sinker test double, shared by the
// http.audio.sink.*_test.go files.
type fakeAudioSinker struct {
	sinks       []audiox.Sink
	sinksErr    error
	currentSink audiox.Sink
	currentErr  error
	activateErr error
	activatedID string
}

func (f *fakeAudioSinker) ListSinks() iterx.Seq[audiox.Sink] {
	return fakeSinkSeq{items: f.sinks, err: f.sinksErr}
}

func (f *fakeAudioSinker) Current() (audiox.Sink, error) { return f.currentSink, f.currentErr }

func (f *fakeAudioSinker) Activate(id string) error {
	f.activatedID = id
	return f.activateErr
}

type fakeSinkSeq struct {
	items []audiox.Sink
	err   error
}

func (s fakeSinkSeq) Each(_ context.Context) iter.Seq[audiox.Sink] { return iterx.From(s.items...) }
func (s fakeSinkSeq) Err() error                                   { return s.err }

func TestHTTPAudioSinkCurrent(t *testing.T) {
	t.Run("gated", func(t *testing.T) {
		routes := mux.NewRouter()
		metaapi.NewHTTPAudioSink(
			metaapi.HTTPAudioSinkOptionSinker(&fakeAudioSinker{}),
			metaapi.HTTPAudioSinkOptionSupported(false),
			metaapi.HTTPAudioSinkOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestContextBytes(context.Background(), http.MethodGet, "/", nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusServiceUnavailable, resp.Code)
	})

	t.Run("current", func(t *testing.T) {
		var result metaapi.AudioSinkCurrentResponse

		routes := mux.NewRouter()
		metaapi.NewHTTPAudioSink(
			metaapi.HTTPAudioSinkOptionSinker(&fakeAudioSinker{
				currentSink: audiox.Sink{ID: "alsa_output.pci-0000_00_1f.3", Name: "Built-in Audio"},
			}),
			metaapi.HTTPAudioSinkOptionSupported(true),
			metaapi.HTTPAudioSinkOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestContextBytes(context.Background(), http.MethodGet, "/", nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, "alsa_output.pci-0000_00_1f.3", result.Sink.Id)
		require.Equal(t, "Built-in Audio", result.Sink.Name)
	})

	t.Run("error", func(t *testing.T) {
		routes := mux.NewRouter()
		metaapi.NewHTTPAudioSink(
			metaapi.HTTPAudioSinkOptionSinker(&fakeAudioSinker{
				currentErr: errors.New("no default sink"),
			}),
			metaapi.HTTPAudioSinkOptionSupported(true),
			metaapi.HTTPAudioSinkOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(testx.Must(uuid.NewV4())(t).String(), jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestContextBytes(context.Background(), http.MethodGet, "/", nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
