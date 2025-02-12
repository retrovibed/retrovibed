-- +goose Up
-- +goose StatementBegin
INSTALL inet;
LOAD inet;

CREATE TABLE torrents_peers (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    next_check TIMESTAMPTZ NOT NULL DEFAULT '-infinity',
    ip INET NOT NULL,
    port USMALLINT NOT NULL,
    bep51 boolean NOT NULL DEFAULT false,
    bep51_available UBIGINT NOT NULL DEFAULT 0,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_peers;
-- +goose StatementEnd
