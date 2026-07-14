package ddisc

import (
	"context"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/ducktype"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/localex"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"golang.org/x/text/language"
)

type DiscoveredOption func(*Discovered)

func DiscoveredOptionNoop(*Discovered) {}
func DiscoveredOptionIndex(b bool) DiscoveredOption {
	if b {
		return func(d *Discovered) {
			d.KnownMediaID = uuid.Max.String()
		}
	}

	return DiscoveredOptionNoop
}

func DiscoveredOptionMimetype(s string) DiscoveredOption {
	return func(d *Discovered) {
		d.Mimetype = langx.FirstNonZero(s, string(mimex.Binary))
		d.Category = mimex.Category(d.Mimetype)
	}
}

func DiscoveredOptionHealth(h uint32) DiscoveredOption {
	return func(d *Discovered) {
		d.Health = h
	}
}

func DiscoveredOptionTitle(s string) DiscoveredOption {
	return func(d *Discovered) {
		d.Title = s
	}
}

func DiscoveredOptionTestDefaults(d *Discovered) {
	d.AudioDefaultLocale = localex.FirstDefined(userx.LocaleLanguage())
	d.SubtitlesDefaultLocale = localex.FirstDefined(userx.LocaleLanguage())
}

func DiscoveredOptionFromTorrentInfo(i *metainfo.Info) DiscoveredOption {
	return func(d *Discovered) {
		d.Title = i.Name
		d.Bytes = uint64(i.TotalLength())
	}
}

// DiscoveredOptionDetectCorrupted detects torrent metadata that cannot be stored
// as-is (e.g. a name that is not valid UTF8) and prevents it from propagating
// to other nodes: the title is sanitized with the unicode replacement glyph and
// the record is excluded from sync by zeroing its sync uid.
func DiscoveredOptionDetectCorrupted(d *Discovered) {
	if utf8.ValidString(d.Title) {
		return
	}

	d.Title = strings.ToValidUTF8(d.Title, "�")
	d.SyncUID = uuid.Nil.String()
}

func DiscoveredOptionKnownMedia(id string) DiscoveredOption {
	return func(d *Discovered) {
		d.KnownMediaID = id
	}
}

func DiscoveredOptionPartitionAuto(partitions *Partition) DiscoveredOption {
	return func(d *Discovered) {
		uid := uuid.FromStringOrNil(d.KnownMediaID)
		if uid.IsZero() || uuid.Max == uid {
			// do nothing
			return
		}

		d.Partition = partitions.Max([]byte(d.KnownMediaID)).String()
	}
}

func DiscoveredOptionPartition(p string) DiscoveredOption {
	return func(d *Discovered) {
		d.Partition = p
	}
}

func DiscoveredOptionFromExtracted(ex Extracted) DiscoveredOption {
	return func(d *Discovered) {
		d.Title = langx.FirstNonZero(
			langx.Autoderef(ex.Music).Title,
			langx.Autoderef(ex.Video).Title,
			d.Title,
		)

		d.Description = langx.FirstNonZero(
			langx.Autoderef(ex.Music).Subtitle,
			langx.Autoderef(ex.Video).Subtitle,
			d.Description,
		)

		d.Collation = langx.FirstNonZero(
			langx.Autoderef(ex.Music).Collation,
			langx.Autoderef(ex.Video).Collation,
			d.Collation,
		)

		DiscoveredOptionMimetype(
			langx.FirstNonZero(
				langx.Autoderef(ex.Music).Mimetype,
				langx.Autoderef(ex.Video).Mimetype,
				d.Mimetype,
			),
		)(d)

		d.ReleasedAt = langx.FirstNonZero(
			langx.Autoderef(ex.Music).Date,
			langx.Autoderef(ex.Video).Date,
			d.ReleasedAt,
		)
	}
}

// the worst possible discoverable candidate.
// used for zero values when performing discovey.
func Worst() Discovered {
	return Discovered{
		PolicyRank: math.MaxUint16,
		Health:     0,
		Bytes:      0,
	}
}

func NewDiscovered(md *int160.T, options ...DiscoveredOption) (m Discovered) {
	r := langx.Clone(Discovered{
		ID:                     torrentx.HashUID(md),
		Infohash:               md.Bytes(),
		KnownMediaID:           uuid.Nil.String(),
		Mimetype:               mimex.Binary,
		Category:               mimex.Application,
		SyncUID:                uuid.Must(uuid.NewV7()).String(),
		Partition:              uuid.Nil.String(),
		AudioDefaultLocale:     language.Und.String(),
		SubtitlesDefaultLocale: language.Und.String(),
	}, options...)
	return r
}

// NewDiscoveredFromKnown builds a Discovered record for a specific known media entity found
// within a torrent, keyed on (infohash, known media id) rather than infohash alone. This allows
// multiple Discovered rows to exist for the same infohash (e.g. one per episode in a season
// pack, or one per track in an album). known must already be resolved (never the Unknown()
// sentinel) - that precondition is the caller's responsibility.
func NewDiscoveredFromKnown(md int160.T, known library.Known, options ...DiscoveredOption) (m Discovered) {
	r := langx.Clone(Discovered{
		ID:                     md5x.FormatUUID(md5x.Digest(md.Bytes(), []byte(known.UID))),
		Infohash:               md.Bytes(),
		KnownMediaID:           known.UID,
		Title:                  known.Title,
		Description:            known.Overview,
		ReleasedAt:             known.Released,
		Adult:                  known.Adult,
		Mimetype:               langx.FirstNonZero(known.Mimetype, mimex.Binary),
		Category:               mimex.Category(langx.FirstNonZero(known.Mimetype, mimex.Application)),
		SyncUID:                uuid.Must(uuid.NewV7()).String(),
		Partition:              uuid.Nil.String(),
		AudioDefaultLocale:     language.Und.String(),
		SubtitlesDefaultLocale: language.Und.String(),
	}, options...)
	return r
}

func DiscoveredQueryNextCheck(r timex.Range) squirrel.Sqlizer {
	return squirrelx.Between("ddisc_media.next_check_at", ducktype.NewNullTime(r.Start), ducktype.NewNullTime(r.End))
}

func DiscoveredQueryKnownMediaID(id string) squirrel.Sqlizer {
	if id == "" {
		return squirrelx.Noop{}
	}
	return squirrel.Eq{"ddisc_media.known_media_id": id}
}

func DiscoveredQueryKnown() squirrel.Sqlizer {
	return squirrel.Expr("ddisc_media.known_media_id NOT IN ('00000000-0000-0000-0000-000000000000', 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF')")
}

func DiscoveredQueryByIDs(ids ...string) squirrel.Sqlizer {
	if len(ids) == 0 {
		return squirrelx.Noop{}
	}
	return squirrel.Eq{"ddisc_media.id": ids}
}

// DiscoveredQueryExplicit toggles whether adult content is allowed in results.
// allow=false restricts results to non-adult content; allow=true permits
// both adult and non-adult content (it does not restrict to adult-only).
func DiscoveredQueryExplicit(allow bool) squirrel.Sqlizer {
	if allow {
		return squirrelx.Noop{}
	}

	return squirrel.Expr("ddisc_media.adult = ?", false)
}

func DiscoveredQueryLanguage(v string) squirrel.Sqlizer {
	if stringsx.Blank(v) {
		return squirrelx.Noop{}
	}

	return squirrel.Expr("ddisc_media.audio_default_locale = ?", v)
}

func DiscoveredQueryText(query string) squirrel.Sqlizer {
	if query == "" {
		return squirrelx.Noop{}
	}
	return lucenex.Query(duckdbx.NewLucene(), query, lucenex.WithDefaultField("description"))
}

// DiscoveredQueryTitle restricts results to candidates whose title contains
// every term in query (bare terms combine with an implicit AND, so a
// multi-word query like "nirvana utero" requires both words present in the
// title) - this is what keeps RankAndSelect from picking a candidate whose
// title has nothing to do with what was actually searched for.
func DiscoveredQueryTitle(query string) squirrel.Sqlizer {
	if query == "" {
		return squirrelx.Noop{}
	}
	return lucenex.Query(duckdbx.NewLucene(), lucenex.Clean(query), lucenex.WithDefaultField("title"))
}

func DiscoveredQueryCategory(mimetype string) squirrel.Sqlizer {
	return squirrel.Expr("ddisc_media.category = ?", mimex.Category(mimetype))
}

func DiscoveredQueryMimetypes(mimetypes ...string) squirrel.Sqlizer {
	mimetypes = slicesx.MapTransform(func(s string) string {
		return md5x.FormatUUID(md5x.Digest(s))
	}, mimetypes...)

	return squirrelx.In("md5(mimetype)::uuid", mimetypes...)
}

func DiscoveredSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) DiscoveredScanner {
	return NewDiscoveredScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func DiscoveredSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(DiscoveredScannerStaticColumns)...).From("ddisc_media")
}
