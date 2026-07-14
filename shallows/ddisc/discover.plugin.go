package ddisc

import (
	"context"
	"iter"
	"log"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

// searchPlugins is the narrow interface PluginStrategy needs from
// *retroapi/searchplugin.Registry.
type searchPlugins interface {
	Search(ctx context.Context, category, query string) iterx.Seq[*ddiscapi.Import]
}

// PluginStrategy runs external search plugins (via a
// *retroapi/searchplugin.Registry) and yields whatever they find. Yielded
// candidates are neither ranked nor persisted here - Discover's central seq
// does both for every candidate regardless of which strategy produced it.
// No-ops if req.Title is empty — plugins can't be usefully queried without
// one.
func PluginStrategy(plugins searchPlugins) DiscoverStrategy {
	return pluginStrategy{plugins: plugins}
}

type pluginStrategy struct {
	plugins searchPlugins
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
		if t.req.Title == "" {
			return
		}

		seq := t.cfg.plugins.Search(ctx, t.req.Category, t.req.Title)
		for imp := range seq.Each(ctx) {
			m, err := metainfo.ParseMagnetURI(imp.Magnet)
			if err != nil {
				log.Println("search plugin returned invalid magnet uri", err)
				continue
			}

			id := int160.FromBytes(m.InfoHash.Bytes())
			d := NewDiscovered(
				&id,
				DiscoveredOptionKnownMedia(t.req.KnownMediaID),
				DiscoveredOptionMimetype(langx.FirstNonZero(imp.Mimetype, t.req.Category)),
				DiscoveredOptionHealth(imp.Health),
				DiscoveredOptionTitle(m.DisplayName),
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
