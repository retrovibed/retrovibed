//go:build duckdb_use_static_lib && !linux && !darwin && !ios

package goosex

// #cgo LDFLAGS: -lduckdb
import "C"
