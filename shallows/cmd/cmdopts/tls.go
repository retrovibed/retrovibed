package cmdopts

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/tlsx"
)

type TLSConfig struct {
	Insecure   bool   `help:"allow unsigned (and therefor insecure) tls certificates to be used" name:"insecure" env:"RETROVIBED_TLS_INSECURE" hidden:"true"`
	Servername string `help:"name the certificate should authorized for" name:"servername" default:"localhost:9998" env:"RETROVIBED_TLS_SERVERNAME"`
	PoolDir    string `help:"directory for trust-on-first-use certificate pool" name:"tls-pool-dir" default:"${vars_user_configuration_directory}/tls.d" env:"RETROVIBED_TLS_POOL_DIR"`
	once       sync.Once
	pool       tlsx.Pool
}

func (t *TLSConfig) AfterApply() error {
	var err error
	t.once.Do(func() {
		t.pool, err = tlsx.NewPool(t.PoolDir)
	})
	return err
}

func (t *TLSConfig) Config() *tls.Config {
	return tlsx.MustClone(t.pool.Config(), tlsx.OptionServerName(t.Servername))
}

func (t *TLSConfig) DefaultClient() *http.Client {
	ctransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       5 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       t.Config(),
	}

	defaultclient := &http.Client{Transport: ctransport}
	defaultclient = httpx.BindRetryTransport(defaultclient, http.StatusTooManyRequests, http.StatusBadGateway, http.StatusInternalServerError, http.StatusRequestTimeout)

	// if env.Boolean(false, eg.EnvLogsNetwork) {
	// 	return httpx.DebugClient(defaultclient)
	// }

	return defaultclient
}
