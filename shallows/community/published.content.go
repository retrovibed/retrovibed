package community

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
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

// PublishedContentOptionTestDefaults normalizes the identifier/temporal fields
// to valid, deterministic values for tests. Combine with testx.Fake so free-form
// fields (Title, Description, MagnetURI, ...) are populated with non-empty fakes
// rather than Go zero values, which would otherwise trip the title CHECK constraint.
func PublishedContentOptionTestDefaults(t *PublishedContent) {
	t.ID = uuid.Nil.String()
	t.CommunityID = uuid.Nil.String()
	t.KnownMediaID = uuid.Nil.String()
	t.LibraryID = uuid.Nil.String()
	t.OAuthGoogleID = uuid.Nil.String()
	t.PublishMode = 0
	t.PublishedAt = timex.Inf()
	t.TombstonedAt = timex.Inf()
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
