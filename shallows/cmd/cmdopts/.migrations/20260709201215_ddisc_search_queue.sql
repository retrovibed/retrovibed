-- +goose Up
-- +goose StatementBegin
CREATE TABLE ddisc_search_queue (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    known_media_id UUID UNIQUE NOT NULL,
    next_check_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    attempts UINTEGER NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ddisc_search_queue;
-- +goose StatementEnd
