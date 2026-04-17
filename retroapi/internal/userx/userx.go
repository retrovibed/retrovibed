package userx

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/retrovibed/retroapi/internal/debugx"
	"github.com/retrovibed/retroapi/internal/envx"
)

func Locale() string {
	return envx.String("", "LANG", "LC_ALL")
}

func Root() user.User {
	return user.User{
		Gid:     "0",
		Uid:     "0",
		HomeDir: "/root",
	}
}

func Zero() user.User {
	return user.User{}
}

// CurrentUserOrDefault returns the current user or the default configured user.
func CurrentUserOrDefault(d user.User) (result *user.User) {
	var (
		err error
	)

	if result, err = user.Current(); err != nil {
		log.Println("failed to retrieve current user, using default", err)
		tmp := d
		return &tmp
	}

	return result
}

// DefaultConfigDir returns the user config directory.
func DefaultConfigDir(rel ...string) string {
	return _configDir(rel...)
}

// DefaultCacheDirectory cache directory for storing data.
func DefaultCacheDirectory(rel ...string) string {
	return _cacheDir(rel...)
}

func DefaultDataDirectory(rel ...string) string {
	return _dataDir(rel...)
}

// DefaultDownloadDirectory returns the user download directory.
func DefaultDownloadDirectory(rel ...string) string {
	return filepath.Join(_downloadDir(), filepath.Join(rel...))
}

// DefaultRuntimeDirectory runtime directory for storing data.
func DefaultRuntimeDirectory(rel ...string) string {
	user := CurrentUserOrDefault(Root())

	if user.Uid == Root().Uid {
		return envx.String(filepath.Join("/", "run"), "RUNTIME_DIRECTORY", "XDG_RUNTIME_DIR")
	}

	defaultdir := filepath.Join("/", "run", "user", user.Uid)
	return filepath.Join(envx.String(defaultdir, "RUNTIME_DIRECTORY", "XDG_RUNTIME_DIR"), filepath.Join(rel...))
}

// return a path relative to the home directory.
func HomeDirectoryRel(rel ...string) (path string, err error) {
	if path, err = os.UserHomeDir(); err != nil {
		return "", err
	}

	return filepath.Join(path, filepath.Join(rel...)), nil
}

// DefaultDirectory finds the first directory root that exists and then returns
// that root directory joined with the relative path provided.
func DefaultDirectory(rel string, roots ...string) (path string) {
	for _, root := range roots {
		path = filepath.Join(root, rel)
		if _, err := os.Stat(root); err == nil {
			return path
		}
	}

	return path
}

// HomeDirectoryOrDefault loads the user home directory or falls back to the provided
// path when an error occurs.
func HomeDirectoryOrDefault(fallback string) (dir string) {
	var (
		err error
	)

	if dir, err = os.UserHomeDir(); err != nil {
		debugx.Println("unable to get user home directory", err)
		return fallback
	}

	return dir
}
