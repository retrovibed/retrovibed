package authn

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/internal/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/httpx"
	"github.com/retrovibed/retrovibed/retroapi/internal/md5x"
	"golang.org/x/crypto/ssh"
)

func Register(ctx context.Context) (*Session, error) {
	c, err := Oauth2DeeppoolHTTPClient(ctx)
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

// AwaitAuthorized polls the authz endpoint until granted returns true for the
// received token or the context is cancelled. Uses exponential backoff capped
// at 30 seconds between attempts.
func AwaitAuthorized(ctx context.Context, c *http.Client, granted func(*Token) bool) error {
	endpoint := fmt.Sprintf("https://%s/m/authz/", deeppool.Deeppool())
	bs := backoffx.New(backoffx.Exponential(500*time.Millisecond), backoffx.Maximum(30*time.Second))

	err := backoffx.Attempt(ctx, bs, func(ctx context.Context) error {
		var authed AuthzResponse

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return err
		}

		resp, err := httpx.AsError(c.Do(req))
		if err != nil {
			log.Println("awaiting authorization", err)
			return err
		}
		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(&authed); err != nil {
			return err
		}

		if authed.Token == nil || !granted(authed.Token) {
			return errors.New("access not yet granted")
		}

		return nil
	})

	return errorsx.Wrap(err, "awaiting authorization failed")
}

func PrintIdentity(w io.Writer, s ssh.Signer, session *Session) error {
	var (
		err3 error
	)
	_, err1 := fmt.Fprintln(w, "fingerprint", ssh.FingerprintSHA256(s.PublicKey()))
	_, err2 := fmt.Fprintln(w, "identity   ", md5x.String(ssh.FingerprintSHA256(s.PublicKey())))
	if session != nil {
		_, err3 = fmt.Fprintln(w, "account    ", session.Account.Id)
	}
	_, err4 := fmt.Fprintln(w, "public     ", strings.TrimSpace(string(ssh.MarshalAuthorizedKey(s.PublicKey()))))
	_, err5 := fmt.Fprintln(w, "base64     ", base64.URLEncoding.EncodeToString(s.PublicKey().Marshal()))
	return errorsx.Compact(err1, err2, err3, err4, err5)
}
