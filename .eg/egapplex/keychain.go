package egapplex

import (
	"context"
	"fmt"
	"os"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// writeFile writes already-decoded secret material to path with 0600 perms.
func writeFile(path string, data []byte, label string) error {
	if len(data) == 0 {
		return fmt.Errorf("%s: missing value", label)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write %s to disk: %w", label, err)
	}

	return nil
}

// keychainSetup returns the shared prefix shell steps for creating a fresh signing keychain:
// create, set the auto-lock timeout, then unlock it.
func keychainSetup(env shell.Command) []shell.Command {
	return []shell.Command{
		env.New("which -a openssl").Lenient(true),
		env.New("openssl version").Lenient(true),
		env.New("security create-keychain -p ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"),
		env.New("security set-keychain-settings -lut 21600 ${APPLE_KEYCHAIN_PATH}"),
		env.New("security unlock-keychain -p ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"),
	}
}

// keychainFinalize returns the shared suffix shell steps for trusting the derived intermediate
// certificate, scoping the key partition list, and registering the keychain for lookup.
func keychainFinalize(env shell.Command, partitions string) []shell.Command {
	return []shell.Command{
		env.New("security import ${APPLE_INTERMEDIATE_CERT} -k ${APPLE_KEYCHAIN_PATH}"),
		env.New(fmt.Sprintf("security set-key-partition-list -S %s -s -k ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}", partitions)),
		env.New("security list-keychains -d user -s ${APPLE_KEYCHAIN_PATH} login.keychain-db"),
	}
}

// UnlockKeychain unlocks the keychain at keychainPath using the current run's password.
func UnlockKeychain(keychainPath string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		return shell.Run(
			ctx,
			shell.Newf("security unlock-keychain -p %s %s", egenv.RunID(), keychainPath),
		)
	}
}

// KeychainP12 creates a temporary keychain and imports the provided .p12 certificate.
func KeychainP12(key []byte, password string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		keypath := egenv.WorkspaceDirectory("apple.p12")
		keychainPath := egenv.WorkspaceDirectory("apple.signing.keychain")
		intermediatepath := egenv.WorkspaceDirectory("apple.intermediate.cer")

		if err := writeFile(keypath, key, "certificate"); err != nil {
			return err
		}

		env := shell.Runtime().
			Environ("APPLE_KEYCHAIN_PASSWORD", egenv.RunID()).
			Environ("APPLE_SIGNING_KEY_PATH", keypath).
			Environ("APPLE_SIGNING_KEY_PASSWORD", password).
			Environ("APPLE_INTERMEDIATE_CERT", intermediatepath).
			Environ("APPLE_KEYCHAIN_PATH", keychainPath)

		cmds := keychainSetup(env)
		cmds = append(
			cmds,
			env.New("security import ${APPLE_SIGNING_KEY_PATH} -P ${APPLE_SIGNING_KEY_PASSWORD} -A -k ${APPLE_KEYCHAIN_PATH}"),
			env.New("curl -fLo ${APPLE_INTERMEDIATE_CERT} $(openssl pkcs12 -legacy -in ${APPLE_SIGNING_KEY_PATH} -nokeys -passin pass:${APPLE_SIGNING_KEY_PASSWORD} | openssl x509 -noout -text | grep \"CA Issuers - URI:\" | cut -d':' -f2- | xargs)"),
		)
		cmds = append(cmds, keychainFinalize(env, "apple-tool:,apple:,codesign:")...)

		return shell.Run(ctx, cmds...)
	}
}

// KeychainPEM creates a temporary keychain and imports a PEM private key and DER certificate.
func KeychainPEM(key, cert []byte) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		keypath := egenv.WorkspaceDirectory("apple.key.pem")
		certpath := egenv.WorkspaceDirectory("apple.cert.der")
		keychainPath := egenv.WorkspaceDirectory("apple.signing.keychain")
		intermediatepath := egenv.WorkspaceDirectory("apple.intermediate.cer")

		if err := writeFile(keypath, key, "key"); err != nil {
			return err
		}

		if err := writeFile(certpath, cert, "certificate"); err != nil {
			return err
		}

		env := shell.Runtime().
			Environ("APPLE_KEYCHAIN_PASSWORD", egenv.RunID()).
			Environ("APPLE_SIGNING_KEY_PATH", keypath).
			Environ("APPLE_SIGNING_CERT_PATH", certpath).
			Environ("APPLE_INTERMEDIATE_CERT", intermediatepath).
			Environ("APPLE_KEYCHAIN_PATH", keychainPath)

		cmds := keychainSetup(env)
		cmds = append(
			cmds,
			env.New("security import ${APPLE_SIGNING_KEY_PATH} -A -k ${APPLE_KEYCHAIN_PATH}"),
			env.New("security import ${APPLE_SIGNING_CERT_PATH} -A -k ${APPLE_KEYCHAIN_PATH}"),
			env.New("curl -fLo ${APPLE_INTERMEDIATE_CERT} $(openssl x509 -inform DER -in ${APPLE_SIGNING_CERT_PATH} -noout -text | grep \"CA Issuers - URI:\" | cut -d':' -f2- | xargs)"),
		)
		cmds = append(cmds, keychainFinalize(env, "apple-tool:,apple:,codesign:")...)

		return shell.Run(ctx, cmds...)
	}
}

// KeychainAppendPEM imports one or more additional certificates into the existing signing
// keychain under a single unlock/partition-list pass, pairing each against the private key
// already imported by an earlier KeychainPEM call. label distinguishes temp files when multiple
// certificates are imported under one unlock/partition-list pass; passing more than one cert is
// for bundling a chain (e.g. leaf + intermediate) for a single identity.
func KeychainAppendPEM(certs ...[]byte) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		keychainPath := egenv.WorkspaceDirectory("apple.signing.keychain")

		env := shell.Runtime().
			Environ("APPLE_KEYCHAIN_PASSWORD", egenv.RunID()).
			Environ("APPLE_KEYCHAIN_PATH", keychainPath)

		cmds := []shell.Command{
			env.New("security unlock-keychain -p ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"),
		}

		for _, cert := range certs {
			certfile, err := os.CreateTemp(egenv.WorkspaceDirectory(), "apple.*.cert.der")
			if err != nil {
				return fmt.Errorf("failed to create temp cert file: %w", err)
			}
			certpath := certfile.Name()
			certfile.Close()

			if err := writeFile(certpath, cert, "certificate"); err != nil {
				return err
			}

			cmds = append(cmds, env.New("security import "+certpath+" -A -k ${APPLE_KEYCHAIN_PATH}"))
		}

		cmds = append(cmds, env.New("security set-key-partition-list -S apple-tool:,apple:,codesign:,productbuild:,productsign: -s -k ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"))

		return shell.Run(ctx, cmds...)
	}
}
