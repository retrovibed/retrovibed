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

// Keychain creates a temporary keychain and imports the provided .p12 certificate.
func Keychain(base64key, keypassword string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		if strings.TrimSpace(base64key) == "" {
			return fmt.Errorf("setting up a temporary keychain requires a base64 encoded certificate")
		}

		keypath := egenv.WorkspaceDirectory("apple.p12")
		keychainPath := egenv.WorkspaceDirectory("apple.signing.keychain")
		intermediatepath := egenv.WorkspaceDirectory("apple.intermediate.cer")

		keyp12, err := base64.URLEncoding.DecodeString(base64key)
		if err != nil {
			return fmt.Errorf("failed to decode base64 certificate: %w", err)
		}

		if err = os.WriteFile(keypath, keyp12, 0600); err != nil {
			return fmt.Errorf("failed to write certificate to disk: %w", err)
		}

		env := shell.Runtime().
			Environ("APPLE_KEYCHAIN_PASSWORD", egenv.RunID()).
			Environ("APPLE_SIGNING_KEY_PATH", keypath).
			Environ("APPLE_SIGNING_KEY_PASSWORD", keypassword).
			Environ("APPLE_INTERMEDIATE_CERT", intermediatepath).
			Environ("APPLE_KEYCHAIN_PATH", keychainPath)

		return shell.Run(
			ctx,
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
		if strings.TrimSpace(base64key) == "" || strings.TrimSpace(base64cert) == "" {
			return fmt.Errorf("setting up a temporary keychain requires a base64 encoded key and certificate")
		}

		keypath := egenv.WorkspaceDirectory("apple.key.pem")
		certpath := egenv.WorkspaceDirectory("apple.cert.der")
		keychainPath := egenv.WorkspaceDirectory("apple.signing.keychain")
		intermediatepath := egenv.WorkspaceDirectory("apple.intermediate.cer")

		keydata, err := base64.URLEncoding.DecodeString(base64key)
		if err != nil {
			return fmt.Errorf("failed to decode base64 key: %w", err)
		}

		if err = os.WriteFile(keypath, keydata, 0600); err != nil {
			return fmt.Errorf("failed to write key to disk: %w", err)
		}

		certdata, err := base64.URLEncoding.DecodeString(base64cert)
		if err != nil {
			return fmt.Errorf("failed to decode base64 certificate: %w", err)
		}

		if err = os.WriteFile(certpath, certdata, 0600); err != nil {
			return fmt.Errorf("failed to write certificate to disk: %w", err)
		}

		env := shell.Runtime().
			Environ("APPLE_KEYCHAIN_PASSWORD", egenv.RunID()).
			Environ("APPLE_SIGNING_KEY_PATH", keypath).
			Environ("APPLE_SIGNING_CERT_PATH", certpath).
			Environ("APPLE_INTERMEDIATE_CERT", intermediatepath).
			Environ("APPLE_KEYCHAIN_PATH", keychainPath)

		return shell.Run(
			ctx,
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

// AuthKey writes the App Store Connect .p8 auth key to ~/.private_keys/AuthKey_<ID>.p8.
func AuthKey(keyid, base64authkey string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		if strings.TrimSpace(base64authkey) == "" || strings.TrimSpace(keyid) == "" {
			return fmt.Errorf("appleAuthKey: missing credentials")
		}

		filename := egenv.EphemeralDirectory("apple.auth.key.p8")

		authkeydata, err := base64.URLEncoding.DecodeString(base64authkey)
		if err != nil {
			return fmt.Errorf("failed to decode base64 auth key: %w", err)
		}

		log.Printf("Writing Auth Key to local sandbox: %s", filename)
		if err = os.WriteFile(filename, authkeydata, 0600); err != nil {
			return fmt.Errorf("failed to write auth key: %w", err)
		}

		return shell.Run(
			ctx,
			shell.Newf("mkdir -p ~/.private_keys"),
			shell.Newf("mv ${APPLE_AUTH_KEY_PATH} ~/.private_keys/AuthKey_${APPLE_API_KEY_ID}.p8").
				Environ("APPLE_API_KEY_ID", keyid).
				Environ("APPLE_AUTH_KEY_PATH", filename),
		)
	}
}

// ProvisioningProfile installs a mobileprovision profile for iOS builds.
func ProvisioningProfile(base64profile string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		if strings.TrimSpace(base64profile) == "" {
			return fmt.Errorf("provisioning profile: missing base64 encoded profile")
		}

		profilepath := egenv.WorkspaceDirectory("apple.mobileprovision")

		profiledata, err := base64.URLEncoding.DecodeString(base64profile)
		if err != nil {
			return fmt.Errorf("failed to decode base64 profile: %w", err)
		}

		if err = os.WriteFile(profilepath, profiledata, 0600); err != nil {
			return fmt.Errorf("failed to write profile to disk: %w", err)
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
