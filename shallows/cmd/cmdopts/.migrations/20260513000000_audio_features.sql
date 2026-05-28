-- +goose Up
-- +goose StatementBegin
CREATE TABLE audio_features (
    media_id UUID PRIMARY KEY NOT NULL,
    features DOUBLE[128] NOT NULL,
    stats_version UINTEGER NOT NULL DEFAULT 0,
    indexed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS audio_features;
-- +goose StatementEnd
