-- +goose Up
-- +goose StatementBegin
ALTER TABLE library_metadata ADD COLUMN IF NOT EXISTS parent_id UUID DEFAULT '00000000-0000-0000-0000-000000000000';
UPDATE library_metadata SET parent_id = DEFAULT;
COMMIT;
BEGIN;
ALTER TABLE library_metadata ALTER COLUMN parent_id SET NOT NULL;
COMMENT ON COLUMN library_metadata.parent_id IS 'the containing folder, itself a library_metadata row with mimetype inode/directory. the zero uuid is the root of the library.';
CREATE INDEX idx_library_metadata_parent ON library_metadata(parent_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_library_metadata_parent;
ALTER TABLE library_metadata DROP COLUMN IF EXISTS parent_id;
-- +goose StatementEnd
