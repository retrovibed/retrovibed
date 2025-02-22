-- +goose Up
-- +goose StatementBegin
INSTALL fts; LOAD fts;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
