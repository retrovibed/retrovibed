package cmdmeta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

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
			fmt.Fprint(w, `{"access_token":"stub-access-token","token_type":"bearer","expires_in":3600}`)
		})
		mux.HandleFunc("/authn/ssh", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, authed) })
		mux.HandleFunc("/m/authz/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, authz) })

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
			`{"profiles":[{"token":"profile-token","account":{"id":"account-id"},"profile":{"id":"profile-id","display":"Jane Doe","email":"jane@example.com"}}]}`,
			`{"bearer":"meta-bearer","token":{"library_read":true,"remote_control":true,"exp":4102444800}}`,
		)

		var buf bytes.Buffer
		require.NoError(t, CloudDisplay{}.run(ctx, &buf, signer, rewrite))

		out := buf.String()
		require.Contains(t, out, "fingerprint")
		require.Contains(t, out, "account     account-id")
		require.Contains(t, out, "display     Jane Doe")
		require.Contains(t, out, "email       jane@example.com")
		require.Contains(t, out, "permissions library_read,remote_control")
	})

	t.Run("authorization failure is returned and nothing is printed", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		rewrite := newServer(
			t,
			`{"profiles":[{"token":"profile-token","account":{"id":"account-id"},"profile":{"id":"profile-id"}}]}`,
			`not json`,
		)

		var buf bytes.Buffer
		err := CloudDisplay{}.run(ctx, &buf, signer, rewrite)
		require.ErrorContains(t, err, "unable to retrieve cloud permissions")
		require.Empty(t, buf.String())
	})
}
