// Command retrodsearch queries a PeerTube instance's (default
// sepiasearch.org, a cross-instance search index) json search api and
// exports results as jsonl. It's both a normal CLI tool for development
// (`go run . --query ubuntu`) and, when invoked as
// `retrodsearch plugin --mimetype <m> --query <q> [--adult]`, a
// searchplugin.Registry plugin (see
// github.com/retrovibed/retrovibed/retroapi/examples/searchplugin-noop for
// the calling contract).
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"golang.org/x/time/rate"

	// autohijack points net.DefaultResolver and http.DefaultTransport at
	// wasinet's virtual sockets when this is built for wasip1 (a no-op on
	// every other platform). Registry.Search mounts the host's real TLS
	// trust store into the sandbox and sets SSL_CERT_DIR before running a
	// plugin, so a plain net/http.Get - TLS verification included - works
	// unmodified once this import is in place.
	_ "github.com/egdaemon/wasinet/wasinet/autohijack"
)

// categoryCodes maps this plugin's genre categories onto PeerTube's numeric
// category ids. "all" is a sentinel meaning "omit categoryOneOf entirely".
var categoryCodes = map[string]string{
	"all":    "",
	"movies": "2",
	"tv":     "10",
	"music":  "1",
	"games":  "7",
}

type peerTubeCategory struct {
	ID int `json:"id"`
}

type peerTubeSearchRow struct {
	UUID          string           `json:"uuid"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	ThumbnailPath string           `json:"thumbnailPath"`
	Views         float64          `json:"views"`
	Category      peerTubeCategory `json:"category"`
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

type cli struct {
	Query      string   `flag:"" name:"query" required:""`
	Categories []string `flag:"" name:"category" default:"all" help:"category(s) to search within (all, movies, tv, music, games); repeatable"`
	Domain     string   `flag:"" name:"domain" default:"https://sepiasearch.org" env:"PEERTUBE_DOMAIN" help:"base url of the PeerTube/SepiaSearch instance"`
	Adult      bool     `flag:"" name:"adult" default:"false" help:"allow adult (NSFW) results"`
	MaxResults uint     `flag:"" name:"max-results" default:"128" help:"maximum results to fetch detail pages for"`
	Attempts   uint     `flag:"" name:"attempts" default:"5" help:"maximum retry attempts per request"`
	Rate       float64  `flag:"" name:"rate" default:"1" help:"maximum requests per second against the target domain"`
	Workers    uint     `flag:"" name:"workers" default:"4" help:"number of concurrent detail-page fetches"`
	Output     string   `flag:"" name:"output" default:"-" help:"output destination; '-' for stdout"`
}

func (t cli) Run(ctx context.Context) (err error) {
	c := &http.Client{
		Timeout: 30 * time.Second,
		Transport: HeaderTransport{
			Headers: http.Header{"User-Agent": []string{"retrodsearch/peertube (+https://github.com/retrovibed/retrovibed/tree/main/retrodsearch)"}},
			Next:    MaybeDebugTransport(http.DefaultTransport),
		},
	}

	out := io.Writer(os.Stdout)
	if t.Output != "" && t.Output != "-" {
		f, err := os.Create(t.Output)
		if err != nil {
			return fmt.Errorf("unable to open output %q: %w", t.Output, err)
		}
		defer f.Close()
		out = f
	}

	limiter := rate.NewLimiter(rate.Limit(t.Rate), 1)
	return t.run(ctx, c, out, limiter)
}

func (t cli) run(ctx context.Context, c *http.Client, out io.Writer, l *rate.Limiter) (err error) {
	categories := make([]string, 0, len(t.Categories))
	for _, category := range t.Categories {
		category = strings.ToLower(category)
		if _, ok := categoryCodes[category]; !ok {
			return fmt.Errorf("unsupported category %q", category)
		}
		categories = append(categories, category)
	}

	var encMu sync.Mutex
	enc := json.NewEncoder(out)

	pool := asynccompute.New(func(ctx context.Context, row peerTubeSearchRow) error {
		target := fmt.Sprintf("%s/api/v1/videos/%s", strings.TrimRight(t.Domain, "/"), url.PathEscape(row.UUID))
		body, ferr := Fetch(ctx, c, l, t.Attempts, target)
		if ferr != nil {
			log.Printf("retrodsearch: failed to fetch video detail %s: %v", row.UUID, ferr)
			return nil
		}

		var video peerTubeVideoResponse
		if ferr := json.Unmarshal(body, &video); ferr != nil {
			log.Printf("retrodsearch: failed to parse video detail %s: %v", row.UUID, ferr)
			return nil
		}

		file, ok := bestFile(video.Files)
		if !ok {
			log.Printf("retrodsearch: no magnet link found for %s", row.UUID)
			return nil
		}

		imp := &ddiscapi.Import{
			Uri:        file.MagnetUri,
			Uritype:    mimex.Magnet,
			Title:      row.Name,
			Overview:   row.Description,
			PosterPath: joinURL(t.Domain, row.ThumbnailPath),
			Popularity: row.Views,
		}

		encMu.Lock()
		encErr := enc.Encode(imp)
		encMu.Unlock()
		if encErr != nil {
			return fmt.Errorf("failed to encode result: %w", encErr)
		}
		return nil
	}, asynccompute.Workers[peerTubeSearchRow](uint16(t.Workers)))
	defer func() {
		if cerr := pool.Close(); err == nil {
			err = cerr
		}
	}()

	for _, category := range categories {
		if err := t.runCategory(ctx, c, l, category, pool); err != nil {
			return err
		}
	}

	return nil
}

// runCategory pages through category's listing via PeerTube's start/count
// pagination, dispatching every result row onto pool, until MaxResults is
// hit or every result has been seen. MaxResults applies per category, not
// to the combined total across every category in t.Categories.
func (t cli) runCategory(ctx context.Context, c *http.Client, l *rate.Limiter, category string, pool *asynccompute.Pool[peerTubeSearchRow]) error {
	const pageSize = 25

	var start uint
	for {
		if t.MaxResults > 0 && start >= t.MaxResults {
			break
		}

		count := uint(pageSize)
		if t.MaxResults > 0 && start+count > t.MaxResults {
			count = t.MaxResults - start
		}

		target, err := searchURL(t.Domain, category, t.Query, t.Adult, start, count)
		if err != nil {
			return err
		}

		body, err := Fetch(ctx, c, l, t.Attempts, target)
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
			break
		}
	}

	return nil
}

// bestFile picks the highest-resolution entry among files that actually
// have a magnet link, since PeerTube videos are typically transcoded into
// several resolutions and retrodsearch only ever emits magnet uris - a
// file without one (e.g. web-seed/download-only) is not a candidate at
// all, not just a lower-priority one.
func bestFile(files []peerTubeFile) (peerTubeFile, bool) {
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

func searchURL(domain, category, query string, adult bool, start, count uint) (string, error) {
	u, err := url.Parse(strings.TrimRight(domain, "/") + "/api/v1/search/videos")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("search", query)
	q.Set("start", strconv.FormatUint(uint64(start), 10))
	q.Set("count", strconv.FormatUint(uint64(count), 10))
	// PeerTube's own api param is "nsfw" - kept as an internal
	// implementation detail here; this plugin's own flag/protocol-facing
	// name is "adult", to match retrovibed's established terminology.
	q.Set("nsfw", strconv.FormatBool(adult))
	if code := categoryCodes[category]; code != "" {
		q.Set("categoryOneOf", code)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func joinURL(domain, path string) string {
	if path == "" {
		return ""
	}
	return strings.TrimRight(domain, "/") + path
}

// Fetch retrieves target, retrying up to attempts times with an exponential
// backoff, and respecting l as a politeness rate limit.
func Fetch(ctx context.Context, c *http.Client, l *rate.Limiter, attempts uint, target string) ([]byte, error) {
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

		if err := l.Wait(ctx); err != nil {
			return nil, err
		}

		attemptCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		body, err := doFetch(attemptCtx, c, target)
		cancel()
		if err != nil {
			lastErr = err
			continue
		}
		return body, nil
	}

	return nil, fmt.Errorf("failed to fetch %s after %d attempts: %w", target, attempts, lastErr)
}

func doFetch(ctx context.Context, c *http.Client, target string) ([]byte, error) {
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

// HeaderTransport sets a fixed set of headers (e.g. User-Agent) on every
// request before delegating to Next.
type HeaderTransport struct {
	Headers http.Header
	Next    http.RoundTripper
}

func (t HeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	cloned := req.Clone(req.Context())
	for k, vs := range t.Headers {
		for _, v := range vs {
			cloned.Header.Set(k, v)
		}
	}
	return t.Next.RoundTrip(cloned)
}

// debugHTTPEnabled gates debugTransport - off by default, set
// RETRODSEARCH_DEBUG_HTTP (any non-empty value) to turn full raw
// request/response logging back on.
func debugHTTPEnabled() bool {
	return os.Getenv("RETRODSEARCH_DEBUG_HTTP") != ""
}

// MaybeDebugTransport wraps delegate in debugTransport when
// debugHTTPEnabled, otherwise returns delegate unchanged.
func MaybeDebugTransport(delegate http.RoundTripper) http.RoundTripper {
	if !debugHTTPEnabled() {
		return delegate
	}
	return debugTransport{delegate: delegate}
}

type debugTransport struct {
	delegate http.RoundTripper
}

func (t debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if raw, err := httputil.DumpRequest(req, true); err == nil {
		log.Println("RAW REQUEST")
		log.Println("Scheme:", req.URL.Scheme)
		log.Println(string(raw))
	}

	resp, err := t.delegate.RoundTrip(req)

	if resp != nil && resp.Body != nil {
		if raw, derr := httputil.DumpResponse(resp, true); derr != nil {
			return resp, derr
		} else {
			log.Println("RAW RESPONSE")
			log.Println(string(raw))
		}
	}

	return resp, err
}

// mimetypeFlag collects every --mimetype occurrence into a slice, since a
// search request can carry several candidate discovery mimetypes and
// stdlib flag has no built-in repeatable-flag type.
type mimetypeFlag []string

func (m *mimetypeFlag) String() string { return strings.Join(*m, ",") }
func (m *mimetypeFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

// ParsePluginArgs parses the searchplugin.Registry calling convention's
// arguments: everything after the leading "plugin" argv token
// (os.Args[2:]), which is always repeatable --mimetype flags, a --query
// flag, and an optional --adult flag.
func ParsePluginArgs(args []string) (mimetypes []string, query string, adult bool) {
	fs := flag.NewFlagSet("plugin", flag.ExitOnError)
	var mf mimetypeFlag
	fs.Var(&mf, "mimetype", "discovery mimetype to search within (repeatable)")
	q := fs.String("query", "", "search text to query")
	a := fs.Bool("adult", false, "allow adult content in results")
	fs.Parse(args)
	return []string(mf), *q, *a
}

// specificSiteCategory maps retrovibed's discovery mimetypes onto the
// genre-based categories this plugin's own categoryCodes actually searches
// within. Deliberately has no entry for the ambiguous "video" umbrella
// mimetype - see resolveCategory below for why.
var specificSiteCategory = map[string]string{
	"application/vnd.retrovibed.discovery.movies": "movies",
	"application/vnd.retrovibed.discovery.tv":     "tv",
	"application/vnd.retrovibed.discovery.music":  "music",
	"application/vnd.retrovibed.discovery.audio":  "music",
	"audio": "music",
}

// resolveCategory resolves mimetypes (a search request's candidate
// discovery mimetypes, most specific first) down to the single genre
// category to search within. If exactly one distinct genre is recognized
// among mimetypes, that genre is used; otherwise this falls back to "all"
// rather than guessing wrong.
func resolveCategory(mimetypes []string) string {
	seen := map[string]bool{}
	for _, m := range mimetypes {
		if cat, ok := specificSiteCategory[m]; ok {
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

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// searchplugin.Registry always invokes a loaded plugin as
	// "<binary> plugin --mimetype <m> [...] --query <query> [--adult]" -
	// argv[1] is literally "plugin", a fixed part of that protocol (see
	// retroapi/examples/searchplugin-noop/main.go's doc comment). Detect
	// that mode before normal kong parsing, since kong doesn't have a
	// bare-vs-named-subcommand default mode that fits both call shapes.
	if len(os.Args) > 1 && os.Args[1] == "plugin" {
		mimetypes, query, adult := ParsePluginArgs(os.Args[2:])
		if err := runPlugin(ctx, mimetypes, query, adult); err != nil {
			log.Println("plugin:", err)
			os.Exit(1)
		}
		return
	}

	var c cli
	kctx := kong.Parse(&c, kong.BindTo(ctx, (*context.Context)(nil)))
	kctx.FatalIfErrorf(kctx.Run())
}

// runPlugin re-parses through kong (picking up its env-var-sourced defaults
// for --domain exactly like direct CLI use would) so per-install domain
// overrides still work under the plugin protocol.
func runPlugin(ctx context.Context, mimetypes []string, query string, adult bool) error {
	var c cli
	parser, err := kong.New(&c, kong.BindTo(ctx, (*context.Context)(nil)))
	if err != nil {
		return err
	}

	kctx, err := parser.Parse([]string{"--category", resolveCategory(mimetypes), "--query", query, "--adult=" + strconv.FormatBool(adult)})
	if err != nil {
		return err
	}

	return kctx.Run()
}
