-- +goose Up
-- +goose StatementBegin
ALTER TABLE authz_meta ADD COLUMN local_only BOOLEAN;
ALTER TABLE authz_meta ALTER COLUMN local_only SET DEFAULT FALSE;
UPDATE authz_meta SET local_only = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE authz_meta ALTER COLUMN local_only SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE authz_meta DROP COLUMN local_only;
-- +goose StatementEnd
