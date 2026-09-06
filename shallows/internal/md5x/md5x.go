package md5x

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"hash"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// digest the provided contents and return the resulting hash.
// if an error occurs during hashing then a nil value is returned.
func Digest[T string | ~[]byte](bs ...T) hash.Hash {
	v := md5.New()

	for _, b := range bs {
		y := []byte(b)
		if n, err := v.Write(y); err != nil || n < len(y) {
			return nil
		}
	}

	return v
}

// deliberately still on encoding/json rather than jsonx: these digests content
// address records that are already stored, and json/v2 marshals differently
// (no html escaping, a nil slice as [] instead of null), which would silently
// rekey every one of them.
func JSON(v any) hash.Hash {
	return Digest(errorsx.Must(json.Marshal(v)))
}

// String to md5 uuid encoded string
func String(s string) string {
	return FormatUUID(Digest(s))
}

// format md5 hash to a uuid encoded string
func FormatUUID(m hash.Hash) string {
	return uuid.FromBytesOrNil(m.Sum(nil)).String()
}

// format md5 hash to a hex encoded string
func FormatHex(m hash.Hash) string {
	return hex.EncodeToString(m.Sum(nil))
}

// format hash to a base64 encoded string
func FormatBase64(m hash.Hash) string {
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}
