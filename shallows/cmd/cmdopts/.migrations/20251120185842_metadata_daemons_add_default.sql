-- +goose Up
-- +goose StatementBegin
ALTER TABLE meta_daemons ADD COLUMN "default" BOOLEAN;
ALTER TABLE meta_daemons ALTER COLUMN "default" SET DEFAULT FALSE;
UPDATE meta_daemons SET "default" = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE meta_daemons ALTER COLUMN "default" SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meta_daemons DROP COLUMN "default";
-- +goose StatementEnd
