package tvdb

import (
	"context"

	"github.com/dashotv/tvdb/openapi"
	"github.com/dashotv/tvdb/openapi/models/operations"
	"github.com/dashotv/tvdb/openapi/models/shared"
)

type Client struct {
	Key   string
	Token string
	sdk   *openapi.SDK
	ctx   context.Context
}

// Login creates a new client by logging in with an apikey
func Login(apikey string) (*Client, error) {
	s := openapi.New()
	c := &Client{
		Key:   apikey,
		Token: "",
		sdk:   s,
		ctx:   context.Background(),
	}

	p := operations.PostLoginRequestBody{
		Apikey: c.Key,
	}

	r, err := c.PostLogin(p)
	if err != nil {
		return nil, err
	}

	c.Token = *r.Data.Token
	return New(apikey, *r.Data.Token), nil
}

// New creates a client with an existing token
func New(apikey, token string) *Client {
	s := openapi.New(openapi.WithSecurity(shared.Security{BearerAuth: token}))
	return &Client{
		Key:   apikey,
		Token: token,
		sdk:   s,
		ctx:   context.Background(),
	}
}
