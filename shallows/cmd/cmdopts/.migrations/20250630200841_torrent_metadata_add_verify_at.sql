-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_metadata ADD COLUMN verify_at TIMESTAMPTZ DEFAULT 'infinity';
ALTER TABLE torrents_metadata ALTER COLUMN verify_at SET NOT NULL;
COMMENT ON COLUMN torrents_metadata.verify_at IS 'timestamp for when to next verify content on disk';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_metadata DROP COLUMN IF EXISTS verify_at;
-- +goose StatementEnd