-- +goose Up
-- +goose StatementBegin
ALTER TABLE library_known_media ADD COLUMN tombstoned_at TIMESTAMPTZ DEFAULT 'infinity';
ALTER TABLE library_known_media ALTER COLUMN tombstoned_at SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE library_known_media DROP COLUMN IF EXISTS tombstoned_at;
-- +goose StatementEnd
