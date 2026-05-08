-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_community_sync_state_community_id;
ALTER TABLE community_sync_state ADD COLUMN sync_feed_at TIMESTAMPTZ DEFAULT 'infinity';
UPDATE community_sync_state SET sync_feed_at = DEFAULT;
COMMIT;
BEGIN;
ALTER TABLE community_sync_state ALTER COLUMN sync_feed_at SET NOT NULL;
CREATE INDEX idx_community_sync_state_community_id ON community_sync_state(community_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE community_sync_state DROP COLUMN IF EXISTS sync_feed_at;
-- +goose StatementEnd
