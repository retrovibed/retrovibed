package metaapi

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
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
)

type HTTPDaemonsOption func(*HTTPDaemons)

func HTTPDaemonsOptionJWTSecret(j jwtx.SecretSource) HTTPDaemonsOption {
	return func(t *HTTPDaemons) {
		t.jwtsecret = j
	}
}

func HTTPDaemonsOptionMDNSDiscovery(enabled bool) HTTPDaemonsOption {
	return func(t *HTTPDaemons) {
		t.mdnsDiscovery = enabled
	}
}

func HTTPDaemonsOptionMDNSLookup(l meta.MDNSLookup) HTTPDaemonsOption {
	return func(t *HTTPDaemons) {
		t.mdnsLookup = l
	}
}

func NewHTTPDaemons(q sqlx.Queryer, options ...HTTPDaemonsOption) *HTTPDaemons {
	svc := langx.Clone(HTTPDaemons{
		q:          q,
		jwtsecret:  env.JWTSecret,
		decoder:    formx.NewDecoder(),
		mdnsLookup: meta.DefaultMDNSLookup,
	}, options...)

	return &svc
}

type HTTPDaemons struct {
	q             sqlx.Queryer
	jwtsecret     jwtx.SecretSource
	decoder       *form.Decoder
	mdnsDiscovery bool
	mdnsLookup    meta.MDNSLookup
}

func (t *HTTPDaemons) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/discover").Methods(http.MethodGet).MatcherFunc(isWebsocketUpgrade).Handler(alice.New(
		httpx.GatedResponse(t.mdnsDiscovery, http.StatusServiceUnavailable),
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.discover))

	r.Path("/latest").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.latest))

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.update))

	r.Path("/{id}").Methods(http.MethodPut).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.touch))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.delete))
}

func (t *HTTPDaemons) discover(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(errorsx.Wrap(err, "failed to accept mdns discovery websocket"))
		return
	}
	defer func() {
		errorsx.Log(c.Close(websocket.StatusNormalClosure, ""))
	}()

	seq := meta.DiscoverOnce(t.q, t.mdnsLookup)
	for p := range seq.Each(r.Context()) {
		// only stream newly discovered peers — ones already known to
		// meta_daemons were touched (updated_at bumped) by the scan but
		// aren't new information for the client.
		if !p.CreatedAt.Equal(p.UpdatedAt) {
			continue
		}

		encoded, err := NewDaemonFromMetaDaemon(p)
		if err != nil {
			log.Println(errorsx.Wrap(err, "response generation failed"))
			errorsx.Log(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "internal service error"))
			return
		}

		data, err := jsonx.Marshal(encoded)
		if err != nil {
			log.Println(errorsx.Wrap(err, "unable to encode discovered peer"))
			errorsx.Log(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "internal service error"))
			return
		}

		if err := c.Write(r.Context(), websocket.MessageText, data); err != nil {
			log.Println(errorsx.Wrap(err, "unable to write discovered peer"))
			return
		}
	}

	if err := seq.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to discover mdns peers"))
		errorsx.Log(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "internal service error"))
		return
	}
}

func (t *HTTPDaemons) search(w http.ResponseWriter, r *http.Request) {
	const resplimit = 128
	var (
		err  error
		resp = DaemonSearchResponse{
			Next: &DaemonSearchRequest{
				Offset: 0,
				Limit:  resplimit,
			},
		}
	)

	if err = t.decoder.Decode(resp.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	resp.Next.Limit = numericx.Min(resp.Next.Limit, resplimit)

	query := meta.DaemonSearchBuilder().
		Offset(resp.Next.Offset * resp.Next.Limit).Limit(resp.Next.Limit)

	s := sqlx.Scan(meta.DaemonSearch(r.Context(), t.q, query))
	for v := range s.Iter() {
		var (
			encoded *Daemon
		)

		if encoded, err = NewDaemonFromMetaDaemon(v); err != nil {
			log.Println(errorsx.Wrap(err, "response generation failed"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		resp.Items = append(resp.Items, encoded)
	}

	if err = s.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "response generation failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	resp.Next.Offset += 1

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDaemons) create(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg DaemonCreateRequest
		v   meta.Daemon
	)

	if err = jsonx.UnmarshalRead(r.Body, &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if v, err = NewMetadaemonFromDaemon(msg.Daemon, meta.DaemonOptionMaybeID, meta.DaemonOptionEnsureDescription, timex.JSONSafeDecodeOption); err != nil {
		log.Println(errorsx.Wrap(err, "converting data failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = meta.DaemonInsertWithDefaults(r.Context(), t.q, v).Scan(&v); err != nil {
		log.Println(errorsx.Wrap(err, "unable to insert record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DaemonCreateResponse{
		Daemon: errorsx.Must(NewDaemonFromMetaDaemon(v)),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDaemons) latest(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		v   meta.Daemon
	)

	if err = meta.DaemonFindDefault(r.Context(), t.q).Scan(&v); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "no daemons known"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "failed to find record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DaemonLookupResponse{
		Daemon: errorsx.Must(NewDaemonFromMetaDaemon(v)),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDaemons) touch(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		v   meta.Daemon
		id  = mux.Vars(r)["id"]
	)

	if err = meta.DaemonTouch(r.Context(), t.q, id).Scan(&v); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "no daemons known"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "failed to find record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DaemonLookupResponse{
		Daemon: errorsx.Must(NewDaemonFromMetaDaemon(v)),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDaemons) download(w http.ResponseWriter, r *http.Request) { //nolint: unused
	var (
		err error
		v   meta.Daemon
		id  = mux.Vars(r)["id"]
	)

	if err = meta.DaemonDownload(r.Context(), t.q, id).Scan(&v); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "no daemons known"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "failed to find record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DaemonLookupResponse{
		Daemon: errorsx.Must(NewDaemonFromMetaDaemon(v)),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDaemons) update(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg DaemonUpdateRequest
		v   meta.Daemon
		id  = mux.Vars(r)["id"]
	)

	if err = jsonx.UnmarshalRead(r.Body, &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if v, err = NewMetadaemonFromDaemon(msg.Daemon, meta.DaemonOptionEnsureDescription, timex.JSONSafeDecodeOption); err != nil {
		log.Println(errorsx.Wrap(err, "converting data failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if v.Downloads {
		var tmp meta.Daemon
		if err = meta.DaemonDownload(r.Context(), t.q, id).Scan(&tmp); err != nil {
			log.Println(errorsx.Wrap(err, "unable to update record"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}
	}

	if err = meta.DaemonUpdateByID(r.Context(), t.q, id, v).Scan(&v); err != nil {
		log.Println(errorsx.Wrap(err, "unable to update record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DaemonUpdateResponse{
		Daemon: errorsx.Must(NewDaemonFromMetaDaemon(v)),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDaemons) delete(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		v   meta.Daemon
		id  = mux.Vars(r)["id"]
	)

	if err = meta.DaemonDeleteByID(r.Context(), t.q, id).Scan(&v); err != nil {
		log.Println(errorsx.Wrap(err, "unable to delete record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DaemonDeleteResponse{
		Daemon: errorsx.Must(NewDaemonFromMetaDaemon(v)),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
