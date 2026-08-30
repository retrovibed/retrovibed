-- +goose Up
-- +goose StatementBegin
CREATE TABLE ddisc_recommendation_defaults (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  profile_id UUID UNIQUE NOT NULL,
  language VARCHAR NOT NULL DEFAULT '',
  adult BOOLEAN NOT NULL DEFAULT false,
  rlimit UINTEGER NOT NULL DEFAULT 5,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON COLUMN ddisc_recommendation_defaults.rlimit IS 'max recommendations to request per mimetype, per run, from a search plugin';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ddisc_recommendation_defaults;
-- +goose StatementEnd
