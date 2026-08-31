-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS plugin_publishers;

CREATE TABLE plugin_publishers (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    path TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    mimetype TEXT NOT NULL DEFAULT ''
);

DROP TABLE IF EXISTS community_publisher;

CREATE TABLE community_publisher (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    community_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    publisher_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_community_publisher_community_publisher ON community_publisher(community_id, publisher_id);
CREATE INDEX IF NOT EXISTS idx_community_publisher_community_id ON community_publisher(community_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS community_publisher;
DROP TABLE IF EXISTS plugin_publishers;
-- +goose StatementEnd
