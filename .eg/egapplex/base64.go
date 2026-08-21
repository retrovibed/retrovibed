package egapplex

import (
	"encoding/base64"
	"log"
	"os"
	"strings"
)

// Base64 reads the first of keys set in the environment and decodes it as URL-safe base64
// (matching how these secrets are actually encoded in CI), returning fallback if none are set or
// none decode successfully.
func Base64(fallback []byte, keys ...string) []byte {
	for _, k := range keys {
		s := strings.TrimSpace(os.Getenv(k))
		if s == "" {
			continue
		}

		decoded, err := base64.URLEncoding.DecodeString(s)
		if err != nil {
			log.Printf("%s stored an invalid base64 value: %v", k, err)
			continue
		}

		return decoded
	}

	return fallback
}
