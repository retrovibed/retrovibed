package media

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/netip"
	"os"
	"sync/atomic"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/anacrolix/missinggo/pubsub"
	"github.com/coder/websocket"
	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/justinas/alice"
	rootenv "github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
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
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/internal/websocketx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type HTTPDiscoveredOption func(*HTTPDiscovered)

func HTTPDiscoveredOptionJWTSecret(j jwtx.SecretSource) HTTPDiscoveredOption {
	return func(t *HTTPDiscovered) {
		t.jwtsecret = j
	}
}

func HTTPDiscoveredOptionRootStorage(vfs fsx.Virtual) HTTPDiscoveredOption {
	return func(t *HTTPDiscovered) {
		t.rootstorage = vfs
	}
}

func HTTPDiscoveredOptionQueryCleaner(mc library.QueryCleaner) HTTPDiscoveredOption {
	return func(t *HTTPDiscovered) {
		t.mediacleaner = mc
	}
}

func NewHTTPDiscovered(q sqlx.Queryer, d *atomic.Pointer[torrent.Client], c storage.ClientImpl, pub *asyncx.Wakeup, options ...HTTPDiscoveredOption) *HTTPDiscovered {
	svc := langx.Clone(HTTPDiscovered{
		q:            q,
		d:            d,
		c:            c,
		pub:          pub,
		jwtsecret:    env.JWTSecret,
		decoder:      formx.NewDecoder(),
		rootstorage:  fsx.DirVirtual(os.TempDir()),
		fts:          duckdbx.NewLucene(),
		mediacleaner: library.QueryCleanerNoop(),
	}, options...)

	errorsx.Log(errorsx.Wrap(fsx.MkDirs(0700, svc.rootstorage.Path(rootenv.TorrentDirName), svc.rootstorage.Path(rootenv.MediaDirName)), "failed to prepare directories"))

	return &svc
}

type HTTPDiscovered struct {
	q            sqlx.Queryer
	d            *atomic.Pointer[torrent.Client]
	c            storage.ClientImpl
	pub          *asyncx.Wakeup
	jwtsecret    jwtx.SecretSource
	decoder      *form.Decoder
	rootstorage  fsx.Virtual
	fts          lucenex.Driver
	mediacleaner library.QueryCleaner
}

func (t *HTTPDiscovered) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.TimeoutRollingWrite(3*time.Second),
	).ThenFunc(t.upload))

	r.Path("/publish").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.TimeoutRollingWrite(3*time.Second),
	).ThenFunc(t.publish))

	r.Path("/magnet").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.magnet))

	r.Path("/available").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/downloading").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.downloading))

	r.Path("/{id}/metadatasync").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.metadatasync))

	r.Path("/{id}/socket").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.websocket))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.metadata))

	r.Path("/{id}").Methods(http.MethodPut).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.update))

	r.Path("/{id}").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.download))

	r.Path("/{id}/tune").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.tune))

	r.Path("/{id}/pause").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.pause))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.reset))
}

func (t *HTTPDiscovered) magnet(w http.ResponseWriter, r *http.Request) {
	var (
		msg MagnetCreateRequest
		dl  torrent.Torrent
	)

	if err := jsonx.UnmarshalRead(r.Body, &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to parse magnet link request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	m, err := metainfo.ParseMagnetURI(msg.Uri)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to parse magnet link"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	lmd := tracking.NewMetadata(
		new(int160.FromByteArray(m.InfoHash)),
		tracking.MetadataOptionFromMagnet(&m),
		tracking.MetadataOptionTrackers(m.Trackers...),
		tracking.MetadataOptionAutoDescription,
		tracking.MetadataOptionAutoHidden,
	)

	if err = tracking.MetadataInsertWithDefaults(r.Context(), t.q, lmd).Scan(&lmd); err != nil {
		log.Println(errorsx.Wrap(err, "unable to record metadata record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	metadata, err := torrent.New(metainfo.Hash(lmd.Infohash), torrent.OptionStorage(t.c), torrentx.OptionTracker(lmd.Tracker), torrent.OptionPublicTrackers(lmd.Private, tracking.PublicTrackers()...))
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to create torrent from metadata %s", lmd.ID))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if dl, _, err = t.d.Load().Start(metadata); err != nil {
		log.Println(errorsx.Wrap(err, "unable to start download"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	go func() {
		var (
			mhash = md5.New()
		)

		errorsx.Log(tracking.DownloadInto(context.Background(), t.q, t.rootstorage, t.mediacleaner, &lmd, dl, mhash, t.pub))
	}()

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &MagnetCreateResponse{
		Download: new(
			langx.Clone(
				Download{},
				DownloadOptionFromTorrentMetadata(langx.Clone(lmd, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) publish(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		decoded PublishedUploadRequest
		copied  = &iox.Copied{Result: new(uint64)}
		mhash   = md5.New()
		buf     bytes.Buffer
	)

	reader, err := r.MultipartReader()
	if err != nil {
		log.Println(errorsx.Wrap(err, "failed creating a multipart form"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	torfile, err := reader.NextPart()
	if err != nil {
		log.Println(errorsx.Wrap(err, "missing torrent file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	defer torfile.Close()

	// we limit torrent files to 128 MB for sanity sake.
	meta, err := metainfo.Load(io.TeeReader(io.LimitReader(torfile, 128*bytesx.MiB), io.MultiWriter(&buf, mhash, copied)))
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to read torrent file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusRequestEntityTooLarge))
		return
	}

	info, err := meta.UnmarshalInfo()
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to read torrent info"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	torimp, err := t.c.OpenTorrent(&info, meta.ID())
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to open torrent storage"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	contents, err := reader.NextPart()
	if err != nil {
		log.Println(errorsx.Wrap(err, "missing torrent contents"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	defer contents.Close()

	n, err := io.Copy(io.NewOffsetWriter(torimp, 0), contents)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to copy torrent to storage"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if n != info.TotalLength() {
		log.Println(errorsx.Errorf("failed to upload all content %d != %d", n, info.TotalLength()))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err := os.WriteFile(t.rootstorage.Path(rootenv.TorrentDirName, fmt.Sprintf("%s.torrent", meta.ID().String())), buf.Bytes(), 0600); err != nil {
		log.Println(errorsx.Errorf("failed to write torrent file: %s", t.rootstorage.Path(rootenv.TorrentDirName, fmt.Sprintf("%s.torrent", meta.ID().String()))))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	metadata, err := reader.NextPart()
	if err != nil {
		log.Println(errorsx.Wrap(err, "missing metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	defer metadata.Close()

	if err = jsonx.UnmarshalRead(metadata, &decoded); err != nil {
		log.Println(errorsx.Wrap(err, "invalid metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	lmd := tracking.NewMetadata(
		new(int160.FromByteArray(meta.HashInfoBytes())),
		tracking.MetadataOptionFromInfo(&info),
		tracking.MetadataOptionAvailable(n),
		tracking.MetadataOptionTrackers(slicesx.Flatten(meta.UpvertedAnnounceList()...)...),
		tracking.MetadataOptionEntropySeed(meta.ID().Bytes(), uuid.FromStringOrNil(decoded.Entropy).Bytes()),
		tracking.MetadataOptionMimetype(decoded.Mimetype),
		tracking.MetadataOptionExpiresAtTTL(time.Duration(decoded.Ttl)),
		tracking.MetadataOptionCompleted,
		tracking.MetadataOptionAutoSeeding,
		tracking.MetadataOptionAutoDescription,
		tracking.MetadataOptionAutoHidden,
	)

	if err = tracking.MetadataInsertWithDefaults(r.Context(), t.q, lmd).Scan(&lmd); err != nil {
		log.Println(errorsx.Wrap(err, "unable to record metadata record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &PublishedUploadResponse{
		Published: new(
			langx.Clone(
				Published{},
				PublishedOptionFromTorrentMetadata(langx.Clone(lmd, timex.JSONSafeEncodeOption)),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) upload(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		dl     torrent.Torrent
		f      multipart.File
		fh     *multipart.FileHeader
		buf    [bytesx.MiB]byte
		copied = &iox.Copied{Result: new(uint64)}
		mhash  = md5.New()
	)

	if f, fh, err = r.FormFile("content"); err != nil {
		log.Println(errorsx.Wrap(err, "content parameter required"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	defer f.Close()

	tmp, err := fsx.CreateTemp(t.rootstorage, "retrovibed.upload.*")
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

	meta, err := metainfo.LoadFromFile(tmp.Name())
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to read temporary file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if info, err := meta.UnmarshalInfo(); err == nil && !langx.Autoderef(info.Private) {
		meta.AnnounceList = append(meta.AnnounceList, tracking.PublicTrackers())
	}

	info, err := meta.UnmarshalInfo()
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to read temporary file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}
	info.Name = fh.Filename

	lmd := tracking.NewMetadata(
		new(int160.FromByteArray(meta.HashInfoBytes())),
		tracking.MetadataOptionFromInfo(&info),
		tracking.MetadataOptionTrackers(slicesx.Flatten(meta.UpvertedAnnounceList()...)...),
		tracking.MetadataOptionAutoDescription,
		tracking.MetadataOptionAutoHidden,
	)

	if err = tracking.MetadataInsertWithDefaults(r.Context(), t.q, lmd).Scan(&lmd); err != nil {
		log.Println(errorsx.Wrap(err, "unable to record metadata record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = os.Rename(tmp.Name(), t.rootstorage.Path(rootenv.TorrentDirName, fmt.Sprintf("%s.torrent", metainfo.Hash(lmd.Infohash).String()))); err != nil {
		log.Println(errorsx.Wrap(err, "unable to failed to record torrent file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	metadata, err := torrent.New(
		metainfo.Hash(lmd.Infohash),
		torrent.OptionStorage(t.c),
		torrent.OptionNodes(meta.NodeList()...),
		torrentx.OptionTracker(lmd.Tracker),
		torrent.OptionWebseeds(meta.UrlList),
		torrent.OptionPublicTrackers(lmd.Private, tracking.PublicTrackers()...),
	)
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to create torrent from metadata %s", lmd.ID))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	go func() {
		var (
			mhash = md5.New()
		)

		if dl, _, err = t.d.Load().Start(metadata); err != nil {
			log.Println(errorsx.Wrap(err, "unable to start download"))
			return
		}

		errorsx.Log(tracking.DownloadInto(context.Background(), t.q, t.rootstorage, t.mediacleaner, &lmd, dl, mhash, t.pub))
	}()

	if err := tracking.MetadataDownloadByID(r.Context(), t.q, lmd.ID).Scan(&lmd); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to track download"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &MediaUploadResponse{
		Media: new(
			langx.Clone(
				Media{},
				MediaOptionFromTorrentMetadata(lmd),
			),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) pause(w http.ResponseWriter, r *http.Request) {
	var (
		md tracking.Metadata
		id = mux.Vars(r)["id"]
	)

	if err := tracking.MetadataFindByID(r.Context(), t.q, id).Scan(&md); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	metadata, err := torrent.New(metainfo.Hash(md.Infohash), torrent.OptionStorage(t.c))
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to create metadata from metadata %s", md.ID))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = t.d.Load().Stop(metadata); err != nil {
		log.Println(errorsx.Wrap(err, "unable to stop download"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = tracking.MetadataPausedByID(r.Context(), t.q, id).Scan(&md); err != nil {
		log.Println(errorsx.Wrap(err, "unable to pause metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &DownloadPauseResponse{
		Download: new(
			langx.Clone(
				Download{},
				DownloadOptionFromTorrentMetadata(langx.Clone(md, timex.JSONSafeEncodeOption))),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) reset(w http.ResponseWriter, r *http.Request) {
	var (
		md tracking.Metadata
		id = mux.Vars(r)["id"]
	)

	if err := tracking.MetadataFindByID(r.Context(), t.q, id).Scan(&md); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	metadata, err := torrent.New(metainfo.Hash(md.Infohash), torrent.OptionStorage(t.c))
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to create metadata from metadata %s", md.ID))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = t.d.Load().Stop(metadata); err != nil {
		log.Println(errorsx.Wrap(err, "unable to stop download"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := tracking.Reset(r.Context(), t.q, t.rootstorage, &md); err != nil {
		log.Println(errorsx.Wrap(err, "unable to delete content"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &DownloadDeleteResponse{
		Download: new(
			langx.Clone(
				Download{},
				DownloadOptionFromTorrentMetadata(langx.Clone(md, timex.JSONSafeEncodeOption))),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) tune(w http.ResponseWriter, r *http.Request) {
	var (
		md    tracking.Metadata
		id    = mux.Vars(r)["id"]
		cl    *torrent.Client
		peers torrent.Tuner
		msg   DownloadTuneRequest
	)

	if cl = t.d.Load(); cl == nil {
		log.Println("torrent client unavailable")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusServiceUnavailable))
		return
	}

	if err := jsonx.UnmarshalRead(r.Body, &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err := tracking.MetadataDownloadByID(r.Context(), t.q, id).Scan(&md); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to track download"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	peersv := make([]torrent.Peer, 0, len(msg.Peers))
	for _, v := range msg.Peers {
		addr, err := netip.ParseAddrPort(v)
		if err != nil {
			log.Println(errorsx.Wrapf(err, "unable to parse peer address: %s", v))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}
		peersv = append(peersv, torrent.NewPeer(int160.Zero(), addr))
	}
	peers = torrent.TunePeers(peersv...)

	go func(meta tracking.Metadata, opts ...torrent.Tuner) {
		md, err := torrent.New(
			metainfo.Hash(meta.Infohash),
			torrent.OptionStorage(t.c),
			torrent.OptionTrackers(meta.Tracker),
			torrent.OptionPublicTrackers(meta.Private, tracking.PublicTrackers()...),
		)
		if err != nil {
			log.Println(errorsx.Wrapf(err, "unable to create metadata from metadata %s", meta.ID))
			return
		}

		if err = cl.Modify(r.Context(), md, opts...); err != nil {
			log.Println(errorsx.Wrap(err, "unable to start download"))
			return
		}
	}(md, peers, torrent.TuneNewConns)

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &DownloadTuneResponse{}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) download(w http.ResponseWriter, r *http.Request) {
	var (
		meta tracking.Metadata
		id   = mux.Vars(r)["id"]
		cl   *torrent.Client
	)

	if cl = t.d.Load(); cl == nil {
		log.Println("torrent client unavailable")
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusServiceUnavailable))
		return
	}

	if err := tracking.MetadataDownloadByID(r.Context(), t.q, id).Scan(&meta); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to track download"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	go func(meta tracking.Metadata) {
		var (
			mhash = md5.New()
			dl    torrent.Torrent
			added bool
		)

		metadata, err := torrent.New(
			metainfo.Hash(meta.Infohash),
			torrent.OptionStorage(t.c),
			torrent.OptionTrackers(meta.Tracker),
			torrent.OptionPublicTrackers(meta.Private, tracking.PublicTrackers()...),
		)
		if err != nil {
			log.Println(errorsx.Wrapf(err, "unable to create metadata from metadata %s", meta.ID))
			return
		}

		if dl, added, err = cl.Start(metadata); err != nil {
			log.Println(errorsx.Wrap(err, "unable to start download"))
			return
		} else if !added {
			log.Println("torrent", meta.ID, meta.Description, "already running")
			return
		}

		errorsx.Log(tracking.DownloadInto(context.Background(), t.q, t.rootstorage, t.mediacleaner, &meta, dl, mhash, t.pub))
	}(meta)

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &DownloadBeginResponse{
		Download: new(
			langx.Clone(
				Download{},
				DownloadOptionFromTorrentMetadata(langx.Clone(meta, timex.JSONSafeEncodeOption))),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) websocket(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("failed to accept websocket", err)
		return
	}
	defer func() {
		errorsx.Log(c.Close(websocket.StatusNormalClosure, ""))
	}()

	var (
		md  tracking.Metadata
		sub pubsub.Subscription
		dl  torrent.Torrent
		buf = bytes.NewBuffer(nil)
		enc = json.NewEncoder(buf)
		id  = mux.Vars(r)["id"]
	)

	if err := tracking.MetadataDownloadByID(r.Context(), t.q, id).Scan(&md); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(errorsx.Wrap(c.Close(websocketx.PrivateStatus(http.StatusNotFound), "not found"), "failed to close websocket"))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to track download"))
		errorsx.Log(errorsx.Wrap(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "internal service error"), "failed to close websocket"))
		return
	}

	ctx := c.CloseRead(context.Background())

	metadata, err := torrent.New(
		metainfo.Hash(md.Infohash),
		torrent.OptionStorage(t.c),
		torrent.OptionTrackers(md.Tracker),
		torrent.OptionPublicTrackers(md.Private, tracking.PublicTrackers()...),
	)

	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to create metadata from %s", md.ID))
		errorsx.Log(errorsx.Wrap(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "internal service error"), "failed to close websocket"))
		return
	}

	tclient := t.d.Load()
	if tclient == nil {
		log.Println("torrent client not initialized")
		errorsx.Log(errorsx.Wrap(c.Close(websocketx.PrivateStatus(http.StatusServiceUnavailable), "torrent client not ready"), "failed to close websocket"))
		return
	}

	if dl, _, err = tclient.Start(metadata, torrent.TuneSubscribe(&sub)); err != nil {
		log.Println(errorsx.Wrap(err, "unable to track download"))
		errorsx.Log(errorsx.Wrap(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "internal service error"), "failed to close websocket"))
		return
	}

	genmsg := func(_ctx context.Context) error {
		buf.Reset()
		msg := new(
			langx.Clone(
				Download{},
				DownloadOptionFromTorrentMetadata(langx.Clone(md, timex.JSONSafeEncodeOption)),
				DownloadOptionFromTorrent(dl),
			),
		)

		if err = enc.Encode(msg); err != nil {
			return errorsx.Wrap(err, "unable to encode status")
		}

		if err = c.Write(_ctx, websocket.MessageBinary, buf.Bytes()); err != nil {
			return errorsx.Wrap(err, "unable to write msg")
		}

		return nil
	}

	if err := genmsg(ctx); err != nil {
		log.Println(errorsx.Wrap(context.Cause(ctx), "websocket failed"))
		errorsx.Log(errorsx.Wrap(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "internal service error"), "failed to close websocket"))
		return
	}

	for {
		select {
		case <-sub.Values:
			if err := genmsg(ctx); err != nil {
				log.Println(err)
				errorsx.Log(errorsx.Wrap(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "internal service error"), "failed to close websocket"))
				return
			}
		case <-ctx.Done():
			log.Println(errorsx.Wrap(context.Cause(ctx), "websocket failed"))
			errorsx.Log(errorsx.Wrap(c.Close(websocketx.PrivateStatus(http.StatusInternalServerError), "internal service error"), "failed to close websocket"))
			return
		}
	}
}

func (t *HTTPDiscovered) update(w http.ResponseWriter, r *http.Request) {
	var (
		msg DownloadUpdateRequest
		md  tracking.Metadata
		id  = mux.Vars(r)["id"]
	)

	if err := jsonx.UnmarshalRead(r.Body, &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode update request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	var err error
	if md, err = MetadataFromDownload(msg.Download); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode update request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	md.ID = id

	if err := tracking.MetadataInsertOrUpdate(r.Context(), t.q, md).Scan(&md); err != nil {
		log.Println(errorsx.Wrap(err, "unable to update metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &DownloadUpdateResponse{
		Download: new(
			langx.Clone(
				Download{},
				DownloadOptionFromTorrentMetadata(langx.Clone(md, timex.JSONSafeEncodeOption))),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) metadata(w http.ResponseWriter, r *http.Request) {
	var (
		meta tracking.Metadata
		id   = mux.Vars(r)["id"]
	)

	if err := tracking.MetadataFindByID(r.Context(), t.q, id).Scan(&meta); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &DownloadMetadataResponse{
		Download: new(
			langx.Clone(
				Download{},
				DownloadOptionFromTorrentMetadata(langx.Clone(meta, timex.JSONSafeEncodeOption))),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) metadatasync(w http.ResponseWriter, r *http.Request) {
	var (
		req Download
		tmd tracking.Metadata
		id  = mux.Vars(r)["id"]
	)

	if err := jsonx.UnmarshalRead(r.Body, &req); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decoded update"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err := tracking.MetadataAssignKnownMediaID(r.Context(), t.q, id, req.Media.KnownMediaId).Scan(&tmd); err != nil {
		log.Println(errorsx.Wrapf(err, "unable to sync media metadata: %s", id))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := sqlx.Discard(sqlx.Scan(library.MetadataSyncKnownMediaIDFromTorrent(r.Context(), t.q, id))); sqlx.IgnoreNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to sync media metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &MetadataSyncResponse{
		Media: new(
			langx.Clone(
				Media{},
				MediaOptionFromTorrentMetadata(langx.Clone(tmd, timex.JSONSafeEncodeOption))),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) downloading(w http.ResponseWriter, r *http.Request) {
	var (
		msg = DownloadSearchResponse{
			Next: &DownloadSearchRequest{
				Limit: 100,
			},
		}
	)

	if err := t.decoder.Decode(msg.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	msg.Next.Limit = numericx.Min(msg.Next.Limit, 100)

	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryCompleted(false),
			tracking.MetadataQueryInitiated(),
			tracking.MetadataQueryNotPaused(),
		},
	).OrderBy("peers > 0 DESC, downloaded < bytes DESC, downloaded/bytes DESC").Offset(msg.Next.Offset * msg.Next.Limit).Limit(msg.Next.Limit)

	qq := sqlx.Scan(tracking.MetadataSearch(r.Context(), t.q, q))
	for p := range qq.Iter() {
		tmp := langx.Clone(Download{}, DownloadOptionFromTorrentMetadata(langx.Clone(p, timex.JSONSafeEncodeOption)))
		msg.Items = append(msg.Items, &tmp)
	}

	if err := qq.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "encoding failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) search(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg = DownloadSearchResponse{
			Next: &DownloadSearchRequest{
				Limit: 100,
			},
		}
	)

	if err = t.decoder.Decode(msg.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	msg.Next.Limit = numericx.Min(msg.Next.Limit, 100)

	q := tracking.MetadataSearchBuilder().Where(squirrel.And{
		tracking.MetadataQueryCompleted(msg.Next.Completed),
		tracking.MetadataQueryNotTombstoned(),
		tracking.MetadataQueryHidden(msg.Next.Hidden),
		lucenex.Query(t.fts, msg.Next.Query, lucenex.WithDefaultField("auto_description")),
	}).OrderBy("created_at DESC").Offset(msg.Next.Offset * msg.Next.Limit).Limit(msg.Next.Limit)

	i := sqlx.Scan(tracking.MetadataSearch(r.Context(), t.q, q))
	for p := range i.Iter() {
		tmp := langx.Clone(Download{}, DownloadOptionFromTorrentMetadata(langx.Clone(p, timex.JSONSafeEncodeOption)))
		msg.Items = append(msg.Items, &tmp)
	}

	if err = i.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "encoding failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
