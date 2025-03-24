package env

import (
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/internal/envx"
	"github.com/retrovibed/retrovibed/internal/userx"
)

const (
	// percentage of requests that should fail.
	ChaosRate = "RETROVIBED_CHAOS_RATE"

	// health code config
	HTTPHealthzProbability = "RETROVIBED_PROBABILITY"
	HTTPHealthzCode        = "RETROVIBED_HEALTHZ_CODE"

	// TLS pem location.
	DaemonTLSPEM = "RETROVIBED_TLS_PEM"
	// JWTSharedSecret shared secret between the applications, used to encrypt data.
	// and sign messages.
	JWTSharedSecret = "RETROVIBED_JWT_SECRET"

	// enable multicast service discovery
	MDNSDisabled  = "RETROVIBED_MDNS_DISABLED"          // enable/disable multicast dns registration, allows for the frontend to automatically find daemons on the local network.
	AutoDiscovery = "RETROVIBED_TORRENT_AUTO_DISCOVERY" // enable/disable automatically discovering torrents from peers.
	AutoBootstrap = "RETROVIBED_TORRENT_AUTO_BOOTSTRAP" // enable/disable the predefined set of public swarms to bootstrap from
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
