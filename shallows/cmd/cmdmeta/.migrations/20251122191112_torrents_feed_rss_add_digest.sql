-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_feed_rss ADD COLUMN "digest" UUID;
ALTER TABLE torrents_feed_rss ALTER COLUMN "digest" SET DEFAULT '00000000-0000-0000-0000-000000000000'::UUID;
UPDATE torrents_feed_rss SET "digest" = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE torrents_feed_rss ALTER COLUMN "digest" SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_feed_rss DROP COLUMN "digest";
-- +goose StatementEnd