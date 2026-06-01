package tlsx_test

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/tlsx"
	"github.com/stretchr/testify/require"
)

func TestExpiredCertificate(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "tls.pem")
		require.True(t, tlsx.ExpiredCertificate(path))
	})

	t.Run("valid cert", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "tls.pem")
		require.NoError(t, tlsx.SelfSignedLocalHostTLSSeeded(rand.Reader, path))
		require.False(t, tlsx.ExpiredCertificate(path))
	})

	t.Run("expired cert", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "tls.pem")
		priv, derbytes := expiredCert(t)
		require.NoError(t, tlsx.WritePEMFile(path, priv, derbytes))
		require.True(t, tlsx.ExpiredCertificate(path))
	})

	t.Run("corrupt content", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "tls.pem")
		require.NoError(t, os.WriteFile(path, []byte("not valid pem"), 0600))
		require.True(t, tlsx.ExpiredCertificate(path))
	})
}
