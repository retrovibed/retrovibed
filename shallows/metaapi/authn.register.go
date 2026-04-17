package metaapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/retrovibed/retroapi/authn"
	"github.com/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

func Register(ctx context.Context) (*Session, error) {
	c, err := authn.Oauth2DeeppoolHTTPClient(ctx)
	if err != nil {
		return nil, err
	}

	bs := backoffx.New(backoffx.Exponential(200*time.Millisecond), backoffx.Maximum(30*time.Second))
	ctx, done := context.WithTimeout(ctx, 5*time.Minute)
	defer done()
	return backoffx.AttemptV(ctx, bs, func(ctx context.Context, attempts uint) (_ *Session, err error) {
		var (
			authed  Authed
			session Session
		)

		defer func() {
			if err == nil {
				return
			}

			log.Println("registration failed", err)
		}()
		resp, err := httpx.AsError(c.Post(fmt.Sprintf("https://%s/authn/ssh", deeppool.Deeppool()), "", nil))
		if err != nil {
			return nil, err
		}

		if err = json.NewDecoder(resp.Body).Decode(&authed); err != nil {
			return nil, err
		}

		switch len(authed.Profiles) {
		case 0:
			// continue
		case 1:
			session := authed.Profiles[0]
			return &Session{
				Token:   session.Token,
				Account: session.Account,
				Profile: session.Profile,
			}, nil
		default:
			return nil, errors.New("multiple profiles not supported yet")
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/authn/signup", deeppool.Deeppool()), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("authorization", fmt.Sprintf("bearer %s", authed.SignupToken))

		resp, err = httpx.AsError(c.Do(req))
		if err != nil {
			return nil, err
		}

		if err = json.NewDecoder(resp.Body).Decode(&session); err != nil {
			return nil, err
		}

		return &session, nil
	})
}
