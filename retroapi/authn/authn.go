package authn

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/internal/debugx"
	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/httpx"
	"github.com/retrovibed/retrovibed/retroapi/internal/md5x"
	"github.com/retrovibed/retrovibed/retroapi/internal/oauth2x"
	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/retrovibed/retrovibed/retroapi/internal/stringsx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
)

func RetryClient(c *http.Client) *http.Client {
	return httpx.BindRetryTransport(c, http.StatusTooManyRequests, http.StatusBadGateway, http.StatusInternalServerError, http.StatusRequestTimeout)
}

func DeeppoolEndpoint() oauth2.Endpoint {
	return EndpointSSHAuth(fmt.Sprintf("https://%s", env.Deeppool()))
}

func NoRedirectFn(req *http.Request, via []*http.Request) error {
	if req.Method == "ENDPOINT" {
		return http.ErrUseLastResponse
	}
	return nil
}

func HTTPClientDefaults() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{ServerName: env.Deeppool(), InsecureSkipVerify: InsecureSkipVerify()},
		},
		CheckRedirect: NoRedirectFn,
	}
}

func HTTPClientLocalDefaults(cfg *tls.Config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       cfg,
		},
	}
}

func oauth2SSHConfig(signer ssh.Signer, otp string, endpoint oauth2.Endpoint) oauth2.Config {
	return oauth2.Config{
		ClientID:     ssh.FingerprintSHA256(signer.PublicKey()),
		ClientSecret: otp,
		Endpoint:     endpoint,
	}
}

func Seeded(ctx context.Context, seed string, force bool, path string) (s ssh.Signer, err error) {
	return sshx.Seeded(ctx, seed, force, path)
}

func Unseeded() error {
	return sshx.Unseeded(env.PrivateKeyPath(userx.DefaultRelRoot()))
}

func PrivateKeyPath() string {
	return env.PrivateKeyPath(userx.DefaultRelRoot())
}

func PublicKeyPath() string {
	return env.PrivateKeyPath(userx.DefaultRelRoot()) + ".pub"
}

func UserDisplayName() string {
	u := userx.CurrentUserOrDefault(userx.Zero())
	return stringsx.FirstNonBlank(u.Name, u.Username)
}

func Oauth2DeeppoolHTTPClient(ctx context.Context, signer ssh.Signer) (*http.Client, error) {
	cfg := oauth2SSHConfig(signer, "", DeeppoolEndpoint())

	c := HTTPClientDefaults()

	token, err := oauth2Bearer(ctx, signer, c, cfg, "", "")
	if err != nil {
		return nil, err
	}

	return cfg.Client(context.WithValue(ctx, oauth2.HTTPClient, c), token), nil
}

func SSHSigner() (ssh.Signer, error) {
	return sshx.AutoCached(sshx.NewKeyGen(), env.PrivateKeyPath(userx.DefaultRelRoot()))
}

func oauth2Bearer(ctx context.Context, signer ssh.Signer, c *http.Client, cfg oauth2.Config, email, displayname string) (*oauth2.Token, error) {
	type exstate struct {
		Entropy   string `json:"uid"`
		PublicKey []byte `json:"pkey"`
		Email     string `json:"email"`
		Display   string `json:"display"`
	}

	c = httpx.BindRetryTransport(c, http.StatusBadGateway, http.StatusTooManyRequests)

	state, err := jwtx.EncodeJSON(exstate{
		Entropy:   errorsx.Must(uuid.NewV4()).String(),
		PublicKey: signer.PublicKey().Marshal(),
		Email:     email,
		Display:   displayname,
	})
	if err != nil {
		return nil, err
	}

	authzuri := cfg.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	)

	exchanged, err := oauth2x.RetrieveAuthCode(ctx, c, authzuri)
	if err != nil {
		return nil, err
	}

	if exchanged.State != state {
		return nil, fmt.Errorf("unexpected state")
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, c)
	token, err := cfg.Exchange(ctx, exchanged.Code, oauth2.AccessTypeOffline)

	return token, errorsx.Wrap(err, "token signature failure")
}

func Oauth2Bearer(ctx context.Context, signer ssh.Signer, c *http.Client, email, displayname string) (*oauth2.Token, error) {
	return oauth2Bearer(ctx, signer, c, oauth2SSHConfig(signer, "", DeeppoolEndpoint()), email, displayname)
}

type FnTokenSource func() (*oauth2.Token, error)

func (t FnTokenSource) Token() (*oauth2.Token, error) {
	return t()
}

func AutomaticTokenSource(jwtsecret func() []byte) (*oauth2.Token, error) {
	signer, err := sshx.AutoCached(sshx.NewKeyGen(), env.PrivateKeyPath(userx.DefaultRelRoot()))
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to read identity")
	}

	id := ssh.FingerprintSHA256(signer.PublicKey())

	claims := jwtx.NewJWTClaims(
		id,
		jwtx.ClaimsOptionAuthnExpiration(),
		jwtx.ClaimsOptionIssuer(id),
	)

	debugx.Println("claims", spew.Sdump(claims))

	bearer, err := jwtx.Signed(jwtsecret(), claims)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create bearer")
	}

	return &oauth2.Token{
		TokenType:   "bearer",
		AccessToken: bearer,
		Expiry:      claims.ExpiresAt.Time,
	}, nil
}

func NewBearer(jwtsecret func() []byte) (string, error) {
	signer, err := sshx.AutoCached(sshx.NewKeyGen(), env.PrivateKeyPath(userx.DefaultRelRoot()))
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

	bearer, err := jwtx.Signed(jwtsecret(), claims)

	return bearer, errorsx.Wrap(err, "token signature failure")
}

func BearerForHost(ctx context.Context, c *http.Client, host string) (*oauth2.Token, error) {
	signer, err := sshx.AutoCached(sshx.NewKeyGen(), env.PrivateKeyPath(userx.DefaultRelRoot()))
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to read identity")
	}

	state, err := AutoTokenState(signer)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to generate authentication state")
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, c)

	endpoint := EndpointSSHAuth(host)

	cfg := OAuth2SSHConfig(signer, "", endpoint)

	authzuri := cfg.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	)

	r, err := RetrieveAuthCode(ctx, c, authzuri)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to retrieve auth code")
	}
	if r.State != state {
		return nil, errorsx.Wrap(err, "invalid state")
	}

	token, err := cfg.Exchange(ctx, r.Code, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to exchange auth code")
	}

	return token, nil
}

// special jwt secret for registration. essentially just the jwt secret hashed
// with the string registration. allows us to basically always return a valid jwt token
// for the oauth2 endpoint but only for the registration endpoints.
func JWTRegistrationSecretFromEnv(jwtsecret func() []byte) func() []byte {
	return func() []byte {
		return md5x.Digest(jwtsecret(), []byte("registration")).Sum(nil)
	}
}

// Bearer extracts the jwt bearer token from a http request.
func BearerToken(req *http.Request) string {
	before, after, _ := strings.Cut(req.Header.Get("authorization"), " ")

	if strings.ToLower(before) != "bearer" {
		return ""
	}

	return after
}

// AuthorizationToken retrieve and decode the authorization token.
func AuthorizationToken(ctx context.Context, secret jwtx.SecretSource, req *http.Request, claims jwt.Claims) (bearer string, err error) {
	// detect authorization header bearer token
	if bearer = BearerToken(req); len(bearer) > 0 {
		if err = jwtx.Validate(secret, bearer, claims); err == nil {
			return bearer, nil
		}
		log.Println("unable to detect jwt bearer token", err)
	}

	return bearer, errors.New("unable to detect session")
}

func EndpointSSHAuth(hostname string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:   fmt.Sprintf("%s/oauth2/ssh/auth", hostname),
		TokenURL:  fmt.Sprintf("%s/oauth2/ssh/token", hostname),
		AuthStyle: oauth2.AuthStyleInHeader,
	}
}

func OAuth2SSHConfig(signer ssh.Signer, otp string, endpoint oauth2.Endpoint) oauth2.Config {
	return oauth2.Config{
		ClientID:     ssh.FingerprintSHA256(signer.PublicKey()),
		ClientSecret: otp,
		Endpoint:     endpoint,
	}
}

func OAuth2SSHToken(ctx context.Context, signer ssh.Signer, endpoint oauth2.Endpoint) (cfg oauth2.Config, tok *oauth2.Token, err error) {
	var (
		sig *ssh.Signature
	)

	password := uuid.Must(uuid.NewV4())

	cfg = OAuth2SSHConfig(signer, password.String(), endpoint)
	if sig, err = signer.Sign(rand.Reader, password.Bytes()); err != nil {
		return cfg, nil, err
	}

	encodedsig := base64.RawURLEncoding.EncodeToString(ssh.Marshal(sig))

	tok, err = cfg.PasswordCredentialsToken(ctx, cfg.ClientID, encodedsig)
	return cfg, tok, err
}

func AutoTokenState(signer ssh.Signer) (encoded string, err error) {
	type reqstate struct {
		ID        string `json:"id"`
		PublicKey []byte `json:"pkey"`
	}

	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	rawstate := reqstate{
		ID:        id.String(),
		PublicKey: signer.PublicKey().Marshal(),
	}

	if encoded, err = jwtx.EncodeJSON(rawstate); err != nil {
		return "", errorsx.Wrap(err, "unable to encode state")
	}

	return encoded, nil
}

func AutoOauth2Client(ctx context.Context, cfg *tls.Config, endpoint oauth2.Endpoint, opts ...SSHTokenSourceOption) (c *http.Client) {
	_c := HTTPClientLocalDefaults(cfg)
	return oauth2.NewClient(context.WithValue(ctx, oauth2.HTTPClient, _c), JWTSSHTokenSource(endpoint, _c, opts...))
}

type AuthResponse struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func RetrieveAuthCode(ctx context.Context, chttp *http.Client, uri string) (r AuthResponse, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return r, err
	}

	resp, err := httpx.AsError(chttp.Do(req))
	if err != nil {
		return r, err
	}
	defer httpx.AutoClose(resp) //nolint:errcheck

	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return r, err
	}

	return r, nil
}
