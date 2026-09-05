package metaapi

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/iox"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/wireguardx"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

type HTTPWireguardOption func(*HTTPWireguard)

func HTTPWireguardOptionJWTSecret(j jwtx.SecretSource) HTTPWireguardOption {
	return func(t *HTTPWireguard) {
		t.jwtsecret = j
	}
}

func NewHTTPWireguard(dir string, q sqlx.Queryer, options ...HTTPWireguardOption) *HTTPWireguard {
	svc := langx.Clone(HTTPWireguard{
		q:         q,
		dir:       fsx.DirVirtual(dir),
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPWireguard struct {
	dir       fsx.Virtual
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPWireguard) Bind(r *mux.Router) {
	r.StrictSlash(false)
	r.Use(httpx.RouteInvoked)
	// r.Use(httpx.DebugRequest)

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

	r.Path("/current").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.current))

	r.Path("/{id}").Methods(http.MethodPatch).Handler(alice.New(
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

func (t *HTTPWireguard) search(w http.ResponseWriter, r *http.Request) {
	const resplimit = 128
	var (
		err  error
		resp = WireguardSearchResponse{
			Next: &WireguardSearchRequest{
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

	query := meta.WireguardSearchBuilder().
		Where(
			squirrel.And{
				squirrel.Expr("'t'"),
				lucenex.Query(duckdbx.NewLucene(), resp.Next.Query, lucenex.WithDefaultField("description")),
			},
		).
		Offset(resp.Next.Offset * resp.Next.Limit).Limit(resp.Next.Limit)

	q := sqlx.Scan(meta.WireguardSearch(r.Context(), t.q, query))
	for p := range q.Iter() {
		var (
			encoded *Wireguard
		)

		if encoded, err = NewWireguardFromMetaWireguard(p); err != nil {
			log.Println(errorsx.Wrap(err, "wireguard generation failed"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		resp.Items = append(resp.Items, encoded)
	}

	if err := q.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to search wireguard"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPWireguard) create(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		f      multipart.File
		fh     *multipart.FileHeader
		buf    [bytesx.MiB]byte
		copied = &iox.Copied{Result: new(uint64)}
		mhash  = md5.New()
		wg     meta.Wireguard
	)

	if f, fh, err = r.FormFile("content"); err != nil {
		log.Println(errorsx.Wrap(err, "content parameter required"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	defer f.Close()

	tmp, err := fsx.CreateTemp(t.dir, "retrovibed.upload.*")
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to create temporary file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}
	defer func() {
		if err == nil {
			return
		}

		log.Println("failure receiving upload, removing attempt", err)
		errorsx.Log(errorsx.Wrap(os.Remove(tmp.Name()), "unable to remove tmp"))
	}()
	defer tmp.Close()

	if _, err = io.CopyBuffer(io.MultiWriter(tmp, mhash, copied), f, buf[:]); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create temporary file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	raw, err := iox.String(tmp)
	if err != nil {
		log.Println(errorsx.Wrap(err, "failed to read configuration"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if _, err = wireguardx.FromWgQuick(string(raw), "retrovibed"); err != nil {
		log.Println(errorsx.Wrap(err, "failed to parse config"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	wg = meta.NewWireguard(md5x.FormatUUID(mhash), meta.WireguardOptionDescription(fh.Filename))

	if err = meta.WireguardInsertWithDefaults(r.Context(), t.q, wg).Scan(&wg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to create temporary file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = os.Rename(tmp.Name(), t.dir.Path(wg.ID)); err != nil {
		log.Println(errorsx.Wrap(err, "failed to rename upload"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	encoded, err := NewWireguardFromMetaWireguard(wg)
	if err != nil {
		log.Println(errorsx.Wrap(err, "failed to encode"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &WireguardUploadResponse{
		Wireguard: encoded,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPWireguard) current(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		wg       meta.Wireguard
		realpath = errorsx.Zero(filepath.EvalSymlinks(t.dir.Path(wireguardx.Current)))
	)

	if err = meta.WireguardCurrent(r.Context(), t.q).Scan(&wg); errors.Is(err, sql.ErrNoRows) {
		log.Println("no wireguard configuration activated")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrapf(err, "failed to read configuration: %s", realpath))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	encoded, err := os.ReadFile(t.dir.Path(wg.ID))
	if err != nil {
		log.Println(errorsx.Wrapf(err, "failed to read configuration: %s", realpath))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	_wg, err := NewWireguardFromMetaWireguard(wg)
	if err != nil {
		log.Println(errorsx.Wrap(err, "failed to encode"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if _, err = wireguardx.FromWgQuick(string(encoded), "retrovibed"); err != nil {
		log.Println(errorsx.Wrap(err, "failed to parse config"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &WireguardCurrentResponse{
		Wireguard: _wg,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPWireguard) update(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		id  = mux.Vars(r)["id"]
		wg  meta.Wireguard
		msg WireguardUpdateRequest
	)

	if err = jsonx.UnmarshalRead(r.Body, &msg); err != nil {
		log.Println("unable to decode request", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if wg, err = NewMetaWireguardFromWireguard(msg.Wireguard, timex.JSONSafeDecodeOption); err != nil {
		log.Println(errorsx.Wrap(err, "failed to decode update"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	wg.ID = id

	if err = meta.WireguardUpdate(r.Context(), t.q, wg).Scan(&wg); err != nil {
		log.Println(errorsx.Wrap(err, "failed to update record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	encoded, err := NewWireguardFromMetaWireguard(wg)
	if err != nil {
		log.Println(errorsx.Wrap(err, "failed to encode"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &WireguardDeleteResponse{
		Wireguard: encoded,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPWireguard) touch(w http.ResponseWriter, r *http.Request) {
	// switches the current config
	var (
		err error
		id  = mux.Vars(r)["id"]
		wg  meta.Wireguard
	)

	if err = os.RemoveAll(t.dir.Path(wireguardx.Current)); fsx.IgnoreIsNotExist(err) != nil {
		log.Println(errorsx.Wrap(err, "failed to remove old config"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = os.Symlink(t.dir.Path(id), t.dir.Path(wireguardx.Current)); err != nil {
		log.Println(errorsx.Wrap(err, "failed to symlink"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = meta.WireguardTouch(r.Context(), t.q, id).Scan(&wg); err != nil {
		log.Println(errorsx.Wrap(err, "failed to update"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	// this shouldnt be necessary, we *should* be able to use a CTE or UPDATE followed by a select.
	// unfortunately duckdb does not support either case.
	if err = meta.WireguardCurrent(r.Context(), t.q).Scan(&wg); errors.Is(err, sql.ErrNoRows) {
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "failed to read new"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	encoded, err := NewWireguardFromMetaWireguard(wg)
	if err != nil {
		log.Println(errorsx.Wrap(err, "failed to encode"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &WireguardTouchResponse{
		Wireguard: encoded,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPWireguard) delete(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		id  = mux.Vars(r)["id"]
		wg  meta.Wireguard
	)

	if err = os.Remove(t.dir.Path(id)); fsx.IgnoreIsNotExist(err) != nil {
		log.Println(errorsx.Wrap(err, "failed to remove config"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = meta.WireguardDeleteByID(r.Context(), t.q, id).Scan(&wg); err != nil {
		log.Println(errorsx.Wrap(err, "failed to remove record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	encoded, err := NewWireguardFromMetaWireguard(wg)
	if err != nil {
		log.Println(errorsx.Wrap(err, "failed to encode"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &WireguardDeleteResponse{
		Wireguard: encoded,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
