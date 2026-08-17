package cmdopts

import (
	"context"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// DeeppoolClient provides an http.Client for deeppool API calls.
type DeeppoolClient interface {
	HTTPClient(ctx context.Context) (*http.Client, error)
}

// DeeppoolClientDefault is the production implementation that registers
// with deeppool and creates an authenticated client.
type DeeppoolClientDefault struct {
	SSHID *SSHID
}

func (t DeeppoolClientDefault) HTTPClient(ctx context.Context) (*http.Client, error) {
	signer, err := t.SSHID.Signer()
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to generate signer id")
	}

	c, err := authn.AutoJWTClient(ctx, signer)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create api client")
	}

	if _, err := authn.Register(ctx, c); err != nil {
		return nil, errorsx.Wrap(err, "unable to register with archival service")
	}

	return c, nil
}

// DeeppoolClientTest is a test implementation that uses a provided http.Client.
type DeeppoolClientTest struct {
	Client *http.Client
}

func (t DeeppoolClientTest) HTTPClient(ctx context.Context) (*http.Client, error) {
	return t.Client, nil
}
