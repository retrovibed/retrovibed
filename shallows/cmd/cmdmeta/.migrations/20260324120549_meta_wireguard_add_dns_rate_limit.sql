-- +goose Up
-- +goose StatementBegin
ALTER TABLE meta_wireguard ADD COLUMN dns_rate_limit UINTEGER;
ALTER TABLE meta_wireguard ALTER COLUMN dns_rate_limit SET DEFAULT 0;
UPDATE meta_wireguard SET dns_rate_limit = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE meta_wireguard ALTER COLUMN dns_rate_limit SET NOT NULL;
COMMENT ON COLUMN meta_wireguard.dns_rate_limit IS 'maximum DNS lookup rate in events per second; 0 means unlimited';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meta_wireguard DROP COLUMN dns_rate_limit;
-- +goose StatementEnd
