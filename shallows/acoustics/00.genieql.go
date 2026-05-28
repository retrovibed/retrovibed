//go:build genieql.generate
// +build genieql.generate

package acoustics

import (
	"context"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

func AudioLSH(gql genieql.Structure) {
	gql.From(gql.Table("audio_lsh"))
}

func AudioLSHScanner(gql genieql.Scanner, pattern func(i AudioLSH)) {
	gql.ColumnNamePrefix("audio_lsh.")
}

func AudioLSHInsert(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a AudioLSH) NewAudioLSHScannerStaticRow,
) {
	gql.Into("audio_lsh").Conflict("ON CONFLICT (table_id, hash, media_id) DO NOTHING")
}

func AudioLSHDeleteByMediaID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, mediaID string) NewAudioLSHScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM audio_lsh WHERE "media_id" = {mediaID} RETURNING ` + AudioLSHScannerStaticColumns)
}

func AudioStats(gql genieql.Structure) {
	gql.From(gql.Table("audio_stats"))
}

func AudioStatsScanner(gql genieql.Scanner, pattern func(i AudioStats)) {
	gql.ColumnNamePrefix("audio_stats.")
}

func AudioStatsInsert(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a AudioStats) NewAudioStatsScannerStaticRow,
) {
	gql.Into("audio_stats").Conflict("ON CONFLICT (dimension) DO UPDATE SET count = EXCLUDED.count, sum = EXCLUDED.sum, sum_sq = EXCLUDED.sum_sq")
}
