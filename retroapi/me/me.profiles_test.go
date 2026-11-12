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

func TestProfiles(t *testing.T) {
	t.Run("lists the profiles of the ssh key", func(t *testing.T) {
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, http.MethodGet, req.Method)
			require.Equal(t, "/me/profile/", req.URL.Path)
			body, _ := json.Marshal(&authn.ProfileSearchResponse{Items: []*authn.Profile{{Id: "p1"}, {Id: "p2"}}})
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})

		profiles, err := me.Profiles(t.Context(), c)
		require.NoError(t, err)
		require.Len(t, profiles, 2)
		require.Equal(t, "p2", profiles[1].Id)
	})
}
