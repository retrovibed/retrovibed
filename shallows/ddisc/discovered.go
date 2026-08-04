package ddisc

import (
	"context"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/ducktype"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/localex"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
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

// acquisition_state is a plain int32 mirroring the AcquisitionState enum
// ordinals (both this package's AcquisitionState and ddisc.discovery.proto's
// generated equivalent). Kept as a bare int32 rather than the AcquisitionState
// type itself because genieql has no per-field type-override mechanism - it
// always generates the DB-driven primitive type for a column.
func DiscoveredOptionAcquisitionState(s AcquisitionState) DiscoveredOption {
	return func(d *Discovered) {
		d.AcquisitionState = int32(s)
	}
}

func DiscoveredOptionInfoHash(i []byte) DiscoveredOption {
	return func(d *Discovered) {
		d.Infohash = i
	}
}

func DiscoveredOptionTitle(s string) DiscoveredOption {
	return func(d *Discovered) {
		d.Title = s
	}
}

func DiscoveredOptionDescription(s string) DiscoveredOption {
	return func(d *Discovered) {
		d.Description = s
	}
}

// DiscoveredOptionPrivate marks a candidate as sourced from a BEP 27 private
// torrent. Private candidates are still persisted (so this node can use them
// locally) but are excluded from every peer-facing sync/search response.
func DiscoveredOptionPrivate(b bool) DiscoveredOption {
	return func(d *Discovered) {
		d.Private = b
	}
}

// DiscoveredOptionContentMime sets ContentMime directly - how the uri is
// fetched (mimex.Bittorrent, mimex.HTTP), not the media's own Mimetype.
func DiscoveredOptionContentMime(s string) DiscoveredOption {
	return func(d *Discovered) {
		d.Contentmime = s
	}
}

func DiscoveredOptionURI(s string) DiscoveredOption {
	return func(d *Discovered) {
		d.URI = s
	}
}

// DiscoveredOptionAutoMagnet synthesizes a minimal magnet uri from d.Infohash
// when uri isn't already set (i.e. NewDiscovered/NewDiscoveredFromKnown was
// called with an empty uri). Callers that only ever have a raw infohash in
// hand (DHT/wire-sync, protobuf bridges, etc.) can use this instead of
// constructing metainfo.Magnet{InfoHash: metainfo.Hash(id.Bytes())}.String()
// themselves at every call site.
func DiscoveredOptionAutoMagnet(d *Discovered) {
	if d.URI != "" {
		return
	}

	d.URI = metainfo.Magnet{InfoHash: metainfo.Hash(d.Infohash)}.String()
	d.Contentmime = mimex.Bittorrent
}

func DiscoveredOptionTestDefaults(d *Discovered) {
	d.AudioDefaultLocale = localex.FirstDefined(userx.LocaleLanguage())
	d.SubtitlesDefaultLocale = localex.FirstDefined(userx.LocaleLanguage())
}

func DiscoveredOptionFromTorrentInfo(i *metainfo.Info) DiscoveredOption {
	return func(d *Discovered) {
		d.Title = i.Name
		d.Bytes = uint64(i.TotalLength())
		d.Private = langx.Autoderef(i.Private)
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

// PolicyRejectionCatalogOnly marks a Discovered synthesized purely from the
// library_known_media catalog (see KnownStrategy) - the catalog has never
// actually resolved a downloadable source for it. Rank's happy path never
// clears an already-set PolicyRejection, so this survives Discover's central
// ranking pass and keeps Select from ever choosing it for download.
const PolicyRejectionCatalogOnly = "catalog-only: no known download source yet"

// DiscoveredOptionCatalogOnly marks d as a catalog-only candidate: not a real
// downloadable source, just a hit against our library_known_media catalog.
func DiscoveredOptionCatalogOnly(d *Discovered) {
	d.PolicyRejection = PolicyRejectionCatalogOnly
	d.PolicyRank = math.MaxUint16
}

// catalogURI synthesizes a non-empty placeholder URI for a catalog-only
// Discovered - ddisc_media.uri is NOT NULL CHECK (uri <> ”), so anything
// that might get persisted (e.g. daemons.SearchQueueBackgroundRun) needs a
// well-formed value even though it's never a real fetchable source.
func catalogURI(knownMediaID string) string {
	return "retrovibed+catalog://" + knownMediaID
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

// defaults content mimetype to bittorrent.
// as that was the initial usecase for uris. the uri was either magnet or a
// http torrent file uri.
// non-torrent results from plugins *must* set their contentmime for resolution
// to work properly.
func NewDiscovered(md *int160.T, options ...DiscoveredOption) (m Discovered) {
	r := langx.Clone(Discovered{
		ID:                     torrentx.HashUID(md),
		Infohash:               md.Bytes(),
		KnownMediaID:           uuid.Nil.String(),
		Mimetype:               mimex.Binary,
		Contentmime:            mimex.Bittorrent,
		Category:               mimex.Application,
		SyncUID:                uuid.Must(uuid.NewV7()).String(),
		Partition:              uuid.Nil.String(),
		AudioDefaultLocale:     language.Und.String(),
		SubtitlesDefaultLocale: language.Und.String(),
		CreatedAt:              timex.NegInf(),
		UpdatedAt:              timex.NegInf(),
		NextCheckAt:            timex.NegInf(),
		TombstonedAt:           timex.Inf(),
		PolicyRank:             math.MaxUint16,
	},
		DiscoveredOptionAcquisitionState(AcquisitionStateEphemeral),
		langx.Compose(options...),
	)
	return r
}

// NewDiscoveredFromKnown builds a Discovered record for a specific known media entity found
// within a torrent, keyed on (infohash, known media id) rather than infohash alone. This allows
// multiple Discovered rows to exist for the same infohash (e.g. one per episode in a season
// pack, or one per track in an album). known must already be resolved (never the Unknown()
// sentinel) - that precondition is the caller's responsibility.
func NewDiscoveredFromKnown(md int160.T, known library.Known, options ...DiscoveredOption) (m Discovered) {
	r := langx.Clone(Discovered{
		Source:                 "retrovibed.media.archive",
		ID:                     md5x.FormatUUID(md5x.Digest(md.Bytes(), []byte(known.UID))),
		Infohash:               md.Bytes(),
		KnownMediaID:           known.UID,
		Title:                  known.Title,
		Description:            known.Overview,
		ReleasedAt:             known.Released,
		Adult:                  known.Adult,
		Mimetype:               langx.FirstNonZero(known.Mimetype, mimex.Binary),
		Contentmime:            mimex.Bittorrent,
		Category:               mimex.Category(langx.FirstNonZero(known.Mimetype, mimex.Application)),
		SyncUID:                uuid.Must(uuid.NewV7()).String(),
		Partition:              uuid.Nil.String(),
		AudioDefaultLocale:     language.Und.String(),
		SubtitlesDefaultLocale: language.Und.String(),
		CreatedAt:              timex.NegInf(),
		UpdatedAt:              timex.NegInf(),
		NextCheckAt:            timex.NegInf(),
		TombstonedAt:           timex.Inf(),
		PolicyRank:             math.MaxUint16,
	},
		DiscoveredOptionAcquisitionState(AcquisitionStateEphemeral),
		langx.Compose(options...),
	)
	return r
}

// NewDiscoveredFromImport builds a Discovered candidate directly from a
// search-plugin result, before its real infohash is known - keyed by
// md5(uri) (same pattern as ddisc.Locate's row id) rather than the real
// infohash NewDiscovered's callers already have in hand. ddisc_media's
// infohash column is NOT NULL and must be exactly 20 bytes, so
// int160.FromHashedBytes(imp.Uri) (SHA1 of the uri) is stored as a
// placeholder until the real infohash is resolved - which only happens
// once this row is actually selected and handed to an importer (see
// daemons.DiscoveredDownload), not for every candidate up front.
func NewDiscoveredFromImport(imp *ddiscapi.Import, options ...DiscoveredOption) (m Discovered) {
	placeholder := int160.FromHashedBytes([]byte(imp.Uri))
	r := langx.Clone(Discovered{
		ID:                     md5x.FormatUUID(md5x.Digest(imp.Uri)),
		Source:                 imp.Source,
		URI:                    imp.Uri,
		Infohash:               placeholder.Bytes(),
		Title:                  imp.Title,
		Health:                 imp.Health,
		Bytes:                  imp.Bytes,
		KnownMediaID:           uuid.Nil.String(),
		Mimetype:               mimex.Binary,
		Contentmime:            mimex.Bittorrent,
		Category:               mimex.Application,
		SyncUID:                uuid.Must(uuid.NewV7()).String(),
		Partition:              uuid.Nil.String(),
		AudioDefaultLocale:     language.Und.String(),
		SubtitlesDefaultLocale: language.Und.String(),
		CreatedAt:              timex.NegInf(),
		UpdatedAt:              timex.NegInf(),
		NextCheckAt:            timex.NegInf(),
		TombstonedAt:           time.Now().Add(3 * time.Hour),
		PolicyRank:             math.MaxUint16,
	},
		DiscoveredOptionAcquisitionState(AcquisitionStateEphemeral),
		langx.Compose(options...),
	)
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
