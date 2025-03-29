package metaapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang-jwt/jwt/v5"
	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/meta"
	"google.golang.org/protobuf/encoding/protojson"
)

type contextKey int

const (
	ckeyAuthz contextKey = iota
)

func FromContext(ctx context.Context) (t *Token, err error) {
	var (
		ok bool
	)

	if t, ok = ctx.Value(ckeyAuthz).(*Token); ok {
		return t, nil
	}

	return nil, jwtx.NotAuthorized()
}

func WithAuthorization(ctx context.Context, b *Token) context.Context {
	return context.WithValue(ctx, ckeyAuthz, b)
}

func AuthzPermUsermanagement(ctx context.Context, cause error) (_ context.Context, token *Token, err error) {
	if cause != nil {
		return ctx, nil, errorsx.Authorization(fmt.Errorf("not authorized"))
	}

	if token, err = FromContext(ctx); err != nil {
		return ctx, token, errorsx.Wrap(errorsx.Authorization(err), "not authorized")
	} else if !token.Usermanagement {
		return ctx, token, errorsx.Authorization(errorsx.WithStack(fmt.Errorf("not authorized: permission denied")))
	}

	return ctx, token, nil
}

func AuthzTokenHTTP(p jwtx.SecretSource, check func(ctx context.Context, cause error) (_ context.Context, token *Token, err error)) func(http.Handler) http.Handler {
	return func(original http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx    context.Context
				err    error
				claims tokenclaims = tokenclaims{Token: &Token{}}
			)

			if err = jwtx.Validate(p, authn.Bearer(r), &claims); err != nil {
				httpx.Unauthorized(w, errorsx.Wrap(err, "invalid token"))
				return
			}

			if ctx, _, err = check(WithAuthorization(r.Context(), claims.Token), errorsx.Wrap(err, "failed to decode token")); err != nil {
				httpx.Unauthorized(w, errorsx.Wrapf(err, "claims: %s", spew.Sdump(&claims)))
				return
			}

			original.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type tokenclaims struct {
	*Token
}

func (t *tokenclaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(t.Expires, 0)), nil
}

func (t *tokenclaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(t.Issued, 0)), nil
}

func (t *tokenclaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(t.NotBefore, 0)), nil
}

func (t *tokenclaims) GetIssuer() (string, error) {
	return "", nil
}

func (t *tokenclaims) GetSubject() (string, error) {
	return t.ProfileId, nil
}
func (t *tokenclaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

func (t *Token) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *Token) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

type TokenOption func(*Token)

func TokenOptionFromAuthz(a meta.Authz) TokenOption {
	return func(t *Token) {
		t.Usermanagement = a.Usermanagement
	}
}

func TokenFromRegisterClaims(claims jwt.RegisteredClaims, options ...TokenOption) *Token {
	return langx.Autoptr(langx.Clone(Token{
		Id:        claims.ID,
		Issuer:    claims.Issuer,
		ProfileId: claims.Subject,
		Issued:    claims.IssuedAt.Unix(),
		Expires:   claims.ExpiresAt.Unix(),
		NotBefore: claims.NotBefore.Unix(),
	}, options...))
}

func NewJWTClaim(t *Token) *tokenclaims {
	return &tokenclaims{Token: t}
}
