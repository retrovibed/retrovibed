package egapplexdep

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// writeBase64File decodes a base64-encoded secret and writes it to path with 0600 perms.
func writeBase64File(path, b64, label string) error {
	if strings.TrimSpace(b64) == "" {
		return fmt.Errorf("%s: missing base64 encoded value", label)
	}

	data, err := base64.URLEncoding.DecodeString(b64)
	if err != nil {
		return fmt.Errorf("failed to decode base64 %s: %w", label, err)
	}

	if err = os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write %s to disk: %w", label, err)
	}

	return nil
}

// KeychainPEM creates a temporary keychain and imports a PEM private key and DER certificate.
func KeychainPEM(base64key, base64cert string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		keypath := egenv.WorkspaceDirectory("apple.key.pem")
		certpath := egenv.WorkspaceDirectory("apple.cert.der")
		keychainPath := egenv.WorkspaceDirectory("apple.signing.keychain")
		intermediatepath := egenv.WorkspaceDirectory("apple.intermediate.cer")

		if err := writeBase64File(keypath, base64key, "key"); err != nil {
			return err
		}

		if err := writeBase64File(certpath, base64cert, "certificate"); err != nil {
			return err
		}

		env := shell.Runtime().
			Environ("APPLE_KEYCHAIN_PASSWORD", egenv.RunID()).
			Environ("APPLE_SIGNING_KEY_PATH", keypath).
			Environ("APPLE_SIGNING_CERT_PATH", certpath).
			Environ("APPLE_INTERMEDIATE_CERT", intermediatepath).
			Environ("APPLE_KEYCHAIN_PATH", keychainPath)

		return shell.Run(
			ctx,
			env.New("which -a openssl").Lenient(true),
			env.New("openssl version").Lenient(true),
			env.New("security create-keychain -p ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"),
			env.New("security set-keychain-settings -lut 21600 ${APPLE_KEYCHAIN_PATH}"),
			env.New("security unlock-keychain -p ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"),
			env.New("security import ${APPLE_SIGNING_KEY_PATH} -A -k ${APPLE_KEYCHAIN_PATH}"),
			env.New("security import ${APPLE_SIGNING_CERT_PATH} -A -k ${APPLE_KEYCHAIN_PATH}"),
			env.New("curl -fLo ${APPLE_INTERMEDIATE_CERT} $(openssl x509 -inform DER -in ${APPLE_SIGNING_CERT_PATH} -noout -text | grep \"CA Issuers - URI:\" | cut -d':' -f2- | xargs)"),
			env.New("security import ${APPLE_INTERMEDIATE_CERT} -k ${APPLE_KEYCHAIN_PATH}"),
			env.New("security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"),
			env.New("security list-keychains -d user -s ${APPLE_KEYCHAIN_PATH} login.keychain-db"),
		)
	}
}

// KeychainAppendPEM imports an additional PEM key and DER certificate into the existing signing keychain.
// label distinguishes temp files when multiple certificates are imported.
func KeychainAppendPEM(label, base64key, base64cert string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		keypath := egenv.WorkspaceDirectory(fmt.Sprintf("apple.%s.key.pem", label))
		certpath := egenv.WorkspaceDirectory(fmt.Sprintf("apple.%s.cert.der", label))
		keychainPath := egenv.WorkspaceDirectory("apple.signing.keychain")

		if err := writeBase64File(keypath, base64key, "installer key"); err != nil {
			return err
		}

		if err := writeBase64File(certpath, base64cert, "installer certificate"); err != nil {
			return err
		}

		env := shell.Runtime().
			Environ("APPLE_KEYCHAIN_PASSWORD", egenv.RunID()).
			Environ("APPLE_INSTALLER_KEY_PATH", keypath).
			Environ("APPLE_INSTALLER_CERT_PATH", certpath).
			Environ("APPLE_KEYCHAIN_PATH", keychainPath)

		return shell.Run(
			ctx,
			env.New("security unlock-keychain -p ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"),
			env.New("security import ${APPLE_INSTALLER_KEY_PATH} -A -k ${APPLE_KEYCHAIN_PATH}"),
			env.New("security import ${APPLE_INSTALLER_CERT_PATH} -A -k ${APPLE_KEYCHAIN_PATH}"),
			env.New("security set-key-partition-list -S apple-tool:,apple:,codesign:,productbuild:,productsign: -s -k ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"),
		)
	}
}
