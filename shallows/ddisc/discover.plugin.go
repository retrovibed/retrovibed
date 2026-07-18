package ddisc

import (
	"context"
	"iter"
	"log"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// searchPlugins is the narrow interface PluginStrategy needs from
// *retroapi/searchplugin.Registry.
type searchPlugins interface {
	Search(ctx context.Context, mimetypes []string, query string, adult bool) iterx.Seq[*ddiscapi.Import]
}

// PluginStrategy runs external search plugins (via a
// *retroapi/searchplugin.Registry) and yields whatever they find. Yielded
// candidates are neither ranked nor persisted here - Discover's central seq
// does both for every candidate regardless of which strategy produced it.
// No-ops if req.Title is empty — plugins can't be usefully queried without
// one. As a side effect, independent of what it yields, it TOFU-records
// (see KnownMediaFromImport) any known-media catalog info a plugin result
// carries — q is used only for that.
func PluginStrategy(q sqlx.Queryer, plugins searchPlugins) DiscoverStrategy {
	return pluginStrategy{q: q, plugins: plugins}
}

type pluginStrategy struct {
	q       sqlx.Queryer
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

		seq := t.cfg.plugins.Search(ctx, t.req.Mimetypes, t.req.Title, t.req.Adult)
		for imp := range seq.Each(ctx) {
			m, err := metainfo.ParseMagnetURI(imp.Uri)
			if err != nil {
				log.Println("search plugin returned an unresolvable uri", err)
				continue
			}

			mimetype := langx.FirstNonZero(imp.Mimetype, langx.FirstNonZero(t.req.Mimetypes...))

			if kid := uuid.FromStringOrNil(imp.KnownMediaId); !uuidx.IsMinMax(kid) {
				known := KnownMediaFromImport(kid, mimetype, imp)
				if err := library.KnownInsertWithDefaultsTOFU(ctx, t.cfg.q, known).Scan(&known); err != nil {
					log.Println("unable to record known media from plugin", err)
				}
			}

			id := int160.FromBytes(m.InfoHash.Bytes())
			d := NewDiscovered(
				&id,
				DiscoveredOptionURI(imp.Uri),
				DiscoveredOptionKnownMedia(langx.FirstNonZero(imp.KnownMediaId, t.req.KnownMediaID)),
				DiscoveredOptionMimetype(Generalize(mimetype)),
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
