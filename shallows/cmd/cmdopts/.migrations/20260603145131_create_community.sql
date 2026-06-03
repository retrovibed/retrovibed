-- +goose Up
-- +goose StatementBegin
CREATE TABLE community (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_sync_at TIMESTAMPTZ NOT NULL DEFAULT '-infinity',
    sync_feed_at TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    auto_download INTEGER NOT NULL DEFAULT 0,
    sync_cursor_published_content UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'
);

INSERT INTO community (id, auto_download, last_sync_at, sync_feed_at, created_at, updated_at)
SELECT
    COALESCE(cs.community_id, css.community_id),
    COALESCE(cs.auto_download, 0),
    GREATEST(COALESCE(cs.last_sync_at, '-infinity'), COALESCE(css.last_sync_at, '-infinity')),
    COALESCE(css.sync_feed_at, 'infinity'),
    LEAST(COALESCE(cs.created_at, now()), COALESCE(css.created_at, now())),
    GREATEST(COALESCE(cs.updated_at, now()), COALESCE(css.updated_at, now()))
FROM community_subscription cs
FULL OUTER JOIN community_sync_state css ON cs.community_id = css.community_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS community;
-- +goose StatementEnd
