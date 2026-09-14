// Package metaapi provides HTTP/system-level concerns for the application
// as a whole rather than any single content domain: authentication,
// identity, authorization, billing, wireguard/network configuration,
// daemon management, and diagnostics. Domain-specific functionality (media,
// discovery, community, etc.) belongs in its own dedicated *api package
// instead.
package metaapi
