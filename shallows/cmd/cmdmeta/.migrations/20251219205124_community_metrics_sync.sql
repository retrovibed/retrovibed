-- +goose Up
-- +goose StatementBegin
CREATE TABLE community_sync_state (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    community_id UUID NOT NULL UNIQUE,
    last_sync_at TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01 00:00:00+00',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_community_sync_state_community_id ON community_sync_state(community_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS community_sync_state;
-- +goose StatementEnd
