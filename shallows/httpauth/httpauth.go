package httpauth

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
)

// Authenticate a session - responds with 401 if unable to locate the token.
func AuthenticateWithToken(p jwtx.SecretSource) func(http.Handler) http.Handler {
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
		})
	}
}

func IssuerSubjectID(ctx context.Context, ss jwtx.SecretSource, req *http.Request) (issuer string, pid string, err error) {
	var (
		b jwt.RegisteredClaims
	)

	if _, err = jwtx.BearerFromHTTPContext(ctx, req, ss, &b); err != nil {
		return "", "", err
	}

	return b.Issuer, b.Subject, nil
}
