-- +goose Up
-- +goose StatementBegin
CREATE TABLE community_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    community_id UUID NOT NULL,
    period_start TIMESTAMPTZ NOT NULL DEFAULT now(),
    period_end TIMESTAMPTZ NOT NULL DEFAULT now(),
    subscribers UINTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (community_id, period_start),
    CHECK (period_start < period_end)
);

CREATE INDEX idx_community_metrics_community_id ON community_metrics(community_id);
CREATE INDEX idx_community_metrics_period ON community_metrics(period_start, period_end);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS community_metrics;
-- +goose StatementEnd
