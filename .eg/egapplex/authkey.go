package egapplex

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// AuthKey writes the App Store Connect .p8 auth key to ~/.private_keys/AuthKey_<keyid>.p8.
func AuthKey(keyid string, key []byte) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		if strings.TrimSpace(keyid) == "" {
			return fmt.Errorf("appleAuthKey: missing key id")
		}

		filename := egenv.EphemeralDirectory("apple.auth.key.p8")

		if err := writeFile(filename, key, "auth key"); err != nil {
			return err
		}

		log.Printf("Writing Auth Key to local sandbox: %s", filename)

		return shell.Run(
			ctx,
			shell.Newf("mkdir -p ~/.private_keys"),
			shell.Newf("mv ${APPLE_AUTH_KEY_PATH} ~/.private_keys/AuthKey_${APPLE_API_KEY_ID}.p8").
				Environ("APPLE_API_KEY_ID", keyid).
				Environ("APPLE_AUTH_KEY_PATH", filename),
		)
	}
}
