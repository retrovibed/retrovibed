-- +goose Up
-- +goose StatementBegin
COMMIT;

BEGIN;
  DROP INDEX IF EXISTS idx_published_content_community_id;
  DROP INDEX IF EXISTS idx_published_content_library_id;
  DROP INDEX IF EXISTS idx_published_content_community_library;
COMMIT;

BEGIN;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
-- +goose StatementEnd
