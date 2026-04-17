package sshx

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"strings"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/internal/cryptox"
	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/fsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/md5x"
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

func NewKeyGenSeeded[T ~string | []byte](seed T) *KeyGen {
	return NewKeyGen(OptionKeyGenRand(cryptox.NewChaCha8([]byte(seed))))
}

func UnsafeNewKeyGen() *KeyGen {
	return NewKeyGen(OptionKeyGenRand(cryptox.NewChaCha8([]byte("unsafe"))))
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

func Seeded(ctx context.Context, seed string, force bool, path string) (s ssh.Signer, err error) {
	ts := time.Now().Unix()
	backup := fmt.Sprintf("%s.%d", path, ts)
	backuppub := fmt.Sprintf("%s.pub.%d", path, ts)

	if exists := fsx.Exists(path); exists && !force {
		return nil, errorsx.Errorf("an identity already exists at %s, use --force to backup and replace", path)
	} else if exists {
		if err := os.Rename(path, backup); errors.Is(err, os.ErrNotExist) {
			log.Println("an identity already existed, renamed to", path, backup)
		} else if err != nil {
			return nil, errorsx.Wrap(err, "unable to rename old key")
		}

		if err := os.Rename(fmt.Sprintf("%s.pub", path), backuppub); errors.Is(err, os.ErrNotExist) {
			log.Println("an identity already existed, renamed to", fmt.Sprintf("%s.pub", path), backuppub)
		} else if err != nil {
			return nil, errorsx.Wrap(err, "unable to rename old public key")
		}
	}

	return AutoCached(NewKeyGenSeeded(seed), path)
}

func Unseeded(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.Remove(path + ".pub"); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
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

type Parsed struct {
	ssh.PublicKey
	Comment string
	Options []string
}

func ParseAuthorizedKeys(encoded []byte) iter.Seq[Parsed] {
	return func(yield func(Parsed) bool) {
		for len(encoded) != 0 {
			var (
				err error
				p   Parsed
			)

			if p.PublicKey, p.Comment, p.Options, encoded, err = ssh.ParseAuthorizedKey(encoded); err != nil {
				if IsNoKeyFound(err) {
					continue
				}
				log.Println(err)
				continue
			}

			if !yield(p) {
				return
			}
		}
	}
}
