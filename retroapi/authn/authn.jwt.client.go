package authn

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/httpx"
	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
)

// AutoJWTClient returns a client authorized against /m/* endpoints via the
// authz layer (requires the identity to already have a profile). Never pass
// its result to Register — Register creates the first profile, but this
// client's authz layer requires one to already exist to fetch a bearer at
// all, including for Register's own /authn/ssh request. Use
// RegistrationJWTClient for that instead.
func AutoJWTClient(ctx context.Context, signer ssh.Signer) (c *http.Client, err error) {
	c, err = Oauth2DeeppoolHTTPClient(ctx, signer)
	if err != nil {
		return nil, errorsx.Wrap(err, "failed to create oauth2 http client")
	}

	return RetryClient(AuthzClient(JWTClientHostname(c, env.Deeppool()))), nil
}

// RegistrationJWTClient returns the client Register must be called with: the
// layer-1 SSH-signed JWT client only, with no dependency on an existing
// profile. See AutoJWTClient's doc comment for why that client can't be used
// here.
func RegistrationJWTClient(ctx context.Context, signer ssh.Signer) (*http.Client, error) {
	return RegistrationJWTClientWithEndpoint(ctx, signer, DeeppoolEndpoint())
}

// RegistrationJWTClientWithEndpoint is RegistrationJWTClient with an explicit
// endpoint instead of env.Deeppool(), for tests driving a fake server.
func RegistrationJWTClientWithEndpoint(ctx context.Context, signer ssh.Signer, endpoint oauth2.Endpoint) (*http.Client, error) {
	return Oauth2DeeppoolHTTPClientWithEndpoint(ctx, signer, endpoint)
}

func JWTClientHostname(oauth2c *http.Client, hostname string) *http.Client {
	return oauth2.NewClient(
		context.WithValue(context.Background(), oauth2.HTTPClient, HTTPClientDefaults()),
		&jwttokensource{
			oauth2c:  oauth2c,
			endpoint: fmt.Sprintf("https://%s", hostname),
		},
	)
}

type jwttokensource struct {
	oauth2c  *http.Client
	endpoint string
}

func (t *jwttokensource) Token() (*oauth2.Token, error) {
	var (
		authed Authed
	)

	ctx, done := context.WithTimeout(context.Background(), 3*time.Second)
	defer done()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/authn/ssh", t.endpoint), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpx.AsError(t.oauth2c.Do(req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&authed); err != nil {
		return nil, err
	}

	switch len(authed.Profiles) {
	case 0:
		return nil, fmt.Errorf("no profiles associated with this identity. this should have been done automatically during boot")
	case 1:
	default:
		return nil, fmt.Errorf("mumlitple profiles per identity not supported")
	}

	return &oauth2.Token{AccessToken: authed.Profiles[0].Token}, err
}

func AuthzClient(oauth2c *http.Client) *http.Client {
	return AuthzClientEndpoint(context.WithValue(context.Background(), oauth2.HTTPClient, HTTPClientDefaults()), oauth2c, fmt.Sprintf("https://%s/m/authz/", env.Deeppool()))
}

func AuthzClientLibrary(tls *tls.Config, oauth2c *http.Client, endpoint string) *http.Client {
	cc := HTTPClientLocalDefaults(tls)
	return AuthzClientEndpoint(context.WithValue(context.Background(), oauth2.HTTPClient, cc), oauth2c, fmt.Sprintf("%s/meta/authz/", endpoint))
}

func AuthzClientEndpoint(ctx context.Context, oauth2c *http.Client, endpoint string) *http.Client {
	return oauth2.NewClient(
		ctx,
		&metatokensource{
			oauth2c:  oauth2c,
			endpoint: endpoint,
		},
	)
}

type metatokensource struct {
	oauth2c  *http.Client
	endpoint string
}

func (t *metatokensource) Token() (*oauth2.Token, error) {
	var (
		authed AuthzResponse
	)

	ctx, done := context.WithTimeout(context.Background(), 3*time.Second)
	defer done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpx.AsError(t.oauth2c.Do(req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&authed); err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken: authed.Bearer,
		Expiry:      time.UnixMilli(authed.Token.Expires * 1000),
		ExpiresIn:   authed.Token.Expires,
	}, err
}
