-- +goose Up
-- +goose StatementBegin
CREATE TABLE torrents_metadata (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    hidden_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    initiated_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    paused_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    next_announce_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    seeding boolean NOT NULL DEFAULT 'false',
    private boolean NOT NULL DEFAULT 'false',
    tracker VARCHAR NOT NULL, -- will convert this to an array later.
    bytes UBIGINT NOT NULL DEFAULT 0,
    downloaded UBIGINT NOT NULL DEFAULT 0,
    uploaded UBIGINT NOT NULL DEFAULT 0,
    peers USMALLINT NOT NULL DEFAULT 0,
    description STRING NOT NULL DEFAULT '',
    auto_description STRING NOT NULL DEFAULT '', -- application managed description
    infohash BINARY NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_metadata;
-- +goose StatementEnd
