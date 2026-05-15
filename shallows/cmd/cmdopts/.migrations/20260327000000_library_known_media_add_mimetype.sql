-- +goose Up
-- +goose StatementBegin
ALTER TABLE library_known_media ADD COLUMN IF NOT EXISTS mimetype TEXT DEFAULT 'application/octet-stream';
UPDATE library_known_media SET mimetype = DEFAULT;
COMMIT;
BEGIN;
ALTER TABLE library_known_media ALTER COLUMN mimetype SET NOT NULL;
COMMENT ON COLUMN library_known_media.mimetype IS 'content category classification using MIME type prefixes (video, audio, image)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE library_known_media DROP COLUMN IF EXISTS mimetype;
-- +goose StatementEnd
