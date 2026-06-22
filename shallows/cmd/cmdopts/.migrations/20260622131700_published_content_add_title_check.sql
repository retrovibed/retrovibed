-- +goose Up
-- +goose StatementBegin

-- 1. Backfill existing rows that would violate the new constraint.
UPDATE published_content SET title = id::TEXT WHERE title = '';

-- 2. Rename existing table
DROP INDEX IF EXISTS idx_published_content_community_library;
ALTER TABLE published_content RENAME TO published_content_old;

-- 3. Recreate the table with the title check constraint
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
    bytes UBIGINT NOT NULL DEFAULT 0,
    mimetype TEXT NOT NULL DEFAULT 'application/octet-stream',
    tombstoned_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    title TEXT NOT NULL DEFAULT '' CHECK (title != ''),
    description TEXT NOT NULL DEFAULT ''
);

-- 4. Move the data
INSERT INTO published_content
SELECT id, community_id, known_media_id, magnet_uri, library_id, published_at,
       created_at, updated_at, publish_mode, oauth_google_id, bytes, mimetype,
       tombstoned_at, title, description
FROM published_content_old;

-- 5. Finalize
DROP TABLE published_content_old;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_published_content_community_library;
ALTER TABLE published_content RENAME TO published_content_old;

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
    bytes UBIGINT NOT NULL DEFAULT 0,
    mimetype TEXT NOT NULL DEFAULT 'application/octet-stream',
    tombstoned_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    title TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT ''
);

INSERT INTO published_content
SELECT id, community_id, known_media_id, magnet_uri, library_id, published_at,
       created_at, updated_at, publish_mode, oauth_google_id, bytes, mimetype,
       tombstoned_at, title, description
FROM published_content_old;

DROP TABLE published_content_old;
CREATE UNIQUE INDEX IF NOT EXISTS idx_published_content_community_library ON published_content(community_id, library_id);
-- +goose StatementEnd
