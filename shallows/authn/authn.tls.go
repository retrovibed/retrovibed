package authn

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"sync"

	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/tlsx"
)

var initpool = sync.OnceValue(func() tlsx.Pool {
	pool, err := tlsx.NewPool(env.TLSPoolDir())
	if err != nil {
		log.Fatalln("failed to initialize tls pool", err)
	}
	return pool
})

// TLSConfig returns a *tls.Config backed by the process-wide TOFU certificate
// pool. The pool is initialized once on first call using the default pool dir.
func TLSConfig() *tls.Config {
	return initpool().Config()
}

// ValidateCertificate validates a DER-encoded certificate for a given hostname
// against the process-wide TOFU pool.
func ValidateCertificate(ctx context.Context, hostname string, der []byte) error {
	cert, err := x509.ParseCertificates(der)
	if err != nil {
		return err
	}

	return initpool().Validate(ctx, hostname, cert)
}
