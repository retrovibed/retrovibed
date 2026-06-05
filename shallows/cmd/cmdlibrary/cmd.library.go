package cmdlibrary

type cmdImports struct {
	Directory importDirectory `cmd:"" help:"import files from a directory into the library"`
	JSONL     importJSONL     `cmd:"" help:"import library records and file data from a JSONL stream"`
}

type Commands struct {
	Import  cmdImports  `cmd:"" help:"import media using various strategies"`
	Export  exportJSONL `cmd:"" help:"export library records and file data as JSONL to stdout"`
	Publish cmdPublish  `cmd:"" help:"publish a library content"`
}
