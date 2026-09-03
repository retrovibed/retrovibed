package backups

import (
	"context"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// Key derives the database encryption key from the account's backup seed and the device's
// private key. the seed is issued once per account and the private key comes back from the
// identity seed, so a backup made on one device is recoverable on another.
func Key(seed string, privatekey []byte) (string, error) {
	var (
		derived [32]byte
	)

	if seed == "" {
		return "", errorsx.String("backup seed is required")
	}

	if len(privatekey) == 0 {
		return "", errorsx.String("private key is required")
	}

	prng := cryptox.NewChaCha8(append([]byte(seed), privatekey...))
	if _, err := prng.Read(derived[:]); err != nil {
		return "", errorsx.Wrap(err, "unable to derive backup key")
	}

	return hex.EncodeToString(derived[:]), nil
}

// ResolveKey fetches the seed from the backup service and derives the key for this device.
func ResolveKey(ctx context.Context, c *http.Client) (string, error) {
	seed, err := deeppool.NewBackups(c).Seed(ctx)
	if err != nil {
		return "", errorsx.Wrap(err, "unable to resolve backup seed")
	}

	privatekey, err := os.ReadFile(env.PrivateKeyPath())
	if err != nil {
		return "", errorsx.Wrap(err, "unable to read identity")
	}

	return Key(seed, privatekey)
}
