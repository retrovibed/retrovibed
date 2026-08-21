package egapplex

const (
	// The single private key backing every code-signing cert across both platforms (darwin's
	// Developer ID, Installer, and Distribution certs, and iOS's Distribution identity) —
	// Apple's CSR-based issuance allows one key to back multiple certificate types, and this
	// module treats that as the standard case (see KeychainAppendPEM). darwin's KeychainPEM
	// reads it as a raw PEM key (EnvDistributionPassword unused); iOS's KeychainP12 reads it
	// as a .p12 blob together with EnvDistributionPassword. One key, two consumption shapes.
	EnvDistributionKey      = "APPLE_DISTRIBUTION_KEY"
	EnvDistributionPassword = "APPLE_DISTRIBUTION_PASSWORD"

	// Certificates — one per distribution channel, each importable independently against
	// EnvDistributionKey (darwin) or bundled into it (iOS's .p12 already contains its cert).
	EnvDeveloperIDCert  = "APPLE_DEVELOPER_ID_CERT"  // darwin-only, direct/notarized distribution
	EnvInstallerCert    = "APPLE_INSTALLER_CERT"     // darwin-only, signs the .pkg installer
	EnvDistributionCert = "APPLE_DISTRIBUTION_CERT"  // App Store submission, darwin's half

	// App Store Connect API key — authenticates notarytool/altool uploads. Names match the
	// AuthKey function's parameters (keyid, key).
	EnvAuthKeyID     = "APPLE_AUTHKEY_ID"
	EnvAuthKeyIssuer = "APPLE_AUTHKEY_ISSUER"
	EnvAuthKey       = "APPLE_AUTHKEY"
)
