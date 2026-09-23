-- +goose Up
-- +goose StatementBegin
ALTER TABLE community RENAME COLUMN last_sync_at TO next_sync_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE community RENAME COLUMN next_sync_at TO last_sync_at;
-- +goose StatementEnd
