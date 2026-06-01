package tlsx_test

import (
	"crypto/ecdsa"
	"crypto/rand"
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

	"github.com/retrovibed/retrovibed/retroapi/cryptox"
	"github.com/retrovibed/retrovibed/retroapi/tlsx"
	"github.com/stretchr/testify/require"
)

func expiredCert(t *testing.T) (*ecdsa.PrivateKey, []byte) {
	t.Helper()
	templ := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             time.Now().Add(-48 * time.Hour),
		NotAfter:              time.Now().Add(-24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	priv, der, err := tlsx.SelfSignedAuto(rand.Reader, templ)
	require.NoError(t, err)
	return priv, der
}

func TestSelfSignedLocalHostTLSSeeded(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "tls.pem")
		require.NoError(t, tlsx.SelfSignedLocalHostTLSSeeded(rand.Reader, path))
		require.FileExists(t, path)

		data, err := os.ReadFile(path)
		require.NoError(t, err)
		cert, err := tlsx.DecodePEMCertificate(data)
		require.NoError(t, err)
		require.True(t, time.Now().Before(cert.NotAfter), "generated cert should be valid")
	})

	t.Run("valid cert not regenerated", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "tls.pem")
		require.NoError(t, tlsx.SelfSignedLocalHostTLSSeeded(rand.Reader, path))

		info1, err := os.Stat(path)
		require.NoError(t, err)

		require.NoError(t, tlsx.SelfSignedLocalHostTLSSeeded(rand.Reader, path))

		info2, err := os.Stat(path)
		require.NoError(t, err)
		require.Equal(t, info1.ModTime(), info2.ModTime(), "valid cert file should not be modified")
	})

	t.Run("expired cert regenerated", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "tls.pem")

		priv, derbytes := expiredCert(t)
		require.NoError(t, tlsx.WritePEMFile(path, priv, derbytes))

		info1, err := os.Stat(path)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		require.NoError(t, tlsx.SelfSignedLocalHostTLSSeeded(rand.Reader, path))

		info2, err := os.Stat(path)
		require.NoError(t, err)
		require.NotEqual(t, info1.ModTime(), info2.ModTime(), "expired cert file should be replaced")

		data, err := os.ReadFile(path)
		require.NoError(t, err)
		cert, err := tlsx.DecodePEMCertificate(data)
		require.NoError(t, err)
		require.True(t, time.Now().Before(cert.NotAfter), "new cert should be valid")
	})

	t.Run("corrupt content regenerated", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "tls.pem")
		require.NoError(t, os.WriteFile(path, []byte("not valid pem"), 0600))

		require.NoError(t, tlsx.SelfSignedLocalHostTLSSeeded(rand.Reader, path))

		data, err := os.ReadFile(path)
		require.NoError(t, err)
		cert, err := tlsx.DecodePEMCertificate(data)
		require.NoError(t, err)
		require.True(t, time.Now().Before(cert.NotAfter), "regenerated cert should be valid")
	})

	t.Run("options forwarded", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "tls.pem")
		require.NoError(t, tlsx.SelfSignedLocalHostTLSSeeded(rand.Reader, path, tlsx.X509OptionHosts("testhost.local")))

		data, err := os.ReadFile(path)
		require.NoError(t, err)
		cert, err := tlsx.DecodePEMCertificate(data)
		require.NoError(t, err)
		require.Contains(t, cert.DNSNames, "testhost.local", "extra host from options should appear in cert SANs")
	})

	t.Run("private key stable across regeneration with same seed", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "tls.pem")

		require.NoError(t, tlsx.SelfSignedLocalHostTLSSeeded(cryptox.NewChaCha8(t.Name()), path))
		data1, err := os.ReadFile(path)
		require.NoError(t, err)
		key1, err := tlsx.DecodePEMPrivateKey(data1)
		require.NoError(t, err)

		// force regeneration by writing an expired cert
		epriv, ederbytes := expiredCert(t)
		require.NoError(t, tlsx.WritePEMFile(path, epriv, ederbytes))

		require.NoError(t, tlsx.SelfSignedLocalHostTLSSeeded(cryptox.NewChaCha8(t.Name()), path))
		data2, err := os.ReadFile(path)
		require.NoError(t, err)
		key2, err := tlsx.DecodePEMPrivateKey(data2)
		require.NoError(t, err)

		require.Equal(t, key1.Bytes, key2.Bytes, "same seed should reproduce the same private key after regeneration")
	})

	t.Run("generated cert can be loaded and served over https", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "tls.pem")
		require.NoError(t, tlsx.SelfSignedLocalHostTLSSeeded(rand.Reader, path))

		cert, err := tls.LoadX509KeyPair(path, path)
		require.NoError(t, err)

		srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, "ok")
		}))
		srv.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
		srv.StartTLS()
		t.Cleanup(srv.Close)

		pool, err := tlsx.NewPool(t.TempDir())
		require.NoError(t, err)
		client := &http.Client{Transport: &http.Transport{TLSClientConfig: pool.Config(), DisableKeepAlives: true}}

		resp, err := client.Get(srv.URL)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
