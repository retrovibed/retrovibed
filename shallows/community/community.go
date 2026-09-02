package community

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
)

// CommunityURLFromDomain builds the standard hosted community url for the given
// domain label, e.g. "myslug" -> "https://myslug.community.retrovibe.space".
func CommunityURLFromDomain(domain string) string {
	return fmt.Sprintf("https://%s.%s", domain, envx.String("community.retrovibe.space", env.CommunityHost))
}

// CommunityDomainFromURL extracts the community subdomain label from a hosted
// community url, e.g. "https://myslug.community.retrovibe.space" -> "myslug".
// Urls that aren't hosted community urls are returned unchanged.
func CommunityDomainFromURL(uri string) string {
	host := envx.String("community.retrovibe.space", env.CommunityHost)

	u, err := url.Parse(strings.TrimSpace(uri))
	if err != nil || u.Host == "" {
		return uri
	}

	if !strings.HasSuffix(u.Host, host) {
		return u.Host
	}

	return langx.FirstNonZero(strings.Split(u.Host, ".")...)
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

// CommunityQueryAccountID matches communities owned by accountID; uuid.Nil
// (i.e. unset) matches every community, applying no filter.
func CommunityQueryAccountID(accountID string) squirrel.Sqlizer {
	if accountID == uuid.Nil.String() {
		return squirrelx.Noop{}
	}

	return squirrel.Eq{"community.account_id": accountID}
}
