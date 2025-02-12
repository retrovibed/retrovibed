package md5x

import (
	"crypto/md5"
	"encoding/hex"
)

// Digest to md5 hex encoded string
func Digest(b []byte) string {
	d := md5.Sum(b)
	return hex.EncodeToString(d[:])
}

// Hex to md5 hex encoded string
func Hex(s string) string {
	return Digest([]byte(s))
}

// Bytes digest byte slice
func Bytes[T ~[]byte | string](b T) []byte {
	d := md5.Sum([]byte(b))
	return d[:]
}
