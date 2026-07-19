package tracking

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/env"
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
	errorsx.Log(errorsx.Wrap(fsx.MkDirs(0700, rootstore.Path(env.TorrentDirName)), "unable to create torrent cache directory"))

	return langx.Clone(URIImport{
		q:         q,
		rootstore: rootstore,
		client:    c,
		l:         rate.NewLimiter(rate.Every(1*time.Second), 1),
	}, options...)
}

// Import resolves uri (see Resolve) and records the result as Metadata.
// options are applied after Resolve's own defaults, so a caller can
// override any of them (e.g. MetadataOptionMimetype,
// MetadataOptionDescription for a per-item title a shared URIImport
// instance has no way to know) or layer on anything else MetadataOption
// supports. Whether to auto-download the result is left to the caller,
// using the returned Metadata's ID (see MetadataAutoDownloadByID) - that's
// not URIImport's job.
func (t URIImport) Import(ctx context.Context, uri string, options ...func(*Metadata)) (meta Metadata, err error) {
	m, err := t.Resolve(ctx, uri)
	if err != nil {
		return meta, err
	}

	m = langx.Clone(m, options...)

	if err = MetadataInsertWithDefaults(ctx, t.q, m).Scan(&meta); err != nil {
		return meta, errorsx.Wrap(err, "unable to record torrent metadata")
	}

	log.Println("uri import recorded", meta.ID, meta.Description)
	return meta, nil
}
