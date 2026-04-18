package env

import (
	"path/filepath"

	"github.com/retrovibed/retrovibed/retroapi/internal/userx"
)

const (
	// used for local dev to change default.
	DeeppoolEndpoint = "RETROVIBED_META_ENDPOINT"
)

func RootStorageDir(rel ...string) string {
	return userx.DefaultDataDirectory(userx.DefaultRelRoot(), filepath.Join(rel...))
}

func MediaDir() string {
	return RootStorageDir("media")
}

func TorrentDir() string {
	return RootStorageDir("torrent")
}

func PrivateKeyPath() string {
	return userx.DefaultConfigDir(userx.DefaultRelRoot(), "id")
}

func TLSPoolDir() string {
	return userx.DefaultConfigDir(userx.DefaultRelRoot(), "tls.d")
}
