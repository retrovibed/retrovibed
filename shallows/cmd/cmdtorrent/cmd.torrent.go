package cmdtorrent

type cmdImports struct {
	Peer      importPeer      `cmd:"" help:"import torrents from a specific peer and a list of magnet urls"`
	Directory importDirectory `cmd:"" help:"import torrents directly from a directory"`
}

type Commands struct {
	Verify     filesystemVerify `cmd:"" help:"verify torrents stored on disk and regenerate bitmaps"`
	Import     cmdImports       `cmd:"" help:"import torrents using various strategies"`
	Export     exportMagnets    `cmd:"" help:"export seeded torrents as magnet urls to stdout or a file"`
	Magnet     cmdMagnet        `cmd:"" help:"insert magnet links for download"`
	KnownMedia cmdKnownMedia    `cmd:"" help:"best effort attempt to assign known media"`
	Download   cmdDownload      `cmd:"" help:"test util for downloading a torrent, primarily used to ensure working torrent library, does not persist anything. use kill -s USR1 to dump statistics"`
	DHT        cmdDHTTest       `cmd:"" help:"test utility for ensuring dht behaviors"`
	Inspect    cmdInspect       `cmd:"" help:"print diagnostic information about torrent file"`
}
