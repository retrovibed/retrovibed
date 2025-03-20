package cmdtorrent

type Commands struct {
	Import importFilesystem `cmd:"" help:"import torrents from a file/directory"`
	Magnet cmdMagnet        `cmd:"" help:"insert magnet links for download"`
}
