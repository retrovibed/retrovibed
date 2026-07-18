-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS ddisc_locate;
CREATE TABLE ddisc_locate (
  id UUID PRIMARY KEY NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  tombstoned_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
  query STRING NOT NULL,
  mimetype STRING NOT NULL CHECK (mimetype <> ''),
  known_media_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  located_torrent_id UUID NOT NULL DEFAULT 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'::uuid,
  autodownload BOOLEAN NOT NULL DEFAULT false,
  adult BOOLEAN NOT NULL DEFAULT false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ddisc_locate;
-- +goose StatementEnd
