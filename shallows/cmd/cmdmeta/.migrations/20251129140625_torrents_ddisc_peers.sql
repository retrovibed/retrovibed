-- +goose Up
-- +goose StatementBegin
CREATE TABLE ddisc_peers (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_syncd TIMESTAMPTZ NOT NULL DEFAULT now(),
    peer BINARY NOT NULL,
    partition UUID NOT NULL,
    syncoffset UUID NOT NULL -- uuid v7 synchronization id.
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ddisc_peers;
-- +goose StatementEnd
