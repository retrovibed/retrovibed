package env

import (
	"path/filepath"

	"github.com/retrovibed/retrovibed/retroapi/userx"
)

const (
	// used for local dev to change default.
	DeeppoolEndpoint = "RETROVIBED_META_ENDPOINT"

	// MediaDirName and TorrentDirName are the on-disk subdirectory names
	// media files and .torrent files/caches are stored under, relative to
	// any rootstore - shared so every consumer of a shared fsx.Virtual root
	// scopes into the same subdirectories.
	MediaDirName   = "media"
	TorrentDirName = "torrent"
)

func RootStorageDir(rel ...string) string {
	return userx.DefaultDataDirectory(userx.DefaultRelRoot(), filepath.Join(rel...))
}

func MediaDir() string {
	return RootStorageDir(MediaDirName)
}

func TorrentDir() string {
	return RootStorageDir(TorrentDirName)
}

func PrivateKeyPath(root string) string {
	return userx.DefaultConfigDir(root, "id")
}

func TLSPoolDir() string {
	return userx.DefaultConfigDir(userx.DefaultRelRoot(), "tls.d")
}
