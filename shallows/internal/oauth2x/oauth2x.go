package oauth2x

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/retrovibed/retrovibed/internal/httpx"
	"golang.org/x/oauth2"
)

type Option func(c *oauth2.Config)

func Scopes(scopes ...string) Option {
	return func(c *oauth2.Config) {
		c.Scopes = scopes
	}
}

func RedirectURL(s string) Option {
	return func(c *oauth2.Config) {
		c.RedirectURL = s
	}
}

func Clone(i *oauth2.Config, options ...Option) *oauth2.Config {
	dup := *i

	for _, opt := range options {
		opt(&dup)
	}

	return &dup
}

type Endpoint interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource
}

// ContstantEndpoint returns a static url with a state param.
type ConstantEndpoint string

func (t ConstantEndpoint) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return fmt.Sprintf(string(t), url.QueryEscape(state))
}

func (t ConstantEndpoint) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return nil, errors.New("not implemented")
}

func (t ConstantEndpoint) TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	return nil
}

type NoopEndpoint string

func (t NoopEndpoint) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return ""
}

func (t NoopEndpoint) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return nil, errors.New("not implemented")
}

func (t NoopEndpoint) TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	return nil
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
	defer httpx.TryClose(resp)

	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return r, err
	}

	return r, nil
}
