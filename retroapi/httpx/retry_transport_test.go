package httpx_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/internal/httptestx"
	"github.com/stretchr/testify/require"
)

func TestRetryTransport(t *testing.T) {
	t.Run("retries once on matching status code", func(t *testing.T) {
		invoked := 0
		var body []byte
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			body, _ = io.ReadAll(req.Body)
			invoked++
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     http.Header{},
			}
		})
		c.Transport = httpx.NewRetryTransport(c.Transport, http.StatusBadGateway)

		req, err := http.NewRequest(http.MethodGet, "http://example.com/", strings.NewReader("Hello World"))
		require.NoError(t, err)
		resp, err := c.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadGateway, resp.StatusCode)
		require.Equal(t, 2, invoked)
		require.Equal(t, []byte("Hello World"), body)
	})

	t.Run("retries on context.DeadlineExceeded", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			time.Sleep(time.Hour)
		}))
		defer s.Close()

		c := &http.Client{}
		c.Transport = httpx.NewRetryTransport(c.Transport, http.StatusBadGateway)

		ctx, cancel := context.WithTimeout(context.Background(), -1*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.URL, strings.NewReader("Hello World"))
		require.NoError(t, err)
		resp, err := c.Do(req)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("retries once on unexpected EOF", func(t *testing.T) {
		invoked := 0
		var body []byte
		s := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			invoked++
			if invoked == 1 {
				conn, _, err := resp.(http.Hijacker).Hijack()
				require.NoError(t, err)
				_, _ = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 5\r\n"))
				_ = conn.Close()
				return
			}

			body, _ = io.ReadAll(req.Body)
			resp.WriteHeader(http.StatusOK)
		}))
		defer s.Close()

		c := &http.Client{}
		c.Transport = httpx.NewRetryTransport(c.Transport, http.StatusBadGateway)

		req, err := http.NewRequest(http.MethodPost, s.URL, strings.NewReader("Hello World"))
		require.NoError(t, err)
		resp, err := c.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, 2, invoked)
		require.Equal(t, []byte("Hello World"), body)
	})

	t.Run("retries with a nil body", func(t *testing.T) {
		invoked := 0
		var body []byte
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			body, _ = io.ReadAll(req.Body)
			invoked++
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     http.Header{},
			}
		})
		c.Transport = httpx.NewRetryTransport(c.Transport, http.StatusBadGateway)

		req, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
		require.NoError(t, err)
		resp, err := c.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadGateway, resp.StatusCode)
		require.Equal(t, 2, invoked)
		require.Equal(t, []byte(""), body)
	})
}
