-- +goose Up
-- +goose StatementBegin
CREATE TABLE library_metadata (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    hidden_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    tombstoned_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    description STRING NOT NULL DEFAULT '',
    auto_description STRING NOT NULL DEFAULT '', -- application managed description
    mimetype STRING NOT NULL DEFAULT 'application/octet-stream',
    image STRING NOT NULL DEFAULT '',
    bytes UBIGINT NOT NULL,
    torrent_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    archive_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS library_metadata;
-- +goose StatementEnd

