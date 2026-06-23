//go:build darwin || ios

package userx

import (
	"os"
	"path/filepath"

	"github.com/retrovibed/retrovibed/retroapi/internal/envx"
	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
)

// returns the relative root that should be used for all well known directories.
func DefaultRelRoot() string {
	return "retrovibed"
}

// platform specific config directory resolution
func _configDir(rel ...string) string {
	homedir := errorsx.Must(os.UserHomeDir())
	defaultdir := filepath.Join(homedir, "Library", "Application Support")
	return filepath.Join(envx.String(defaultdir, "XDG_CONFIG_HOME"), filepath.Join(rel...))
}

func _cacheDir(rel ...string) string {
	cachedir := errorsx.Must(os.UserCacheDir())
	return filepath.Join(envx.String(cachedir, "XDG_CACHE_HOME"), filepath.Join(rel...))
}

func _dataDir(rel ...string) string {
	return _configDir(rel...)
}

func _downloadDir() string {
	homedir := errorsx.Must(os.UserHomeDir())
	return envx.String(filepath.Join(homedir, "Downloads"), "XDG_DOWNLOAD_DIR")
}
