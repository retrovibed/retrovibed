package torrenttestx

import (
	"log"
	"testing"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func QuickDHT(t testing.TB, options ...dht.Option) *dht.Server {
	dhts, err := dht.NewServer(
		32,
		options...,
	)
	require.NoError(t, err)
	return dhts
}

func QuickClient(t testing.TB, options ...torrent.ClientConfigOption) *torrent.Client {
	cdir := t.TempDir()
	return Client(t, autobind.NewLoopback(
		autobind.EnableDHT(QuickDHT(t)), // dht should be optional, but refactor made it not for now default to a basic one
	), torrent.NewMetadataCache(cdir), storage.NewFile(cdir), options...)
}

func QuickClientWithDHT(t testing.TB, dhts *dht.Server, options ...torrent.ClientConfigOption) *torrent.Client {
	cdir := t.TempDir()
	return Client(t, autobind.NewLoopback(autobind.EnableDHT(dhts)), torrent.NewMetadataCache(cdir), storage.NewFile(cdir), options...)
}

func Client(t testing.TB, binder autobind.Autobind, mdcache torrent.MetadataStore, scache storage.ClientImpl, options ...torrent.ClientConfigOption) *torrent.Client {
	cdir := t.TempDir()

	return testx.Must(binder.Bind(
		torrent.NewClient(
			torrent.NewDefaultClientConfig(
				mdcache,
				scache,
				torrent.ClientConfigCacheDirectory(cdir),
				torrent.ClientConfigSeed(true),
				torrent.ClientConfigInfoLogger(log.New(log.Writer(), "[torrent] ", log.Flags())),
				torrent.ClientConfigDebugLogger(log.New(log.Writer(), "[torrent] ", log.Flags())),
				torrent.ClientConfigCompose(options...),
				torrent.ClientConfigDialPoolSize(1),
				torrent.ClientConfigUploadLimit(rate.NewLimiter(rate.Inf, 10)),
			),
		),
	))(t)
}
