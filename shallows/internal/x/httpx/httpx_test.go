package httpx_test

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	. "github.com/james-lawrence/deeppool/internal/x/httpx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("httpx", func() {
	DescribeTable("RedirectHTTPRequest",
		func(inURL *url.URL, cs *tls.ConnectionState, inIP net.IP, expectedURL, defaultPort string) {
			req := &http.Request{Host: inURL.Host, URL: inURL, TLS: cs}
			Expect(RedirectHTTPRequest(req, inIP.String(), defaultPort).String()).To(Equal(expectedURL))
		},
		Entry(
			"it should use the port provided in the host field",
			&url.URL{Scheme: "http", Host: "www.example.com:123", Path: "tallachat/details"},
			nil,
			net.ParseIP("127.0.0.1"),
			"http://127.0.0.1:123/tallachat/details",
			"456",
		),
		Entry(
			"it should use scheme provided by the url",
			&url.URL{Scheme: "https", Host: "www.example.com:123", Path: "tallachat/details"},
			&tls.ConnectionState{},
			net.ParseIP("127.0.0.1"),
			"https://127.0.0.1:123/tallachat/details",
			"456",
		),
		Entry(
			"it should use the default port if no port is provided in the host field",
			&url.URL{Scheme: "http", Host: "www.example.com", Path: "tallachat/details"},
			nil,
			net.ParseIP("127.0.0.1"),
			"http://127.0.0.1:456/tallachat/details",
			"456",
		),
	)

	Describe("NewUpload", func() {
		It("should be able to copy contents", func() {
			_, out, err := NewUpload("foo", "bar.txt", strings.NewReader("hello world"))
			Expect(err).To(Succeed())
			decoded, err := io.ReadAll(out)
			Expect(err).To(Succeed())
			Expect(string(decoded)).To(ContainSubstring("hello world"))
		})
	})
})
