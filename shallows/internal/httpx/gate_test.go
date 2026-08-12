package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/stretchr/testify/require"
)

func TestGatedResponseDisabledShortCircuits(t *testing.T) {
	called := false
	s := httptest.NewServer(alice.New(
		httpx.GatedResponse(false, http.StatusPreconditionFailed),
	).ThenFunc(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		called = true
	})))
	defer s.Close()

	resp, err := http.Get(s.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusPreconditionFailed, resp.StatusCode)
	require.False(t, called)
}

func TestGatedResponseEnabledPassesThrough(t *testing.T) {
	called := false
	s := httptest.NewServer(alice.New(
		httpx.GatedResponse(true, http.StatusPreconditionFailed),
	).ThenFunc(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		called = true
	})))
	defer s.Close()

	resp, err := http.Get(s.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, called)
}
