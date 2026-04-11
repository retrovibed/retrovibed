-- +goose Up
-- +goose StatementBegin
CREATE TABLE meta_wireguard (
    id UUID PRIMARY KEY NOT NULL, -- md5 of the contents
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    "default" BOOLEAN DEFAULT 'f' NOT NULL,
    port USMALLINT NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS meta_wireguard;
-- +goose StatementEnd