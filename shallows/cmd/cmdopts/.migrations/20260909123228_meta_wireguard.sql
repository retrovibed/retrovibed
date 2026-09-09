-- +goose Up
-- +goose StatementBegin
-- 1. Ensure the old table exists. On a fresh database (a future state) the
--    original meta_wireguard table may never have been created; this is a
--    no-op there but guarantees the rename below always has a table to act on.
CREATE TABLE IF NOT EXISTS meta_wireguard (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    "default" BOOLEAN DEFAULT 'f' NOT NULL,
    port USMALLINT NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    dns_rate_limit UINTEGER NOT NULL DEFAULT 0,
    maximum_connections UBIGINT NOT NULL DEFAULT (~0::UBIGINT)::UBIGINT,
    outbound_rate_limit UINTEGER NOT NULL DEFAULT 0
);

-- 2. Rename existing table
ALTER TABLE IF EXISTS meta_wireguard RENAME TO meta_wireguard_old;
CREATE TABLE meta_wireguard (
    id UUID PRIMARY KEY NOT NULL, -- md5 of the contents
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    nettype UINT32 NOT NULL DEFAULT 0 CHECK (nettype IN (0, 1, 2)),
    port USMALLINT NOT NULL DEFAULT 0,
    rate_limit_dns UINTEGER NOT NULL DEFAULT 0,
    rate_limit_outbound UINTEGER NOT NULL DEFAULT 0,
    maximum_connections UBIGINT NOT NULL DEFAULT (~0::UBIGINT)::UBIGINT,
    description TEXT NOT NULL DEFAULT ''
);

COMMENT ON COLUMN meta_wireguard.nettype IS 'defines how this wireguard configuration is used within the application. distributionnet is for content distribution. socialnet is for social activities like publishing. 0 = unspecified, 1 = distribution, 2 = social';
COMMENT ON COLUMN meta_wireguard.maximum_connections IS 'maximum number of connections';
COMMENT ON COLUMN meta_wireguard.rate_limit_dns IS 'maximum DNS lookup rate in events per second; 0 means unlimited';
COMMENT ON COLUMN meta_wireguard.rate_limit_outbound IS 'maximum outbound dial rate in connections per second; 0 means use the unspecified';

-- 3. Copy rows across, mapping the renamed columns; netsocial and
-- maximum_connections are new so they take their column defaults.
INSERT INTO meta_wireguard (
    id,
    created_at,
    updated_at,
    nettype,
    port,
    rate_limit_dns,
    rate_limit_outbound,
    description
)
SELECT
    id,
    created_at,
    updated_at,
    CASE WHEN "default" THEN 1 ELSE 0 END,
    port,
    dns_rate_limit,
    outbound_rate_limit,
    description
FROM meta_wireguard_old;

-- 5. Finalize
DROP TABLE meta_wireguard_old;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- IRREVERSIBLE: The old schema used a "default" boolean and separate rate-limit
-- columns; the new schema uses a nettype enum and renamed columns. If a rollback
-- is ever needed, implement it manually at that time.
SELECT CAST('IRREVERSIBLE: manual rollback required' AS INTEGER);
-- +goose StatementEnd
