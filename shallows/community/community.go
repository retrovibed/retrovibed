package community

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
)

// CommunityURLFromDomain builds the standard hosted community url for the given
// domain label, e.g. "myslug" -> "https://myslug.community.retrovibe.space".
func CommunityURLFromDomain(domain string) string {
	return fmt.Sprintf("https://%s.community.retrovibe.space", domain)
}

func CommunitySearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) CommunityScanner {
	return NewCommunityScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func CommunitySearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(CommunityScannerStaticColumns)...).From("community")
}

func CommunityQueryNotTombstoned() squirrel.Sqlizer {
	return squirrel.Expr("community.tombstoned_at = 'infinity'")
}
