//go:build !duckdb_use_static_lib

package goosex

const inetSQL = "INSTALL inet; LOAD inet;"
const vssSQL = "INSTALL vss; LOAD vss;"
