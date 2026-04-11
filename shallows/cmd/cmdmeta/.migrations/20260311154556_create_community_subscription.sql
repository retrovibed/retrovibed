-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS community_subscription;
CREATE TABLE community_subscription (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    community_id UUID NOT NULL UNIQUE,
    auto_download INTEGER NOT NULL DEFAULT 0,
    last_sync_at TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01 00:00:00+00',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS community_subscription;
-- +goose StatementEnd
