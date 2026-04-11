-- +goose Up
-- +goose StatementBegin
ALTER TABLE ddisc_media ADD COLUMN "partition" UUID DEFAULT '00000000-0000-0000-0000-000000000000'::UUID;
COMMENT ON COLUMN ddisc_media.partition IS 'uuid used to track what partition discovered media belongs to once its been identified';
UPDATE ddisc_media SET "partition" = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE ddisc_media ALTER COLUMN "partition" SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ddisc_media DROP COLUMN "partition";
-- +goose StatementEnd
