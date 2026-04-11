-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_feed_rss ADD COLUMN last_built_at TIMESTAMPTZ DEFAULT 'infinity';
ALTER TABLE torrents_feed_rss ALTER COLUMN last_built_at SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_feed_rss DROP COLUMN last_built_at;
-- +goose StatementEnd
