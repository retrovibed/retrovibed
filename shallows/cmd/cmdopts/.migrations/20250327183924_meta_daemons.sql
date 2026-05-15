-- +goose Up
-- +goose StatementBegin
CREATE TABLE meta_daemons (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    hostname TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS meta_daemons;
-- +goose StatementEnd
