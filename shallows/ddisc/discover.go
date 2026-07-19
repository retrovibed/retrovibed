package ddisc

import (
	"context"
	"errors"
	"iter"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// DiscoverRequest carries whatever a DiscoverStrategy might need to find a
// match for a known-media-id.
type DiscoverRequest struct {
	KnownMediaID string
	Title        string   // optional; needed by the plugin strategy, empty means "skip it"
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
		Title:        known.Title,
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

// Discover tries each strategy in order, stopping at the first that yields a
// result. Earlier strategies' side effects (e.g. a fired-and-forgotten peer
// query) still happen even when that strategy doesn't yield anything itself.
// Every yielded candidate is ranked with policy before being handed to the
// caller, so callers observe the real policy rank rather than the unranked
// sentinel - but candidates are not persisted here. Persisting is the
// caller's job: only the caller knows whether a candidate is worth keeping
// (e.g. daemons.DiscoveredDownload persists only the winner, once import has
// resolved its real infohash) or whether every candidate should be recorded
// regardless (e.g. daemons.SearchQueueBackgroundRun, which exists purely to
// populate ddisc_media for background browsing).
func Discover(ctx context.Context, policy Policy, req DiscoverRequest, strategies ...DiscoverStrategy) iterx.Seq[Discovered] {
	return &discoverSeq{policy: policy, req: req, strategies: strategies}
}

type discoverSeq struct {
	policy     Policy
	req        DiscoverRequest
	strategies []DiscoverStrategy
	err        error
}

func (t *discoverSeq) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {
		for _, strategy := range t.strategies {
			seq := iterx.NotFound(strategy.Discover(ctx, t.req))

			found := false
			for d := range seq.Each(ctx) {
				found = true

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
				t.err = err
				return
			}

			if found {
				return
			}
		}
	}
}

func (t *discoverSeq) Err() error {
	return t.err
}

// SyncStrategies returns the synchronous (non-DHT) discovery strategies:
// external search plugins (if plugins is non-nil) and the local ddisc_media
// fallback (if knownMediaID resolves to something other than uuid.Nil).
func SyncStrategies(q sqlx.Queryer, plugins searchplugin.T, knownMediaID string) []DiscoverStrategy {
	strategies := []DiscoverStrategy{}
	if plugins != nil {
		strategies = append(strategies, PluginStrategy(q, plugins))
	}
	if knownMediaID != uuid.Nil.String() {
		strategies = append(strategies, LocalStrategy(q))
	}
	return strategies
}
