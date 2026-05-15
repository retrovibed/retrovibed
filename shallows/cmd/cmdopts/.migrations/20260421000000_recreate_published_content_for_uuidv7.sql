-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS published_content;
CREATE TABLE published_content (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    community_id UUID NOT NULL,
    known_media_id UUID NOT NULL,
    magnet_uri TEXT NOT NULL,
    library_id UUID NOT NULL,
    published_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    publish_mode INTEGER NOT NULL DEFAULT 0,
    oauth_google_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    bytes UBIGINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS published_content;
CREATE TABLE published_content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    community_id UUID NOT NULL,
    known_media_id UUID NOT NULL,
    magnet_uri TEXT NOT NULL,
    library_id UUID NOT NULL,
    published_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    publish_mode INTEGER NOT NULL DEFAULT 0,
    oauth_google_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    bytes UBIGINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd
