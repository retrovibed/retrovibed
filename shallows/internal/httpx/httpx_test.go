package httpx_test

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/stretchr/testify/require"
)

func TestRedirectHTTPRequest(t *testing.T) {
	cases := []struct {
		name        string
		inURL       *url.URL
		cs          *tls.ConnectionState
		inIP        net.IP
		defaultPort string
		expectedURL string
	}{
		{
			name:        "uses port provided in the host field",
			inURL:       &url.URL{Scheme: "http", Host: "www.example.com:123", Path: "foo/bar"},
			cs:          nil,
			inIP:        net.ParseIP("127.0.0.1"),
			defaultPort: "456",
			expectedURL: "http://127.0.0.1:123/foo/bar",
		},
		{
			name:        "uses scheme provided by the url",
			inURL:       &url.URL{Scheme: "https", Host: "www.example.com:123", Path: "foo/bar"},
			cs:          &tls.ConnectionState{},
			inIP:        net.ParseIP("127.0.0.1"),
			defaultPort: "456",
			expectedURL: "https://127.0.0.1:123/foo/bar",
		},
		{
			name:        "uses default port if no port provided in host field",
			inURL:       &url.URL{Scheme: "http", Host: "www.example.com", Path: "foo/bar"},
			cs:          nil,
			inIP:        net.ParseIP("127.0.0.1"),
			defaultPort: "456",
			expectedURL: "http://127.0.0.1:456/foo/bar",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := &http.Request{Host: tc.inURL.Host, URL: tc.inURL, TLS: tc.cs}
			got := httpx.RedirectHTTPRequest(req, tc.inIP.String(), tc.defaultPort)
			require.Equal(t, tc.expectedURL, got.String())
		})
	}
}

func TestNewUpload(t *testing.T) {
	_, out, err := httpx.NewUpload("foo", "bar.txt", strings.NewReader("hello world"))
	require.NoError(t, err)
	decoded, err := io.ReadAll(out)
	require.NoError(t, err)
	require.Contains(t, string(decoded), "hello world")
}
