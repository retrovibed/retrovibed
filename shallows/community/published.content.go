package community

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
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
	prototype.KnownMediaID = langx.FirstNonZero(prototype.KnownMediaID, uuidx.Nil().String())
	prototype.OAuthGoogleID = langx.FirstNonZero(prototype.OAuthGoogleID, uuidx.Nil().String())
	prototype.PublishedAt = langx.FirstNonZero(prototype.PublishedAt, timex.Inf())
	return prototype
}

// PublishedContentOptionTestDefaults normalizes the identifier/temporal fields
// to valid, deterministic values for tests. Combine with testx.Fake so free-form
// fields (Title, Description, MagnetURI, ...) are populated with non-empty fakes
// rather than Go zero values, which would otherwise trip the title CHECK constraint.
func PublishedContentOptionTestDefaults(t *PublishedContent) {
	t.ID = uuidx.Nil().String()
	t.CommunityID = uuidx.Nil().String()
	t.KnownMediaID = uuidx.Nil().String()
	t.LibraryID = uuidx.Nil().String()
	t.OAuthGoogleID = uuidx.Nil().String()
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
