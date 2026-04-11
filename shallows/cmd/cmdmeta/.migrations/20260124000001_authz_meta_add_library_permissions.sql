-- +goose Up
-- +goose StatementBegin
ALTER TABLE authz_meta ADD COLUMN library_read BOOLEAN;
ALTER TABLE authz_meta ALTER COLUMN library_read SET DEFAULT FALSE;
UPDATE authz_meta SET library_read = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE authz_meta ALTER COLUMN library_read SET NOT NULL;

ALTER TABLE authz_meta ADD COLUMN library_modify BOOLEAN;
ALTER TABLE authz_meta ALTER COLUMN library_modify SET DEFAULT FALSE;
UPDATE authz_meta SET library_modify = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE authz_meta ALTER COLUMN library_modify SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE authz_meta DROP COLUMN library_read;
ALTER TABLE authz_meta DROP COLUMN library_modify;
-- +goose StatementEnd
