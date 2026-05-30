//go:build duckdb_use_static_lib && (darwin || ios)

package goosex

// Static linking on darwin is handled via -Wl,-force_load in CGO_LDFLAGS at build time.

import "C"
