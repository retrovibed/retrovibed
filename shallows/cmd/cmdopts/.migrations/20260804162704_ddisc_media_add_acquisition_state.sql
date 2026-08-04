-- +goose Up
-- +goose StatementBegin
ALTER TABLE ddisc_media RENAME TO ddisc_media_old;
CREATE TABLE ddisc_media (
  id UUID PRIMARY KEY NOT NULL, -- md5 of infohash
  infohash BINARY NOT NULL CHECK (octet_length(infohash) = 20),
  uri STRING NOT NULL CHECK (uri <> ''),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  next_check_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  attempts UINTEGER NOT NULL DEFAULT 0,
  hidden_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
  tombstoned_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
  released_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
  title STRING NOT NULL DEFAULT '',
  source STRING NOT NULL DEFAULT '',
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
  adult BOOLEAN NOT NULL DEFAULT false,
  private BOOLEAN NOT NULL DEFAULT false,
  health UINTEGER NOT NULL DEFAULT 0,
  policy_rank USMALLINT NOT NULL DEFAULT 65535,
  policy_rejection STRING NOT NULL DEFAULT '',
  acquisition_state INTEGER NOT NULL DEFAULT 0
);
COMMENT ON COLUMN ddisc_media.source IS 'tracks what the specific subsystem the ddisc media as sourced from';
COMMENT ON COLUMN ddisc_media.uri IS 'magnet uri for this torrent, when known; infohash remains the canonical identity/dedup key';
COMMENT ON COLUMN ddisc_media.sync_uid IS 'uuid v7 used to track sync status between peers, peers track the maximum sync_uid of the records they have received';
COMMENT ON COLUMN ddisc_media.partition IS 'uuid used to track what partition discovered media belongs to once its been identified';
COMMENT ON COLUMN ddisc_media.health IS 'seed/health count as reported by whatever source discovered this candidate (e.g. a search plugin); 0 means unknown, not necessarily zero seeds';
COMMENT ON COLUMN ddisc_media.policy_rank IS 'lower is better; computed by a ranking Policy, 65535 means unranked or hard-rejected';
COMMENT ON COLUMN ddisc_media.policy_rejection IS 'non-empty when a ranking Policy hard-rejected this candidate, set to the matched dealbreaker term';
COMMENT ON COLUMN ddisc_media.acquisition_state IS 'mirrors ddisc.discovery.proto AcquisitionState ordinals: 0 unknown, 1 ephemeral, 2 downloading, 3 completed - see shallows/ddisc.DiscoveredOptionAcquisition* options';
INSERT INTO ddisc_media (
  id, infohash, uri, created_at, updated_at, next_check_at, attempts, hidden_at, tombstoned_at, released_at,
  title, source, "collation", description, mimetype, bytes, video_resolution, video_runtime, audio_bitrate,
  audio_default_locale, subtitles_default_locale, known_media_id, sync_uid, "partition", category, adult,
  private, health, policy_rank, policy_rejection
)
SELECT
  id, infohash, uri, created_at, updated_at, next_check_at, attempts, hidden_at, tombstoned_at, released_at,
  title, source, "collation", description, mimetype, bytes, video_resolution, video_runtime, audio_bitrate,
  audio_default_locale, subtitles_default_locale, known_media_id, sync_uid, "partition", category, adult,
  private, health, policy_rank, policy_rejection
FROM ddisc_media_old;
DROP TABLE IF EXISTS ddisc_media_old;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ddisc_media RENAME TO ddisc_media_old;
CREATE TABLE ddisc_media (
  id UUID PRIMARY KEY NOT NULL, -- md5 of infohash
  infohash BINARY NOT NULL CHECK (octet_length(infohash) = 20),
  uri STRING NOT NULL CHECK (uri <> ''),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  next_check_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  attempts UINTEGER NOT NULL DEFAULT 0,
  hidden_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
  tombstoned_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
  released_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
  title STRING NOT NULL DEFAULT '',
  source STRING NOT NULL DEFAULT '',
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
  adult BOOLEAN NOT NULL DEFAULT false,
  private BOOLEAN NOT NULL DEFAULT false,
  health UINTEGER NOT NULL DEFAULT 0,
  policy_rank USMALLINT NOT NULL DEFAULT 65535,
  policy_rejection STRING NOT NULL DEFAULT ''
);
COMMENT ON COLUMN ddisc_media.source IS 'tracks what the specific subsystem the ddisc media as sourced from';
COMMENT ON COLUMN ddisc_media.uri IS 'magnet uri for this torrent, when known; infohash remains the canonical identity/dedup key';
COMMENT ON COLUMN ddisc_media.sync_uid IS 'uuid v7 used to track sync status between peers, peers track the maximum sync_uid of the records they have received';
COMMENT ON COLUMN ddisc_media.partition IS 'uuid used to track what partition discovered media belongs to once its been identified';
COMMENT ON COLUMN ddisc_media.health IS 'seed/health count as reported by whatever source discovered this candidate (e.g. a search plugin); 0 means unknown, not necessarily zero seeds';
COMMENT ON COLUMN ddisc_media.policy_rank IS 'lower is better; computed by a ranking Policy, 65535 means unranked or hard-rejected';
COMMENT ON COLUMN ddisc_media.policy_rejection IS 'non-empty when a ranking Policy hard-rejected this candidate, set to the matched dealbreaker term';
INSERT INTO ddisc_media (
  id, infohash, uri, created_at, updated_at, next_check_at, attempts, hidden_at, tombstoned_at, released_at,
  title, source, "collation", description, mimetype, bytes, video_resolution, video_runtime, audio_bitrate,
  audio_default_locale, subtitles_default_locale, known_media_id, sync_uid, "partition", category, adult,
  private, health, policy_rank, policy_rejection
)
SELECT
  id, infohash, uri, created_at, updated_at, next_check_at, attempts, hidden_at, tombstoned_at, released_at,
  title, source, "collation", description, mimetype, bytes, video_resolution, video_runtime, audio_bitrate,
  audio_default_locale, subtitles_default_locale, known_media_id, sync_uid, "partition", category, adult,
  private, health, policy_rank, policy_rejection
FROM ddisc_media_old;
DROP TABLE IF EXISTS ddisc_media_old;
-- +goose StatementEnd
