package release

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
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

// Keychain creates a temporary keychain and imports the provided .p12 certificate.
func Keychain(base64key, keypassword string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		keypath := egenv.WorkspaceDirectory("apple.p12")
		keychainPath := egenv.WorkspaceDirectory("apple.signing.keychain")
		intermediatepath := egenv.WorkspaceDirectory("apple.intermediate.cer")

		if err := writeBase64File(keypath, base64key, "certificate"); err != nil {
			return err
		}

		env := shell.Runtime().
			Environ("APPLE_KEYCHAIN_PASSWORD", egenv.RunID()).
			Environ("APPLE_SIGNING_KEY_PATH", keypath).
			Environ("APPLE_SIGNING_KEY_PASSWORD", keypassword).
			Environ("APPLE_INTERMEDIATE_CERT", intermediatepath).
			Environ("APPLE_KEYCHAIN_PATH", keychainPath)

		return shell.Run(
			ctx,
			env.New("which -a openssl").Lenient(true),
			env.New("openssl version").Lenient(true),
			env.New("security create-keychain -p ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"),
			env.New("security set-keychain-settings -lut 21600 ${APPLE_KEYCHAIN_PATH}"),
			env.New("security unlock-keychain -p ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"),
			env.New("security import ${APPLE_SIGNING_KEY_PATH} -P ${APPLE_SIGNING_KEY_PASSWORD} -A -k ${APPLE_KEYCHAIN_PATH}"),
			env.New("curl -fLo ${APPLE_INTERMEDIATE_CERT} $(openssl pkcs12 -in ${APPLE_SIGNING_KEY_PATH} -nokeys -passin pass:${APPLE_SIGNING_KEY_PASSWORD} | openssl x509 -noout -text | grep \"CA Issuers - URI:\" | cut -d':' -f2- | xargs)"),
			env.New("security import ${APPLE_INTERMEDIATE_CERT} -k ${APPLE_KEYCHAIN_PATH}"),
			env.New("security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k ${APPLE_KEYCHAIN_PASSWORD} ${APPLE_KEYCHAIN_PATH}"),
			env.New("security list-keychains -d user -s ${APPLE_KEYCHAIN_PATH} login.keychain-db"),
		)
	}
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

// AuthKey writes the App Store Connect .p8 auth key to ~/.private_keys/AuthKey_<ID>.p8.
func AuthKey(keyid, base64authkey string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		if strings.TrimSpace(keyid) == "" {
			return fmt.Errorf("appleAuthKey: missing key id")
		}

		filename := egenv.EphemeralDirectory("apple.auth.key.p8")

		if err := writeBase64File(filename, base64authkey, "auth key"); err != nil {
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

// EmbedProvisioningProfile decodes a provisioning profile and writes it to destpath.
func EmbedProvisioningProfile(base64profile, destpath string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		return writeBase64File(destpath, base64profile, "embedded provisioning profile")
	}
}

// ProvisioningProfile installs a mobileprovision profile for iOS builds.
func ProvisioningProfile(base64profile string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		profilepath := egenv.WorkspaceDirectory("apple.mobileprovision")

		if err := writeBase64File(profilepath, base64profile, "provisioning profile"); err != nil {
			return err
		}

		env := shell.Runtime().Environ("APPLE_PROFILE_PATH", profilepath)

		return shell.Run(
			ctx,
			env.New("mkdir -p ~/Library/MobileDevice/Provisioning\\ Profiles"),
			env.New("security cms -D -i ${APPLE_PROFILE_PATH} -o /tmp/apple_profile_decoded.plist && "+
				"UUID=$(/usr/libexec/PlistBuddy -c 'Print :UUID' /tmp/apple_profile_decoded.plist) && "+
				"cp ${APPLE_PROFILE_PATH} ~/Library/MobileDevice/Provisioning\\ Profiles/${UUID}.mobileprovision"),
			env.New("ls -lha ~/Library/MobileDevice/Provisioning\\ Profiles/"),
		)
	}
}
