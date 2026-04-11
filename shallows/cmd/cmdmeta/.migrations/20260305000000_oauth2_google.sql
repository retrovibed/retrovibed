-- +goose Up
-- +goose StatementBegin
CREATE TABLE oauth2_google (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    access_token VARCHAR NOT NULL DEFAULT '',
    refresh_token VARCHAR NOT NULL DEFAULT '',
    token_type VARCHAR NOT NULL DEFAULT 'Bearer',
    expiry TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01 00:00:00+00',
    scopes VARCHAR NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS oauth2_google;
-- +goose StatementEnd
