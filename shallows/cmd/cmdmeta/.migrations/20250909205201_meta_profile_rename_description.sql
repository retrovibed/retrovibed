-- +goose Up
-- +goose StatementBegin
ALTER TABLE meta_profiles RENAME COLUMN description TO display;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meta_profiles RENAME COLUMN display TO description;
-- +goose StatementEnd
