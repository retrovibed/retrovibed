-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_metadata ADD COLUMN imported_at TIMESTAMPTZ DEFAULT 'infinity';
ALTER TABLE torrents_metadata ALTER COLUMN imported_at SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_metadata DROP COLUMN IF EXISTS imported_at;
-- +goose StatementEnd
