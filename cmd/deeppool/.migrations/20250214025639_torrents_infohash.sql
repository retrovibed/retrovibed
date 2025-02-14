-- +goose Up
-- +goose StatementBegin
CREATE TABLE torrents_unknown_infohashes (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    infohash BINARY NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_unknown_infohashes;
-- +goose StatementEnd
