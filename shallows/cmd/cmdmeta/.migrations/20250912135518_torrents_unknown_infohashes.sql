-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_unknown_infohashes;
CREATE TABLE torrents_unknown_infohashes (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
    next_check TIMESTAMPTZ NOT NULL DEFAULT '-infinity',
    attempts UBIGINT NOT NULL DEFAULT 0,
    peer BINARY NOT NULL,
    ip INET NOT NULL DEFAULT '::',
    port USMALLINT NOT NULL DEFAULT 0,
    infohash BINARY NOT NULL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_unknown_infohashes;
-- +goose StatementEnd