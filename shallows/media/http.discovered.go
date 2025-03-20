package media

import (
	"context"
	"crypto/md5"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/x/bytesx"
	"github.com/retrovibed/retrovibed/internal/x/errorsx"
	"github.com/retrovibed/retrovibed/internal/x/formx"
	"github.com/retrovibed/retrovibed/internal/x/fsx"
	"github.com/retrovibed/retrovibed/internal/x/httpx"
	"github.com/retrovibed/retrovibed/internal/x/iox"
	"github.com/retrovibed/retrovibed/internal/x/jwtx"
	"github.com/retrovibed/retrovibed/internal/x/langx"
	"github.com/retrovibed/retrovibed/internal/x/numericx"
	"github.com/retrovibed/retrovibed/internal/x/slicesx"
	"github.com/retrovibed/retrovibed/internal/x/sqlx"
	"github.com/retrovibed/retrovibed/internal/x/sqlxx"
	"github.com/retrovibed/retrovibed/tracking"
)

type HTTPDiscoveredOption func(*HTTPDiscovered)

func HTTPDiscoveredOptionJWTSecret(j jwtx.JWTSecretSource) HTTPDiscoveredOption {
	return func(t *HTTPDiscovered) {
		t.jwtsecret = j
	}
}

type download interface {
	Start(t torrent.Metadata) (dl torrent.Torrent, added bool, err error)
	Stop(t torrent.Metadata) (err error)
}

func NewHTTPDiscovered(q sqlx.Queryer, d download, c storage.ClientImpl, options ...HTTPDiscoveredOption) *HTTPDiscovered {
	svc := langx.Clone(HTTPDiscovered{
		q:            q,
		d:            d,
		c:            c,
		jwtsecret:    env.JWTSecret,
		decoder:      formx.NewDecoder(),
		mediastorage: fsx.DirVirtual(os.TempDir()),
	}, options...)

	return &svc
}

type HTTPDiscovered struct {
	q            sqlx.Queryer
	d            download
	c            storage.ClientImpl
	jwtsecret    jwtx.JWTSecretSource
	decoder      *form.Decoder
	mediastorage fsx.Virtual
}

func (t *HTTPDiscovered) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.TimeoutRollingWrite(3*time.Second),
	).ThenFunc(t.upload))

	r.Path("/available").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/downloading").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.downloading))

	r.Path("/{id}").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.download))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		// httpauth.AuthenticateWithToken(t.jwtsecret),
		// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.pause))
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

	tmp, err := fsx.CreateTemp(t.mediastorage, "retrovibed.upload.*")
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
	if info, err := meta.UnmarshalInfo(); err == nil && !info.Private {
		meta.AnnounceList = append(meta.AnnounceList, tracking.PublicTrackers())
	}

	info, err := meta.UnmarshalInfo()
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to read temporary file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}
	info.Name = fh.Filename

	lmd := tracking.NewMetadata(langx.Autoptr(meta.HashInfoBytes()),
		tracking.MetadataOptionFromInfo(&info),
		tracking.MetadataOptionTrackers(slicesx.Flatten(meta.AnnounceList...)...),
	)

	if err = tracking.MetadataInsertWithDefaults(r.Context(), t.q, lmd).Scan(&lmd); err != nil {
		log.Println(errorsx.Wrap(err, "unable to record metadata record"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	metadata, err := torrent.New(metainfo.Hash(lmd.Infohash), torrent.OptionStorage(t.c), torrent.OptionNodes(meta.NodeList()...), torrent.OptionTrackers(meta.AnnounceList...), torrent.OptionWebseeds(meta.UrlList))
	if err != nil {
		log.Println(errorsx.Wrapf(err, "unable to create torrent from metadata %s", lmd.ID))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if dl, _, err = t.d.Start(metadata); err != nil {
		log.Println(errorsx.Wrap(err, "unable to start download"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	go func() {
		errorsx.Log(tracking.Download(context.Background(), t.q, t.mediastorage, &lmd, dl))
	}()

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &MediaUploadResponse{
		Media: langx.Autoptr(
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

	if err = t.d.Stop(metadata); err != nil {
		log.Println(errorsx.Wrap(err, "unable to stop download"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = tracking.MetadataPausedByID(r.Context(), t.q, id).Scan(&md); err != nil {
		log.Println(errorsx.Wrap(err, "unable to pause metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &DownloadBeginResponse{
		Download: langx.Autoptr(
			langx.Clone(
				Download{},
				DownloadOptionFromTorrentMetadata(langx.Clone(md, tracking.MetadataOptionJSONSafeEncode))),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) download(w http.ResponseWriter, r *http.Request) {
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

	if _, _, err := t.d.Start(metadata); err != nil {
		log.Println(errorsx.Wrap(err, "unable to start download"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := tracking.MetadataDownloadByID(r.Context(), t.q, id).Scan(&md); err != nil {
		log.Println(errorsx.Wrap(err, "unable to track download"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &DownloadBeginResponse{
		Download: langx.Autoptr(
			langx.Clone(
				Download{},
				DownloadOptionFromTorrentMetadata(langx.Clone(md, tracking.MetadataOptionJSONSafeEncode))),
		),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) downloading(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg DownloadSearchResponse = DownloadSearchResponse{
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

	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryInitiated(),
			tracking.MetadataQueryIncomplete(),
			tracking.MetadataQueryNotPaused(),
		},
	).OrderBy("created_at DESC").Offset(msg.Next.Offset * msg.Next.Limit).Limit(msg.Next.Limit)

	err = sqlxx.ScanEach(tracking.MetadataSearch(r.Context(), t.q, q), func(p *tracking.Metadata) error {
		tmp := langx.Clone(Download{}, DownloadOptionFromTorrentMetadata(langx.Clone(*p, tracking.MetadataOptionJSONSafeEncode)))
		msg.Items = append(msg.Items, &tmp)
		return nil
	})

	if err != nil {
		log.Println(errorsx.Wrap(err, "encoding failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPDiscovered) search(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg MediaSearchResponse = MediaSearchResponse{
			Next: &MediaSearchRequest{
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
		tracking.MetadataQueryNotInitiated(),
		tracking.MetadataQuerySearch(msg.Next.Query, "description"),
	}).OrderBy("created_at DESC").Offset(msg.Next.Offset * msg.Next.Limit).Limit(msg.Next.Limit)

	err = sqlxx.ScanEach(tracking.MetadataSearch(r.Context(), t.q, q), func(p *tracking.Metadata) error {
		tmp := langx.Clone(Media{}, MediaOptionFromTorrentMetadata(langx.Clone(*p, tracking.MetadataOptionJSONSafeEncode)))
		msg.Items = append(msg.Items, &tmp)
		return nil
	})

	if err != nil {
		log.Println(errorsx.Wrap(err, "encoding failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
