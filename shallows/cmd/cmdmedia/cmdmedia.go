package cmdmedia

type Commands struct {
	Import `cmd:"" help:"provides functionality to import torrent files and directories"`
}

type Import struct {
	Files    importFilesystem `cmd:"" help:"import files from a directory tree"`
	Torrents importTorrents   `cmd:"" help:"import torrents from a directory tree, should be run after files have been imported"`
}
