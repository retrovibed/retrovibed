package env

import (
	"sync"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
)

func MediaDir() string                    { return env.MediaDir() }
func TorrentDir() string                  { return env.TorrentDir() }
func PrivateKeyPath() string              { return env.PrivateKeyPath() }
func TLSPoolDir() string                  { return env.TLSPoolDir() }
func RootStorageDir(rel ...string) string { return env.RootStorageDir(rel...) }

const (
	// percentage of requests that should fail.
	ChaosRate = "RETROVIBED_CHAOS_RATE"

	// health code config
	HTTPHealthzProbability = "RETROVIBED_PROBABILITY"
	HTTPHealthzCode        = "RETROVIBED_HEALTHZ_CODE"

	// TLS pem location.
	DaemonTLSPEM = "RETROVIBED_TLS_PEM"
	DaemonSocket = "RETROVIBED_DAEMON_SOCKET" // specify the socket/port to listen on in the form: tcp://:9998

	// used for local dev to change default.
	DeeppoolEndpoint = "RETROVIBED_META_ENDPOINT"

	// JWTSharedSecret used to create jwt tokens
	JWTSharedSecret = "RETROVIBED_JWT_SECRET"

	AutoDownloadMetadata    = "RETROVIBED_METADATA_AUTODOWNLOAD"            // enable/disable automatically downloading metadata.
	AutoMDNS                = "RETROVIBED_AUTO_MDNS"                        // enable/disable multicast dns registration, allows for the frontend to automatically find daemons on the local network.
	AutoArchive             = "RETROVIBED_MEDIA_AUTO_ARCHIVE"               // enable/disable automatic archiving of eligible media.
	AutoReclaim             = "RETROVIBED_MEDIA_AUTO_RECLAIM"               // enable/disable automatic reclaiming of disk space for media that has been archived.
	AutoIdentifyMedia       = "RETROVIBED_MEDIA_AUTO_IDENTIFY"              // enable/disable automatically identified media.
	AutoLocateMedia         = "RETROVIBED_MEDIA_AUTO_LOCATE"                // enable/disable automatically locate and download media.
	AutoDiscovery           = "RETROVIBED_TORRENT_AUTO_DISCOVERY"           // enable/disable automatically discovering torrents from peers.
	AutoBootstrap           = "RETROVIBED_TORRENT_AUTO_BOOTSTRAP"           // enable/disable the predefined set of public swarms to bootstrap from.
	TorrentPort             = "RETROVIBED_TORRENT_PORT"                     // specify the port to listen to torrents on
	TorrentPublicIP4        = "RETROVIBED_TORRENT_PUBLIC_IP4"               // specify the public ipv4 address the torrent service.
	TorrentPublicIP6        = "RETROVIBED_TORRENT_PUBLIC_IP6"               // specify the public ipv6 address the torrent service.
	TorrentLogging          = "RETROVIBED_TORRENT_LOGGING"                  // enable/disable torrent logging. strconv.ParseBool
	TorrentDebug            = "RETROVIBED_TORRENT_DEBUG"                    // enable/disable torrent debug logging. strconv.ParseBool
	TorrentDisableWireguard = "RETROVIBED_TORRENT_DISABLE_WIREGUARD"        // enable/disable torrent vpn functionality. strconv.ParseBool
	TorrentPEX              = "RETROVIBED_TORRENT_PEX"                      // enable/disable torrent pex functionality. strconv.ParseBool
	TorrentAllowSeeding     = "RETROVIBED_TORRENT_ALLOW_SEEDING"            // enable/disable torrent allow the daemon to seed. strconv.ParseBool
	TorrentDownloadStats    = "RETROVIBED_TORRENT_DOWNLOAD_STATS_FREQUENCY" // adjust the frequency at which download stats are logged. time.Duration format.
	TorrentPrivateNetwork   = "RETROVIBED_TORRENT_PRIVATE"                  // specify that the torrent system should firewall to only private networks.
	TorrentDirectoryWatch   = "RETROVIBED_TORRENT_DIRECTORY_WATCH"          // specify the list of directories to watch for torrent files.
	TorrentVerifyFrequency  = "RETROVIBED_TORRENT_VERIFY_FREQUENCY"         // specify how frequently to check for torrents to verify.
	DHTDebug                = "RETROVIBED_DHT_DEBUG"                        // enable dht debug logging
	SelfSignedHosts         = "RETROVIBED_SELF_SIGNED_HOSTS"                // list of hosts to add to the self signed certificate.
	DDiscFrequency          = "RETROVIBED_DDISC_ANNOUNCE_FREQUENCY"         // how frequently to announce partitions frequency
	LoggingVerbosity        = "RETROVIBED_LOGGING_VERBOSITY"                // controls logging verbosity level
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
