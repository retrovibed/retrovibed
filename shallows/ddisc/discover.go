package ddisc

import (
	"context"
	"database/sql"
	"errors"
	"iter"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// DiscoverRequest carries whatever a DiscoverStrategy might need to find a
// match for a known-media-id.
type DiscoverRequest struct {
	KnownMediaID string
	Query        string   // optional; needed by the plugin strategy, empty means "skip it"
	Mimetypes    []string // optional; discovery mimetypes, needed by the plugin strategy
	Adult        bool     // optional; enable adult content.
}

// DiscoverRequestFromKnown builds a DiscoverRequest from a library.Known
// catalog entry, deriving Mimetypes from its Mimetype — the shape every
// caller that already has a Known row needs, so the field mapping lives in
// one place.
func DiscoverRequestFromKnown(known library.Known) DiscoverRequest {
	return DiscoverRequest{
		KnownMediaID: known.UID,
		Query:        known.Title,
		Mimetypes:    Category(known.Mimetype),
		Adult:        known.Adult,
	}
}

// Category expands mime (a coarse mimex category, e.g. "video"/"audio") into
// the ordered list of specific Retrovibed discovery mimetypes a search
// plugin should be queried with — this is discovery-specific policy (which
// coarse category maps to which discovery mimetypes), not a general
// mimetype fact, so it lives here rather than in mimex. Anything that isn't
// a recognized coarse category (mime is already specific, or unrelated)
// passes through unchanged.
func Category(mime string) []string {
	switch mimex.Category(mime) {
	case mimex.Video:
		return []string{mimex.RetrovibedDiscoveryMovies, mimex.RetrovibedDiscoveryTV, mimex.RetrovibedDiscoveryVideo}
	case mimex.Audio:
		return []string{mimex.RetrovibedDiscoveryMusic, mimex.RetrovibedDiscoveryAudio}
	default:
		return []string{mime}
	}
}

// Generalize inverts Category: maps one of the specific discovery
// mimetypes back to the coarse mimex category it was expanded from, so
// persisted Discovered/ddisc_media rows keep storing "video"/"audio"
// (which mimex.Category can still coarsen correctly) rather than a
// discovery-specific mimetype it can't parse. Anything else passes through
// unchanged.
func Generalize(mime string) string {
	switch mime {
	case mimex.RetrovibedDiscoveryMovies, mimex.RetrovibedDiscoveryTV, mimex.RetrovibedDiscoveryVideo:
		return mimex.Video
	case mimex.RetrovibedDiscoveryMusic, mimex.RetrovibedDiscoveryAudio:
		return mimex.Audio
	default:
		return mime
	}
}

// DiscoverStrategy is a single way to look for media matching a
// DiscoverRequest. A strategy may have side effects (e.g. firing a
// fire-and-forget peer query) beyond what it yields synchronously from
// Discover.
type DiscoverStrategy interface {
	Discover(ctx context.Context, req DiscoverRequest) iterx.Seq[Discovered]
}

// UnimplementedStrategy is a safe default DiscoverStrategy: every Discover
// yields nothing and never errors, so a disabled/unavailable strategy slot
// (e.g. peertube search turned off, or a test that doesn't care about it)
// can use this instead of a nil DiscoverStrategy - mirrors
// searchplugin.Unimplemented.
type UnimplementedStrategy struct{}

func (UnimplementedStrategy) Discover(ctx context.Context, req DiscoverRequest) iterx.Seq[Discovered] {
	return unimplementedDiscoverSeq{}
}

type unimplementedDiscoverSeq struct{}

func (unimplementedDiscoverSeq) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {}
}

func (unimplementedDiscoverSeq) Err() error { return nil }

// DiscoverOption customizes a Discover call.
type DiscoverOption func(*discoverSeq)

// DiscoverOptionFilter adds pred as a filter every candidate from every
// strategy must pass before it's ranked and yielded. Omitted entirely,
// every candidate passes through unfiltered (see discoverNoopFilter).
func DiscoverOptionFilter(pred func(context.Context, Discovered) (bool, error)) DiscoverOption {
	return func(t *discoverSeq) {
		t.filter = pred
	}
}

// discoverNoopFilter is discoverSeq's default filter: every candidate
// passes. Keeps discoverSeq.filter always callable, no nil check needed.
func discoverNoopFilter(context.Context, Discovered) (bool, error) {
	return true, nil
}

// DiscoverOptionDetectMedia sets detectMedia as the transform discoverSeq
// runs every strategy's candidate seq through, letting it fill in
// KnownMediaID on any Discovered candidate that doesn't already carry one
// (e.g. a plugin/peertube hit for a free-text query) - see
// KnownMediaDetector for the concrete clean+match implementation. Omitted
// entirely, no candidate is touched (see discoverNoopDetectMedia).
func DiscoverOptionDetectMedia(detectMedia func(iterx.Seq[Discovered]) iterx.Seq[Discovered]) DiscoverOption {
	return func(t *discoverSeq) {
		t.detectMedia = detectMedia
	}
}

// discoverNoopDetectMedia is discoverSeq's default detectMedia transform:
// every candidate passes through unchanged. Keeps discoverSeq.detectMedia
// always callable, no nil check needed.
func discoverNoopDetectMedia(s iterx.Seq[Discovered]) iterx.Seq[Discovered] {
	return s
}

// Discover tries every strategy in order, regardless of whether an earlier
// strategy already yielded something - so a caller sees candidates from
// every strategy, not just the first one to produce a hit. A genuine error
// from any strategy stops the whole chain; a strategy simply coming up empty
// (iterx.ErrNotFound) does not. Every yielded candidate passes the filter
// (see DiscoverOptionFilter) and is ranked with policy before being handed
// to the caller, so callers observe the real policy rank rather than the
// unranked sentinel - but candidates are not persisted here. Persisting is
// the caller's job: only the caller knows whether a candidate is worth
// keeping (e.g. daemons.DiscoveredDownload persists only the winner, once
// import has resolved its real infohash) or whether every candidate should
// be recorded regardless (e.g. daemons.SearchQueueBackgroundRun, which
// exists purely to populate ddisc_media for background browsing).
func Discover(ctx context.Context, policy Policy, req DiscoverRequest, options []DiscoverOption, strategies ...DiscoverStrategy) iterx.Seq[Discovered] {
	d := &discoverSeq{policy: policy, req: req, strategies: strategies, filter: discoverNoopFilter, detectMedia: discoverNoopDetectMedia}
	for _, opt := range options {
		opt(d)
	}
	return d
}

type discoverSeq struct {
	policy      Policy
	req         DiscoverRequest
	strategies  []DiscoverStrategy
	filter      func(context.Context, Discovered) (bool, error)
	detectMedia func(iterx.Seq[Discovered]) iterx.Seq[Discovered]
	err         error
}

func (t *discoverSeq) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {
		for _, strategy := range t.strategies {
			seq := strategy.Discover(ctx, t.req)
			seq = iterx.Filter(seq, t.filter)
			seq = iterx.NotFound(seq)
			seq = t.detectMedia(seq)

			for d := range seq.Each(ctx) {
				if err := t.policy.Rank(&d); err != nil {
					t.err = err
					return
				}

				if !yield(d) {
					return
				}
			}

			if err := seq.Err(); errors.Is(err, iterx.ErrNotFound) {
				continue
			} else if err != nil {
				t.err = errorsx.Wrapf(err, "%T failed", strategy)
				return
			}
		}
	}
}

func (t *discoverSeq) Err() error {
	return t.err
}

// TitleFilter matches Discovered candidates whose Title satisfies req's
// lucene query, run via duckdbx.Search against the single candidate title.
type TitleFilter struct {
	q   sqlx.Queryer
	req DiscoverRequest
}

func NewTitleFilter(q sqlx.Queryer, req DiscoverRequest) TitleFilter {
	return TitleFilter{q: q, req: req}
}

func (t TitleFilter) Match(ctx context.Context, d Discovered) (bool, error) {
	var matched bool
	err := duckdbx.Search(lucenex.Clean(t.req.Query), d.Title).RunWith(t.q).QueryRowContext(ctx).Scan(&matched)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// ExternalStrategies returns the discovery strategies that reach outside
// this node's own database: external wasm search plugins (if plugins is
// non-nil) and the in-process PeerTube/SepiaSearch strategy (if peertube is
// non-nil). Unlike SyncStrategies, this never includes LocalStrategy - a
// caller that only wants "what can the outside world tell us" (e.g. a
// background queue drain populating ddisc_media from scratch) would gain
// nothing from a strategy that just reads ddisc_media back.
func ExternalStrategies(q sqlx.Queryer, plugins searchplugin.T, peertube DiscoverStrategy) []DiscoverStrategy {
	strategies := []DiscoverStrategy{}
	if plugins != nil {
		strategies = append(strategies, PluginStrategy(q, plugins))
	}
	if peertube != nil {
		strategies = append(strategies, peertube)
	}
	return strategies
}

// SyncStrategies returns the synchronous (non-DHT) discovery strategies:
// ExternalStrategies plus the local ddisc_media scan (if knownMediaID
// resolves to something other than uuid.Nil). All of them always run
// together via Discover - the local scan is not gated on the external
// strategies coming up empty.
func SyncStrategies(q sqlx.Queryer, plugins searchplugin.T, peertube DiscoverStrategy, knownMediaID string) []DiscoverStrategy {
	strategies := append([]DiscoverStrategy{KnownStrategy(q)}, ExternalStrategies(q, plugins, peertube)...)
	if knownMediaID != uuid.Nil.String() {
		strategies = append(strategies, LocalStrategy(q))
	}
	return strategies
}
