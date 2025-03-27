package httpauth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
)

// Authenticate a session - responds with 401 if unable to locate the token.
func AuthenticateWithToken(p jwtx.JWTSecretSource) func(http.Handler) http.Handler {
	return func(original http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			var (
				err error
				b   jwt.RegisteredClaims
			)

			if err = jwtx.Validate(p, authn.Bearer(req), &b); err != nil {
				httpx.Unauthorized(resp, errorsx.Wrap(err, "failed to decode token"))
				return
			}

			original.ServeHTTP(resp, req)
			// original.ServeHTTP(resp, req.WithContext(authn.WithLoginToken(req.Context(), b)))
		})
	}
}
