-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
ALTER TABLE published_content ADD COLUMN IF NOT EXISTS bytes UBIGINT DEFAULT 0;
UPDATE published_content SET bytes = DEFAULT;
COMMIT;
BEGIN;
ALTER TABLE published_content ALTER COLUMN bytes SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
ALTER TABLE published_content DROP COLUMN IF EXISTS bytes;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd
