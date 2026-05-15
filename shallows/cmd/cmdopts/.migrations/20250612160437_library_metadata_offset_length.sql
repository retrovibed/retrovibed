-- +goose Up
-- +goose StatementBegin
ALTER TABLE library_metadata ADD COLUMN disk_offset UBIGINT DEFAULT 0;
ALTER TABLE library_metadata ALTER COLUMN disk_offset SET NOT NULL;
COMMENT ON COLUMN library_metadata.disk_offset IS 'specifies the offset into the media block storage layer';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE library_metadata DROP COLUMN IF EXISTS disk_offset;
-- +goose StatementEnd