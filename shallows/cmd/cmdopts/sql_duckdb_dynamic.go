//go:build !duckdb_use_static_lib

package cmdopts

const extensionSQL = "INSTALL icu; LOAD icu; INSTALL inet; LOAD inet; INSTALL vss; LOAD vss;"
