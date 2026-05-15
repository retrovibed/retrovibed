-- +goose Up
-- +goose StatementBegin

-- 1. Rename existing table
ALTER TABLE torrents_metadata RENAME TO torrents_metadata_old;

-- 2. Create the new table with the 20-byte constraint
CREATE TABLE torrents_metadata (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    hidden_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT 'infinity',
    initiated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT 'infinity',
    paused_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT 'infinity',
    next_announce_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT '-infinity',
    seeding BOOLEAN NOT NULL DEFAULT false,
    private BOOLEAN NOT NULL DEFAULT false,
    tracker VARCHAR NOT NULL,
    bytes UBIGINT NOT NULL DEFAULT 0,
    downloaded UBIGINT NOT NULL DEFAULT 0,
    uploaded UBIGINT NOT NULL DEFAULT 0,
    peers USMALLINT NOT NULL DEFAULT 0,
    description VARCHAR NOT NULL DEFAULT '',
    auto_description VARCHAR NOT NULL DEFAULT '',
    infohash BLOB NOT NULL CHECK (octet_length(infohash) = 20),
    known_media_id UUID NOT NULL DEFAULT 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'::UUID,
    encryption_seed UUID NOT NULL DEFAULT gen_random_uuid(),
    archivable BOOLEAN NOT NULL DEFAULT false,
    verify_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT 'infinity',
    mimetype VARCHAR NOT NULL DEFAULT 'application/x-bittorrent',
    imported_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT 'infinity',
    tombstoned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT 'infinity',
    completed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT 'infinity',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT 'infinity'
);

-- 3. Clean up the old data before migration
DELETE FROM torrents_metadata_old WHERE octet_length(infohash) != 20;

-- 4. Move the data (Listing columns for safety)
INSERT INTO torrents_metadata 
SELECT 
    id, created_at, updated_at, hidden_at, initiated_at, paused_at, next_announce_at,
    seeding, private, tracker, bytes, downloaded, uploaded, peers, 
    description, auto_description, infohash, known_media_id, encryption_seed, 
    archivable, verify_at, mimetype, imported_at, tombstoned_at, completed_at, expires_at 
FROM torrents_metadata_old;

-- 5. Finalize
DROP TABLE torrents_metadata_old;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- nothing to do here.
-- +goose StatementEnd
