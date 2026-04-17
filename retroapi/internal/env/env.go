package env

import (
	"path/filepath"
	"sync"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retroapi/internal/envx"
	"github.com/retrovibed/retroapi/internal/userx"
)

const (
	// used for local dev to change default.
	DeeppoolEndpoint = "RETROVIBED_META_ENDPOINT"

	// JWTSharedSecret used to create jwt tokens
	JWTSharedSecret = "RETROVIBED_JWT_SECRET"
)

var v = sync.OnceValue(func() []byte {
	return []byte(envx.String(
		uuid.Must(uuid.NewV4()).String(),
		JWTSharedSecret,
	))
})

func JWTSecret() []byte {
	return v()
}

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
