package httpx_test

import (
	"net/http"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsError(t *testing.T) {
	t.Run("returns nil for a successful response", func(t *testing.T) {
		resp, err := httpx.AsError(&http.Response{StatusCode: http.StatusOK, Header: http.Header{}}, nil)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("returns an *Error for a 4xx/5xx response", func(t *testing.T) {
		_, err := httpx.AsError(&http.Response{StatusCode: http.StatusNotFound, Status: "404 Not Found", Header: http.Header{}}, nil)
		require.Error(t, err)

		cause, ok := httpx.ErrorWithCode(err, http.StatusNotFound)
		require.True(t, ok)
		require.Equal(t, http.StatusNotFound, cause.Code)
	})

	t.Run("passes through the original error unchanged", func(t *testing.T) {
		orig := &http.Response{}
		_, err := httpx.AsError(orig, assert.AnError)
		require.Equal(t, assert.AnError, err)
	})
}

func TestErrorWithCode(t *testing.T) {
	t.Run("matches when the error code is one of the given codes", func(t *testing.T) {
		_, err := httpx.AsError(&http.Response{StatusCode: http.StatusConflict, Status: "409 Conflict", Header: http.Header{}}, nil)

		cause, ok := httpx.ErrorWithCode(err, http.StatusNotFound, http.StatusConflict)
		require.True(t, ok)
		require.Equal(t, http.StatusConflict, cause.Code)
	})

	t.Run("does not match a different code", func(t *testing.T) {
		_, err := httpx.AsError(&http.Response{StatusCode: http.StatusInternalServerError, Status: "500 Internal Server Error", Header: http.Header{}}, nil)

		_, ok := httpx.ErrorWithCode(err, http.StatusNotFound)
		require.False(t, ok)
	})

	t.Run("does not match a plain error", func(t *testing.T) {
		_, ok := httpx.ErrorWithCode(assert.AnError, http.StatusNotFound)
		require.False(t, ok)
	})

	t.Run("does not match a nil error", func(t *testing.T) {
		_, ok := httpx.ErrorWithCode(nil, http.StatusNotFound)
		require.False(t, ok)
	})
}

func TestIgnoreError(t *testing.T) {
	t.Run("true when the error code is one of the given codes", func(t *testing.T) {
		_, err := httpx.AsError(&http.Response{StatusCode: http.StatusNotFound, Status: "404 Not Found", Header: http.Header{}}, nil)
		require.True(t, httpx.IgnoreError(err, http.StatusNotFound))
	})

	t.Run("false when the error code does not match", func(t *testing.T) {
		_, err := httpx.AsError(&http.Response{StatusCode: http.StatusInternalServerError, Status: "500 Internal Server Error", Header: http.Header{}}, nil)
		require.False(t, httpx.IgnoreError(err, http.StatusNotFound))
	})

	t.Run("false for a plain error", func(t *testing.T) {
		require.False(t, httpx.IgnoreError(assert.AnError, http.StatusNotFound))
	})
}

func TestErrorCode(t *testing.T) {
	t.Run("nil for a successful response", func(t *testing.T) {
		require.NoError(t, httpx.ErrorCode(&http.Response{StatusCode: http.StatusOK, Header: http.Header{}}))
	})

	t.Run("an *Error for a 4xx/5xx response", func(t *testing.T) {
		err := httpx.ErrorCode(&http.Response{StatusCode: http.StatusNotFound, Status: "404 Not Found", Header: http.Header{}})
		cause, ok := httpx.ErrorWithCode(err, http.StatusNotFound)
		require.True(t, ok)
		require.Equal(t, http.StatusNotFound, cause.Code)
	})
}
