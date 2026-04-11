package ddisc

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/mimex"
	"github.com/retrovibed/retrovibed/internal/slicesx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/squirrelx"
	"github.com/retrovibed/retrovibed/internal/torrentx"
)

func DiscoveredOptionNoop(*Discovered) {}
func DiscoveredOptionIndex(b bool) func(*Discovered) {
	if b {
		return func(d *Discovered) {
			d.KnownMediaID = uuid.Max.String()
		}
	}

	return DiscoveredOptionNoop
}

func DiscoveredOptionMimetype(s string) func(*Discovered) {
	return func(d *Discovered) {
		d.Mimetype = langx.FirstNonZero(s, string(mimex.Binary))
	}
}

func DiscoveredOptionFromTorrentInfo(i *metainfo.Info) func(*Discovered) {
	return func(d *Discovered) {
		d.Title = i.Name
		d.Bytes = uint64(i.TotalLength())
	}
}

func DiscoveredOptionKnownMedia(id string) func(*Discovered) {
	return func(d *Discovered) {
		d.KnownMediaID = id
	}
}

func DiscoveredOptionPartitionAuto(partitions *Partition) func(*Discovered) {
	return func(d *Discovered) {
		uid := uuid.FromStringOrNil(d.KnownMediaID)
		if uid.IsZero() || uuid.Max == uid {
			// do nothing
			return
		}

		d.Partition = partitions.Max([]byte(d.KnownMediaID)).String()
	}
}

func DiscoveredOptionPartition(p string) func(*Discovered) {
	return func(d *Discovered) {
		d.Partition = p
	}
}

func DiscoveredOptionFromExtracted(ex Extracted) func(*Discovered) {
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

		d.Mimetype = langx.FirstNonZero(
			langx.Autoderef(ex.Music).Mimetype,
			langx.Autoderef(ex.Video).Mimetype,
			d.Mimetype,
		)

		d.ReleasedAt = langx.FirstNonZero(
			langx.Autoderef(ex.Music).Date,
			langx.Autoderef(ex.Video).Date,
			d.ReleasedAt,
		)
	}
}

func NewDiscovered(md *int160.T, options ...func(*Discovered)) (m Discovered) {
	r := langx.Clone(Discovered{
		ID:           torrentx.HashUID(md),
		Infohash:     md.Bytes(),
		KnownMediaID: uuid.Nil.String(),
		Mimetype:     mimex.Bittorrent,
		SyncUID:      uuid.Must(uuid.NewV7()).String(),
		Partition:    uuid.Nil.String(),
	}, options...)
	return r
}

func DiscoveredQueryNeedsCheck() squirrel.Sqlizer {
	return squirrel.Expr("ddisc_media.next_check_at < NOW()")
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
