-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS ddisc_peers;
ALTER TABLE torrents_peers RENAME TO torrents_peers_old;
CREATE TABLE torrents_peers (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    tombstoned_at TIMESTAMPTZ NOT NULL,
    next_check TIMESTAMPTZ NOT NULL DEFAULT '-infinity',
    description TEXT NOT NULL DEFAULT '',
    peer BINARY NOT NULL,
    network TEXT NOT NULL,
    ip INET NOT NULL DEFAULT '::',
    port USMALLINT NOT NULL,
    bep51 boolean NOT NULL DEFAULT false,
    bep51_ttl USMALLINT NOT NULL DEFAULT 0,
    bep51_available UBIGINT NOT NULL DEFAULT 0,
    ddisc boolean NOT NULL DEFAULT false,
    ddisc_partition UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    ddisc_syncoffset UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000' -- uuid v7 synchronization id.
);
INSERT INTO torrents_peers (
    id,
    created_at,
    updated_at,
    tombstoned_at,
    next_check,
    description,
    peer,
    network,
    ip,
    port,
    bep51,
    bep51_ttl,
    bep51_available,
    ddisc,
    ddisc_partition,
    ddisc_syncoffset
)
SELECT
    id,
    created_at,
    updated_at,
    tombstoned_at,
    next_check,
    '', -- Default for description
    peer,
    network,
    ip,
    port,
    bep51,
    bep51_ttl,
    bep51_available,
    FALSE AS ddisc, -- Default for new column
    '00000000-0000-0000-0000-000000000000'::UUID AS ddisc_partition, -- Default for new column
    '00000000-0000-0000-0000-000000000000'::UUID AS ddisc_syncoffset -- Default for new column
FROM torrents_peers_old;
DROP TABLE IF EXISTS torrents_peers_old;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE torrents_peers DROP COLUMN "ddisc";
ALTER TABLE torrents_peers DROP COLUMN "ddisc_partition";
ALTER TABLE torrents_peers DROP COLUMN "ddisc_syncoffset";
-- +goose StatementEnd