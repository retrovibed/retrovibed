-- +goose Up
-- +goose StatementBegin
CREATE TABLE library_autoimport_directories (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    path TEXT NOT NULL CHECK (path <> ''),
    description TEXT NOT NULL DEFAULT '',
    debounce INTERVAL NOT NULL DEFAULT INTERVAL '1 hour' CHECK (debounce > INTERVAL '0 seconds'),
    mode UINTEGER NOT NULL DEFAULT 0 CHECK (mode IN (0, 1)),
    library_directory_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    last_scanned_at TIMESTAMPTZ NOT NULL DEFAULT '-infinity'
);

COMMENT ON COLUMN library_autoimport_directories.path IS 'absolute path of the directory to monitor, only the top level is scanned';
COMMENT ON COLUMN library_autoimport_directories.debounce IS 'files are imported once they have not been modified for this long';
COMMENT ON COLUMN library_autoimport_directories.mode IS '0 = copy (keep the original), 1 = move (remove the original after import)';
COMMENT ON COLUMN library_autoimport_directories.library_directory_id IS 'library virtual directory (library_metadata.directory_id) imported files are placed into';

CREATE UNIQUE INDEX IF NOT EXISTS idx_library_autoimport_directories_path ON library_autoimport_directories(path);

CREATE TABLE library_autoimport_files (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    directory_id UUID NOT NULL,
    name TEXT NOT NULL CHECK (name <> ''),
    import_at TIMESTAMPTZ NOT NULL,
    last_imported_at TIMESTAMPTZ NOT NULL DEFAULT '-infinity',
    library_metadata_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    attempts UINTEGER NOT NULL DEFAULT 0,
    error TEXT NOT NULL DEFAULT ''
);

COMMENT ON COLUMN library_autoimport_files.name IS 'file name relative to the monitored directory';
COMMENT ON COLUMN library_autoimport_files.import_at IS 'file modification time + debounce, or the retry backoff after a failure. only moves forward.';
COMMENT ON COLUMN library_autoimport_files.last_imported_at IS 'the file is pending while last_imported_at < import_at';

CREATE UNIQUE INDEX IF NOT EXISTS idx_library_autoimport_files_directory_name ON library_autoimport_files(directory_id, name);
CREATE INDEX IF NOT EXISTS idx_library_autoimport_files_import_at ON library_autoimport_files(import_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS library_autoimport_files;
DROP TABLE IF EXISTS library_autoimport_directories;
-- +goose StatementEnd
