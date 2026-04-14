package httpauth

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/retrovibed/retrovibed/shallows/authn"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jwtx"
)

// Authenticate a session - responds with 401 if unable to locate the token.
func AuthenticateWithToken(p jwtx.SecretSource) func(http.Handler) http.Handler {
	return func(original http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			var (
				err     error
				b       jwt.RegisteredClaims
				encoded = authn.Bearer(req)
			)

			if err = jwtx.Validate(p, encoded, &b); err != nil {
				httpx.Unauthorized(resp, errorsx.Wrapf(err, "failed to verify token: %s", encoded))
				return
			}

			original.ServeHTTP(resp, req)
		})
	}
}

func BearerFromQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tok := r.URL.Query().Get("token"); tok != "" && r.Header.Get("Authorization") == "" {
			r.Header.Set("Authorization", "bearer "+tok)
		}
		next.ServeHTTP(w, r)
	})
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
