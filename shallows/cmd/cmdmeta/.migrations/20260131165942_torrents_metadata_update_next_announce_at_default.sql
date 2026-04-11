-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_metadata ALTER COLUMN "next_announce_at" SET DEFAULT '-infinity';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_metadata ALTER COLUMN "next_announce_at" SET DEFAULT 'infinity';
-- +goose StatementEnd
