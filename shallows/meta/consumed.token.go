package meta

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/retrovibed/retrovibed/internal/md5x"
)

func ConsumedTokenFromJWTClaims(raw string, c jwt.RegisteredClaims) ConsumedToken {
	return ConsumedToken{
		ID:           md5x.String(raw),
		TombstonedAt: c.ExpiresAt.Time.Add(time.Minute),
		Token:        raw,
	}
}

func ConsumedTokenString(raw string, tombstone time.Time) ConsumedToken {
	return ConsumedToken{
		ID:           md5x.String(raw),
		TombstonedAt: tombstone,
		Token:        raw,
	}
}
