package community

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

type PublishedContentOption func(*PublishedContent)

func NewPublishedContent(prototype PublishedContent, options ...PublishedContentOption) PublishedContent {
	for _, o := range options {
		o(&prototype)
	}
	prototype.KnownMediaID = langx.FirstNonZero(prototype.KnownMediaID, uuid.Nil.String())
	prototype.OAuthGoogleID = langx.FirstNonZero(prototype.OAuthGoogleID, uuid.Nil.String())
	prototype.PublishedAt = langx.FirstNonZero(prototype.PublishedAt, timex.Inf())
	return prototype
}

// PublishedContentOptionFromDB converts a database model to proto options.
func PublishedContentOptionFromDB(pc PublishedContent) func(*meta.PublishedContent) {
	return func(p *meta.PublishedContent) {
		p.Id = pc.ID
		p.Title = pc.Title
		p.Description = pc.Description
		p.CommunityId = pc.CommunityID
		p.KnownMediaId = pc.KnownMediaID
		p.MagnetUri = pc.MagnetURI
		p.LibraryId = pc.LibraryID
		p.OauthGoogleId = pc.OAuthGoogleID
		p.PublishedAt = grpcx.EncodeTime(pc.PublishedAt)
		p.CreatedAt = grpcx.EncodeTime(pc.CreatedAt)
		p.UpdatedAt = grpcx.EncodeTime(pc.UpdatedAt)
		p.Bytes = pc.Bytes
	}
}

func PublishedContentSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) PublishedContentScanner {
	return NewPublishedContentScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func PublishedContentSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(PublishedContentScannerStaticColumns)...).From("published_content")
}

func PublishedContentQueryCommunityID(cid string) squirrel.Sqlizer {
	return squirrel.Expr("published_content.community_id = ?", cid)
}

func PublishedContentQueryNotTombstoned() squirrel.Sqlizer {
	return squirrel.Expr("published_content.tombstoned_at = 'infinity'")
}
