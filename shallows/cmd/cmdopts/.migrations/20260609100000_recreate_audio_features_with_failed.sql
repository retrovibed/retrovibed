-- +goose Up
-- +goose StatementBegin
SET hnsw_enable_experimental_persistence = true;

DROP INDEX IF EXISTS audio_features_hnsw;
DROP TABLE IF EXISTS audio_features;

CREATE TABLE audio_features (
    media_id UUID PRIMARY KEY NOT NULL,
    features FLOAT[128] NOT NULL,
    failed BOOLEAN NOT NULL DEFAULT false,
    stats_version UINTEGER NOT NULL DEFAULT 0,
    indexed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX audio_features_hnsw ON audio_features
USING HNSW (features) WITH (metric = 'cosine');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET hnsw_enable_experimental_persistence = true;

DROP INDEX IF EXISTS audio_features_hnsw;
DROP TABLE IF EXISTS audio_features;

CREATE TABLE audio_features (
    media_id UUID PRIMARY KEY NOT NULL,
    features FLOAT[128] NOT NULL,
    stats_version UINTEGER NOT NULL DEFAULT 0,
    indexed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX audio_features_hnsw ON audio_features
USING HNSW (features) WITH (metric = 'cosine');
-- +goose StatementEnd
