-- +goose Up
-- +goose StatementBegin
ALTER TABLE authz_meta ADD COLUMN remote_control BOOLEAN;
ALTER TABLE authz_meta ALTER COLUMN remote_control SET DEFAULT FALSE;
UPDATE authz_meta SET remote_control = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE authz_meta ALTER COLUMN remote_control SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE authz_meta DROP COLUMN remote_control;
-- +goose StatementEnd
