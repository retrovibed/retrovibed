-- +goose Up
-- +goose StatementBegin
ALTER TABLE meta_wireguard ADD COLUMN maximum_connections UBIGINT;
ALTER TABLE meta_wireguard ALTER COLUMN maximum_connections SET DEFAULT (~0::UBIGINT)::UBIGINT;
UPDATE meta_wireguard SET maximum_connections = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE meta_wireguard ALTER COLUMN maximum_connections SET NOT NULL;
COMMENT ON COLUMN meta_wireguard.maximum_connections IS 'maximum concurrent connections';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meta_wireguard DROP COLUMN maximum_connections;
-- +goose StatementEnd
