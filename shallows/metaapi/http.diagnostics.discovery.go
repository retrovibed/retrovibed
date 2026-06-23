package metaapi

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type HTTPDiscoveryDiagnosticsOption func(*HTTPDiscoveryDiagnostics)

func HTTPDiscoveryDiagnosticsOptionJWTSecret(j jwtx.SecretSource) HTTPDiscoveryDiagnosticsOption {
	return func(t *HTTPDiscoveryDiagnostics) { t.jwtsecret = j }
}

type HTTPDiscoveryDiagnostics struct {
	db                *sql.DB
	discoverySnapshot func() (ddisc.Snapshot, error)
	jwtsecret         jwtx.SecretSource
}

func NewHTTPDiscoveryDiagnostics(db *sql.DB, discoverySnapshot func() (ddisc.Snapshot, error), options ...HTTPDiscoveryDiagnosticsOption) *HTTPDiscoveryDiagnostics {
	d := &HTTPDiscoveryDiagnostics{
		db:                db,
		discoverySnapshot: discoverySnapshot,
		jwtsecret:         env.JWTSecret,
	}
	for _, o := range options {
		o(d)
	}
	return d
}

func (t *HTTPDiscoveryDiagnostics) Bind(r *mux.Router) {
	r.StrictSlash(false)
	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.get))
}

func (t *HTTPDiscoveryDiagnostics) get(w http.ResponseWriter, r *http.Request) {
	s, err := t.discoverySnapshot()
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to read unknown hash count"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	var peers, peersddisc, peersbep51, unknown int64
	if err := tracking.PeerTotals(r.Context(), t.db).Scan(&peers, &peersddisc, &peersbep51); err != nil {
		log.Println(errorsx.Wrap(err, "unable to read peer totals"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := tracking.UnknownHashCount(r.Context(), t.db).Scan(&unknown); err != nil {
		log.Println(errorsx.Wrap(err, "unable to read unknown hash count"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	resp := DiscoveryMetricsResponse{
		Discovery: &DiscoveryDiagnostics{
			Enabled:        s.Enabled,
			Ratio:          s.Ratio,
			Partitions:     s.Partitions,
			Workloads:      s.Workloads,
			LocalPartition: s.LocalPartition,
			Peers:          uint64(peers),
			PeersDdisc:     uint64(peersddisc),
			PeersBep51:     uint64(peersbep51),
			UnknownHashes:  uint64(unknown),
		},
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write diagnostics response"))
	}
}
