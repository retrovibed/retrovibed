package community

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
)

func PluginPublisherSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) PluginPublisherScanner {
	return NewPluginPublisherScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func PluginPublisherSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(PluginPublisherScannerStaticColumns)...).From("plugin_publishers")
}
