package env

import (
	"sync"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
)

func MediaDir() string                    { return env.MediaDir() }
func TorrentDir() string                  { return env.TorrentDir() }
func PrivateKeyPath() string              { return env.PrivateKeyPath(userx.DefaultRelRoot()) }
func TLSPoolDir() string                  { return env.TLSPoolDir() }
func RootStorageDir(rel ...string) string { return env.RootStorageDir(rel...) }

const (
	// smoke test exit after initialization
	Smoke = "RETROVIBED_SMOKE"

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
	JWTSharedSecret            = "RETROVIBED_JWT_SECRET"
	Endpoint                   = "RETROVIBED_HTTP_ENDPOINT"                    // http address for the retrovibed daemon
	NeuralMediaID              = "RETROVIBED_NEURALS_MEDIA_ID"                 // model for extracting media metadata from essentially file names. used in identification of media.
	AutoDownloadNeurals        = "RETROVIBED_NEURALS_AUTODOWNLOAD"             // enable/disable automatically downloading ai models.
	AutoDownloadMetadata       = "RETROVIBED_METADATA_AUTODOWNLOAD"            // enable/disable automatically downloading metadata.
	MDNSAdvertise              = "RETROVIBED_MDNS_ADVERTISE"                   // enable/disable multicast dns registration, allows for the frontend to automatically find daemons on the local network.
	MDNSDiscovery              = "RETROVIBED_MDNS_DISCOVERY"                   // enable/disable api driven mdns lan peer discovery.
	AutoArchive                = "RETROVIBED_MEDIA_AUTO_ARCHIVE"               // enable/disable automatic archiving of eligible media.
	AutoReclaim                = "RETROVIBED_MEDIA_AUTO_RECLAIM"               // enable/disable automatic reclaiming of disk space for media that has been archived.
	AutoIdentifyMedia          = "RETROVIBED_MEDIA_AUTO_IDENTIFY"              // enable/disable automatically identified media.
	AutoLocateMedia            = "RETROVIBED_MEDIA_AUTO_LOCATE"                // enable/disable automatically locate and download media.
	AutoDiscovery              = "RETROVIBED_TORRENT_AUTO_DISCOVERY"           // enable/disable automatically discovering torrents from peers.
	AutoPeerTube               = "RETROVIBED_AUTO_PEERTUBE"                    // enable/disable the built-in PeerTube/SepiaSearch discovery strategy.
	AutoSubscriptions          = "RETROVIBED_AUTO_SUBSCRIPTIONS"               // enable/disable auto subscription setup, primarily used for retrokiosk since they are low powered devices
	PeerTubeDomain             = "PEERTUBE_DOMAIN"                             // base url of the PeerTube/SepiaSearch instance to search.
	AutoBootstrap              = "RETROVIBED_TORRENT_AUTO_BOOTSTRAP"           // enable/disable the predefined set of public swarms to bootstrap from.
	TorrentPort                = "RETROVIBED_TORRENT_PORT"                     // specify the port to listen to torrents on
	TorrentPublicIP4           = "RETROVIBED_TORRENT_PUBLIC_IP4"               // specify the public ipv4 address the torrent service.
	TorrentPublicIP6           = "RETROVIBED_TORRENT_PUBLIC_IP6"               // specify the public ipv6 address the torrent service.
	TorrentLogging             = "RETROVIBED_TORRENT_LOGGING"                  // enable/disable torrent logging. strconv.ParseBool
	TorrentDebug               = "RETROVIBED_TORRENT_DEBUG"                    // enable/disable torrent debug logging. strconv.ParseBool
	TorrentDisableWireguard    = "RETROVIBED_TORRENT_DISABLE_WIREGUARD"        // enable/disable torrent vpn functionality. strconv.ParseBool
	WireguardAutohealFrequency = "RETROVIBED_WIREGUARD_AUTOHEAL_FREQUENCY"     // base interval for wireguard autoheal exponential backoff. time.Duration format.
	WireguardAutohealMax       = "RETROVIBED_WIREGUARD_AUTOHEAL_MAX"           // maximum interval between wireguard autoheal recovery attempts. time.Duration format.
	TorrentPEX                 = "RETROVIBED_TORRENT_PEX"                      // enable/disable torrent pex functionality. strconv.ParseBool
	TorrentAllowSeeding        = "RETROVIBED_TORRENT_ALLOW_SEEDING"            // enable/disable torrent allow the daemon to seed. strconv.ParseBool
	TorrentDownloadStats       = "RETROVIBED_TORRENT_DOWNLOAD_STATS_FREQUENCY" // adjust the frequency at which download stats are logged. time.Duration format.
	TorrentPrivateNetwork      = "RETROVIBED_TORRENT_PRIVATE"                  // specify that the torrent system should firewall to only private networks.
	TorrentDirectoryWatch      = "RETROVIBED_TORRENT_DIRECTORY_WATCH"          // specify the list of directories to watch for torrent files.
	TorrentVerifyFrequency     = "RETROVIBED_TORRENT_VERIFY_FREQUENCY"         // specify how frequently to check for torrents to verify.
	DHTDebug                   = "RETROVIBED_DHT_DEBUG"                        // enable dht debug logging
	SelfSignedHosts            = "RETROVIBED_SELF_SIGNED_HOSTS"                // list of hosts to add to the self signed certificate.
	DDiscFrequency             = "RETROVIBED_DDISC_ANNOUNCE_FREQUENCY"         // how frequently to announce partitions frequency.
	DDiscP2PLocate             = "RETROVIBED_DDISC_P2P_LOCATE"                 // enable the discovery of media from peers.
	DDiscIndexRatio            = "RETROVIBED_DDISC_INDEX_RATIO"                // what percentage of discovered media to attempt to index.
	DDiscBackgroundFrequency   = "RETROVIBED_DDISC_BACKGROUND_FREQUENCY"       // how frequently to push work into the queue.
	DDiscBackgroundWorkers     = "RETROVIBED_DDISC_BACKGROUND_WORKERS"         // number of identifying workers
	LoggingVerbosity           = "RETROVIBED_LOGGING_VERBOSITY"                // controls logging verbosity level
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
