-- +goose Up
-- +goose StatementBegin
ALTER TABLE torrents_ddisc_media RENAME TO ddisc_media;
ALTER TABLE ddisc_media ADD COLUMN "sync_uid" UUID DEFAULT '00000000-0000-0000-0000-000000000000'::UUID;
COMMENT ON COLUMN ddisc_media.sync_uid IS 'uuid v7 used to track sync status between peers, peers track the maximum sync_uid of the records they have received';
UPDATE ddisc_media SET "sync_uid" = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE ddisc_media ALTER COLUMN "sync_uid" SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ddisc_media DROP COLUMN "sync_uid";
ALTER TABLE ddisc_media RENAME TO torrents_ddisc_media;
-- +goose StatementEnd