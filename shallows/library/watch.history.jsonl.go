package library

import (
	"context"
	"io"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// WatchHistoryJSONLEncode streams library_watch_history rows matching b as JSONL to w.
func WatchHistoryJSONLEncode(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder, w io.Writer) error {
	scanner := WatchHistorySearch(ctx, q, b)
	defer scanner.Close()

	enc := jsonl.NewEncoder(w)
	for scanner.Next() {
		var row WatchHistory
		if err := scanner.Scan(&row); err != nil {
			return err
		}

		rec := WatchHistoryRecord{
			Id:       row.ID,
			MediaId:  row.MediaID,
			Duration: uint64(row.Duration.Milliseconds()),
		}
		if err := enc.Encode(&rec); err != nil {
			return err
		}
	}

	return scanner.Err()
}
