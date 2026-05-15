-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_metadata ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ DEFAULT 'infinity';
UPDATE torrents_metadata SET completed_at = DEFAULT;
COMMIT;
BEGIN;
ALTER TABLE torrents_metadata ALTER COLUMN completed_at SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_metadata DROP COLUMN IF EXISTS completed_at;
-- +goose StatementEnd
