package sshx

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/retrovibed/retrovibed/internal/cryptox"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"golang.org/x/crypto/ssh"
)

// IsNoKeyFound check if ssh key is not found.
func IsNoKeyFound(err error) bool {
	return err.Error() == "ssh: no key found"
}

// Comment adds comment to the ssh public key.
func Comment(encoded []byte, comment string) []byte {
	if strings.TrimSpace(comment) == "" {
		return encoded
	}

	comment = " " + comment + "\r\n"
	return append(bytes.TrimSpace(encoded), []byte(comment)...)
}

type option func(*KeyGen)

func OptionKeyGenRand(src io.Reader) option {
	return func(kg *KeyGen) {
		kg.rand = src
	}
}

func NewKeyGenSeeded(seed string) *KeyGen {
	return NewKeyGen(OptionKeyGenRand(cryptox.NewPRNGSHA512([]byte(seed))))
}

func UnsafeNewKeyGen() *KeyGen {
	return NewKeyGen(OptionKeyGenRand(cryptox.NewPRNGSHA512([]byte("unsafe"))))
}

func NewKeyGen(options ...option) *KeyGen {
	kg := KeyGen{
		rand: nil, // if nil crypto packages use crypto/rand
	}

	for _, opt := range options {
		opt(&kg)
	}

	return &kg
}

type KeyGen struct {
	rand io.Reader
}

func (t KeyGen) Generate() (epriv, epub []byte, err error) {
	var (
		priv   ed25519.PrivateKey
		pub    ed25519.PublicKey
		pubkey ssh.PublicKey
		mpriv  []byte
	)

	if pub, priv, err = ed25519.GenerateKey(t.rand); err != nil {
		return nil, nil, err
	}

	if pubkey, err = ssh.NewPublicKey(pub); err != nil {
		return nil, nil, err
	}

	if mpriv, err = x509.MarshalPKCS8PrivateKey(priv); err != nil {
		return nil, nil, err
	}

	pemKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: mpriv,
	}

	return pem.EncodeToMemory(pemKey), ssh.MarshalAuthorizedKey(pubkey), nil
}

type keygen interface {
	Generate() (epriv, epub []byte, err error)
}

func loadcached(path string) (s ssh.Signer, err error) {
	var (
		privencoded []byte
	)

	if privencoded, err = os.ReadFile(path); err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(privencoded)
}

func SignerFromGenerator(kg keygen) (s ssh.Signer, err error) {
	var (
		privencoded []byte
	)

	if privencoded, _, err = kg.Generate(); err != nil {
		return nil, err
	}

	if s, err = ssh.ParsePrivateKey(privencoded); err != nil {
		return nil, err
	}

	return s, nil
}

func AutoCached(kg keygen, path string) (s ssh.Signer, err error) {
	var (
		privencoded, pubencoded []byte
	)

	if s, err = loadcached(path); err == nil {
		return s, nil
	}

	if privencoded, pubencoded, err = kg.Generate(); err != nil {
		return nil, err
	}

	if s, err = ssh.ParsePrivateKey(privencoded); err != nil {
		return nil, err
	}

	if err = os.WriteFile(path, privencoded, 0600); err != nil {
		return nil, err
	}

	if err = os.WriteFile(fmt.Sprintf("%s.pub", path), pubencoded, 0600); err != nil {
		return nil, err
	}

	return s, err
}

// ensure the public key exists
func EnsurePublicKey(s ssh.Signer, path string) error {
	return os.WriteFile(fmt.Sprintf("%s.pub", path), ssh.MarshalAuthorizedKey(s.PublicKey()), 0600)
}

func EncodeBase64PublicKey(pub ssh.PublicKey) string {
	return base64.RawURLEncoding.EncodeToString(pub.Marshal())
}

func DecodeBase64PublicKey(s string) (pub ssh.PublicKey, err error) {
	var (
		encoded []byte
	)

	if encoded, err = base64.RawURLEncoding.DecodeString(s); err != nil {
		return nil, err
	}

	return ssh.ParsePublicKey(encoded)
}

func FingerprintMD5(pub ssh.PublicKey) string {
	return md5x.FormatUUID(md5x.Digest(pub.Marshal()))
}
