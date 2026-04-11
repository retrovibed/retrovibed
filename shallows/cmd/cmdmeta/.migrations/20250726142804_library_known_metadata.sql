-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS library_known_media;
CREATE TABLE library_known_media (
    uid UUID PRIMARY KEY NOT NULL,
    md5 UUID UNIQUE NOT NULL,
    md5_lower UBIGINT NOT NULL,
    duplicates BIGINT NOT NULL DEFAULT 0,
    source VARCHAR NOT NULL DEFAULT '',
    id VARCHAR NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    released TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    title VARCHAR NOT NULL DEFAULT '',
    popularity DOUBLE NOT NULL DEFAULT 0.0,
    adult BOOLEAN NOT NULL DEFAULT FALSE,
    backdrop_path VARCHAR NOT NULL DEFAULT '',
    poster_path VARCHAR NOT NULL DEFAULT '',
    original_language VARCHAR NOT NULL DEFAULT '',
    original_title VARCHAR NOT NULL DEFAULT '',
    overview VARCHAR NOT NULL DEFAULT '',
    auto_description VARCHAR NOT NULL DEFAULT '' -- application managed for full text search
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS library_known_media;
-- +goose StatementEnd
