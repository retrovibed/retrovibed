-- +goose Up
-- +goose StatementBegin
ALTER TABLE library_known_media ADD COLUMN sync_id UUID DEFAULT '00000000-0000-0000-0000-000000000000';
UPDATE library_known_media SET sync_id = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE library_known_media ALTER COLUMN sync_id SET NOT NULL;
COMMENT ON COLUMN library_known_media.sync_id IS 'deeppool golden copy sync cursor (uuidv7) for this record; nil uuid when the record did not come from the golden copy';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE library_known_media DROP COLUMN IF EXISTS sync_id;
-- +goose StatementEnd
