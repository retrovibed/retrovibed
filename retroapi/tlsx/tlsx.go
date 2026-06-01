package tlsx

import (
	"crypto"
	"crypto/ed25519"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
)

type privatekey interface {
	Public() crypto.PublicKey
	Equal(x crypto.PrivateKey) bool
}

// X509Option ...
type X509Option func(*x509.Certificate)

// X509OptionSubject subject for the cert
func X509OptionSubject(s pkix.Name) X509Option {
	return func(t *x509.Certificate) {
		t.Subject = s
	}
}

// X509OptionCA enables the certificate as a ca.
func X509OptionCA(t *x509.Certificate) {
	t.IsCA = true
	t.KeyUsage = t.KeyUsage | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
	t.ExtKeyUsage = append(t.ExtKeyUsage, x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth)
}

// X509OptionHosts set the hosts
func X509OptionHosts(names ...string) X509Option {
	return func(t *x509.Certificate) {
		for _, h := range names {
			if ip := net.ParseIP(h); ip != nil {
				t.IPAddresses = append(t.IPAddresses, ip)
			} else {
				t.DNSNames = append(t.DNSNames, h)
			}
		}
	}
}

// X509OptionUsage set the usage options for the certificate.
func X509OptionUsage(u x509.KeyUsage) X509Option {
	return func(t *x509.Certificate) {
		t.KeyUsage = t.KeyUsage | u
	}
}

// X509OptionUsageExt set the usage extension bits.
func X509OptionUsageExt(u ...x509.ExtKeyUsage) X509Option {
	return func(t *x509.Certificate) {
		t.ExtKeyUsage = u
	}
}

// X509OptionTimeWindow where the certificate is valid.
// clock can be nil.
func X509OptionTimeWindow(c clock, d time.Duration) X509Option {
	return func(cert *x509.Certificate) {
		cert.NotBefore = c.Now()
		cert.NotAfter = cert.NotBefore.Add(d)
	}
}

// X509OptionCompose
func X509OptionCompose(options ...X509Option) X509Option {
	return func(cert *x509.Certificate) {
		for _, opt := range options {
			opt(cert)
		}
	}
}

type clock interface {
	Now() time.Time
}

type stdlibclock struct{}

func (t stdlibclock) Now() time.Time {
	return time.Now()
}

// DefaultClock that can be used to generate a template cert.
func DefaultClock() clock {
	return stdlibclock{}
}

// X509TemplateRand generate a template using the provided random source.
// the clock is allowed to be nil.
func X509TemplateRand(r io.Reader, d time.Duration, c clock, options ...X509Option) (template x509.Certificate, err error) {
	var (
		serialNumber *big.Int
	)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)

	if serialNumber, err = rand.Int(r, serialNumberLimit); err != nil {
		return template, errorsx.WithStack(err)
	}

	orgHash := md5.New()
	if _, err = io.CopyN(orgHash, r, 1024); err != nil {
		return template, errorsx.WithStack(err)
	}

	template = x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{hex.EncodeToString(orgHash.Sum(nil))},
		},
		KeyUsage:              0,
		ExtKeyUsage:           nil,
		BasicConstraintsValid: true,
	}

	// ensure there is always a valid window.
	X509OptionTimeWindow(stdlibclock{}, d)(&template)

	for _, opt := range options {
		opt(&template)
	}

	return template, errorsx.WithStack(err)
}

// X509Template ...
func X509Template(d time.Duration, options ...X509Option) (template x509.Certificate, err error) {
	return X509TemplateRand(rand.Reader, d, stdlibclock{}, options...)
}

func SelfSignedED25519(r io.Reader, template *x509.Certificate) (priv ed25519.PrivateKey, derBytes []byte, err error) {
	if _, priv, err = ed25519.GenerateKey(r); err != nil {
		return priv, derBytes, err
	}

	return SelfSigned(priv, template)
}

// SelfSigned signs its own certificate ..
func SelfSigned[T privatekey](priv T, template *x509.Certificate) (_ T, derBytes []byte, err error) {
	return SelfSignedRand(rand.Reader, priv, template)
}

// SelfSignedRand signs its own certificate ..
func SelfSignedRand[T privatekey](r io.Reader, priv T, template *x509.Certificate) (_ T, derBytes []byte, err error) {
	return SignedRand(r, priv, template, template)
}

// SignedRand signs its own certificate ..
func SignedRand[T privatekey](r io.Reader, priv T, template, parent *x509.Certificate) (_ T, derBytes []byte, err error) {
	if derBytes, err = x509.CreateCertificate(r, template, parent, priv.Public(), priv); err != nil {
		return priv, derBytes, errorsx.WithStack(err)
	}

	return priv, derBytes, nil
}

// WritePEMFile ...
func WritePEMFile(path string, key ed25519.PrivateKey, derBytes []byte) (err error) {
	var (
		dst *os.File
	)

	if dst, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600); err != nil {
		return err
	}

	return errorsx.Compact(WritePEM(dst, key, derBytes), dst.Close())
}

// WritePEM ...
func WritePEM(dst io.Writer, key ed25519.PrivateKey, derBytes []byte) (err error) {
	if err = pem.Encode(dst, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return errorsx.WithStack(err)
	}

	pkcs8, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return errorsx.Wrap(err, "failed to marshal ed25519 private key")
	}

	if err = pem.Encode(dst, &pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8}); err != nil {
		return errorsx.WithStack(err)
	}

	return nil
}

// WritePrivateKey ...
func WritePrivateKey(dst io.Writer, key ed25519.PrivateKey) error {
	pkcs8, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return errorsx.Wrap(err, "failed to marshal ed25519 private key")
	}
	return errorsx.WithStack(pem.Encode(dst, &pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8}))
}

// WritePrivateKeyFile ...
func WritePrivateKeyFile(path string, key ed25519.PrivateKey) (err error) {
	var (
		dst *os.File
	)

	if dst, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600); err != nil {
		return err
	}

	return errorsx.Compact(WritePrivateKey(dst, key), dst.Close())
}

// WriteCertificate ...
func WriteCertificate(dst io.Writer, cert []byte) error {
	return errorsx.WithStack(pem.Encode(dst, &pem.Block{Type: "CERTIFICATE", Bytes: cert}))
}

// WriteCertificateFile ...
func WriteCertificateFile(path string, cert []byte) (err error) {
	var (
		dst *os.File
	)

	if dst, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600); err != nil {
		return err
	}

	return errorsx.Compact(WriteCertificate(dst, cert), dst.Sync(), dst.Close())
}

// Option tls config options
type Option func(*tls.Config) error

// OptionVerifyClientIfGiven see tls.VerifyClientCertIfGiven
func OptionVerifyClientIfGiven(c *tls.Config) error {
	c.ClientAuth = tls.VerifyClientCertIfGiven
	return nil
}

// OptionNoClientCert see tls.NoClientCert
func OptionNoClientCert(c *tls.Config) error {
	c.ClientAuth = tls.NoClientCert
	return nil
}

// OptionServerName see tls.Config.ServerName
func OptionServerName(name string) Option {
	return func(c *tls.Config) error {
		c.ServerName = name
		return nil
	}
}

// OptionInsecureSkipVerify see tls.Config.InsecureSkipVerify
func OptionInsecureSkipVerify(c *tls.Config) error {
	c.InsecureSkipVerify = true
	return nil
}

// OptionNextProtocols ALPN see tls.NextProtos
func OptionNextProtocols(protocols ...string) Option {
	return func(c *tls.Config) error {
		c.NextProtos = append(c.NextProtos, protocols...)
		return nil
	}
}

// Clone ...
func Clone(c *tls.Config, options ...Option) (updated *tls.Config, err error) {
	updated = c.Clone()

	for _, opt := range options {
		if err = opt(updated); err != nil {
			return updated, err
		}
	}

	return updated, nil
}

// MustClone ...
func MustClone(c *tls.Config, options ...Option) *tls.Config {
	updated, err := Clone(c, options...)
	if err != nil {
		panic(err)
	}
	return updated
}

// DecodePEMCertificate decode a pem encoded x509 certiciate.
func DecodePEMCertificate(encoded []byte) (cert *x509.Certificate, err error) {
	var (
		p *pem.Block
	)

	if p, _ = pem.Decode(encoded); p == nil {
		return cert, errorsx.Wrap(err, "unable to decode pem certificate")
	}

	if cert, err = x509.ParseCertificate(p.Bytes); err != nil {
		return cert, errorsx.Wrap(err, "failed parse certificate")
	}

	return cert, nil
}

// DecodePEMPrivateKey returns the first PEM block whose type ends in "PRIVATE KEY",
// skipping any leading non-key blocks (e.g. certificates). The caller is responsible
// for parsing the block bytes based on block.Type.
func DecodePEMPrivateKey(encoded []byte) (*pem.Block, error) {
	rest := encoded
	for {
		var p *pem.Block
		p, rest = pem.Decode(rest)
		if p == nil {
			return nil, errorsx.New("no private key block found in PEM data")
		}

		if strings.HasSuffix(p.Type, "PRIVATE KEY") {
			return p, nil
		}
	}
}

// NewDialer for tls configurations.
func NewDialer(c *tls.Config, options ...Option) *tls.Dialer {
	return &tls.Dialer{
		Config: MustClone(c, options...),
		NetDialer: &net.Dialer{
			Timeout: 5 * time.Second,
		},
	}
}

// ExpiredCertificate reports whether the PEM file at path needs to be (re)generated.
// Returns true if the file is missing, unreadable, corrupt, or contains an expired cert.
// Returns false only when the file contains a valid, unexpired certificate.
func ExpiredCertificate(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return true
	}
	cert, err := DecodePEMCertificate(data)
	return err != nil || cert == nil || time.Now().After(cert.NotAfter)
}

// SelfSignedLocalHostTLS generates or renews a self-signed localhost cert at path.
func SelfSignedLocalHostTLS(path string, options ...X509Option) error {
	return SelfSignedLocalHostTLSSeeded(rand.Reader, path, options...)
}

// SelfSignedLocalHostTLSSeeded generates or renews a self-signed localhost cert at path using r as
// the random source. If the file already exists and the cert inside is still valid, it is left
// untouched. If the cert is expired (or unreadable), the file is removed and a new cert is written.
func SelfSignedLocalHostTLSSeeded(r io.Reader, path string, options ...X509Option) error {
	if !ExpiredCertificate(path) {
		return nil
	}

	certtempl, err := X509TemplateRand(
		r,
		365*24*time.Hour,
		stdlibclock{},
		X509OptionCA,
		X509OptionUsage(
			x509.KeyUsageKeyEncipherment|
				x509.KeyUsageDigitalSignature|
				x509.KeyUsageDataEncipherment,
		),
		X509OptionUsageExt(x509.ExtKeyUsageAny),
		X509OptionSubject(pkix.Name{
			CommonName: "localhost",
		}),
		X509OptionHosts(
			"localhost",
			"127.0.0.1",
			"[::1]",
		),
		X509OptionCompose(
			options...,
		),
	)

	if err != nil {
		return err
	}

	priv, derbytes, err := SelfSignedED25519(r, &certtempl)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	if err = WritePEMFile(path, priv, derbytes); err != nil {
		return err
	}

	return nil
}
