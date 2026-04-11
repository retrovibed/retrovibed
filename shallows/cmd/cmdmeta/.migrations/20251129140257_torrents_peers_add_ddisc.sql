-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_peers ADD COLUMN "ddisc" UUID DEFAULT '00000000-0000-0000-0000-000000000000'::UUID;
COMMENT ON COLUMN torrents_peers.ddisc IS 'distributed discovery partition this peer maintains, two special values. "00000000-0000-0000-0000-000000000000" = doesnt maintain any partitions. "11111111-1111-1111-1111-111111111111" maintains all partitions';
UPDATE torrents_peers SET "ddisc" = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE torrents_peers ALTER COLUMN "ddisc" SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_peers DROP COLUMN "ddisc";
-- +goose StatementEnd
