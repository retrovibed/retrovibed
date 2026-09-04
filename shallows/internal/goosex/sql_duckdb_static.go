//go:build duckdb_use_static_lib

package goosex

const extensionSQL = "LOAD icu; LOAD inet; LOAD vss; LOAD httpfs;"
