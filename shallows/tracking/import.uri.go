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
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"golang.org/x/time/rate"
)

// URIImportOption configures a URIImport at construction; see NewURIImport.
type URIImportOption func(*URIImport)

func URIImportOptionHTTPClient(c *http.Client) URIImportOption {
	return func(t *URIImport) { t.client = c }
}

// URIImport resolves a uri - a magnet: uri, or an http(s) url that serves a
// raw .torrent file when fetched - persists the .torrent file to disk for
// the fetch case, and records a Metadata row for it either way. It
// consolidates what used to be duplicated per-item logic (see the old
// handlehttp/handlemagnet closures in torrent.rss.go). Everything about the
// recorded Metadata beyond the resolved uri/mimetype itself (known-media-id,
// autoarchive, encryption seed, description, ...) is the caller's job via
// Import's options, since URIImport has no way to know it up front.
// Deciding whether to auto-download the result is likewise the caller's
// job, using the returned Metadata's ID - not URIImport's.
type URIImport struct {
	q         sqlx.Queryer
	rootstore fsx.Virtual
	client    *http.Client
	l         *rate.Limiter
}

func NewURIImport(q sqlx.Queryer, c *http.Client, rootstore fsx.Virtual, options ...URIImportOption) URIImport {
	return langx.Clone(URIImport{
		q:         q,
		rootstore: rootstore,
		client:    c,
		l:         rate.NewLimiter(rate.Every(1*time.Second), 1),
	}, options...)
}

// Import resolves uri and records it as Metadata: a "magnet:" uri is parsed
// directly, anything else is treated as an http(s) url to fetch and parse
// as a raw .torrent file. options are applied after URIImport's own
// defaults, so a caller can override any of them (e.g.
// MetadataOptionMimetype, MetadataOptionDescription for a per-item title a
// shared URIImport instance has no way to know) or layer on anything else
// MetadataOption supports. Whether to auto-download the result is left to
// the caller, using the returned Metadata's ID (see
// MetadataAutoDownloadByID) - that's not URIImport's job.
func (t URIImport) Import(ctx context.Context, uri string, options ...func(*Metadata)) (meta Metadata, err error) {
	if strings.HasPrefix(uri, "magnet:") {
		meta, err = t.importMagnet(ctx, uri, options...)
	} else {
		meta, err = t.importHTTP(ctx, uri, options...)
	}
	if err != nil {
		return meta, err
	}

	log.Println("uri import recorded", meta.ID, meta.Description)
	return meta, nil
}

func (t URIImport) importMagnet(ctx context.Context, uri string, options ...func(*Metadata)) (meta Metadata, err error) {
	md, err := metainfo.ParseMagnetURI(uri)
	if err != nil {
		return meta, errorsx.Wrap(err, "unable to parse magnet link")
	}

	m := NewMetadata(
		new(int160.FromByteArray(md.InfoHash)),
		MetadataOptionFromMagnet(&md),
		langx.Compose(options...),
		MetadataOptionAutoDescription,
		MetadataOptionAutoHidden,
	)

	if err = MetadataInsertWithDefaults(ctx, t.q, m).Scan(&meta); err != nil {
		return meta, errorsx.Wrap(err, "unable to record torrent metadata")
	}

	return meta, nil
}

func (t URIImport) importHTTP(ctx context.Context, uri string, options ...func(*Metadata)) (meta Metadata, err error) {
	if err := t.l.Wait(ctx); err != nil {
		return meta, errorsx.Wrap(err, "rate limited")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return meta, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return meta, errorsx.Wrap(err, "unable to retrieve uri")
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return meta, errorsx.Wrap(err, "unable to read response")
	}

	md, err := metainfo.Load(bytes.NewReader(buf))
	if err != nil {
		return meta, errorsx.Wrap(err, "unable to read metainfo from response")
	}

	mi, err := md.UnmarshalInfo()
	if err != nil {
		return meta, errorsx.Wrap(err, "unable to read info from metadata")
	}

	if err = os.WriteFile(t.rootstore.Path("torrent", fmt.Sprintf("%s.torrent", md.HashInfoBytes().String())), buf, 0600); err != nil {
		return meta, errorsx.Wrap(err, "unable to persist torrent to disk")
	}

	m := NewMetadata(
		new(md.ID()),
		MetadataOptionFromInfo(&mi),
		MetadataOptionTrackers(md.Announce),
		langx.Compose(options...),
		MetadataOptionAutoDescription,
		MetadataOptionAutoHidden,
	)

	if err = MetadataInsertWithDefaults(ctx, t.q, m).Scan(&meta); err != nil {
		return meta, errorsx.Wrap(err, "unable to record torrent metadata")
	}

	return meta, nil
}
