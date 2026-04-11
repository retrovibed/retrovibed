-- +goose Up
-- +goose StatementBegin
CREATE TABLE torrents_feed_rss (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    disabled_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    next_check TIMESTAMPTZ NOT NULL DEFAULT '-infinity',
    autodownload boolean NOT NULL DEFAULT 'false',
    autoarchive boolean NOT NULL DEFAULT 'false',
    contributing boolean NOT NULL DEFAULT 'false',
    ttl INTERVAL NOT NULL DEFAULT INTERVAL '0 seconds',
    description VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_feed_rss;
-- +goose StatementEnd