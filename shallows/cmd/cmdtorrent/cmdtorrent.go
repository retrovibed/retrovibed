package cmdtorrent

type Commands struct {
	Magnet cmdMagnet `cmd:"" help:"insert magnet links for download"`
}
