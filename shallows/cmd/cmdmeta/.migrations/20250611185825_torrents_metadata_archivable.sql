-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_metadata ADD COLUMN archivable boolean DEFAULT 'false';
ALTER TABLE torrents_metadata ALTER COLUMN archivable SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_metadata DROP COLUMN IF EXISTS archivable;
-- +goose StatementEnd