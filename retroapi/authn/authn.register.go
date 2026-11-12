package authn

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/internal/md5x"
	"github.com/retrovibed/retrovibed/retroapi/internal/timex"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"golang.org/x/crypto/ssh"
)

// AutoRegistration builds the client Register needs from signer and
// registers in one call. failures are retried, see AutoRegistrationOnce for the single pass version.
// account selects the profile when the key is associated with multiple.
func AutoRegistration(ctx context.Context, signer ssh.Signer, account string, options ...httpx.ClientOption) (*Session, error) {
	c, err := RegistrationJWTClient(ctx, signer, options...)
	if err != nil {
		return nil, err
	}

	return Register(ctx, c, account)
}

// AutoRegistrationOnce is the single pass version of AutoRegistration for cli commands: a failure is returned
// immediately instead of being retried. account selects the profile when the key is associated with multiple, when
// join is set the identity instead joins the existing account rather than signing up for a new one; the profile
// created is pending approval by an administrator of the account. when the key is associated with multiple
// profiles and no account was selected choose is asked to select one. the returned session is the identity
// to act as, pass its account to the clients created afterwards.
func AutoRegistrationOnce(ctx context.Context, signer ssh.Signer, account string, join bool, choose Chooser, options ...httpx.ClientOption) (*Session, error) {
	c, err := RegistrationJWTClient(ctx, signer, options...)
	if err != nil {
		return nil, err
	}

	return RegisterOnce(ctx, c, HTTPClientDefaults(options...), signer, account, join, choose)
}

// Chooser selects the profile to use when the key is associated with multiple. returning no profile (and no error)
// means a selection could not be made, e.g. when not running interactively.
type Chooser func(ctx context.Context, profiles []*Authn) (*Authn, error)

// permanent marks errors retrying cannot resolve.
type permanent struct{ error }

func (t permanent) Unwrap() error { return t.error }

// Register registers the identity, signing up for a new account when the key is not associated with any.
// account selects the profile when the key is associated with multiple.
func Register(ctx context.Context, c *http.Client, account string) (*Session, error) {
	return RegisterAccount(ctx, c, nil, nil, account, false)
}

// RegisterOnce is a single pass of RegisterAccount, failures are returned instead of retried. when the key
// is associated with multiple profiles and no account was selected choose is asked to select one.
func RegisterOnce(ctx context.Context, c *http.Client, chttp *http.Client, signer ssh.Signer, account string, join bool, choose Chooser) (*Session, error) {
	body, err := LoginOptionsToken(ctx, c, fmt.Sprintf("https://%s", env.Deeppool()), account)
	if err != nil {
		return nil, err
	}

	return registerPass(ctx, c, chttp, signer, account, join, body, choose)
}

// RegisterAccount registers the identity. a blank account signs up for a new account, an account selects the
// profile when the key is associated with multiple, and when join is set the identity joins the account instead.
// failures are retried until the registration succeeds, the server rejects it, or the context is done.
// c is the oauth2 client of the identity. joining an account sends a jwt signed by signer with chttp,
// a plain client, because the oauth2 transport overwrites the authorization header.
func RegisterAccount(ctx context.Context, c *http.Client, chttp *http.Client, signer ssh.Signer, account string, join bool) (*Session, error) {
	body, err := LoginOptionsToken(ctx, c, fmt.Sprintf("https://%s", env.Deeppool()), account)
	if err != nil {
		return nil, err
	}

	bs := backoffx.New(backoffx.Exponential(200*time.Millisecond), backoffx.Maximum(30*time.Second))
	ctx, done := context.WithTimeout(ctx, 5*time.Minute)
	defer done()
	return backoffx.AttemptV(ctx, bs, func(ctx context.Context, attempts uint) (*Session, error) {
		session, err := registerPass(ctx, c, chttp, signer, account, join, body, nil)
		if err == nil {
			return session, nil
		}

		log.Println("registration failed", err)

		var perm permanent
		if errors.As(err, &perm) {
			return nil, errors.Join(backoffx.ErrStopAttempts, err)
		}

		return nil, err
	})
}

// registerPass a single registration attempt, errors retrying cannot resolve are marked permanent.
func registerPass(ctx context.Context, c *http.Client, chttp *http.Client, signer ssh.Signer, account string, join bool, body []byte, choose Chooser) (_ *Session, err error) {
	var (
		authed  Authed
		session Session
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/authn/ssh", env.Deeppool()), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return nil, err
	}

	if err = jsonx.UnmarshalRead(resp.Body, &authed); err != nil {
		return nil, err
	}

	switch len(authed.Profiles) {
	case 0:
		// the selected account must be joined explicitly, never silently sign up for another account.
		if !join && account != "" {
			return nil, permanent{errorsx.Errorf("identity is not associated with account %s, join it with --account", account)}
		}
		// continue to registration
	case 1:
		session := authed.Profiles[0]
		return &Session{
			Token:   session.Token,
			Account: session.Account,
			Profile: session.Profile,
		}, nil
	default:
		if choose != nil {
			selected, err := choose(ctx, authed.Profiles)
			if err != nil {
				return nil, permanent{err}
			}

			if selected != nil {
				return &Session{
					Token:   selected.Token,
					Account: selected.Account,
					Profile: selected.Profile,
				}, nil
			}
		}

		return nil, permanent{ambiguousProfiles(&authed)}
	}

	endpoint, bearer, client := "signup", authed.SignupToken, c
	if join {
		// joining an account mirrors eg: the identity signs a jwt with its ssh key, the issuer is the account.
		endpoint, client = "register", chttp
		if bearer, err = jwt.NewWithClaims(
			jwtx.NewSSHSigner(),
			jwtx.NewJWTClaims(md5x.String(ssh.FingerprintSHA256(signer.PublicKey())), jwtx.ClaimsOptionAuthnExpiration(), jwtx.ClaimsOptionIssuer(account)),
		).SignedString(signer); err != nil {
			return nil, errorsx.Wrap(err, "unable to sign registration")
		}
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/authn/%s", env.Deeppool(), endpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", fmt.Sprintf("bearer %s", bearer))

	resp, err = httpx.AsError(client.Do(req))
	var herr *httpx.Error
	if errors.As(err, &herr) && herr.Code >= 400 && herr.Code < 500 && herr.Code != http.StatusTooManyRequests {
		// the server rejected the request, retrying cannot succeed.
		return nil, permanent{err}
	} else if err != nil {
		return nil, err
	}

	if err = jsonx.UnmarshalRead(resp.Body, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// LoginOptionsToken fetches the signed login options token selecting the account, it is sent as the body of
// the /authn/ssh request. a blank account selects nothing and returns no options.
func LoginOptionsToken(ctx context.Context, c *http.Client, endpoint string, account string) (_ []byte, err error) {
	if account == "" {
		return nil, nil
	}

	if uuidx.FromStringOrNil(account).IsZero() {
		return nil, errorsx.Errorf("account must be an account id: %s", account)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/authn/s/?%s", endpoint, url.Values{"account_id": {account}}.Encode()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to fetch login options")
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// ambiguousProfiles the error returned when the key is associated with multiple profiles and no account was selected.
func ambiguousProfiles(authed *Authed) error {
	accounts := make([]string, 0, len(authed.Profiles))
	for _, p := range authed.Profiles {
		if p.Account != nil {
			accounts = append(accounts, fmt.Sprintf("%s (%s)", p.Account.Id, p.Account.Description))
		}
	}

	return errorsx.Errorf("identity has %d profiles: %s, select one with --account-id", len(authed.Profiles), strings.Join(accounts, ", "))
}

// AwaitAuthorized polls the authz endpoint until granted returns true for the
// received token or the context is cancelled. Uses exponential backoff capped
// at 30 seconds between attempts.
func AwaitAuthorized(ctx context.Context, c *http.Client, granted func(*Token) bool) error {
	endpoint := fmt.Sprintf("https://%s/m/authz/", env.Deeppool())
	bs := backoffx.New(backoffx.Exponential(500*time.Millisecond), backoffx.Maximum(30*time.Second))

	err := backoffx.Attempt(ctx, bs, func(ctx context.Context, attempts uint) error {
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

		if err = jsonx.UnmarshalRead(resp.Body, &authed); err != nil {
			return err
		}

		if authed.Token == nil || !granted(authed.Token) {
			return errors.New("access not yet granted")
		}

		return nil
	})

	return errorsx.Wrap(err, "awaiting authorization failed")
}

// ProfileState is enabled when none of the disabled timestamps of the profile have passed, otherwise disabled.
// a timestamp the server did not send is treated as not passed.
func ProfileState(p *Profile) string {
	if timex.Min(slices.Collect(timex.Parser(time.RFC3339Nano, p.DisabledAt, p.DisabledManuallyAt, p.DisabledPendingApprovalAt))...).Equal(timex.Inf()) {
		return "enabled"
	}

	return "disabled"
}

func PrintIdentity(w io.Writer, s ssh.Signer, session *Session) error {
	var (
		err3 error
	)
	_, err1 := fmt.Fprintln(w, "fingerprint", ssh.FingerprintSHA256(s.PublicKey()))
	_, err2 := fmt.Fprintln(w, "identity   ", md5x.String(ssh.FingerprintSHA256(s.PublicKey())))
	if session != nil {
		if session.Account != nil {
			_, err3 = fmt.Fprintln(w, "account    ", session.Account.Id)
		}
		if session.Profile != nil {
			_, err6 := fmt.Fprintln(w, "display    ", session.Profile.Display)
			_, err7 := fmt.Fprintln(w, "email      ", session.Profile.Email)
			_, err8 := fmt.Fprintln(w, "status     ", ProfileState(session.Profile))
			err3 = errorsx.Compact(err3, err6, err7, err8)
		}
	}
	_, err4 := fmt.Fprintln(w, "public     ", strings.TrimSpace(string(ssh.MarshalAuthorizedKey(s.PublicKey()))))
	_, err5 := fmt.Fprintln(w, "base64     ", base64.URLEncoding.EncodeToString(s.PublicKey().Marshal()))
	return errorsx.Compact(err1, err2, err3, err4, err5)
}
