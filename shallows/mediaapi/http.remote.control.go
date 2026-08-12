package mediaapi

import (
	"context"
	"crypto/rand"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/websocketx"
)

// generated once per process, never persisted to disk - dies with the
// process so a leaked/logged listen token can never outlive it or be
// replayed against a different run.
var remoteControlSecret = sync.OnceValue(func() []byte {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(errorsx.Wrap(err, "unable to generate remote control listen secret"))
	}
	return b
})

// RemoteControlListenToken mints a bearer authorized to connect to the local
// /rc/listen endpoint. Only this process can ever produce or validate one,
// since it's signed with remoteControlListenSecret.
func RemoteControlListenToken() (string, error) {
	claims := jwtx.NewJWTClaims(
		"remotecontrol",
		jwtx.ClaimsOptionAuthnExpiration(),
		jwtx.ClaimsOptionIssuer("remotecontrol"),
	)

	bearer, err := jwtx.Signed(remoteControlSecret(), claims)
	return bearer, errorsx.Wrap(err, "unable to sign remote control listen token")
}

type HTTPRemoteControlOption func(*HTTPRemoteControl)

func HTTPRemoteControlOptionJWTSecret(j jwtx.SecretSource) HTTPRemoteControlOption {
	return func(t *HTTPRemoteControl) {
		t.jwtsecret = j
	}
}

func NewHTTPRemoteControl(enabled bool, options ...HTTPRemoteControlOption) *HTTPRemoteControl {
	svc := &HTTPRemoteControl{
		enabled:   enabled,
		jwtsecret: env.JWTSecret,
		connects:  make(map[*websocket.Conn]struct{}),
	}

	for _, opt := range options {
		opt(svc)
	}

	return svc
}

// HTTPRemoteControl brokers playback commands between exactly one local
// "listen" socket (this device's own frontend) and any number of "connect"
// sockets (remote clients belonging to the account). connect is
// fire-and-forget: commands are relayed to listen without correlation, and
// anything listen writes back is broadcast to every open connect socket.
type HTTPRemoteControl struct {
	enabled   bool
	jwtsecret jwtx.SecretSource

	mu       sync.Mutex // guards everything below, including all conn writes
	listener *websocket.Conn
	connects map[*websocket.Conn]struct{}
}

func (t *HTTPRemoteControl) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/listen").Methods(http.MethodGet).Handler(alice.New(
		httpx.GatedResponse(t.enabled, http.StatusForbidden),
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(remoteControlSecret),
		httpx.Timeout2s(),
	).ThenFunc(t.listen))

	r.Path("/connect").Methods(http.MethodGet).Handler(alice.New(
		httpx.GatedResponse(t.enabled, http.StatusForbidden),
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.connect))
}

func (t *HTTPRemoteControl) listen(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("failed to accept remote control listen websocket", err)
		return
	}

	t.mu.Lock()
	if t.listener != nil {
		errorsx.Log(t.listener.Close(websocketx.PrivateStatus(http.StatusConflict), "replaced by another listener"))
	}
	t.listener = c
	t.mu.Unlock()

	defer func() {
		t.mu.Lock()
		if t.listener == c {
			t.listener = nil
		}
		t.mu.Unlock()
		errorsx.Log(c.Close(websocket.StatusNormalClosure, ""))
	}()

	ctx := context.Background()
	for {
		_, data, err := c.Read(ctx)
		if err != nil {
			log.Println(errorsx.Wrap(err, "remote control listen socket closed"))
			return
		}

		var msg Stream
		if err := protojson.Unmarshal(data, &msg); err != nil {
			log.Println(errorsx.Wrap(err, "unable to decode remote control frame"))
			continue
		}

		t.mu.Lock()
		for conn := range t.connects {
			if err := conn.Write(ctx, websocket.MessageBinary, data); err != nil {
				log.Println(errorsx.Wrap(err, "unable to broadcast remote control frame"))
			}
		}
		t.mu.Unlock()
	}
}

func (t *HTTPRemoteControl) connect(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("failed to accept remote control connect websocket", err)
		return
	}

	t.mu.Lock()
	t.connects[c] = struct{}{}
	t.mu.Unlock()

	defer func() {
		t.mu.Lock()
		delete(t.connects, c)
		t.mu.Unlock()
		errorsx.Log(c.Close(websocket.StatusNormalClosure, ""))
	}()

	ctx := context.Background()
	for {
		_, data, err := c.Read(ctx)
		if err != nil {
			log.Println(errorsx.Wrap(err, "remote control connect socket closed"))
			return
		}

		var msg Stream
		if err := protojson.Unmarshal(data, &msg); err != nil {
			log.Println(errorsx.Wrap(err, "unable to decode remote control command"))
			continue
		}

		t.mu.Lock()
		listener := t.listener
		if listener != nil {
			err = listener.Write(ctx, websocket.MessageBinary, data)
		}
		t.mu.Unlock()

		if listener == nil {
			errorsx.Log(c.Close(websocket.StatusTryAgainLater, "no remote control listener attached"))
			return
		} else if err != nil {
			log.Println(errorsx.Wrap(err, "unable to relay remote control command"))
		}
	}
}
