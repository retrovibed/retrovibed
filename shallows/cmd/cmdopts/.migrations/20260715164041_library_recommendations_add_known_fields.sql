-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS library_recommendations;
CREATE TABLE library_recommendations (
  id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
  source UUID NOT NULL,
  recommendations UINTEGER NOT NULL DEFAULT 0,
  content_id UUID UNIQUE NOT NULL,
  known_media_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::UUID,
  mimetype TEXT NOT NULL DEFAULT 'application/octet-stream',
  adult BOOLEAN NOT NULL DEFAULT FALSE,
  language VARCHAR NOT NULL DEFAULT '',
  title TEXT NOT NULL DEFAULT '',
  overview TEXT NOT NULL DEFAULT '',
  image TEXT NOT NULL DEFAULT '',
  popularity DOUBLE NOT NULL DEFAULT 0.0,
  released TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  tombstone_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity'
);
CREATE INDEX idx_library_recommendations_updated_at ON library_recommendations(updated_at DESC);
CREATE INDEX idx_library_recommendations_tombstone_at ON library_recommendations(tombstone_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_library_recommendations_tombstone_at;
DROP INDEX IF EXISTS idx_library_recommendations_updated_at;
DROP TABLE IF EXISTS library_recommendations;
-- +goose StatementEnd
