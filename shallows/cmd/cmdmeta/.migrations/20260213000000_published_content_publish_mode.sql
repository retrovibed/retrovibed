-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
ALTER TABLE published_content ADD COLUMN IF NOT EXISTS publish_mode INTEGER DEFAULT 0;
UPDATE published_content SET publish_mode = 0;
COMMIT;
BEGIN;
ALTER TABLE published_content ALTER COLUMN publish_mode SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
ALTER TABLE published_content DROP COLUMN IF EXISTS publish_mode;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd
