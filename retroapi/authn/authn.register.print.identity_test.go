package authn_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/internal/timex"
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
			Profile: &authn.Profile{Id: "profile-1", Display: "Jane Doe", Email: "jane@example.com"},
		}
		require.NoError(t, authn.PrintIdentity(&out, signer, session))
		require.Contains(t, out.String(), "account     account-1")
		require.Contains(t, out.String(), "profile     profile-1")
		require.Contains(t, out.String(), "display     Jane Doe")
		require.Contains(t, out.String(), "email       jane@example.com")
	})

	t.Run("every value starts in the same column", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		signer, err := authn.Seeded(ctx, "print-identity", false, filepath.Join(t.TempDir(), "id"))
		require.NoError(t, err)

		var out bytes.Buffer
		session := &authn.Session{
			Account: &authn.Account{Id: "account-1"},
			Profile: &authn.Profile{Id: "profile-1", Display: "Jane Doe", Email: "jane@example.com"},
		}
		require.NoError(t, authn.PrintIdentity(&out, signer, session))

		for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
			require.Equal(t, " ", line[11:12], line)
			require.NotEqual(t, " ", line[:1], line)
		}
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
		require.NotContains(t, out.String(), "status")
	})

	t.Run("status is enabled when no disabled timestamp has passed", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		signer, err := authn.Seeded(ctx, "print-identity", false, filepath.Join(t.TempDir(), "id"))
		require.NoError(t, err)

		inf := timex.RFC3339Inf().Format(time.RFC3339Nano)
		var out bytes.Buffer
		session := &authn.Session{Profile: &authn.Profile{DisabledAt: inf, DisabledManuallyAt: inf, DisabledPendingApprovalAt: inf}}
		require.NoError(t, authn.PrintIdentity(&out, signer, session))
		require.Contains(t, out.String(), "status      enabled")
	})

	t.Run("status is enabled when the server does not send the timestamps", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		signer, err := authn.Seeded(ctx, "print-identity", false, filepath.Join(t.TempDir(), "id"))
		require.NoError(t, err)

		var out bytes.Buffer
		require.NoError(t, authn.PrintIdentity(&out, signer, &authn.Session{Profile: &authn.Profile{}}))
		require.Contains(t, out.String(), "status      enabled")
	})

	t.Run("status is disabled when the profile was disabled manually", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		signer, err := authn.Seeded(ctx, "print-identity", false, filepath.Join(t.TempDir(), "id"))
		require.NoError(t, err)

		inf := timex.RFC3339Inf().Format(time.RFC3339Nano)
		var out bytes.Buffer
		session := &authn.Session{Profile: &authn.Profile{DisabledAt: inf, DisabledManuallyAt: time.Now().Add(-time.Hour).Format(time.RFC3339Nano), DisabledPendingApprovalAt: inf}}
		require.NoError(t, authn.PrintIdentity(&out, signer, session))
		require.Contains(t, out.String(), "status      disabled")
	})

	t.Run("status is disabled when the profile is pending approval", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		signer, err := authn.Seeded(ctx, "print-identity", false, filepath.Join(t.TempDir(), "id"))
		require.NoError(t, err)

		inf := timex.RFC3339Inf().Format(time.RFC3339Nano)
		var out bytes.Buffer
		session := &authn.Session{Profile: &authn.Profile{DisabledAt: inf, DisabledManuallyAt: inf, DisabledPendingApprovalAt: timex.RFC3339NegInf().Format(time.RFC3339Nano)}}
		require.NoError(t, authn.PrintIdentity(&out, signer, session))
		require.Contains(t, out.String(), "status      disabled")
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
