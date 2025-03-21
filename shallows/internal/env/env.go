package env

import (
	"github.com/gofrs/uuid"
	"github.com/retrovibed/retrovibed/internal/envx"
	"github.com/retrovibed/retrovibed/internal/userx"
)

const (
	// percentage of requests that should fail.
	ChaosRate = "SHALLOWS_CHAOS_RATE"

	// health code config
	HTTPHealthzProbability = "SHALLOWS_PROBABILITY"
	HTTPHealthzCode        = "SHALLOWS_HEALTHZ_CODE"

	// TLS pem location.
	DaemonTLSPEM = "SHALLOWS_TLS_PEM"
	// JWTSharedSecret shared secret between the applications, used to encrypt data.
	// and sign messages.
	JWTSharedSecret = "SHALLOWS_JWT_SECRET"
)

func JWTSecret() []byte {
	return []byte(envx.String(uuid.Must(uuid.NewV4()).String(), JWTSharedSecret))
}

func MediaDir() string {
	return userx.DefaultDataDirectory(userx.DefaultRelRoot(), "media")
}

func TorrentDir() string {
	return userx.DefaultDataDirectory(userx.DefaultRelRoot(), "torrent")
}

func PrivateKeyPath() string {
	return userx.DefaultConfigDir(userx.DefaultRelRoot(), "id")
}
