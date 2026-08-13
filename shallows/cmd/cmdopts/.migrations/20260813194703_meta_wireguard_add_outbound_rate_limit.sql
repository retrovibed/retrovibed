-- +goose Up
-- +goose StatementBegin
ALTER TABLE meta_wireguard ADD COLUMN outbound_rate_limit UINTEGER;
ALTER TABLE meta_wireguard ALTER COLUMN outbound_rate_limit SET DEFAULT 0;
UPDATE meta_wireguard SET outbound_rate_limit = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE meta_wireguard ALTER COLUMN outbound_rate_limit SET NOT NULL;
COMMENT ON COLUMN meta_wireguard.outbound_rate_limit IS 'maximum torrent dial rate in connections per second; 0 means use the torrent daemon default';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meta_wireguard DROP COLUMN outbound_rate_limit;
-- +goose StatementEnd
