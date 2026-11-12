package cmdmeta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestCloudDisplay(t *testing.T) {
	signer, err := sshx.SignerFromGenerator(sshx.UnsafeNewKeyGen())
	require.NoError(t, err)

	newServer := func(t *testing.T, authed string, authz string) httpx.ClientOption {
		t.Helper()

		mux := http.NewServeMux()
		mux.HandleFunc("/oauth2/ssh/auth", func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode(map[string]string{"code": "stub-code", "state": r.URL.Query().Get("state")}))
		})
		mux.HandleFunc("/oauth2/ssh/token", func(w http.ResponseWriter, r *http.Request) {
			errorsx.Must(fmt.Fprint(w, `{"access_token":"stub-access-token","token_type":"bearer","expires_in":3600}`))
		})
		// the login options token is the selected account id, /authn/ssh only returns the profiles of that account.
		mux.HandleFunc("/authn/s/", func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode(map[string]string{"token": r.URL.Query().Get("account_id")}))
		})
		mux.HandleFunc("/authn/ssh", func(w http.ResponseWriter, r *http.Request) {
			var options struct {
				Token string `json:"token"`
			}
			_ = json.NewDecoder(r.Body).Decode(&options)
			if options.Token == "" {
				errorsx.Must(fmt.Fprint(w, authed))
				return
			}

			var all struct {
				Profiles []map[string]any `json:"profiles"`
			}
			require.NoError(t, json.Unmarshal([]byte(authed), &all))

			var selected []map[string]any
			for _, p := range all.Profiles {
				if p["account"].(map[string]any)["id"] == options.Token {
					selected = append(selected, p)
				}
			}
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"profiles": selected}))
		})
		mux.HandleFunc("/m/authz/", func(w http.ResponseWriter, r *http.Request) { errorsx.Must(fmt.Fprint(w, authz)) })

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			mux.ServeHTTP(w, r)
		}))
		t.Cleanup(srv.Close)

		dst, err := url.Parse(srv.URL)
		require.NoError(t, err)

		return func(c *http.Client) {
			c.Transport = httpx.RewriteHostTransport(dst, c.Transport)
		}
	}

	t.Run("prints identity, display, email and permissions", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		rewrite := newServer(
			t,
			`{"profiles":[{"token":"profile-token","account":{"id":"6d1f5c9e-2f2c-4d2f-9a4b-6a3a3f0a7c11"},"profile":{"id":"profile-id","display":"Jane Doe","email":"jane@example.com","disabled_at":"9999-12-31T23:59:59.999Z","disabled_manually_at":"9999-12-31T23:59:59.999Z","disabled_pending_approval_at":"9999-12-31T23:59:59.999Z"}}]}`,
			`{"bearer":"meta-bearer","token":{"library_read":true,"remote_control":true,"exp":4102444800}}`,
		)

		var buf bytes.Buffer
		require.NoError(t, CloudDisplay{}.run(ctx, &buf, signer, "", nil, rewrite))

		out := buf.String()
		require.Contains(t, out, "fingerprint")
		require.Contains(t, out, "account     6d1f5c9e-2f2c-4d2f-9a4b-6a3a3f0a7c11")
		require.Contains(t, out, "profile     profile-id")
		require.Contains(t, out, "display     Jane Doe")
		require.Contains(t, out, "email       jane@example.com")
		require.Contains(t, out, "status      enabled")
		require.Contains(t, out, "permissions library_read,remote_control")
	})

	t.Run("multiple profiles are chosen interactively", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		rewrite := newServer(
			t,
			`{"profiles":[`+
				`{"token":"t1","account":{"id":"6d1f5c9e-2f2c-4d2f-9a4b-6a3a3f0a7c11","description":"acme"},"profile":{"id":"p1","display":"first"}},`+
				`{"token":"t2","account":{"id":"6d1f5c9e-2f2c-4d2f-9a4b-6a3a3f0a7c22","description":"globex"},"profile":{"id":"p2","display":"second"}}]}`,
			`{"bearer":"meta-bearer","token":{"library_read":true,"exp":4102444800}}`,
		)

		var offered []string
		choose := func(ctx context.Context, profiles []*authn.Authn) (*authn.Authn, error) {
			for _, p := range profiles {
				offered = append(offered, p.Account.Id)
			}
			return profiles[1], nil
		}

		var buf bytes.Buffer
		require.NoError(t, CloudDisplay{}.run(ctx, &buf, signer, "", choose, rewrite))
		require.Equal(t, []string{"6d1f5c9e-2f2c-4d2f-9a4b-6a3a3f0a7c11", "6d1f5c9e-2f2c-4d2f-9a4b-6a3a3f0a7c22"}, offered)
		require.Contains(t, buf.String(), "account     6d1f5c9e-2f2c-4d2f-9a4b-6a3a3f0a7c22")
		require.Contains(t, buf.String(), "display     second")
	})

	t.Run("multiple profiles without a selection is an ambiguity error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		rewrite := newServer(
			t,
			`{"profiles":[{"token":"t1","account":{"id":"account-1"},"profile":{"id":"p1"}},{"token":"t2","account":{"id":"account-2"},"profile":{"id":"p2"}}]}`,
			`{"bearer":"meta-bearer","token":{"library_read":true,"exp":4102444800}}`,
		)

		var buf bytes.Buffer
		err := CloudDisplay{}.run(ctx, &buf, signer, "", func(context.Context, []*authn.Authn) (*authn.Authn, error) { return nil, nil }, rewrite)
		require.ErrorContains(t, err, "--account-id")
		require.Empty(t, buf.String())
	})

	t.Run("authorization failure is returned and nothing is printed", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		rewrite := newServer(
			t,
			`{"profiles":[{"token":"profile-token","account":{"id":"6d1f5c9e-2f2c-4d2f-9a4b-6a3a3f0a7c11"},"profile":{"id":"profile-id"}}]}`,
			`not json`,
		)

		var buf bytes.Buffer
		err := CloudDisplay{}.run(ctx, &buf, signer, "", nil, rewrite)
		require.ErrorContains(t, err, "unable to retrieve cloud permissions")
		require.Empty(t, buf.String())
	})
}
