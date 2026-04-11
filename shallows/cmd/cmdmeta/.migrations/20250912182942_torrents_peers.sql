-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_peers;
CREATE TABLE torrents_peers (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    tombstoned_at TIMESTAMPTZ NOT NULL,
    next_check TIMESTAMPTZ NOT NULL DEFAULT '-infinity',
    peer BINARY NOT NULL,
    network VARCHAR NOT NULL,
    ip INET NOT NULL DEFAULT '::',
    port USMALLINT NOT NULL,
    bep51 boolean NOT NULL DEFAULT false,
    bep51_ttl USMALLINT NOT NULL DEFAULT 0,
    bep51_available UBIGINT NOT NULL DEFAULT 0,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_peers;
-- +goose StatementEnd