-- +goose Up
-- +goose StatementBegin
ALTER TABLE library_metadata ADD COLUMN IF NOT EXISTS disk_usage UBIGINT DEFAULT 0;
ALTER TABLE library_metadata ALTER COLUMN disk_usage SET NOT NULL;
UPDATE library_metadata SET disk_usage = bytes;

ALTER TABLE library_metadata ADD COLUMN IF NOT EXISTS quota_usage UBIGINT DEFAULT 0;
ALTER TABLE library_metadata ALTER COLUMN quota_usage SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE library_metadata DROP COLUMN IF EXISTS disk_usage;
ALTER TABLE library_metadata DROP COLUMN IF EXISTS quota_usage;
-- +goose StatementEnd
