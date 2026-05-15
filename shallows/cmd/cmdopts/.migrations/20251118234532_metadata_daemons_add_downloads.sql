-- +goose Up
-- +goose StatementBegin
ALTER TABLE meta_daemons ADD COLUMN downloads BOOLEAN;
ALTER TABLE meta_daemons ALTER COLUMN downloads SET DEFAULT FALSE;
UPDATE meta_daemons SET downloads = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE meta_daemons ALTER COLUMN downloads SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meta_daemons DROP COLUMN downloads;
-- +goose StatementEnd
