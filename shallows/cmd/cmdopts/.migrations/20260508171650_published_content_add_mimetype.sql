-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
ALTER TABLE published_content ADD COLUMN IF NOT EXISTS mimetype TEXT DEFAULT 'application/octet-stream';
UPDATE published_content SET mimetype = library_metadata.mimetype FROM library_metadata WHERE published_content.library_id = library_metadata.id;
COMMIT;
BEGIN;
ALTER TABLE published_content ALTER COLUMN mimetype SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE published_content DROP COLUMN IF EXISTS mimetype;
-- +goose StatementEnd
