package envx

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
)

// Base64 read value as a base64 encoded string
func Base64(fallback []byte, keys ...string) []byte {
	enc := base64.StdEncoding.WithPadding('=')
	return envval(fallback, os.Getenv, func(s string) ([]byte, error) {
		decoded, err := enc.DecodeString(s)
		if err != nil {
			return decoded, fmt.Errorf("invalid base64 encoded data '%s' - %w", s, err)
		}
		return decoded, nil
	}, keys...)
}

func envval[T any](fallback T, m func(string) string, parse func(string) (T, error), keys ...string) T {
	for _, k := range keys {
		s := strings.TrimSpace(m(k))
		if s == "" {
			continue
		}

		decoded, err := parse(s)
		if err != nil {
			log.Printf("%s stored an invalid value %v\n", k, err)
			continue
		}

		return decoded
	}

	return fallback
}
