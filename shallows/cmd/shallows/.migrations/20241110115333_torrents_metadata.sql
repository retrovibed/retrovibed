-- +goose Up
-- +goose StatementBegin
CREATE TABLE torrents_metadata (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    hidden_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    initiated_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    paused_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    tracker VARCHAR NOT NULL, -- will convert this to an array later.
    bytes UBIGINT NOT NULL,
    downloaded UBIGINT NOT NULL,
    peers USMALLINT NOT NULL,
    description STRING NOT NULL DEFAULT '',
    infohash BINARY NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_metadata;
-- +goose StatementEnd
-- ALTER TABLE torrents_metadata ADD COLUMN peers_pending USMALLINT DEFAULT 0;
-- ALTER TABLE torrents_metadata ADD COLUMN tracker VARCHAR DEFAULT '';
