-- +goose Up
-- +goose StatementBegin
CREATE TABLE torrents_feed_rss (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    next_check TIMESTAMPTZ NOT NULL DEFAULT '-infinity',
    autodownload boolean NOT NULL DEFAULT 'false',
    description VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_feed_rss;
-- +goose StatementEnd
