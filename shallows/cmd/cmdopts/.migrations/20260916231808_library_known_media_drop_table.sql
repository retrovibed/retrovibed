-- +goose Up
-- +goose StatementBegin
-- library_known_media now lives in cache.db, see .migrations.cache; it is
-- rebuildable from external sources (tmdb/tvdb/musicbrainz/deeppool imports,
-- ddisc discovery), so no data is carried over.
DROP TABLE IF EXISTS library_known_media;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
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
    auto_description VARCHAR NOT NULL DEFAULT '',
    mimetype TEXT NOT NULL DEFAULT 'application/octet-stream',
    tombstoned_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    "collation" UINTEGER NOT NULL DEFAULT 0,
    subtitle VARCHAR NOT NULL DEFAULT '',
    parent_uid UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'
);
-- +goose StatementEnd
