-- +goose Up
-- +goose StatementBegin
CREATE TABLE library_locate (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    known_media_id uuid UNIQUE NOT NULL,
    located_torrent_id uuid NOT NULL DEFAULT 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'::uuid,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS library_locate;
-- +goose StatementEnd
