package ddiscapi

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"net/netip"

	"github.com/Masterminds/squirrel"
	"github.com/coder/websocket"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/websocketx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type HTTPDiscoveryOption func(*HTTPDiscovery)

func HTTPDiscoveryOptionJWTSecret(j jwtx.SecretSource) HTTPDiscoveryOption {
	return func(t *HTTPDiscovery) {
		t.jwtsecret = j
	}
}

func NewHTTPDiscovery(q sqlx.Queryer, plugins searchplugin.T, options ...HTTPDiscoveryOption) *HTTPDiscovery {
	svc := langx.Clone(HTTPDiscovery{
		q:         q,
		plugins:   plugins,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPDiscovery struct {
	q         sqlx.Queryer
	plugins   searchplugin.T
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPDiscovery) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/locate").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.websocket))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, AuthzPermPeer),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		metaapi.AuthzTokenHTTP(t.jwtsecret, AuthzPermPeer),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))

	r.Path("/{id}").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool1024(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.download))

}

func (t *HTTPDiscovery) download(w http.ResponseWriter, r *http.Request) {
	var (
		disc ddisc.Discovered
		md   tracking.Metadata
		id   = mux.Vars(r)["id"]
	)

	if err := ddisc.DiscoveredFindByID(r.Context(), t.q, id).Scan(&disc); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find discovered"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find discovered"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := tracking.MetadataFindByInfohash(r.Context(), t.q, hex.EncodeToString(disc.Infohash)).Scan(&md); err != nil {
		log.Println(errorsx.Wrap(err, "unable to find tracking data"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := tracking.MetadataAutoDownloadByID(r.Context(), t.q, md.ID).Scan(&md); err != nil {
		log.Println(errorsx.Wrap(err, "unable to autodownload tracking"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), DiscoveryDownloadResponse{
		Discovery: NewDiscoveryFromDiscovered(disc),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to find response"))
		return
	}
}

// websocket runs a live, synchronous-only (no DHT) discovery search and
// streams each ranked candidate to the client as it's found; the very last
// message sent before the connection closes is always the best candidate
// found (or nothing was found at all, if no messages were sent). Fully
// ephemeral - unlike POST /, nothing is queued for a background daemon to
// keep searching after the socket closes.
func (t *HTTPDiscovery) websocket(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("failed to accept websocket", err)
		return
	}
	defer func() {
		errorsx.Log(c.Close(websocket.StatusNormalClosure, ""))
	}()

	var req Locate
	if err := t.decoder.Decode(&req, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(errorsx.Wrap(c.Close(websocketx.PrivateStatus(http.StatusBadRequest), "invalid request"), "failed to close websocket"))
		return
	}

	if req.Query == "" && req.KnownMediaId == "" {
		log.Println("locate socket requires a query or known_media_id")
		errorsx.Log(errorsx.Wrap(c.Close(websocketx.PrivateStatus(http.StatusBadRequest), "query or known_media_id required"), "failed to close websocket"))
		return
	}

	ctx := c.CloseRead(context.Background())

	discreq := ddisc.DiscoverRequest{
		KnownMediaID: req.KnownMediaId,
		Title:        req.Query,
		Mimetypes:    ddisc.Category(req.Mimetype),
		Adult:        req.Adult,
	}

	strategies := ddisc.SyncStrategies(t.q, t.plugins, discreq.KnownMediaID)
	seq := ddisc.Discover(ctx, ddisc.DefaultPolicy(), discreq, strategies...)

	var (
		buf   = bytes.NewBuffer(nil)
		enc   = json.NewEncoder(buf)
		best  = ddisc.Worst()
		found bool
	)

	write := func(d ddisc.Discovered) error {
		buf.Reset()
		if err := enc.Encode(NewDiscoveryFromDiscovered(d)); err != nil {
			return errorsx.Wrap(err, "unable to encode discovery")
		}
		return c.Write(ctx, websocket.MessageBinary, buf.Bytes())
	}

	for d := range seq.Each(ctx) {
		if d.PolicyRejection != "" {
			continue
		}

		if err := write(d); err != nil {
			log.Println(errorsx.Wrap(err, "unable to write discovery"))
			return
		}

		if ddisc.Compare(d, best) < 0 {
			best, found = d, true
		}
	}

	if err := seq.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "locate search failed"))
		errorsx.Log(errorsx.Wrap(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "search failed"), "failed to close websocket"))
		return
	}

	if found {
		if err := write(best); err != nil {
			log.Println(errorsx.Wrap(err, "unable to write best discovery"))
			return
		}
	}
}

func (t *HTTPDiscovery) search(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp = DiscoverySearchResponse{
			Next: &DiscoverySearchRequest{
				Offset:      0,
				Limit:       100,
				NextCheck:   meta.NewDateRange(timex.NewRangeEverything()),
				AttemptsMax: math.MaxInt64,
			},
		}
	)

	if err = t.decoder.Decode(resp.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	resp.Next.Limit = numericx.Min(resp.Next.Limit, 1024)

	query := tracking.UnknownSearchBuilder().
		Offset(resp.Next.Offset * resp.Next.Limit).Limit(resp.Next.Limit).
		Where(squirrel.And{
			tracking.UnknownHashQueryAttemptsRange(resp.Next.AttemptsMin, resp.Next.AttemptsMax),
			tracking.UnknownHashQueryByIDs(resp.Next.Id...),
			tracking.UnknownHashQueryNextCheck(meta.TimexRange(resp.Next.NextCheck, timex.NewRangeEverything())),
		})

	q := sqlx.Scan(tracking.UnknownSearch(r.Context(), t.q, query))
	for uh := range q.Iter() {
		resp.Items = append(resp.Items, NewDiscoveryFromTrackingUnknownHash(uh))
	}

	if err = q.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "discovery generation failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovery) create(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg DiscoveryCreateRequest
		uh  tracking.UnknownHash
	)

	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	uh = tracking.NewUnknownHash(
		int160.FromBytesOrZero(msg.Discovery.GetInfohash()),
		func(uh *tracking.UnknownHash) { uh.IP = netip.IPv6Unspecified() },
	)

	if err = tracking.UnknownHashInsertWithDefaults(r.Context(), t.q, uh).Scan(&uh); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create discovery entry"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DiscoveryCreateResponse{
		Discovery: NewDiscoveryFromTrackingUnknownHash(uh),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovery) delete(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		vars = mux.Vars(r)
		uh   tracking.UnknownHash
	)

	if err := tracking.UnknownHashDeleteByID(r.Context(), t.q, vars["id"]).Scan(&uh); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DiscoveryDeleteResponse{
		Discovery: NewDiscoveryFromTrackingUnknownHash(uh),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
