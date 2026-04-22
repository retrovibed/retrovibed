package cmdmedia

type Commands struct {
	Import  importFilesystem `cmd:"" help:"import files and directories"`
	Export  exportFilesystem `cmd:"" help:"export media to a directory"`
	Reindex reindex          `cmd:"" help:"run the indexing process on media contents, this can take a bit"`
	Known   Known            `cmd:"" help:"functionality for managing known media"`
	Inspect Inspect          `cmd:"" help:"inspect a file to identify its metadata"`
}

type Known struct {
	Env         knownenv       `cmd:"" help:"extract environment infornation for known media import, specifically dates from and existing database"`
	Tarchive    tarchiveexport `cmd:"" help:"export known media from a directory of tar.gz archives and writes to stdout in known media jsonl format for importing"`
	Duckdb      duckdbexport   `cmd:"" help:"export known media from a duckdb database and writes to stdout in known media jsonl format for importing"`
	TMDB        tmdbimport     `cmd:"" help:"import known media from tmdb and writes to stdout in known media jsonl format for importing"`
	TVDB        tvdbimport     `cmd:"" help:"import known media from tvdb and writes to stdout in known media jsonl format for importing"`
	MusicBrainz mbimport       `cmd:"" help:"import known media from music brainz and writes to stdout in known media jsonl format for importing"`
	Deeppool    deeppoolimport `cmd:"" help:"import known media from deeppool published content and writes to stdout in known media jsonl format for importing"`
	Query       knownquery     `cmd:"" help:"run a query against known media"`
	Import      knownimport    `cmd:"" help:"processes a file or stdin to import media metadata records directly into the database"`
	Archive     knownarchive   `cmd:"" help:"processes stdin and creates a directory of files of media metadata"`
}

type Inspect struct {
	File mediaInspectFile `cmd:"" help:"inspect a file to identify its metadata"`
	// Media mediaInspect `cmd:"" help:"download media from a library and inspect its medata"`
}
