-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_metadata ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ DEFAULT 'infinity';
UPDATE torrents_metadata SET expires_at = DEFAULT;
COMMIT;
BEGIN;
ALTER TABLE torrents_metadata ALTER COLUMN expires_at SET NOT NULL;
COMMENT ON COLUMN torrents_metadata.expires_at IS 'the timestamp when the torrent is no longer considered valid, useful for things like distro install images, media metadata, anything with a temporal lifespan';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_metadata DROP COLUMN IF EXISTS expires_at;
-- +goose StatementEnd
