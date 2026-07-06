-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS ddisc_media;
CREATE TABLE ddisc_media (
  id UUID PRIMARY KEY NOT NULL, -- md5 of infohash
  infohash BINARY NOT NULL CHECK (octet_length(infohash) = 20),
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
  bytes UBIGINT NOT NULL,
  video_resolution STRING NOT NULL DEFAULT '',
  video_runtime INTERVAL NOT NULL DEFAULT INTERVAL '0 second',
  audio_bitrate UINTEGER NOT NULL DEFAULT 0,
  audio_default_locale STRING NOT NULL DEFAULT 'und' CHECK (audio_default_locale <> ''),
  subtitles_default_locale STRING NOT NULL DEFAULT 'und' CHECK (subtitles_default_locale <> ''),
  known_media_id UUID NOT NULL, -- Default to Nil uuid when unknown and not going to be indexed by this node. Max uuid when it needs to be processed.
  sync_uid UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::UUID,
  "partition" UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::UUID,
  category STRING NOT NULL DEFAULT '',
  adult BOOLEAN NOT NULL DEFAULT false
);
COMMENT ON COLUMN ddisc_media.sync_uid IS 'uuid v7 used to track sync status between peers, peers track the maximum sync_uid of the records they have received';
COMMENT ON COLUMN ddisc_media.partition IS 'uuid used to track what partition discovered media belongs to once its been identified';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ddisc_media;
-- +goose StatementEnd
