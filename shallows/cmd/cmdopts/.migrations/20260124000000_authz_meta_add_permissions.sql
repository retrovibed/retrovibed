-- +goose Up
-- +goose StatementBegin
ALTER TABLE authz_meta ADD COLUMN billing_read BOOLEAN;
ALTER TABLE authz_meta ALTER COLUMN billing_read SET DEFAULT FALSE;
UPDATE authz_meta SET billing_read = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE authz_meta ALTER COLUMN billing_read SET NOT NULL;

ALTER TABLE authz_meta ADD COLUMN billing_modify BOOLEAN;
ALTER TABLE authz_meta ALTER COLUMN billing_modify SET DEFAULT FALSE;
UPDATE authz_meta SET billing_modify = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE authz_meta ALTER COLUMN billing_modify SET NOT NULL;

ALTER TABLE authz_meta ADD COLUMN community_modify BOOLEAN;
ALTER TABLE authz_meta ALTER COLUMN community_modify SET DEFAULT FALSE;
UPDATE authz_meta SET community_modify = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE authz_meta ALTER COLUMN community_modify SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE authz_meta DROP COLUMN billing_read;
ALTER TABLE authz_meta DROP COLUMN billing_modify;
ALTER TABLE authz_meta DROP COLUMN community_modify;
-- +goose StatementEnd
