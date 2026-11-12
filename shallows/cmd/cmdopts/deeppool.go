package cmdopts

import (
	"context"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
)

// DeeppoolClient provides an http.Client for deeppool API calls.
type DeeppoolClient interface {
	HTTPClient(ctx context.Context) (*http.Client, error)
}

// DeeppoolClientDefault is the production implementation that registers
// with deeppool and creates an authenticated client.
type DeeppoolClientDefault struct {
	SSHID *SSHID
	// Global provides the account (--account-id) that selects the profile when the ssh key is associated
	// with multiple. a pointer because the flags are parsed after the client is constructed.
	Global *Global
}

func (t DeeppoolClientDefault) HTTPClient(ctx context.Context) (*http.Client, error) {
	signer, err := t.SSHID.Signer()
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to generate signer id")
	}

	session, err := authn.AutoRegistrationOnce(ctx, signer, t.Global.AccountID, false, ChooseProfile)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to register with archival service")
	}

	// act as the account the registration resolved.
	c, err := authn.AutoJWTClient(ctx, signer, session.Account.GetId())
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create api client")
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
