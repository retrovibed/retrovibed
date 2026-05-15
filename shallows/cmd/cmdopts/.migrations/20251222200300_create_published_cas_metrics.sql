-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS published_cas_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    published_content_id UUID NOT NULL,
    period_start TIMESTAMPTZ NOT NULL DEFAULT now(),
    period_end TIMESTAMPTZ NOT NULL DEFAULT now(),
    archivers UINTEGER NOT NULL DEFAULT 0,
    bytes BIGINT NOT NULL DEFAULT 0,
    revenue BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (published_content_id, period_start)
);

CREATE INDEX IF NOT EXISTS idx_published_cas_metrics_published_content_id ON published_cas_metrics(published_content_id);
CREATE INDEX IF NOT EXISTS idx_published_cas_metrics_period ON published_cas_metrics(period_start, period_end);

DROP TABLE IF EXISTS published_content_metrics;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS published_cas_metrics;
-- +goose StatementEnd
