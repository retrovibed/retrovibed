package tracking

import (
	"context"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/deeppool/internal/x/duckdbx"
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/internal/x/squirrelx"
	"github.com/james-lawrence/deeppool/internal/x/timex"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
)

func MetadataOptionFromInfo(i *metainfo.Info) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = i.Name
		m.Bytes = uint64(i.TotalLength())
	}
}

func MetadataOptionDescription(d string) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = d
	}
}

func MetadataOptionJSONSafeEncode(p *Metadata) {
	p.CreatedAt = timex.RFC3339NanoEncode(p.CreatedAt)
	p.UpdatedAt = timex.RFC3339NanoEncode(p.UpdatedAt)
}

func NewMetadata(md *metainfo.Hash, options ...func(*Metadata)) (m Metadata) {
	r := langx.Clone(Metadata{
		ID:       HashUID(md),
		Infohash: md.Bytes(),
	}, options...)
	return r
}

func MetadataQueryNotInitiated() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.initiated_at = 'infinity'")
}

func MetadataQueryInitiated() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.initiated_at < NOW()")
}

func MetadataQueryIncomplete() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.downloaded < torrents_metadata.bytes")
}

func MetadataQueryNotPaused() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.paused_at = 'infinity'")
}

func MetadataSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) MetadataScanner {
	return NewMetadataScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func MetadataQuerySearch(q string) squirrel.Sqlizer {
	return duckdbx.FTSSearch("fts_main_torrents_metadata", q)
}

func MetadataSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(MetadataScannerStaticColumns)...).From("torrents_metadata")
}

func DownloadProgress(ctx context.Context, q sqlx.Queryer, md Metadata, dl torrent.Torrent) {
	for range time.Tick(time.Second) {
		current := uint64(dl.BytesCompleted())
		if md.Downloaded == current {
			continue
		}

		if err := MetadataProgressByID(ctx, q, md.ID, current).Scan(&md); err != nil {
			log.Println("failed to update progress", err)
		} else {
			log.Println(md.ID, "updated", md.Downloaded/md.Bytes, md.Downloaded, "/", md.Bytes)
		}
	}
}
