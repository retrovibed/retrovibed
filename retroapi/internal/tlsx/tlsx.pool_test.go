package tlsx_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/retrovibed/retroapi/internal/tlsx"
	"github.com/stretchr/testify/require"
)

// newCAAndSignedCert builds a CA cert and a leaf cert signed by it, valid for the given window.
// The leaf has DNSNames: ["localhost.lan"] so x509.Verify passes with DNSName: "localhost.lan".
func newCAAndSignedCert(t *testing.T, notBefore, notAfter time.Time) (*x509.CertPool, *tls.Certificate) {
	t.Helper()

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test-ca"},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	require.NoError(t, err)
	caCert, err := x509.ParseCertificate(caDER)
	require.NoError(t, err)

	leafKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	leafTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "localhost.lan"},
		DNSNames:              []string{"localhost.lan"},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, caCert, &leafKey.PublicKey, caKey)
	require.NoError(t, err)
	leafParsed, err := x509.ParseCertificate(leafDER)
	require.NoError(t, err)

	caPool := x509.NewCertPool()
	caPool.AddCert(caCert)

	return caPool, &tls.Certificate{Certificate: [][]byte{leafDER}, PrivateKey: leafKey, Leaf: leafParsed}
}

// newTLSCertHost builds a self-signed *tls.Certificate for the given hostname and validity window.
func newTLSCertHost(t *testing.T, hostname string, notBefore, notAfter time.Time) *tls.Certificate {
	t.Helper()

	templ := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: hostname},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	der, err := x509.CreateCertificate(rand.Reader, templ, templ, &priv.PublicKey, priv)
	require.NoError(t, err)

	leaf, err := x509.ParseCertificate(der)
	require.NoError(t, err)

	return &tls.Certificate{
		Certificate: [][]byte{der},
		PrivateKey:  priv,
		Leaf:        leaf,
	}
}

// newTLSCert builds a self-signed *tls.Certificate with the given validity window.
// CommonName is "localhost.lan" so connectionHostname falls back to it for IP connections.
func newTLSCert(t *testing.T, notBefore, notAfter time.Time) *tls.Certificate {
	return newTLSCertHost(t, "localhost.lan", notBefore, notAfter)
}

// tlsServer starts a TLS httptest server serving the given certificate.
// Certificates is pre-populated so httptest does not inject its own cert, and
// Go's TLS stack serves ours unconditionally (no SNI required for IP clients).
const tlsServerBody = "ok"

func tlsServer(t *testing.T, cert *tls.Certificate) *httptest.Server {
	t.Helper()
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, tlsServerBody)
	}))
	srv.TLS = &tls.Config{
		Certificates: []tls.Certificate{*cert},
	}
	srv.StartTLS()
	t.Cleanup(srv.Close)
	return srv
}

// freshClient returns an *http.Client with the pool's TLS config and no
// connection reuse, so each Get triggers a full TLS handshake.
func freshClient(t *testing.T, pool tlsx.Pool) *http.Client {
	t.Helper()
	return &http.Client{Transport: &http.Transport{TLSClientConfig: pool.Config(), DisableKeepAlives: true}}
}

func requireBody(t *testing.T, resp *http.Response) {
	t.Helper()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, tlsServerBody, string(body))
}

func hostnamePath(dir, hostname string) string {
	return filepath.Join(dir, tlsx.HostnameKey(hostname))
}

func TestPool_TLS(t *testing.T) {
	now := time.Now()

	t.Run("CA-signed cert is accepted without TOFU enrollment", func(t *testing.T) {
		dir := t.TempDir()

		caPool, leafCert := newCAAndSignedCert(t, now.Add(-time.Hour), now.Add(time.Hour))
		pool, err := tlsx.NewPool(dir, tlsx.WithRoots(caPool))
		require.NoError(t, err)

		srv := tlsServer(t, leafCert)
		resp, err := freshClient(t, pool).Get(srv.URL)
		require.NoError(t, err, "CA-signed cert should be accepted via root pool")
		requireBody(t, resp)

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Empty(t, entries, "CA-signed cert should not be stored in the TOFU store")
	})

	t.Run("localhost cert is accepted without TOFU enrollment", func(t *testing.T) {
		dir := t.TempDir()
		pool, err := tlsx.NewPool(dir)
		require.NoError(t, err)

		cert := newTLSCertHost(t, "localhost", now.Add(-time.Hour), now.Add(time.Hour))
		srv := tlsServer(t, cert)

		resp, err := freshClient(t, pool).Get(srv.URL)
		require.NoError(t, err, "localhost cert should be accepted via bypass")
		requireBody(t, resp)

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Empty(t, entries, "localhost cert should not be stored in the TOFU store")
	})

	t.Run("first request enrolls cert", func(t *testing.T) {
		dir := t.TempDir()
		pool, err := tlsx.NewPool(dir)
		require.NoError(t, err)

		cert := newTLSCert(t, now.Add(-time.Hour), now.Add(time.Hour))
		srv := tlsServer(t, cert)

		resp, err := freshClient(t, pool).Get(srv.URL)
		require.NoError(t, err)
		requireBody(t, resp)

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, entries, 1, "cert should be persisted after first request")
	})

	t.Run("second request with same cert succeeds", func(t *testing.T) {
		dir := t.TempDir()
		pool, err := tlsx.NewPool(dir)
		require.NoError(t, err)

		cert := newTLSCert(t, now.Add(-time.Hour), now.Add(time.Hour))
		srv := tlsServer(t, cert)
		client := freshClient(t, pool)

		resp, err := client.Get(srv.URL)
		require.NoError(t, err)
		requireBody(t, resp)

		resp, err = client.Get(srv.URL)
		require.NoError(t, err)
		requireBody(t, resp)
	})

	// Both certs have CommonName "localhost.lan" so they share the same TOFU key.
	// After enrolling cert1, presenting cert2 for the same host must be rejected.
	t.Run("changed cert for same host is rejected", func(t *testing.T) {
		dir := t.TempDir()
		pool, err := tlsx.NewPool(dir)
		require.NoError(t, err)

		cert1 := newTLSCert(t, now.Add(-time.Hour), now.Add(time.Hour))
		cert2 := newTLSCert(t, now.Add(-time.Hour), now.Add(time.Hour))

		srv1 := tlsServer(t, cert1)
		srv2 := tlsServer(t, cert2)
		client := freshClient(t, pool)

		// Enroll cert1 under key "localhost.lan".
		resp, err := client.Get(srv1.URL)
		require.NoError(t, err)
		requireBody(t, resp)

		// cert2 shares the same key but has different bytes — mismatch.
		_, err = client.Get(srv2.URL)
		require.Error(t, err, "changed cert for the same host should be rejected")
	})

	t.Run("enrolled cert survives pool reload from disk", func(t *testing.T) {
		dir := t.TempDir()

		cert := newTLSCert(t, now.Add(-time.Hour), now.Add(time.Hour))
		srv := tlsServer(t, cert)

		pool1, err := tlsx.NewPool(dir)
		require.NoError(t, err)
		resp, err := freshClient(t, pool1).Get(srv.URL)
		require.NoError(t, err)
		requireBody(t, resp)

		pool2, err := tlsx.NewPool(dir)
		require.NoError(t, err)
		resp, err = freshClient(t, pool2).Get(srv.URL)
		require.NoError(t, err)
		requireBody(t, resp)
	})

	t.Run("renewed cert is accepted after old cert expired", func(t *testing.T) {
		dir := t.TempDir()
		pool, err := tlsx.NewPool(dir)
		require.NoError(t, err)

		expired := newTLSCert(t, now.Add(-25*time.Hour), now.Add(-time.Hour))
		renewed := newTLSCert(t, now.Add(-time.Minute), now.Add(time.Hour))

		// Pre-populate the TOFU store with the expired cert (simulating prior enrollment).
		p := hostnamePath(dir, "localhost.lan")
		require.NoError(t, os.MkdirAll(dir, 0700))
		require.NoError(t, os.WriteFile(p, expired.Leaf.Raw, 0600))

		// Connecting to the expired cert removes the stored file.
		srvExpired := tlsServer(t, expired)
		client := freshClient(t, pool)
		_, err = client.Get(srvExpired.URL)
		require.Error(t, err, "expired cert should be rejected")

		_, statErr := os.Stat(p)
		require.True(t, os.IsNotExist(statErr), "stored .der file should be removed after expiry rejection")

		// The renewed cert is now accepted via TOFU and the connection succeeds.
		srvRenewed := tlsServer(t, renewed)
		resp, err := client.Get(srvRenewed.URL)
		require.NoError(t, err, "renewed cert should be accepted after old cert expired")
		requireBody(t, resp)

		_, statErr = os.Stat(p)
		require.NoError(t, statErr, "renewed cert should be persisted to TOFU store")
	})

	t.Run("renewed cert is accepted without presenting old expired cert first", func(t *testing.T) {
		dir := t.TempDir()
		pool, err := tlsx.NewPool(dir)
		require.NoError(t, err)

		expired := newTLSCert(t, now.Add(-25*time.Hour), now.Add(-time.Hour))
		renewed := newTLSCert(t, now.Add(-time.Minute), now.Add(time.Hour))

		// Pre-populate the TOFU store with the expired cert (simulating prior enrollment).
		p := hostnamePath(dir, "localhost.lan")
		require.NoError(t, os.MkdirAll(dir, 0700))
		require.NoError(t, os.WriteFile(p, expired.Leaf.Raw, 0600))

		// The renewed cert is accepted directly — stored expired cert is replaced via TOFU.
		srvRenewed := tlsServer(t, renewed)
		resp, err := freshClient(t, pool).Get(srvRenewed.URL)
		require.NoError(t, err, "renewed cert should be accepted when stored cert is expired")
		requireBody(t, resp)

		_, statErr := os.Stat(p)
		require.NoError(t, statErr, "renewed cert should be persisted to TOFU store")
	})

	t.Run("expired cert is rejected and stored file is removed", func(t *testing.T) {
		dir := t.TempDir()
		pool, err := tlsx.NewPool(dir)
		require.NoError(t, err)

		expired := newTLSCert(t, now.Add(-25*time.Hour), now.Add(-time.Hour))
		srv := tlsServer(t, expired)
		_ = srv

		// The TOFU key for a cert with CommonName "localhost.lan" and no SANs is "localhost.lan".
		p := hostnamePath(dir, "localhost.lan")
		require.NoError(t, os.MkdirAll(dir, 0700))
		require.NoError(t, os.WriteFile(p, expired.Leaf.Raw, 0600))

		_, err = freshClient(t, pool).Get(srv.URL)
		require.Error(t, err, "expired cert should be rejected")

		_, statErr := os.Stat(p)
		require.True(t, os.IsNotExist(statErr), "stored .der file should be removed after expiry")
	})
}
