-- +goose Up
-- +goose StatementBegin
CREATE TABLE meta_profiles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    description text DEFAULT '' NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    disabled_at timestamp with time zone DEFAULT 'infinity'::timestamp with time zone NOT NULL,
    disabled_manually_at timestamp with time zone DEFAULT 'infinity'::timestamp with time zone NOT NULL,
    disabled_pending_approval_at timestamp with time zone DEFAULT '-infinity'::timestamp with time zone NOT NULL,
    session_watermark uuid DEFAULT gen_random_uuid() NOT NULL
);

CREATE TABLE meta_consumed_tokens (
    id uuid PRIMARY KEY,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    tombstoned_at timestamp with time zone NOT NULL,
    token text DEFAULT '' NOT NULL
);

CREATE TABLE authz_meta (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    profile_id uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid NOT NULL UNIQUE,
    usermanagement boolean DEFAULT 'f' NOT NULL
);

CREATE TABLE meta_sso_identity_ssh (
    id uuid PRIMARY KEY, -- md5 of the public key
    disabled_at timestamp with time zone DEFAULT 'infinity'::timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    profile_id uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid NOT NULL,
    public_key TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS authz_meta;
DROP TABLE IF EXISTS meta_consumed_tokens;
DROP TABLE IF EXISTS meta_sso_identity_ssh;
DROP TABLE IF EXISTS meta_profiles;
-- +goose StatementEnd
