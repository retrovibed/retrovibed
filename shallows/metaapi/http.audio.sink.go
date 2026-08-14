package metaapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/coder/websocket"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/audiox"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/websocketx"
)

// Sinker abstracts audio output device enumeration/selection so
// HTTPAudioSink can be tested without a real PulseAudio server.
type Sinker interface {
	ListSinks() iterx.Seq[audiox.Sink]
	Current() (audiox.Sink, error)
	Activate(id string) error
}

// defaultSinker wraps the audiox package functions.
type defaultSinker struct{}

func (defaultSinker) ListSinks() iterx.Seq[audiox.Sink] { return audiox.ListSinks() }
func (defaultSinker) Current() (audiox.Sink, error)     { return audiox.Current() }
func (defaultSinker) Activate(id string) error          { return audiox.Activate(id) }

func newAudioSink(s audiox.Sink) *AudioSink {
	return &AudioSink{
		Id:   s.ID,
		Name: s.Name,
	}
}

type HTTPAudioSinkOption func(*HTTPAudioSink)

func HTTPAudioSinkOptionSinker(s Sinker) HTTPAudioSinkOption {
	return func(t *HTTPAudioSink) {
		t.sink = s
	}
}

func HTTPAudioSinkOptionJWTSecret(j jwtx.SecretSource) HTTPAudioSinkOption {
	return func(t *HTTPAudioSink) {
		t.jwtsecret = j
	}
}

func HTTPAudioSinkOptionSupported(supported bool) HTTPAudioSinkOption {
	return func(t *HTTPAudioSink) {
		t.supported = supported
	}
}

func NewHTTPAudioSink(options ...HTTPAudioSinkOption) *HTTPAudioSink {
	svc := langx.Clone(HTTPAudioSink{
		sink:      defaultSinker{},
		jwtsecret: env.JWTSecret,
		supported: audiox.Supported,
	}, options...)

	return &svc
}

type HTTPAudioSink struct {
	sink      Sinker
	jwtsecret jwtx.SecretSource
	supported bool
}

func isWebsocketUpgrade(r *http.Request, _ *mux.RouteMatch) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket")
}

func (t *HTTPAudioSink) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).MatcherFunc(isWebsocketUpgrade).Handler(alice.New(
		httpx.GatedResponse(t.supported, http.StatusServiceUnavailable),
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.listen))

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.GatedResponse(t.supported, http.StatusServiceUnavailable),
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.current))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.GatedResponse(t.supported, http.StatusServiceUnavailable),
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.activate))
}

func (t *HTTPAudioSink) listen(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("failed to accept audio sink listen websocket", err)
		return
	}
	defer func() {
		errorsx.Log(c.Close(websocket.StatusNormalClosure, ""))
	}()

	seq := t.sink.ListSinks()
	resp := AudioSinkSearchResponse{}
	for s := range seq.Each(r.Context()) {
		resp.Items = append(resp.Items, newAudioSink(s))
	}

	if err := seq.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to list sinks"))
		errorsx.Log(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "internal service error"))
		return
	}

	encoded, err := json.Marshal(&resp)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to encode sinks"))
		errorsx.Log(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "internal service error"))
		return
	}

	if err := c.Write(r.Context(), websocket.MessageBinary, encoded); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write sinks"))
		return
	}

	ctx := c.CloseRead(context.Background())
	<-ctx.Done()
}

func (t *HTTPAudioSink) current(w http.ResponseWriter, r *http.Request) {
	sink, err := t.sink.Current()
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to determine current sink"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &AudioSinkCurrentResponse{
		Sink: newAudioSink(sink),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPAudioSink) activate(w http.ResponseWriter, r *http.Request) {
	var msg AudioSinkTouchRequest

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err := t.sink.Activate(msg.Id); err != nil {
		log.Println(errorsx.Wrap(err, "unable to activate sink"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	sink, err := t.sink.Current()
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to determine current sink"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &AudioSinkTouchResponse{
		Sink: newAudioSink(sink),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
