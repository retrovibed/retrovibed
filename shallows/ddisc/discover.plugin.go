package ddisc

import (
	"context"
	"iter"

	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

// PluginStrategy runs external search plugins (via a
// *retroapi/searchplugin.Registry) and yields whatever they find. Yielded
// candidates are neither ranked nor persisted here - Discover's central seq
// does both for every candidate regardless of which strategy produced it,
// including TOFU-recording known-media catalog info (see KnownMediaTOFU).
// No-ops if req.Title is empty — plugins can't be usefully queried without
// one.
func PluginStrategy(plugins searchplugin.T) DiscoverStrategy {
	return pluginStrategy{plugins: plugins}
}

type pluginStrategy struct {
	plugins searchplugin.T
}

func (t pluginStrategy) Discover(ctx context.Context, req DiscoverRequest) iterx.Seq[Discovered] {
	return &pluginSeq{cfg: t, req: req}
}

type pluginSeq struct {
	cfg pluginStrategy
	req DiscoverRequest
	err error
}

func (t *pluginSeq) Each(ctx context.Context) iter.Seq[Discovered] {
	return func(yield func(Discovered) bool) {
		if t.req.Query == "" {
			return
		}

		seq := t.cfg.plugins.Search(ctx, t.req.Mimetypes, t.req.Query, t.req.Adult)
		for imp := range seq.Each(ctx) {
			if imp.Uri == "" {
				continue
			}

			mimetype := langx.FirstNonZero(imp.Mimetype, langx.FirstNonZero(t.req.Mimetypes...))

			d := NewDiscoveredFromImport(
				imp,
				DiscoveredOptionKnownMedia(langx.FirstNonZero(imp.KnownMediaId, t.req.KnownMediaID)),
				DiscoveredOptionMimetype(Generalize(mimetype)),
			)

			if !yield(d) {
				return
			}
		}

		if err := seq.Err(); err != nil {
			t.err = err
		}
	}
}

func (t *pluginSeq) Err() error {
	return t.err
}
