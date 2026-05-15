-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_metadata ADD COLUMN tombstoned_at TIMESTAMPTZ DEFAULT 'infinity';
ALTER TABLE torrents_metadata ALTER COLUMN tombstoned_at SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_metadata DROP COLUMN IF EXISTS tombstoned_at;
-- +goose StatementEnd