package authn

import (
	"context"
	"net/http"
	"time"

	"github.com/retrovibed/retroapi/internal/env"
	"github.com/retrovibed/retroapi/internal/errorsx"
	"github.com/retrovibed/retroapi/internal/langx"
	"github.com/retrovibed/retroapi/internal/sshx"
	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
)

func JWTSSHTokenSource(endpoint oauth2.Endpoint, c *http.Client, opts ...SSHTokenSourceOption) *sshtokensource {
	return langx.Autoptr(langx.Clone(sshtokensource{
		endpoint: endpoint,
		c:        c,
	}, langx.Compose(opts...), SSHTokenSourceOptionAutoSigner))
}

type sshtokensource struct {
	endpoint oauth2.Endpoint
	c        *http.Client
	signer   ssh.Signer
}

type SSHTokenSourceOption func(*sshtokensource)

func SSHTokenSourceOptionSigner(v ssh.Signer) SSHTokenSourceOption {
	return func(s *sshtokensource) {
		s.signer = v
	}
}

func SSHTokenSourceOptionAutoSigner(s *sshtokensource) {
	if s.signer != nil {
		return
	}
	s.signer = errorsx.Must(sshx.AutoCached(sshx.NewKeyGen(), env.PrivateKeyPath()))
}

func (t *sshtokensource) Token() (*oauth2.Token, error) {
	state, err := AutoTokenState(t.signer)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to generate authentication state")
	}

	ctx, done := context.WithTimeout(context.Background(), time.Minute)
	defer done()

	ctx = context.WithValue(ctx, oauth2.HTTPClient, t.c)

	cfg := OAuth2SSHConfig(t.signer, "", t.endpoint)

	authzuri := cfg.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	)

	r, err := RetrieveAuthCode(ctx, t.c, authzuri)
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
