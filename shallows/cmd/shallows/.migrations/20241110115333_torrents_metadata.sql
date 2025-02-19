-- +goose Up
-- +goose StatementBegin
CREATE TABLE torrents_metadata (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    bytes UBIGINT NOT NULL,
    description STRING NOT NULL DEFAULT '',
    infohash BINARY NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_metadata;
-- +goose StatementEnd
