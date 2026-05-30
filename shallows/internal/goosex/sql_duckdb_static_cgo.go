//go:build duckdb_use_static_lib && linux

package goosex

// -Bstatic/-Bdynamic wrapping forces duckdb to link statically without affecting other libs.
// Done here rather than CGO_LDFLAGS env to avoid Go repeating the flags once per CGo module.

// #cgo LDFLAGS: -Wl,-Bstatic -lduckdb -Wl,-Bdynamic
import "C"
