package authn_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/internal/httptestx"
	"github.com/stretchr/testify/require"
)

func TestRegisterOnce(t *testing.T) {
	newClient := func(calls *int) *http.Client {
		return httptestx.NewTestClient(func(req *http.Request) *http.Response {
			*calls++
			require.Equal(t, "/authn/ssh", req.URL.Path)
			body, _ := json.Marshal(&authn.Authed{Profiles: []*authn.Authn{
				{Token: "t1", Account: &authn.Account{Id: "account-1", Description: "acme"}, Profile: &authn.Profile{Id: "p1"}},
				{Token: "t2", Account: &authn.Account{Id: "account-2", Description: "globex"}, Profile: &authn.Profile{Id: "p2"}},
			}})
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header)}
		})
	}

	t.Run("the chooser selects among multiple profiles and the selection is returned", func(t *testing.T) {
		var calls int

		session, err := authn.RegisterOnce(t.Context(), newClient(&calls), nil, nil, "", false, func(ctx context.Context, profiles []*authn.Authn) (*authn.Authn, error) {
			require.Len(t, profiles, 2)
			return profiles[1], nil
		})
		require.NoError(t, err)
		require.Equal(t, "t2", session.Token)
		require.Equal(t, "account-2", session.Account.Id)
		require.Equal(t, 1, calls)
	})

	t.Run("a chooser that makes no selection reports the ambiguity once without retrying", func(t *testing.T) {
		var calls int

		_, err := authn.RegisterOnce(t.Context(), newClient(&calls), nil, nil, "", false, func(context.Context, []*authn.Authn) (*authn.Authn, error) {
			return nil, nil
		})
		require.ErrorContains(t, err, "--account-id")
		require.Equal(t, 1, calls)
	})

	t.Run("a chooser error is returned without retrying", func(t *testing.T) {
		var calls int

		_, err := authn.RegisterOnce(t.Context(), newClient(&calls), nil, nil, "", false, func(context.Context, []*authn.Authn) (*authn.Authn, error) {
			return nil, context.Canceled
		})
		require.ErrorIs(t, err, context.Canceled)
		require.Equal(t, 1, calls)
	})

	t.Run("no chooser reports the ambiguity once without retrying", func(t *testing.T) {
		var calls int

		_, err := authn.RegisterOnce(t.Context(), newClient(&calls), nil, nil, "", false, nil)
		require.ErrorContains(t, err, "account-1 (acme)")
		require.Equal(t, 1, calls)
	})
}
