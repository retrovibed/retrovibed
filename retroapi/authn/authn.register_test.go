package authn_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/internal/httptestx"
	"github.com/stretchr/testify/require"
)

func authzResponse(t *authn.Token) string {
	body, _ := json.Marshal(&authn.AuthzResponse{Token: t})
	return string(body)
}

func newAuthzClient(fn httptestx.RoundTripFunc) *http.Client {
	return httptestx.NewTestClient(fn)
}

func TestAwaitAuthorized(t *testing.T) {
	granted := func(tok *authn.Token) bool {
		return tok != nil && tok.LibraryRead
	}

	t.Run("succeeds immediately when access is already granted", func(t *testing.T) {
		c := newAuthzClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(authzResponse(&authn.Token{LibraryRead: true}))),
				Header:     make(http.Header),
			}
		})

		err := authn.AwaitAuthorized(t.Context(), c, granted)
		require.NoError(t, err)
	})

	t.Run("retries and succeeds when access is granted on a later attempt", func(t *testing.T) {
		var calls atomic.Int32
		c := newAuthzClient(func(req *http.Request) *http.Response {
			n := calls.Add(1)
			tok := &authn.Token{}
			if n >= 2 {
				tok.LibraryRead = true
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(authzResponse(tok))),
				Header:     make(http.Header),
			}
		})

		err := authn.AwaitAuthorized(t.Context(), c, granted)
		require.NoError(t, err)
		require.GreaterOrEqual(t, calls.Load(), int32(2))
	})

	t.Run("returns error when context is cancelled before access is granted", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
		defer cancel()

		c := newAuthzClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(authzResponse(&authn.Token{}))),
				Header:     make(http.Header),
			}
		})

		err := authn.AwaitAuthorized(ctx, c, granted)
		require.Error(t, err)
	})

	t.Run("retries on HTTP error and returns error when context is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
		defer cancel()

		c := newAuthzClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
				Header:     make(http.Header),
			}
		})

		err := authn.AwaitAuthorized(ctx, c, granted)
		require.Error(t, err)
	})

	t.Run("nil token in response is treated as not granted", func(t *testing.T) {
		var calls atomic.Int32
		c := newAuthzClient(func(req *http.Request) *http.Response {
			n := calls.Add(1)
			var body string
			if n >= 2 {
				body = authzResponse(&authn.Token{LibraryRead: true})
			} else {
				body = `{"bearer":""}`
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}
		})

		err := authn.AwaitAuthorized(t.Context(), c, granted)
		require.NoError(t, err)
		require.GreaterOrEqual(t, calls.Load(), int32(2))
	})
}
