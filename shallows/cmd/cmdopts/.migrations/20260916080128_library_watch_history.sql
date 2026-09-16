-- +goose Up
-- +goose StatementBegin
CREATE TABLE library_watch_history (
    id UUID PRIMARY KEY,
    profile_id UUID NOT NULL,
    media_id UUID NOT NULL,
    duration INTERVAL NOT NULL DEFAULT INTERVAL '0 seconds',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS library_watch_history;
-- +goose StatementEnd
