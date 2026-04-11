package tlsx

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/retrovibed/retrovibed/internal/debugx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/netx"
)

// Pool is a directory-based trust-on-first-use certificate pool.
// Each hostname's trusted certificate is stored in a file named after the
// SHA-256 of the hostname, so cert changes for a given host are detected.
type Pool struct {
	dir   string
	roots *x509.CertPool
}

// PoolOption configures a Pool.
type PoolOption func(*Pool)

// WithRoots sets the CA root pool used for verifying CA-signed certificates.
// Passing nil restores the default (system roots).
func WithRoots(roots *x509.CertPool) PoolOption {
	return func(p *Pool) { p.roots = roots }
}

// NewPool returns a Pool backed by dir.
func NewPool(dir string, opts ...PoolOption) (p Pool, err error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return p, err
	}

	return langx.Clone(Pool{dir: dir, roots: pool}, opts...), nil
}

// Config returns a *tls.Config that performs trust-on-first-use during
// the TLS handshake. InsecureSkipVerify is set internally so that unknown certs
// reach VerifyConnection; that callback does the TOFU check keyed by hostname.
func (t Pool) Config() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true, //nolint:gosec // TOFU: custom verification via VerifyConnection
		VerifyConnection: func(cs tls.ConnectionState) error {
			if len(cs.PeerCertificates) == 0 {
				return errorsx.New("no certificates provided by peer")
			}
			hostname := connectionHostname(cs)
			return t.Validate(context.Background(), hostname, cs.PeerCertificates)
		},
	}
}

// connectionHostname returns the best available hostname for keying a TOFU entry.
// For named connections cs.ServerName is set; for raw IP connections it falls
// back to the cert's IP SANs, DNS SANs, then Subject CommonName.
func connectionHostname(cs tls.ConnectionState) string {
	cert := cs.PeerCertificates[0]
	return langx.FirstNonZero(
		cs.ServerName,
		netx.IPString(netx.FirstNonZeroIP(cert.IPAddresses...)),
		langx.FirstNonZero(cert.DNSNames...),
		cert.Subject.CommonName,
	)
}

// HostnameKey returns the filename key for a TOFU entry keyed by hostname.
func HostnameKey(hostname string) string {
	h := sha256.Sum256([]byte(hostname))
	return hex.EncodeToString(h[:])
}

// Validate checks whether cert is trusted for the given hostname.
// System CA-signed certificates are verified against the system root pool and
// accepted without TOFU enrollment. For self-signed certificates, trust-on-first-use
// applies: the DER bytes are stored on first encounter and must match on subsequent
// calls. Expired certificates are rejected and their stored file is removed so a
// renewed certificate can be enrolled on the next connection.
func (t Pool) Validate(ctx context.Context, hostname string, certs []*x509.Certificate) error {
	if len(certs) == 0 {
		return errorsx.New("no certificates provided")
	}

	cert, intermediates := certs[0], certs[1:]
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return errorsx.Errorf("certificate is not yet valid (valid from %s)", cert.NotBefore)
	}

	path := filepath.Join(t.dir, HostnameKey(hostname))

	if now.After(cert.NotAfter) {
		// Remove stale entry so a renewed cert can be trusted on the next connection.
		errorsx.Log(errorsx.Wrapf(fsx.IgnoreIsNotExist(os.Remove(path)), "failed to cleanup expired certificate for %s", hostname))
		return errorsx.Errorf("certificate has expired (expired %s)", cert.NotAfter)
	}

	opts := x509.VerifyOptions{
		DNSName:       hostname,
		CurrentTime:   now,
		Roots:         t.roots,
		Intermediates: x509.NewCertPool(), // Start empty
	}

	// Just add everything from the slice index 1 onwards.
	// If it's empty, nothing happens.
	for _, icert := range intermediates {
		opts.Intermediates.AddCert(icert)
	}

	// System CA-signed certs don't need TOFU enrollment.
	if _, err := cert.Verify(opts); err == nil {
		return nil
	} else {
		debugx.Println("failed to use system certificates to validate", hostname, err)
	}

	// ignore localhost, this is a problem but our plan is to replace tofu completely eventually.
	if strings.TrimSpace(hostname) == "localhost" {
		return nil
	}

	if err := os.MkdirAll(t.dir, 0700); err != nil {
		return errorsx.WithStack(err)
	}

	existing, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// Trust on first use: persist.
		if err := os.WriteFile(path, cert.Raw, 0600); err != nil {
			return errorsx.WithStack(err)
		}
		return nil
	}

	if err != nil {
		return errorsx.WithStack(err)
	}

	if bytes.Equal(existing, cert.Raw) {
		return nil
	}

	// If the stored cert is expired, allow the new cert to replace it via TOFU.
	if stored, err := x509.ParseCertificate(existing); err == nil && now.After(stored.NotAfter) {
		if err := os.WriteFile(path, cert.Raw, 0600); err != nil {
			return errorsx.WithStack(err)
		}
		return nil
	}

	return errorsx.Errorf("certificate mismatch: certificate has changed since first use: %s", hostname)
}
