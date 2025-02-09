-- +goose Up
-- +goose StatementBegin
CREATE TABLE torrents_peers (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    bep51 boolean NOT NULL DEFAULT false,
    bep51_next_available TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    bep51_available UINTEGER NOT NULL DEFAULT 0,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS torrents_peers;
-- +goose StatementEnd
