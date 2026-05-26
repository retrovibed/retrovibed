-- +goose Up
-- +goose StatementBegin
-- cleanup previous version.
DROP TABLE IF EXISTS library_recent_sessions;

CREATE TABLE library_recent_sessions (
    id UUID PRIMARY KEY NOT NULL,
    mimetype TEXT NOT NULL DEFAULT 'application/octet-stream',
    media_id UUID NOT NULL,
    position INTERVAL NOT NULL DEFAULT INTERVAL '0 seconds',
    duration INTERVAL NOT NULL DEFAULT INTERVAL '0 seconds',
    query BINARY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_played_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_rs_updated_at ON library_recent_sessions(updated_at DESC);
CREATE INDEX idx_rs_last_played ON library_recent_sessions(last_played_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_rs_last_played;
DROP INDEX IF EXISTS idx_rs_updated_at;
DROP INDEX IF EXISTS idx_rs_media_id;
DROP TABLE IF EXISTS library_recent_sessions;
-- +goose StatementEnd

