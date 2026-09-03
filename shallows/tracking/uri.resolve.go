package tracking

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

// Resolve resolves uri - a magnet: uri, or an http(s) url that serves a raw
// .torrent file when fetched - into an in-memory Metadata carrying its real
// infohash and torrent info, without persisting anything to the database.
// options are applied before the auto (AutoDescription/AutoHidden) options,
// so a caller-supplied MetadataOptionMimetype/MetadataOptionDescription is
// what AutoHidden/AutoDescription actually see. The http(s) case still
// writes the fetched .torrent bytes to t.rootstore's cache (so a later
// Import of the same uri doesn't need to re-fetch) and is still
// rate-limited by t.l - only the database insert is skipped, so callers
// that only need the real infohash (e.g. background search, which checks
// many candidates it has no intention of downloading) aren't forced to
// create a torrents_metadata row for every one of them.
func (t URIImport) Resolve(ctx context.Context, uri string, options ...func(*Metadata)) (meta Metadata, err error) {
	if strings.HasPrefix(uri, "magnet:") {
		return t.resolveMagnet(uri, options...)
	}
	return t.resolveHTTP(ctx, uri, options...)
}

func (t URIImport) resolveMagnet(uri string, options ...func(*Metadata)) (meta Metadata, err error) {
	md, err := metainfo.ParseMagnetURI(uri)
	if err != nil {
		return meta, errorsx.Wrap(err, "unable to parse magnet link")
	}

	return NewMetadata(
		new(int160.FromByteArray(md.InfoHash)),
		MetadataOptionFromMagnet(&md),
		langx.Compose(options...),
		MetadataOptionAutoDescription,
		MetadataOptionAutoHidden,
	), nil
}

func (t URIImport) resolveHTTP(ctx context.Context, uri string, options ...func(*Metadata)) (meta Metadata, err error) {
	log.Println("DERP DERP initiated")
	defer log.Println("DERP DERP completed")

	if err := t.l.Wait(ctx); err != nil {
		return meta, errorsx.Wrap(err, "rate limited")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return meta, err
	}

	log.Println("checkpoint 0")
	resp, err := httpx.AsError(t.client.Do(req))
	if err != nil {
		return meta, errorsx.Wrap(err, "unable to retrieve uri")
	}
	defer resp.Body.Close()
	log.Println("checkpoint 1")
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return meta, errorsx.Wrap(err, "unable to read response")
	}

	log.Println("checkpoint 2")
	md, err := metainfo.Load(bytes.NewBuffer(buf))
	if err != nil {
		log.Println("FAILURE", string(buf))
		return meta, errorsx.Wrap(err, "unable to read metainfo from response")
	}
	log.Println("checkpoint 3")
	mi, err := md.UnmarshalInfo()
	if err != nil {
		return meta, errorsx.Wrap(err, "unable to read info from metadata")
	}

	torrentvfs := fsx.DirVirtual(t.rootstore.Path(env.TorrentDirName))
	if err = os.WriteFile(torrentvfs.Path(fmt.Sprintf("%s.torrent", md.HashInfoBytes().String())), buf, 0600); err != nil {
		return meta, errorsx.Wrap(err, "unable to persist torrent to disk")
	}

	return NewMetadata(
		new(md.ID()),
		MetadataOptionFromInfo(&mi),
		MetadataOptionTrackers(md.Announce),
		langx.Compose(options...),
		MetadataOptionAutoDescription,
		MetadataOptionAutoHidden,
	), nil
}
