-- +goose Up
-- +goose StatementBegin
ALTER TABLE meta_wireguard ADD COLUMN dht BOOLEAN DEFAULT true;
UPDATE meta_wireguard SET dht = DEFAULT;
ALTER TABLE meta_wireguard ADD COLUMN rate_limit_inbound UINTEGER DEFAULT 0;
UPDATE meta_wireguard SET rate_limit_inbound = DEFAULT;
COMMIT;
BEGIN;
ALTER TABLE meta_wireguard ALTER COLUMN dht SET NOT NULL;
ALTER TABLE meta_wireguard ALTER COLUMN rate_limit_inbound SET NOT NULL;
COMMENT ON COLUMN meta_wireguard.dht IS 'whether the dht participates in the network when this configuration is used; false runs the dht passively (no dht traffic, node identity retained)';
COMMENT ON COLUMN meta_wireguard.rate_limit_inbound IS 'maximum inbound accept rate in connections per second; 0 means use the torrent settings';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meta_wireguard DROP COLUMN rate_limit_inbound;
ALTER TABLE meta_wireguard DROP COLUMN dht;
-- +goose StatementEnd
