//go:build !darwin

package userx

import (
	"os"
	"path/filepath"

	"github.com/retrovibed/retrovibed/retroapi/internal/envx"
)

// returns the relative root that should be used for all well known directories.
func DefaultRelRoot() string {
	return filepath.Base(os.Args[0])
}

// platform specific config directory resolution
func _configDir(rel ...string) string {
	user := CurrentUserOrDefault(Root())
	defaultdir := filepath.Join(user.HomeDir, ".config")
	return filepath.Join(envx.String(defaultdir, "XDG_CONFIG_HOME"), filepath.Join(rel...))
}

func _cacheDir(rel ...string) string {
	user := CurrentUserOrDefault(Root())
	if user.Uid == Root().Uid {
		return envx.String(filepath.Join("/", "var", "cache"), "CACHE_DIRECTORY")
	}

	defaultdir := filepath.Join(user.HomeDir, ".cache")
	return filepath.Join(envx.String(defaultdir, "CACHE_DIRECTORY", "XDG_CACHE_HOME"), filepath.Join(rel...))
}

func _dataDir(rel ...string) string {
	user := CurrentUserOrDefault(Root())
	defaultdir := filepath.Join(user.HomeDir, ".local", "share")
	return filepath.Join(envx.String(defaultdir, "XDG_DATA_HOME"), filepath.Join(rel...))
}

func _downloadDir() string {
	user := CurrentUserOrDefault(Root())
	return envx.String(filepath.Join(user.HomeDir, "Downloads"), "XDG_DOWNLOAD_DIR")
}
