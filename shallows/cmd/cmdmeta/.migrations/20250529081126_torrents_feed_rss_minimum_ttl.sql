-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_feed_rss ADD COLUMN ttl_minimum INTERVAL DEFAULT INTERVAL '0 seconds';
ALTER TABLE torrents_feed_rss ALTER COLUMN ttl_minimum SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_feed_rss DROP COLUMN ttl_minimum;
-- +goose StatementEnd
