// Package egapplex provides shared Apple code-signing and release tooling used by both the
// darwin and ios release flows.
//
// It covers: creating and populating a signing keychain (KeychainP12, KeychainPEM,
// KeychainAppendPEM, UnlockKeychain), installing an App Store Connect API auth key (AuthKey),
// installing a provisioning profile (Provision), codesigning (Sign), and submitting a build to
// App Store Connect for upload or notarization (Upload, Notarize).
//
// Every function decodes secret material as already-decoded []byte rather than base64 strings —
// callers extract secrets from the environment via egenv.Base64/egenv.String and decide which
// function's encoding (PEM, DER, .p12) applies; egapplex itself is format-agnostic about what an
// env var is named.
package egapplex