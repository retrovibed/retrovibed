package cmdopts

import (
	"context"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

// DeeppoolClient provides an http.Client for deeppool API calls.
type DeeppoolClient interface {
	HTTPClient(ctx context.Context) (*http.Client, error)
}

// DeeppoolClientDefault is the production implementation that registers
// with deeppool and creates an authenticated client.
type DeeppoolClientDefault struct{}

func (t DeeppoolClientDefault) HTTPClient(ctx context.Context) (*http.Client, error) {
	if _, err := metaapi.Register(ctx); err != nil {
		return nil, errorsx.Wrap(err, "unable to register with archival service")
	}

	c, err := authn.AutoJWTClient(ctx)
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
