package cmdmeta

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/stretchr/testify/require"
)

func TestCloudMeLs(t *testing.T) {
	t.Run("renders the current profile first, marked, with the status of each profile", func(t *testing.T) {
		inf := "9999-12-31T23:59:59.999Z"
		past := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339Nano)

		var buf bytes.Buffer
		require.NoError(t, CloudMeLs{}.render(&buf, "b", []*authn.Profile{
			{Id: "a", AccountId: "acct-a", Display: "first", DisabledAt: inf, DisabledManuallyAt: past, DisabledPendingApprovalAt: inf},
			{Id: "b", AccountId: "acct-b", Display: "second", Email: "b@example.com", DisabledAt: inf, DisabledManuallyAt: inf, DisabledPendingApprovalAt: inf},
		}))

		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		require.Len(t, lines, 3)
		require.Contains(t, lines[0], "STATUS")
		require.Regexp(t, `^\*\s+b\s+acct-b\s+second\s+b@example.com\s+enabled$`, lines[1])
		require.Regexp(t, `^\s+a\s+acct-a\s+first\s+disabled$`, lines[2])
	})
}
