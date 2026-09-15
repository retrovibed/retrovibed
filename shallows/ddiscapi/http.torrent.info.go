package ddiscapi

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/justinas/alice"
	rootenv "github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type HTTPTorrentInfoOption func(*HTTPTorrentInfo)

func HTTPTorrentInfoOptionJWTSecret(j jwtx.SecretSource) HTTPTorrentInfoOption {
	return func(t *HTTPTorrentInfo) { t.jwtsecret = j }
}

func HTTPTorrentInfoOptionRootStorage(vfs fsx.Virtual) HTTPTorrentInfoOption {
	return func(t *HTTPTorrentInfo) { t.rootstorage = vfs }
}

type HTTPTorrentInfo struct {
	db          *sql.DB
	rootstorage fsx.Virtual
	jwtsecret   jwtx.SecretSource
}

func NewHTTPTorrentInfo(db *sql.DB, options ...HTTPTorrentInfoOption) *HTTPTorrentInfo {
	t := &HTTPTorrentInfo{
		db:          db,
		rootstorage: fsx.DirVirtual(os.TempDir()),
		jwtsecret:   env.JWTSecret,
	}
	for _, o := range options {
		o(t)
	}
	return t
}

func (t *HTTPTorrentInfo) Bind(r *mux.Router) {
	r.StrictSlash(false)
	r.Use(httpx.RouteInvoked)

	r.Path("/{id}/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.get))
}

func (t *HTTPTorrentInfo) get(w http.ResponseWriter, r *http.Request) {
	var (
		meta tracking.Metadata
		id   = mux.Vars(r)["id"]
	)

	if err := tracking.MetadataFindByID(r.Context(), t.db, id).Scan(&meta); sqlx.ErrNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	} else if err != nil {
		log.Println(errorsx.Wrap(err, "unable to find metadata"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	path := t.rootstorage.Path(rootenv.TorrentDirName, metainfo.Hash(meta.Infohash).String())
	mi, err := metainfo.LoadFromFile(path + tracking.TorrentSuffix)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to load torrent file"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to unmarshal torrent info"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	announceList := make([]string, 0, len(mi.UpvertedAnnounceList()))
	for tracker := range mi.UpvertedAnnounceList().DistinctValues() {
		announceList = append(announceList, tracker)
	}

	files := make([]*metaapi.TorrentFile, 0, len(info.Files))
	for f := range metainfo.Files(&info) {
		files = append(files, &metaapi.TorrentFile{
			Name:   f.Path,
			Length: f.Length,
			Path:   f.Path,
		})
	}

	resp := metaapi.TorrentInfoResponse{
		Meta: &metaapi.TorrentMeta{
			Comment:      mi.Comment,
			Encoding:     mi.Encoding,
			CreatedBy:    mi.CreatedBy,
			CreationDate: uint64(mi.CreationDate),
			AnnounceList: announceList,
			UrlList:      []string(mi.UrlList),
		},
		Details: &metaapi.TorrentDetails{
			Name:    info.Name,
			Length:  uint64(info.TotalLength()),
			Source:  info.Source,
			Private: info.Private != nil && *info.Private,
		},
		Files: files,
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
	}
}
