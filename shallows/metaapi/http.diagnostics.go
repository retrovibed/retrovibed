package metaapi

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/netmonx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/wireguardx"
)

type HTTPDiagnosticsOption func(*HTTPDiagnostics)

func HTTPDiagnosticsOptionJWTSecret(j jwtx.SecretSource) HTTPDiagnosticsOption {
	return func(t *HTTPDiagnostics) { t.jwtsecret = j }
}

type HTTPDiagnostics struct {
	wgSnapshot func() (wireguardx.Statistics, error)
	jwtsecret  jwtx.SecretSource
}

func NewHTTPDiagnostics(wgSnapshot func() (wireguardx.Statistics, error), options ...HTTPDiagnosticsOption) *HTTPDiagnostics {
	d := &HTTPDiagnostics{
		wgSnapshot: wgSnapshot,
		jwtsecret:  env.JWTSecret,
	}
	for _, o := range options {
		o(d)
	}
	return d
}

func (t *HTTPDiagnostics) Bind(r *mux.Router) {
	r.StrictSlash(false)
	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.get))
}

func (t *HTTPDiagnostics) get(w http.ResponseWriter, r *http.Request) {
	resp := NetworkMetricsResponse{
		Wireguard: t.wireguardMetrics(),
		Network:   t.networkMetrics(),
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write diagnostics response"))
	}
}

func (t *HTTPDiagnostics) wireguardMetrics() *WireguardDiagnostics {
	s, err := t.wgSnapshot()
	if err != nil || s.PeerKey == "" {
		return &WireguardDiagnostics{Status: "Inactive"}
	}

	status := "Healthy"
	if s.LastHandshakeSec == 0 {
		status = "No Handshake"
	} else {
		elapsed := time.Now().Unix() - s.LastHandshakeSec
		switch {
		case s.TXBytes > 0 && s.RXBytes == 0:
			status = "Unbalanced Pipe"
		case elapsed > 180:
			status = "Stale Handshake"
		}
	}

	return &WireguardDiagnostics{
		PeerKey:           s.PeerKey,
		KeepaliveInterval: s.KeepaliveInterval,
		TxBytes:           s.TXBytes,
		RxBytes:           s.RXBytes,
		LastHandshakeSec:  s.LastHandshakeSec,
		Status:            status,
	}
}

func (t *HTTPDiagnostics) networkMetrics() *Network {
	mon := netmonx.Global()
	if mon == nil {
		return &Network{}
	}

	state := mon.Current()
	if state == nil {
		return &Network{}
	}

	n := &Network{
		HaveV4:           state.HaveV4,
		HaveV6:           state.HaveV6,
		DefaultInterface: state.DefaultRouteInterface,
	}

	for _, nd := range state.Networks {
		iface := NetworkInterface{
			Name:    nd.Name,
			Ip:      nd.IP.String(),
			Metered: nd.Metered,
		}
		n.Interfaces = append(n.Interfaces, &iface)
	}

	return n
}
