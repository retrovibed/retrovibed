package library

import (
	"context"
	"log"
	"time"
	"unicode"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/runesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

type MetadataOption = func(*Metadata)

func MetadataOptionDescription(d string) MetadataOption {
	return func(m *Metadata) {
		m.Description = d
	}
}

func MetadataOptionKnownMediaID(d string) MetadataOption {
	return func(m *Metadata) {
		m.KnownMediaID = d
	}
}

func MetadataOptionAutoDescription(d string) MetadataOption {
	return func(m *Metadata) {
		m.AutoDescription = d
	}
}

func MetadataOptionTorrentID(d string) MetadataOption {
	return func(m *Metadata) {
		m.TorrentID = d
	}
}

func MetadataOptionEncryptionSeed(d string) MetadataOption {
	return func(m *Metadata) {
		m.EncryptionSeed = d
	}
}

func MetadataOptionArchivable(b bool) MetadataOption {
	return func(m *Metadata) {
		if b {
			m.ArchiveID = uuid.Max.String()
		}
	}
}

func MetadataOptionArchiveID(id string) MetadataOption {
	return func(m *Metadata) {
		m.ArchiveID = id
	}
}

func MetadataOptionBytes(d uint64) MetadataOption {
	return func(m *Metadata) {
		m.Bytes = d
	}
}

func MetadataOptionOffset(d uint64) MetadataOption {
	return func(m *Metadata) {
		m.DiskOffset = d
	}
}

func MetadataOptionParentID(id string) MetadataOption {
	return func(m *Metadata) {
		m.ParentID = id
	}
}

func MetadataOptionMimetype(s string) MetadataOption {
	return func(m *Metadata) {
		m.Mimetype = s
	}
}

func MetadataOptionHidden(b bool) MetadataOption {
	return func(m *Metadata) {
		if b {
			m.HiddenAt = time.Now()
		} else {
			m.HiddenAt = timex.Inf()
		}
	}
}

func MetadataOptionCompose(options ...func(*Metadata)) MetadataOption {
	return func(m *Metadata) {
		for _, opt := range options {
			opt(m)
		}
	}
}

func MetadataOptionTestDefaults(p *Metadata) {
	p.ID = uuid.Nil.String()
	p.ArchiveID = uuid.Nil.String()
	p.TorrentID = uuid.Nil.String()
	p.KnownMediaID = uuid.Nil.String()
	p.EncryptionSeed = uuid.Nil.String()
	p.ParentID = uuid.Nil.String()
	p.HiddenAt = timex.Inf()
	p.Bytes = 16 * bytesx.KiB
	p.Mimetype = mimex.Binary
	p.DiskOffset = 0
}

func MetadataOptionTestID(id string) MetadataOption {
	return func(p *Metadata) {
		p.ID = id
	}
}

func MetadataOptionTestRandomID(p *Metadata) {
	p.ID = uuid.Must(uuid.NewV4()).String()
}

func NewMetadata(id string, options ...func(*Metadata)) (m Metadata) {
	r := langx.Clone(Metadata{
		ID:             id,
		TorrentID:      uuid.Nil.String(),
		ArchiveID:      uuid.Nil.String(),
		KnownMediaID:   uuid.Nil.String(),
		EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		ParentID:       uuid.Nil.String(),
		HiddenAt:       timex.Inf(),
	}, options...)

	return r
}

func MetadataQueryRandomAfter(id string) squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.id >= ?", id)
}

func MetadataQueryRandomBefore(id string) squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.id < ?", id)
}

func MetadataQueryMimetypes(mimes ...string) squirrel.Sqlizer {
	return squirrelx.In("library_metadata.mimetype", mimes...)
}

func MetadataQueryUpdatedBetween(r timex.Range) squirrel.Sqlizer {
	return squirrelx.Between("library_metadata.updated_at", r.Start, r.End)
}

func MetadataQueryNotTombstoned() squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.tombstoned_at = 'infinity'")
}

func MetadataQueryTombstoned() squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.tombstoned_at < 'infinity'")
}

func MetadataQueryByTorrentID(tid string) squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.torrent_id = ?", tid)
}

func MetadataQueryParent(id string) squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.parent_id = ?", id)
}

func MetadataQueryDirectory(b bool) squirrel.Sqlizer {
	if b {
		return squirrel.Expr("library_metadata.mimetype = ?", mimex.Directory)
	}

	return squirrel.Expr("library_metadata.mimetype != ?", mimex.Directory)
}

func MetadataQueryHidden(b bool) squirrel.Sqlizer {
	if !b {
		return squirrel.Expr("library_metadata.hidden_at = 'infinity'")
	}

	return squirrel.Expr("library_metadata.hidden_at < NOW()")
}

func MetadataQueryArchived() squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.archive_id != '00000000-0000-0000-0000-000000000000' AND library_metadata.archive_id != 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'")
}

func MetadataQueryArchivable() squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.archive_id = 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'")
}

func MetadataQueryShared() squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.torrent_id != '00000000-0000-0000-0000-000000000000' AND library_metadata.archive_id != 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'")
}

func MetadataQueryNotIndexed() squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.auto_description == ''")
}

func MetadataQueryNeedsKnownMediaID() squirrel.Sqlizer {
	return squirrel.Expr("library_metadata.known_media_id = 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'")
}

func MetadataQueryHasKnownMedia(has bool) squirrel.Sqlizer {
	if has {
		return squirrel.Expr("library_metadata.known_media_id != '00000000-0000-0000-0000-000000000000' AND library_metadata.known_media_id != 'ffffffff-ffff-ffff-ffff-ffffffffffff'")
	}
	return squirrel.Expr("library_metadata.known_media_id IN ('00000000-0000-0000-0000-000000000000', 'ffffffff-ffff-ffff-ffff-ffffffffffff')")
}

func MetadataQueryHasTorrent(has bool) squirrel.Sqlizer {
	if has {
		return squirrel.Expr("library_metadata.torrent_id != '00000000-0000-0000-0000-000000000000'")
	}
	return squirrel.Expr("library_metadata.torrent_id = '00000000-0000-0000-0000-000000000000'")
}

func MetadataQueryByIDs(ids ...string) squirrel.Sqlizer {
	if len(ids) == 0 {
		return squirrelx.Noop{}
	}

	return squirrelx.In("library_metadata.id", ids...)
}

func MetadataSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) MetadataScanner {
	return NewMetadataScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func MetadataSearchOne(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) MetadataScanner {
	return NewMetadataScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func MetadataSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(MetadataScannerStaticColumns)...).From("library_metadata")
}

func MetadataDiskStorageUsage(ctx context.Context, q sqlx.Queryer) sqlx.IntRowScanner {
	return sqlx.NewIntRowScanner(q.QueryRowContext(ctx, "SELECT SUM(bytes) FROM library_metadata"))
}

// normalizes the text input. if n error occurs it'll log and return the input unmodified.
func NormalizedDescription(s string) string {
	transformer := transform.Chain(runes.ReplaceIllFormed(), runesx.Replace(unicode.IsPunct, ' '))
	o, _, err := transform.String(transformer, s)
	if err == nil {
		return stringsx.CompactWhitespace(o)
	}

	log.Println("unable to normalize description", err)
	return s
}
