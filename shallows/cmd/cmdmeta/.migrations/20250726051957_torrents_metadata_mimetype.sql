-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_metadata ADD COLUMN mimetype text DEFAULT 'application/x-bittorrent';
ALTER TABLE torrents_metadata ALTER COLUMN mimetype SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_metadata DROP COLUMN IF EXISTS mimetype;
-- +goose StatementEnd