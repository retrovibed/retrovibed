// Package egapplexdep holds Apple keychain functions that still need a private key per
// certificate rather than one shared key across a keychain (see egapplex.KeychainAppendPEM,
// which assumes a single shared key — an assumption that doesn't hold for these particular CI
// secrets: the installer and appstore certs are each issued from their own CSR with their own
// distinct key).
//
// This is a staging package: a home for the still-needed per-cert-key variants while
// egapplex.KeychainPEM/KeychainAppendPEM are redesigned to support that case too. Once that
// redesign lands, this package's callers should migrate back to egapplex and this package
// removed.
package egapplexdep
