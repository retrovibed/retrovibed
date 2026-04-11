-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
ALTER TABLE published_content ADD COLUMN IF NOT EXISTS oauth_google_id UUID DEFAULT '00000000-0000-0000-0000-000000000000';
UPDATE published_content SET oauth_google_id = DEFAULT;
COMMIT;
BEGIN;
ALTER TABLE published_content ALTER COLUMN oauth_google_id SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
ALTER TABLE published_content DROP COLUMN IF EXISTS oauth_google_id;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd
