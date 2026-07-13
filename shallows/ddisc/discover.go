package ddisc

import (
	"context"
	"errors"
	"iter"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// DiscoverRequest carries whatever a DiscoverStrategy might need to find a
// match for a known-media-id.
type DiscoverRequest struct {
	KnownMediaID string
	Title        string // optional; needed by the plugin strategy, empty means "skip it"
	Category     string // optional; mimex category, needed by the plugin strategy
}

// DiscoverRequestFromKnown builds a DiscoverRequest from a library.Known
// catalog entry, deriving Category from its Mimetype — the shape every
// caller that already has a Known row needs, so the field mapping lives in
// one place.
func DiscoverRequestFromKnown(known library.Known) DiscoverRequest {
	return DiscoverRequest{
		KnownMediaID: known.UID,
		Title:        known.Title,
		Category:     mimex.Category(known.Mimetype),
	}
}

// DiscoverStrategy is a single way to look for media matching a
// DiscoverRequest. A strategy may have side effects (e.g. firing a
// fire-and-forget peer query, or persisting a plugin result into ddisc_media)
// beyond what it yields synchronously from Discover.
type DiscoverStrategy interface {
	Discover(ctx context.Context, req DiscoverRequest) iterx.Seq[Discovered]
}

// Discover tries each strategy in order, stopping at the first that yields a
// result. Earlier strategies' side effects (e.g. a fired-and-forgotten peer
// query, or a plugin result persisted along the way) still happen even when
// that strategy doesn't yield anything itself.
func Discover(ctx context.Context, req DiscoverRequest, strategies ...DiscoverStrategy) iterx.Seq[Discovered] {
	return &discoverSeq{req: req, strategies: strategies}
}

type discoverSeq struct {
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
