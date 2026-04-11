package community

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/squirrelx"
)

func CommunityMetricAggregateSearch(ctx context.Context, q sqlx.Queryer, communityID string, periodStart, periodEnd time.Time) (summary CommunityMetric, err error) {
	b := squirrelx.PSQL.Select(
		"COALESCE(MAX(subscribers), 0)",
	).From("community_metrics").Where(squirrel.And{
		squirrel.Eq{"community_id": communityID},
		squirrel.GtOrEq{"period_start": periodStart},
		squirrel.LtOrEq{"period_start": periodEnd},
	})

	row := b.RunWith(q).QueryRowContext(ctx)
	err = row.Scan(&summary.Subscribers)

	summary.CommunityID = communityID
	summary.PeriodStart = periodStart
	summary.PeriodEnd = periodEnd

	return summary, err
}

func PublishedCASMetricAggregateSearch(ctx context.Context, q sqlx.Queryer, communityID string, periodStart, periodEnd time.Time) (totalArchivers int32, err error) {
	query := `
		SELECT COALESCE(SUM(latest.archivers), 0)
		FROM (
			SELECT published_cas_metrics.published_content_id,
			       MAX(published_cas_metrics.archivers) as archivers
			FROM published_cas_metrics
			JOIN published_content ON published_cas_metrics.published_content_id = published_content.id
			WHERE published_content.community_id = $1
			  AND published_cas_metrics.period_start BETWEEN $2 AND $3
			GROUP BY published_cas_metrics.published_content_id
		) latest
	`

	row := q.QueryRowContext(ctx, query, communityID, periodStart, periodEnd)
	err = row.Scan(&totalArchivers)

	return totalArchivers, err
}

func PublishedCASMetricPerContentSearch(ctx context.Context, q sqlx.Queryer, communityID string, periodStart, periodEnd time.Time) (results []PublishedCASMetric, err error) {
	var (
		rows *sql.Rows
		m    PublishedCASMetric
	)

	query := `
		SELECT CAST(published_cas_metrics.published_content_id AS TEXT),
		       MAX(published_cas_metrics.archivers) as archivers,
		       MAX(published_cas_metrics.bytes) as bytes,
		       MAX(published_cas_metrics.revenue) as revenue
		FROM published_cas_metrics
		JOIN published_content ON published_cas_metrics.published_content_id = published_content.id
		WHERE published_content.community_id = $1
		  AND published_cas_metrics.period_start BETWEEN $2 AND $3
		GROUP BY published_cas_metrics.published_content_id
	`

	if rows, err = q.QueryContext(ctx, query, communityID, periodStart, periodEnd); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&m.PublishedContentID, &m.Archivers, &m.Bytes, &m.Revenue); err != nil {
			return nil, err
		}
		results = append(results, m)
	}

	return results, rows.Err()
}
