package tracking

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/internal/x/squirrelx"
	"github.com/james-lawrence/deeppool/internal/x/timex"
	"github.com/james-lawrence/torrent/metainfo"
)

func MetadataOptionFromInfo(i *metainfo.Info) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = i.Name
		m.Bytes = uint64(i.TotalLength())
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

func MetadataSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) MetadataScanner {
	return NewMetadataScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func MetadataSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(MetadataScannerStaticColumns)...).From("torrents_metadata")
}
