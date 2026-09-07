-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_metadata ADD COLUMN IF NOT EXISTS available UBIGINT DEFAULT 0;
UPDATE torrents_metadata SET available = DEFAULT;
COMMIT;
BEGIN;
ALTER TABLE torrents_metadata ALTER COLUMN available SET NOT NULL;
COMMENT ON COLUMN torrents_metadata.available IS 'bytes verified present on disk, regardless of source (peer transfer, local authoring, resume verification). distinct from downloaded, which is real bytes fetched from peers.';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_metadata DROP COLUMN IF EXISTS available;
-- +goose StatementEnd
