package authn_test

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/internal/httptestx"
	"github.com/retrovibed/retrovibed/retroapi/internal/md5x"
	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestRegisterAccount(t *testing.T) {
	const aid = "6d1f5c9e-2f2c-4d2f-9a4b-6a3a3f0a7c11"

	signer, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded(t.Name()))
	require.NoError(t, err)

	t.Run("blank account signs up without requesting login options", func(t *testing.T) {
		var paths []string
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			paths = append(paths, req.URL.Path)
			body, _ := json.Marshal(&authn.Authed{SignupToken: "signup-token"})
			if req.URL.Path == "/authn/signup" {
				body, _ = json.Marshal(&authn.Session{Token: "session"})
			}
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})

		session, err := authn.RegisterAccount(t.Context(), c, nil, nil, "", false)
		require.NoError(t, err)
		require.Equal(t, "session", session.Token)
		require.Equal(t, []string{"/authn/ssh", "/authn/signup"}, paths)
	})

	t.Run("account is sent as login options and registers into it with an ssh signed jwt", func(t *testing.T) {
		var (
			paths   []string
			sshbody string
			bearer  string
		)
		oauthc := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			paths = append(paths, req.URL.Path)
			var body []byte
			switch req.URL.Path {
			case "/authn/s/":
				require.Equal(t, aid, req.URL.Query().Get("account_id"))
				body, _ = json.Marshal(&authn.LoginOptionsResponse{Token: "options-token"})
			case "/authn/ssh":
				b, _ := io.ReadAll(req.Body)
				sshbody = string(b)
				body, _ = json.Marshal(&authn.Authed{SignupToken: "signup-token"})
			default:
				t.Fatalf("unexpected oauth2 request %s", req.URL.Path)
			}
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})
		chttp := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			paths = append(paths, req.URL.Path)
			require.Equal(t, "/authn/register", req.URL.Path)
			bearer = strings.TrimPrefix(req.Header.Get("authorization"), "bearer ")
			body, _ := json.Marshal(&authn.Session{Token: "pending"})
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})

		session, err := authn.RegisterAccount(t.Context(), oauthc, chttp, signer, aid, true)
		require.NoError(t, err)
		require.Equal(t, "pending", session.Token)
		require.Equal(t, []string{"/authn/s/", "/authn/ssh", "/authn/register"}, paths)
		require.Contains(t, sshbody, "options-token")

		// the ssh signing alg is unknown to the jwt parser, read the payload segment directly.
		segments := strings.Split(bearer, ".")
		require.Len(t, segments, 3)
		payload, err := base64.RawURLEncoding.DecodeString(segments[1])
		require.NoError(t, err)
		var claims jwt.RegisteredClaims
		require.NoError(t, json.Unmarshal(payload, &claims))
		require.Equal(t, aid, claims.Issuer)
		require.Equal(t, md5x.String(ssh.FingerprintSHA256(signer.PublicKey())), claims.Subject)
	})

	t.Run("a rejected registration is not retried indefinitely", func(t *testing.T) {
		var registrations int
		oauthc := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			var body []byte
			switch req.URL.Path {
			case "/authn/s/":
				body, _ = json.Marshal(&authn.LoginOptionsResponse{Token: "options-token"})
			case "/authn/ssh":
				body, _ = json.Marshal(&authn.Authed{SignupToken: "signup-token"})
			}
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})
		chttp := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			registrations++
			return &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden", Body: io.NopCloser(strings.NewReader("")), Header: make(http.Header)}
		})

		_, err := authn.RegisterAccount(t.Context(), oauthc, chttp, signer, aid, true)
		require.Error(t, err)
		require.LessOrEqual(t, registrations, 2)
	})

	t.Run("account that is not an id is rejected", func(t *testing.T) {
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			t.Fatalf("unexpected request %s", req.URL.Path)
			return nil
		})

		_, err := authn.RegisterAccount(t.Context(), c, c, signer, "acme", true)
		require.ErrorContains(t, err, "account must be an account id")
	})

	t.Run("multiple profiles without a selected account is an ambiguity error listing the accounts", func(t *testing.T) {
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, "/authn/ssh", req.URL.Path)
			body, _ := json.Marshal(&authn.Authed{Profiles: []*authn.Authn{
				{Token: "t1", Account: &authn.Account{Id: "account-1", Description: "acme"}, Profile: &authn.Profile{Id: "p1"}},
				{Token: "t2", Account: &authn.Account{Id: "account-2", Description: "globex"}, Profile: &authn.Profile{Id: "p2"}},
			}})
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})

		_, err := authn.RegisterAccount(t.Context(), c, nil, nil, "", false)
		require.ErrorContains(t, err, "account-1 (acme)")
		require.ErrorContains(t, err, "account-2 (globex)")
		require.ErrorContains(t, err, "--account-id")
	})

	t.Run("the selected account picks the profile", func(t *testing.T) {
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			var body []byte
			switch req.URL.Path {
			case "/authn/s/":
				require.Equal(t, aid, req.URL.Query().Get("account_id"))
				body, _ = json.Marshal(&authn.LoginOptionsResponse{Token: "options-token"})
			case "/authn/ssh":
				b, _ := io.ReadAll(req.Body)
				require.Contains(t, string(b), "options-token")
				body, _ = json.Marshal(&authn.Authed{Profiles: []*authn.Authn{
					{Token: "t2", Account: &authn.Account{Id: aid}, Profile: &authn.Profile{Id: "p2"}},
				}})
			default:
				t.Fatalf("unexpected request %s", req.URL.Path)
			}
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})

		session, err := authn.RegisterAccount(t.Context(), c, nil, nil, aid, false)
		require.NoError(t, err)
		require.Equal(t, "t2", session.Token)
	})

	t.Run("a selected account the key is not associated with is not silently signed up", func(t *testing.T) {
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			var body []byte
			switch req.URL.Path {
			case "/authn/s/":
				body, _ = json.Marshal(&authn.LoginOptionsResponse{Token: "options-token"})
			case "/authn/ssh":
				body, _ = json.Marshal(&authn.Authed{SignupToken: "signup-token"})
			default:
				t.Fatalf("unexpected request %s", req.URL.Path)
			}
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})

		_, err := authn.RegisterAccount(t.Context(), c, nil, nil, aid, false)
		require.ErrorContains(t, err, "not associated with account")
	})
}
