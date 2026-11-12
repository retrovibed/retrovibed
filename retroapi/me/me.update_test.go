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

func TestUpdateProfile(t *testing.T) {
	t.Run("sends only the display and email", func(t *testing.T) {
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, http.MethodPatch, req.Method)
			require.Equal(t, "/me/", req.URL.Path)

			var sent authn.ProfileUpdateRequest
			require.NoError(t, json.NewDecoder(req.Body).Decode(&sent))
			require.Equal(t, "Jane", sent.Profile.Display)
			require.Equal(t, "jane@example.com", sent.Profile.Email)
			require.Empty(t, sent.Profile.Id)

			body, _ := json.Marshal(&authn.ProfileUpdateResponse{Profile: &authn.Profile{Id: "profile-1", Display: "Jane", Email: "jane@example.com"}})
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})

		p, err := me.UpdateProfile(t.Context(), c, "Jane", "jane@example.com")
		require.NoError(t, err)
		require.Equal(t, "Jane", p.Display)
	})

	t.Run("an invalid email is a bad request", func(t *testing.T) {
		c := httptestx.NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{StatusCode: http.StatusBadRequest, Status: "400 Bad Request", Body: io.NopCloser(strings.NewReader("")), Header: make(http.Header)}
		})

		_, err := me.UpdateProfile(t.Context(), c, "Jane", "nope")
		require.True(t, me.BadRequest(err))
	})
}
