//go:build duckdb_use_static_lib

package goosex

// -Bstatic/-Bdynamic wrapping forces duckdb to link statically without affecting other libs.
// Done here rather than CGO_LDFLAGS env to avoid Go repeating the flags once per CGo module.
// Darwin/iOS static linking is handled via -Wl,-force_load in CGO_LDFLAGS at build time.

// #cgo linux LDFLAGS: -Wl,-Bstatic -lduckdb -Wl,-Bdynamic
// #cgo !linux,!darwin,!ios LDFLAGS: -lduckdb
import "C"
