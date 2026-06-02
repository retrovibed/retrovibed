-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
ALTER TABLE published_content ADD COLUMN title TEXT DEFAULT '';
UPDATE published_content SET title = DEFAULT;
ALTER TABLE published_content ADD COLUMN description TEXT DEFAULT '';
UPDATE published_content SET description = DEFAULT;
COMMIT;
BEGIN;
ALTER TABLE published_content ALTER COLUMN title SET NOT NULL;
ALTER TABLE published_content ALTER COLUMN description SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE published_content DROP COLUMN IF EXISTS title;
ALTER TABLE published_content DROP COLUMN IF EXISTS description;
-- +goose StatementEnd
