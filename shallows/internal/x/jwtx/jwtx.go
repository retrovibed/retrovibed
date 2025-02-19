package jwtx

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/stringsx"
	"google.golang.org/grpc/metadata"
)

const (
	mdkey = "authorization"
)

type JWTSecretSource func() []byte

func NotAuthorized() error {
	return errorsx.New("not authorized")
}

type Option func(*jwt.RegisteredClaims)

func ClaimsOptionAuthnRefreshExpiration() Option {
	return func(rc *jwt.RegisteredClaims) {
		rc.ExpiresAt = nil
	}
}

func ClaimsOptionAuthzExpiration() Option {
	return ClaimsOptionExpiration(time.Hour)
}

func ClaimsOptionAuthnExpiration() Option {
	return ClaimsOptionExpiration(7 * 24 * time.Hour)
}

func ClaimsOptionExpiration(d time.Duration) Option {
	return func(rc *jwt.RegisteredClaims) {
		rc.ExpiresAt = jwt.NewNumericDate(rc.IssuedAt.Time.Add(d))
	}
}

func ClaimsOptionIssuer(s string) Option {
	return func(rc *jwt.RegisteredClaims) {
		rc.Issuer = s
	}
}

func ClaimsOptionID(s string) Option {
	return func(rc *jwt.RegisteredClaims) {
		rc.ID = s
	}
}

func ClaimsOptionAudience(s ...string) Option {
	return func(rc *jwt.RegisteredClaims) {
		rc.Audience = s
	}
}

func ClaimsOptionComposed(opts ...Option) Option {
	return func(rc *jwt.RegisteredClaims) {
		for _, opt := range opts {
			opt(rc)
		}
	}
}

func NewJWTClaims(subject string, options ...Option) (c jwt.RegisteredClaims) {
	ts := time.Now()
	c = jwt.RegisteredClaims{
		ID:        uuid.Must(uuid.NewV4()).String(),
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(ts),
		NotBefore: jwt.NewNumericDate(ts),
		ExpiresAt: jwt.NewNumericDate(ts.Add(time.Minute)),
	}

	for _, opt := range options {
		opt(&c)
	}

	return c
}

// UnsafeEncodeLoginToken generates a login token using the provided source of entropy. generally code should not
// use this method.
func UnsafeSigned(ran io.Reader, jwtsecret []byte, t jwt.Claims) (signed string, err error) {
	if err = t.Valid(); err != nil {
		return signed, errorsx.Wrap(err, "invalid login token")
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS512, t).SignedString(jwtsecret)
}

// Signed generates a signed token
func Signed(jwtsecret []byte, t jwt.Claims) (signed string, err error) {
	// since we're passing the crypto rand to the unsafe method it becomes safe.
	return UnsafeSigned(rand.Reader, jwtsecret, t)
}

func BearerFromGRPCContext(ctx context.Context, jwtsecret JWTSecretSource, t jwt.Claims) (_ string, err error) {
	var (
		ok   bool
		md   metadata.MD
		vals []string
	)

	if md, ok = metadata.FromIncomingContext(ctx); !ok {
		return "", errorsx.Authorization(errorsx.New("missing authorization data"))
	}

	if vals = md.Get(mdkey); len(vals) != 1 {
		return "", errorsx.Authorization(errorsx.New("missing authorization data"))
	}

	_, err = Validate(jwtsecret, vals[0], t)
	return vals[0], err
}

func BearerFromHTTPContext(ctx context.Context, r *http.Request, jwtsecret JWTSecretSource, t jwt.Claims) (_ string, err error) {
	token, err := Validate(jwtsecret, r.Header.Get(mdkey), t)
	if err != nil {
		return "", err
	}
	return token.Raw, nil
}

func Validate(jwtsecret JWTSecretSource, encoded string, t jwt.Claims) (token *jwt.Token, err error) {
	encoded = strings.NewReplacer(
		"bearer ",
		"",
		"BEARER ",
		"",
		"Bearer ",
		"",
	).Replace(encoded)

	token, err = jwt.ParseWithClaims(string(encoded), t, func(token *jwt.Token) (interface{}, error) {
		return jwtsecret(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))

	if err != nil {
		return nil, errorsx.Wrap(err, "unable to parse jwt token")
	}

	if !token.Valid {
		return nil, errorsx.Errorf("invalid token %s", t)
	}

	return token, nil
}

func EncodeJSON(v any) (_ string, err error) {
	var (
		encoded []byte
	)

	if encoded, err = json.Marshal(v); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(encoded), nil
}

func DecodeJSON(s string, v any) (err error) {
	var (
		decoded []byte
	)

	if decoded, err = base64.URLEncoding.DecodeString(s); err != nil {
		return err
	}

	return json.Unmarshal(decoded, v)
}

type RedirectClaim string

func (t RedirectClaim) Valid() error {
	return nil
}

func SecureRedirect(uri *url.URL, s JWTSecretSource, fallback string, available ...string) (string, error) {
	uri.Host = enabledHosts(uri.Host, fallback, available...)
	redirect := RedirectClaim(uri.String())
	return Signed(s(), redirect)
}

func DecodeRedirect(s JWTSecretSource, encoded, fallback string) string {
	var (
		decoded RedirectClaim
	)

	if stringsx.Blank(encoded) {
		return fallback
	}

	if _, err := Validate(s, encoded, &decoded); err != nil {
		log.Println("redirect invalid using fallback", err)
		return fallback
	}

	return string(decoded)
}

func enabledHosts(s, fallback string, available ...string) string {
	for _, allowed := range available {
		if allowed == s {
			return allowed
		}
	}

	return fallback
}

// EncodeSegment encodes a JWT specific base64url encoding with padding stripped
func EncodeSegment(seg []byte) string {
	return base64.RawURLEncoding.EncodeToString(seg)
}

// DecodeSegment decodes a JWT specific base64url encoding with padding stripped
func DecodeSegment(seg string) ([]byte, error) {
	encoding := base64.RawURLEncoding

	if jwt.DecodePaddingAllowed {
		if l := len(seg) % 4; l > 0 {
			seg += strings.Repeat("=", 4-l)
		}
		encoding = base64.URLEncoding
	}

	if jwt.DecodeStrict {
		encoding = encoding.Strict()
	}
	return encoding.DecodeString(seg)
}
