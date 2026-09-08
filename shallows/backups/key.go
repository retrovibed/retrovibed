package backups

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"net/http"
	"os"
	"slices"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

const saltSize = 16

// Key is the account's half and the device's half of the backup secret. the seed is issued
// once per account and the private key comes back from the identity seed, so a backup made
// on one device is recoverable on another.
type Key struct {
	seed       string
	privatekey []byte
}

func NewKey(seed string, privatekey []byte) (Key, error) {
	if seed == "" {
		return Key{}, errorsx.String("backup seed is required")
	}

	if len(privatekey) == 0 {
		return Key{}, errorsx.String("private key is required")
	}

	return Key{seed: seed, privatekey: privatekey}, nil
}

// ResolveKey fetches the seed from the backup service and pairs it with this device's
// identity.
func ResolveKey(ctx context.Context, c *http.Client) (Key, error) {
	seed, err := deeppool.NewBackups(c).Seed(ctx)
	if err != nil {
		return Key{}, errorsx.Wrap(err, "unable to resolve backup seed")
	}

	privatekey, err := os.ReadFile(env.PrivateKeyPath())
	if err != nil {
		return Key{}, errorsx.Wrap(err, "unable to read identity")
	}

	return NewKey(seed, privatekey)
}

// Encrypt is the archive cipher under a salt written ahead of the ciphertext, so every
// backup runs under its own keystream.
func (t Key) Encrypt(src io.Reader) (_ io.Reader, err error) {
	var (
		salt [saltSize]byte
	)

	if _, err = rand.Read(salt[:]); err != nil {
		return nil, errorsx.Wrap(err, "unable to generate salt")
	}

	r, err := cryptox.NewReaderChaCha20(t.chacha8(salt[:]), src)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to initialize cipher")
	}

	return io.MultiReader(bytes.NewReader(salt[:]), r), nil
}

// Decrypt reads the salt Encrypt wrote and returns the plaintext of what follows.
func (t Key) Decrypt(src io.Reader) (_ io.Reader, err error) {
	var (
		salt [saltSize]byte
	)

	if _, err = io.ReadFull(src, salt[:]); err != nil {
		return nil, errorsx.Wrap(err, "unable to read salt")
	}

	r, err := cryptox.NewReaderChaCha20(t.chacha8(salt[:]), src)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to initialize cipher")
	}

	return r, nil
}

func (t Key) chacha8(salt []byte) io.Reader {
	return cryptox.NewChaCha8(slices.Concat([]byte(t.seed), t.privatekey, salt))
}
