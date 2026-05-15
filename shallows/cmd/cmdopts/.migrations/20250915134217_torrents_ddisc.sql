-- +goose Up
-- +goose StatementBegin
CREATE TABLE torrents_ddisc_media (
    id UUID PRIMARY KEY NOT NULL, -- md5 of infohash
    infohash BINARY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    next_check_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    attempts UINTEGER NOT NULL DEFAULT 0,
    hidden_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    tombstoned_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    released_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    title STRING NOT NULL DEFAULT '',
    "collation" STRING NOT NULL DEFAULT '',
    description STRING NOT NULL DEFAULT '',
    mimetype STRING NOT NULL DEFAULT 'application/octet-stream',
    -- mimetype_md5 UUID NOT NULL DEFAULT md5(mimetype),
    bytes UBIGINT NOT NULL,
    video_resolution STRING NOT NULL DEFAULT '',
    video_runtime INTERVAL NOT NULL DEFAULT INTERVAL '0 second',
    audio_bitrate UINTEGER NOT NULL DEFAULT 0, -- kb/s?
    audio_default_locale STRING NOT NULL DEFAULT '',
    subtitles_default_locale STRING NOT NULL DEFAULT '',
    known_media_id uuid NOT NULL, -- Default to Nil uuid when unknown and not going to be indexed by this node. Max uuid when it needs to be processed.
);

-- CREATE INDEX torrents_ddisc_media_md5_index ON torrents_ddisc_media (mimetype_md5);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_ddisc_media;
-- +goose StatementEnd

-- ALTER TABLE torrents_ddisc_media RENAME TO torrents_ddisc_media_old;
-- INSERT INTO torrents_ddisc_media (id, infohash, created_at, updated_at, next_check_at, attempts, hidden_at, tombstoned_at, released_at, title, "collation", description, mimetype, bytes, video_resolution, video_runtime, audio_bitrate, audio_default_locale, subtitles_default_locale, known_media_id) SELECT id, infohash, created_at, updated_at, next_check_at, attempts, hidden_at, tombstoned_at, released_at, title, "collation", description, mimetype, bytes, video_resolution, COALESCE(video_runtime, INTERVAL '0 second'), audio_bitrate, audio_default_locale, subtitles_default_locale, known_media_id FROM torrents_ddisc_media_old;