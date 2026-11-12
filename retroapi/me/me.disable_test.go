package me_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/internal/httptestx"
	"github.com/retrovibed/retrovibed/retroapi/me"
	"github.com/stretchr/testify/require"
)

func TestDisableProfile(t *testing.T) {
	t.Run("deletes the profile and returns it", func(t *testing.T) {
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, http.MethodDelete, req.Method)
			require.Equal(t, "/me/profile/profile-1", req.URL.Path)
			body, _ := json.Marshal(&authn.ProfileUpdateResponse{Profile: &authn.Profile{Id: "profile-1"}})
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})

		p, err := me.DisableProfile(t.Context(), c, "profile-1")
		require.NoError(t, err)
		require.Equal(t, "profile-1", p.Id)
	})

	t.Run("the current profile cannot be disabled", func(t *testing.T) {
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden", Body: io.NopCloser(strings.NewReader("")), Header: make(http.Header)}
		})

		_, err := me.DisableProfile(t.Context(), c, "profile-1")
		require.True(t, me.Forbidden(err))
		require.False(t, me.NotFound(err))
	})

	t.Run("a profile of another ssh key is not found", func(t *testing.T) {
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{StatusCode: http.StatusNotFound, Status: "404 Not Found", Body: io.NopCloser(strings.NewReader("")), Header: make(http.Header)}
		})

		_, err := me.DisableProfile(t.Context(), c, "profile-1")
		require.True(t, me.NotFound(err))
		require.False(t, me.Forbidden(err))
	})
}
