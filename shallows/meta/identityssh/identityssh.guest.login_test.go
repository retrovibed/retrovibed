package identityssh_test

import (
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/meta/identityssh"
	"github.com/stretchr/testify/require"
)

func TestGuestLogin(t *testing.T) {
	t.Run("bootstraps exactly one local_only profile", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)
		keypath := authn.SeededOptionPath(filepath.Join(t.TempDir(), "id"))

		require.NoError(t, identityssh.GuestLogin(ctx, q, keypath))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM authz_meta"))(t))

		localOnly, err := sqlx.Value[bool](ctx, q, "SELECT local_only FROM authz_meta")
		require.NoError(t, err)
		require.True(t, localOnly)
	})
}
