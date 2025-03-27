package authn

import (
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/retrovibed/retrovibed/internal/debugx"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/envx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/sshx"
	"golang.org/x/crypto/ssh"
)

func PublicKeyPath() string {
	return env.PrivateKeyPath() + ".pub"
}

func NewBearer() (string, error) {
	signer, err := sshx.AutoCached(sshx.NewKeyGen(), env.PrivateKeyPath())
	if err != nil {
		return "", errorsx.Wrap(err, "unable to read identity")
	}

	id := ssh.FingerprintSHA256(signer.PublicKey())

	claims := jwtx.NewJWTClaims(
		id,
		jwtx.ClaimsOptionAuthnExpiration(),
		jwtx.ClaimsOptionIssuer(id),
	)

	debugx.Println("claims", spew.Sdump(claims))

	bearer, err := jwt.NewWithClaims(
		jwtx.NewSSHSigner(),
		claims,
	).SignedString(signer)
	return bearer, errorsx.Wrap(err, "token signature failure")
}

func JWTSecretFromEnv() []byte {
	return []byte(envx.String(uuid.Must(uuid.NewV4()).String(), env.JWTSharedSecret))
}

// Bearer extracts the jwt bearer token from a http request.
func Bearer(req *http.Request) string {
	before, after, _ := strings.Cut(req.Header.Get("authorization"), " ")

	if strings.ToLower(before) != "bearer" {
		return ""
	}

	return after
}
