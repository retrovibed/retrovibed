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

func TestEnableProfile(t *testing.T) {
	t.Run("puts the profile and returns it", func(t *testing.T) {
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, http.MethodPut, req.Method)
			require.Equal(t, "/me/profile/profile-1", req.URL.Path)
			body, _ := json.Marshal(&authn.ProfileUpdateResponse{Profile: &authn.Profile{Id: "profile-1"}})
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})

		p, err := me.EnableProfile(t.Context(), c, "profile-1")
		require.NoError(t, err)
		require.Equal(t, "profile-1", p.Id)
	})
}
