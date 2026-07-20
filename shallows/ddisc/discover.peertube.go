package ddisc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"golang.org/x/time/rate"
)

// peerTubeCategoryCodes maps this strategy's genre categories onto
// PeerTube's numeric category ids. "all" is a sentinel meaning "omit
// categoryOneOf entirely".
var peerTubeCategoryCodes = map[string]string{
	"all":    "",
	"movies": "2",
	"tv":     "10",
	"music":  "1",
	"games":  "7",
}

// peerTubeMimetypeCategory maps retrovibed's discovery mimetypes onto the
// genre category this strategy actually searches within. Deliberately has
// no entry for the ambiguous "video" umbrella mimetype - see
// resolvePeerTubeCategory for why.
var peerTubeMimetypeCategory = map[string]string{
	mimex.RetrovibedDiscoveryMovies: "movies",
	mimex.RetrovibedDiscoveryTV:     "tv",
	mimex.RetrovibedDiscoveryMusic:  "music",
	mimex.RetrovibedDiscoveryAudio:  "music",
	mimex.Audio:                     "music",
}

// resolvePeerTubeCategory resolves mimetypes (a search request's candidate
// discovery mimetypes, most specific first) down to the single genre
// category to search within. If exactly one distinct genre is recognized
// among mimetypes, that genre is used; otherwise this falls back to "all"
// rather than guessing wrong.
func resolvePeerTubeCategory(mimetypes []string) string {
	seen := map[string]bool{}
	for _, m := range mimetypes {
		if cat, ok := peerTubeMimetypeCategory[m]; ok {
			seen[cat] = true
		}
	}

	if len(seen) == 1 {
		for cat := range seen {
			return cat
		}
	}

	return "all"
}

// peerTubeCategoryMimetype maps a resolved genre category back onto a coarse
// mimex category for DiscoveredOptionMimetype, since PeerTube's categories
// don't carry a mimetype of their own.
func peerTubeCategoryMimetype(category string) string {
	switch category {
	case "movies", "tv", "games":
		return mimex.Video
	case "music":
		return mimex.Audio
	default:
		return mimex.Binary
	}
}

type peerTubeSearchRow struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type peerTubeSearchResponse struct {
	Total uint                `json:"total"`
	Data  []peerTubeSearchRow `json:"data"`
}

type peerTubeResolution struct {
	ID int `json:"id"`
}

type peerTubeFile struct {
	Resolution peerTubeResolution `json:"resolution"`
	MagnetUri  string             `json:"magnetUri"`
}

type peerTubeVideoResponse struct {
	Files []peerTubeFile `json:"files"`
}

// bestPeerTubeFile picks the highest-resolution entry among files that
// actually have a magnet link - PeerTube videos are typically transcoded
// into several resolutions and this strategy only ever emits magnet uris, so
// a file without one (e.g. web-seed/download-only) is not a candidate at
// all, not just a lower-priority one.
func bestPeerTubeFile(files []peerTubeFile) (peerTubeFile, bool) {
	var best peerTubeFile
	found := false
	for _, f := range files {
		if f.MagnetUri == "" {
			continue
		}
		if !found || f.Resolution.ID > best.Resolution.ID {
			best = f
			found = true
		}
	}
	return best, found
}

func peerTubeSearchURL(domain, category, query string, adult bool, start, count uint) (string, error) {
	u, err := url.Parse(strings.TrimRight(domain, "/") + "/api/v1/search/videos")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("search", query)
	q.Set("start", strconv.FormatUint(uint64(start), 10))
	q.Set("count", strconv.FormatUint(uint64(count), 10))
	// PeerTube's own api param is "nsfw" - kept as an internal
	// implementation detail here; DiscoverRequest's own field is "Adult",
	// to match retrovibed's established terminology.
	q.Set("nsfw", strconv.FormatBool(adult))
	if code := peerTubeCategoryCodes[category]; code != "" {
		q.Set("categoryOneOf", code)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// peerTubeHeaderTransport sets a descriptive User-Agent on every request
// before delegating to next.
type peerTubeHeaderTransport struct {
	next http.RoundTripper
}

func (t peerTubeHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	cloned := req.Clone(req.Context())
	cloned.Header.Set("User-Agent", "retrovibed/peertube-discovery (+https://github.com/retrovibed/retrovibed)")
	return t.next.RoundTrip(cloned)
}

type peerTubeStrategy struct {
	client     *http.Client
	domain     string
	limiter    *rate.Limiter
	maxResults uint
	attempts   uint
	workers    uint
}

type PeerTubeOption func(*peerTubeStrategy)

// PeerTubeOptionMaxResults caps how many search results (per Discover call)
// this strategy pages through before it stops fetching detail pages. 0
// means unbounded.
func PeerTubeOptionMaxResults(n uint) PeerTubeOption {
	return func(t *peerTubeStrategy) { t.maxResults = n }
}

// PeerTubeOptionAttempts sets the maximum retry attempts per HTTP request.
func PeerTubeOptionAttempts(n uint) PeerTubeOption {
	return func(t *peerTubeStrategy) { t.attempts = n }
}

// PeerTubeOptionRate sets the maximum requests per second issued against the
// configured domain, shared across every concurrent Discover call this
// strategy serves for the lifetime of the process.
func PeerTubeOptionRate(requestsPerSecond float64) PeerTubeOption {
	return func(t *peerTubeStrategy) { t.limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), 1) }
}

// PeerTubeOptionWorkers sets the number of concurrent detail-page fetches.
func PeerTubeOptionWorkers(n uint) PeerTubeOption {
	return func(t *peerTubeStrategy) { t.workers = n }
}

// PeerTubeStrategy searches a PeerTube/SepiaSearch instance's public JSON
// search api directly in-process via c - unlike PluginStrategy, which
// sandboxes external, untrusted plugins through wazero, this is trusted,
// first-party code with no wasm compile/install step, so it runs the same
// way LocalStrategy does.
func PeerTubeStrategy(c *http.Client, domain string, options ...PeerTubeOption) DiscoverStrategy {
	client := *c
	client.Transport = peerTubeHeaderTransport{next: langx.FirstNonZero(client.Transport, http.DefaultTransport)}

	t := &peerTubeStrategy{
		client:     &client,
		domain:     domain,
		limiter:    rate.NewLimiter(1, 1),
		maxResults: 128,
		attempts:   5,
		workers:    4,
	}

	for _, opt := range options {
		opt(t)
	}

	return t
}

func (t *peerTubeStrategy) Discover(ctx context.Context, req DiscoverRequest) iterx.Seq[Discovered] {
	return &peerTubeSeq{cfg: t, req: req}
}

type peerTubeSeq struct {
	cfg *peerTubeStrategy
	req DiscoverRequest
	err error
}

func (t *peerTubeSeq) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {
		if t.req.Title == "" {
			return
		}

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		results := make(chan Discovered)
		done := make(chan error, 1)

		go func() {
			defer close(results)
			done <- t.cfg.run(ctx, t.req, results)
		}()

		for d := range results {
			if !yield(d) {
				cancel()
				break
			}
		}

		t.err = <-done
	}
}

func (t *peerTubeSeq) Err() error {
	return t.err
}

// run pages the search listing for req, fanning detail-page fetches out
// across a worker pool, and sends every candidate with a usable magnet link
// to results.
func (t *peerTubeStrategy) run(ctx context.Context, req DiscoverRequest, results chan<- Discovered) error {
	category := resolvePeerTubeCategory(req.Mimetypes)

	pool := asynccompute.New(func(ctx context.Context, row peerTubeSearchRow) error {
		d, ok, err := t.resolveRow(ctx, category, row)
		if err != nil {
			log.Printf("ddisc: peertube: failed to resolve %s: %v", row.UUID, err)
			return nil
		}
		if !ok {
			return nil
		}

		select {
		case results <- d:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}, asynccompute.Workers[peerTubeSearchRow](uint16(t.workers)))

	serr := t.search(ctx, category, req.Title, req.Adult, pool)
	cerr := pool.Close()
	return errors.Join(serr, cerr)
}

// search pages through category's listing via PeerTube's start/count
// pagination, dispatching every result row onto pool, until maxResults is
// hit or every result has been seen.
func (t *peerTubeStrategy) search(ctx context.Context, category, query string, adult bool, pool *asynccompute.Pool[peerTubeSearchRow]) error {
	const pageSize = 25

	var start uint
	for {
		if t.maxResults > 0 && start >= t.maxResults {
			return nil
		}

		count := uint(pageSize)
		if t.maxResults > 0 && start+count > t.maxResults {
			count = t.maxResults - start
		}

		target, err := peerTubeSearchURL(t.domain, category, query, adult, start, count)
		if err != nil {
			return err
		}

		body, err := t.fetch(ctx, target)
		if err != nil {
			return fmt.Errorf("failed to fetch search results: %w", err)
		}

		var resp peerTubeSearchResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return fmt.Errorf("failed to parse search results: %w", err)
		}

		for _, row := range resp.Data {
			if err := pool.Run(ctx, row); err != nil {
				return err
			}
		}

		start += uint(len(resp.Data))
		if len(resp.Data) == 0 || start >= resp.Total {
			return nil
		}
	}
}

// resolveRow fetches row's video detail page, picks the best magnet-bearing
// file, and builds a Discovered keyed by that magnet's real infohash.
// Returns ok=false (no error) for videos with no usable magnet - not every
// PeerTube result is a candidate.
func (t *peerTubeStrategy) resolveRow(ctx context.Context, category string, row peerTubeSearchRow) (Discovered, bool, error) {
	target := fmt.Sprintf("%s/api/v1/videos/%s", strings.TrimRight(t.domain, "/"), url.PathEscape(row.UUID))
	body, err := t.fetch(ctx, target)
	if err != nil {
		return Discovered{}, false, err
	}

	var video peerTubeVideoResponse
	if err := json.Unmarshal(body, &video); err != nil {
		return Discovered{}, false, err
	}

	file, ok := bestPeerTubeFile(video.Files)
	if !ok {
		return Discovered{}, false, nil
	}

	magnet, err := metainfo.ParseMagnetURI(file.MagnetUri)
	if err != nil {
		return Discovered{}, false, fmt.Errorf("unable to parse magnet uri: %w", err)
	}

	infohash := int160.FromByteArray(magnet.InfoHash)
	d := NewDiscovered(
		&infohash,
		DiscoveredOptionURI(file.MagnetUri),
		DiscoveredOptionTitle(row.Name),
		DiscoveredOptionDescription(row.Description),
		DiscoveredOptionMimetype(peerTubeCategoryMimetype(category)),
		DiscoveredOptionDetectCorrupted,
	)

	return d, true, nil
}

// fetch retrieves target, retrying up to t.attempts times with an
// exponential backoff, and respecting t.limiter as a politeness rate limit
// against the configured domain.
func (t *peerTubeStrategy) fetch(ctx context.Context, target string) ([]byte, error) {
	attempts := t.attempts
	if attempts == 0 {
		attempts = 1
	}

	var lastErr error
	for attempt := uint(0); attempt < attempts; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(time.Duration(attempt) * time.Second):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		if err := t.limiter.Wait(ctx); err != nil {
			return nil, err
		}

		attemptCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		body, err := doPeerTubeFetch(attemptCtx, t.client, target)
		cancel()
		if err != nil {
			lastErr = err
			continue
		}
		return body, nil
	}

	return nil, fmt.Errorf("failed to fetch %s after %d attempts: %w", target, attempts, lastErr)
}

func doPeerTubeFetch(ctx context.Context, c *http.Client, target string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, target)
	}

	return body, nil
}
