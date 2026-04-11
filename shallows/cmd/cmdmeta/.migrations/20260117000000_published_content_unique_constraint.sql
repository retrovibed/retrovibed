-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
-- +goose StatementEnd
