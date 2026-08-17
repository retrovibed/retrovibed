package identityssh_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta/identityssh"
	"github.com/stretchr/testify/require"
)

func TestInitializeGuest(t *testing.T) {
	t.Run("sets the display name to the given hostname and marks the profile local_only", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)

		id, err := sshx.SignerFromGenerator(sshx.UnsafeNewKeyGen())
		require.NoError(t, err)

		require.NoError(t, identityssh.InitializeGuest(ctx, q, id.PublicKey(), "example-host"))

		display, err := sqlx.String(ctx, q, "SELECT display FROM meta_profiles")
		require.NoError(t, err)
		require.Equal(t, "example-host", display)

		localOnly, err := sqlx.Value[bool](ctx, q, "SELECT local_only FROM authz_meta")
		require.NoError(t, err)
		require.True(t, localOnly)
	})
}
