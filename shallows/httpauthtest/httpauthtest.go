package httpauthtest

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
)

// only utilized in testing, must not be used anywhere else
func UnsafeJWTSecretSource() []byte {
	return []byte("unsafe")
}

func UnsafeToken(c jwt.Claims, s jwtx.SecretSource) string {
	return errorsx.Must(jwtx.Signed(s(), c))
}

func UnsafeClaimsToken(c jwt.Claims, s jwtx.SecretSource) string {
	return fmt.Sprintf("BEARER %s", UnsafeToken(c, s))
}
