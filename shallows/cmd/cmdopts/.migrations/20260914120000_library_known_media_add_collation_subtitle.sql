-- +goose Up
-- +goose StatementBegin
ALTER TABLE library_known_media ADD COLUMN collation UINTEGER DEFAULT 0;
UPDATE library_known_media SET collation = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE library_known_media ALTER COLUMN collation SET NOT NULL;
COMMENT ON COLUMN library_known_media.collation IS 'ordering key: 0 = standalone/overall item; TV episode = season(hi 16 bits, 0xFFFF for specials)/episode(lo 16 bits); album = plain track sequence';

ALTER TABLE library_known_media ADD COLUMN subtitle VARCHAR DEFAULT '';
UPDATE library_known_media SET subtitle = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE library_known_media ALTER COLUMN subtitle SET NOT NULL;
COMMENT ON COLUMN library_known_media.subtitle IS 'episode name; blank for non-episodic rows (movies, albums, the series row itself)';

ALTER TABLE library_known_media ADD COLUMN parent_uid UUID DEFAULT '00000000-0000-0000-0000-000000000000';
UPDATE library_known_media SET parent_uid = DEFAULT;
COMMIT; BEGIN;
ALTER TABLE library_known_media ALTER COLUMN parent_uid SET NOT NULL;
COMMENT ON COLUMN library_known_media.parent_uid IS 'uid of the parent record (e.g. an episode''s show); nil uuid when not a child record';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE library_known_media DROP COLUMN IF EXISTS collation;
ALTER TABLE library_known_media DROP COLUMN IF EXISTS subtitle;
ALTER TABLE library_known_media DROP COLUMN IF EXISTS parent_uid;
-- +goose StatementEnd
