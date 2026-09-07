package metaapi

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type HTTPTorrentDiagnosticsOption func(*HTTPTorrentDiagnostics)

func HTTPTorrentDiagnosticsOptionJWTSecret(j jwtx.SecretSource) HTTPTorrentDiagnosticsOption {
	return func(t *HTTPTorrentDiagnostics) { t.jwtsecret = j }
}

type HTTPTorrentDiagnostics struct {
	db        *sql.DB
	jwtsecret jwtx.SecretSource
}

func NewHTTPTorrentDiagnostics(db *sql.DB, options ...HTTPTorrentDiagnosticsOption) *HTTPTorrentDiagnostics {
	d := &HTTPTorrentDiagnostics{
		db:        db,
		jwtsecret: env.JWTSecret,
	}
	for _, o := range options {
		o(d)
	}
	return d
}

func (t *HTTPTorrentDiagnostics) Bind(r *mux.Router) {
	r.StrictSlash(false)
	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.get))
}

func (t *HTTPTorrentDiagnostics) get(w http.ResponseWriter, r *http.Request) {
	var (
		total, seeding, bytes, downloaded, available, uploaded, peers int64
	)

	if err := tracking.MetadataDiagnostics(r.Context(), t.db).Scan(&total, &seeding, &bytes, &downloaded, &available, &uploaded, &peers); err != nil {
		log.Println(errorsx.Wrap(err, "unable to read torrent metadata totals"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	resp := TorrentMetricsResponse{
		Torrent: &TorrentDiagnostics{
			Total:      uint64(total),
			Seeding:    uint64(seeding),
			Bytes:      uint64(bytes),
			Downloaded: uint64(downloaded),
			Uploaded:   uint64(uploaded),
			Peers:      uint64(peers),
		},
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write diagnostics response"))
	}
}
