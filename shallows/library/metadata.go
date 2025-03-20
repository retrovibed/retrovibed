package library

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
	"github.com/retrovibed/retrovibed/internal/x/duckdbx"
	"github.com/retrovibed/retrovibed/internal/x/langx"
	"github.com/retrovibed/retrovibed/internal/x/sqlx"
	"github.com/retrovibed/retrovibed/internal/x/squirrelx"
	"github.com/retrovibed/retrovibed/internal/x/timex"
)

type MetadataOption = func(*Metadata)

func MetadataOptionDescription(d string) MetadataOption {
	return func(m *Metadata) {
		m.Description = d
	}
}

func MetadataOptionTorrentID(d string) MetadataOption {
	return func(m *Metadata) {
		m.TorrentID = d
	}
}

func MetadataOptionBytes(d uint64) MetadataOption {
	return func(m *Metadata) {
		m.Bytes = d
	}
}

func MetadataOptionMimetype(s string) MetadataOption {
	return func(m *Metadata) {
		m.Mimetype = s
	}
}

func MetadataOptionCompose(options ...func(*Metadata)) MetadataOption {
	return func(m *Metadata) {
		for _, opt := range options {
			opt(m)
		}
	}
}

func MetadataOptionJSONSafeEncode(p *Metadata) {
	p.CreatedAt = timex.RFC3339NanoEncode(p.CreatedAt)
	p.UpdatedAt = timex.RFC3339NanoEncode(p.UpdatedAt)
	p.HiddenAt = timex.RFC3339NanoEncode(p.HiddenAt)
	p.TombstonedAt = timex.RFC3339NanoEncode(p.TombstonedAt)
}

func NewMetadata(id string, options ...func(*Metadata)) (m Metadata) {
	r := langx.Clone(Metadata{
		ID:        id,
		TorrentID: uuid.Nil.String(),
		ArchiveID: uuid.Nil.String(),
	}, options...)

	return r
}

func MetadataQueryVisible() squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.hidden_at = 'infinity'")
}

func MetadataQueryHidden() squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.hidden_at < NOW()")
}

func MetadataQueryArchived() squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.achived_id != '00000000-0000-0000-0000-000000000000'")
}

func MetadataQueryShared() squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.torrent_id != '00000000-0000-0000-0000-000000000000'")
}

func MetadataSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) MetadataScanner {
	return NewMetadataScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func MetadataQuerySearch(q string, columns ...string) squirrel.Sqlizer {
	return duckdbx.FTSSearch("fts_main_library_metadata", q, columns...)
}

func MetadataSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(MetadataScannerStaticColumns)...).From("library_metadata")
}
