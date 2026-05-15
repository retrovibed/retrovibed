-- +goose Up
-- +goose StatementBegin
ALTER TABLE library_metadata ADD COLUMN encryption_seed uuid DEFAULT gen_random_uuid();
ALTER TABLE library_metadata ALTER COLUMN encryption_seed SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE library_metadata DROP COLUMN IF EXISTS encryption_seed;
-- +goose StatementEnd
