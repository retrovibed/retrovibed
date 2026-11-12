package authn_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/stretchr/testify/require"
)

func TestPrintIdentity(t *testing.T) {
	t.Run("includes display and email when the session has a profile", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		signer, err := authn.Seeded(ctx, "print-identity", false, filepath.Join(t.TempDir(), "id"))
		require.NoError(t, err)

		var out bytes.Buffer
		session := &authn.Session{
			Account: &authn.Account{Id: "account-1"},
			Profile: &authn.Profile{Display: "Jane Doe", Email: "jane@example.com"},
		}
		require.NoError(t, authn.PrintIdentity(&out, signer, session))
		require.Contains(t, out.String(), "account-1")
		require.Contains(t, out.String(), "Jane Doe")
		require.Contains(t, out.String(), "jane@example.com")
	})

	t.Run("omits display and email without a session", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		signer, err := authn.Seeded(ctx, "print-identity", false, filepath.Join(t.TempDir(), "id"))
		require.NoError(t, err)

		var out bytes.Buffer
		require.NoError(t, authn.PrintIdentity(&out, signer, nil))
		require.NotContains(t, out.String(), "display")
		require.NotContains(t, out.String(), "email")
	})

	t.Run("does not panic when the session has no account", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		signer, err := authn.Seeded(ctx, "print-identity", false, filepath.Join(t.TempDir(), "id"))
		require.NoError(t, err)

		var out bytes.Buffer
		session := &authn.Session{Profile: &authn.Profile{Display: "Jane Doe", Email: "jane@example.com"}}
		require.NoError(t, authn.PrintIdentity(&out, signer, session))
		require.NotContains(t, out.String(), "account")
		require.Contains(t, out.String(), "jane@example.com")
	})
}
