package metaapi

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

type HTTPDHTDiagnosticsOption func(*HTTPDHTDiagnostics)

func HTTPDHTDiagnosticsOptionJWTSecret(j jwtx.SecretSource) HTTPDHTDiagnosticsOption {
	return func(t *HTTPDHTDiagnostics) { t.jwtsecret = j }
}

type HTTPDHTDiagnostics struct {
	dhtSnapshot func() (dht.ServerStats, error)
	jwtsecret   jwtx.SecretSource
}

func NewHTTPDHTDiagnostics(dhtSnapshot func() (dht.ServerStats, error), options ...HTTPDHTDiagnosticsOption) *HTTPDHTDiagnostics {
	d := &HTTPDHTDiagnostics{
		dhtSnapshot: dhtSnapshot,
		jwtsecret:   env.JWTSecret,
	}
	for _, o := range options {
		o(d)
	}
	return d
}

func (t *HTTPDHTDiagnostics) Bind(r *mux.Router) {
	r.StrictSlash(false)
	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.get))
}

func (t *HTTPDHTDiagnostics) get(w http.ResponseWriter, r *http.Request) {
	s, err := t.dhtSnapshot()
	if err != nil {
		log.Println("unable to generate snapshot", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	resp := DHTMetricsResponse{
		Dht: &DHTDiagnostics{
			GoodNodes:                             int32(s.GoodNodes),
			Nodes:                                 int32(s.Nodes),
			OutstandingTransactions:               int32(s.OutstandingTransactions),
			SuccessfulOutboundAnnouncePeerQueries: s.SuccessfulOutboundAnnouncePeerQueries,
			BadNodes:                              uint32(s.BadNodes),
			OutboundQueriesAttempted:              s.OutboundQueriesAttempted,
		},
	}
	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write diagnostics response"))
	}
}
