-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
ALTER TABLE published_content ADD COLUMN tombstoned_at TIMESTAMPTZ DEFAULT 'infinity';
UPDATE published_content SET tombstoned_at = DEFAULT;
COMMIT;
BEGIN;
ALTER TABLE published_content ALTER COLUMN tombstoned_at SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE published_content DROP COLUMN IF EXISTS tombstoned_at;
-- +goose StatementEnd
