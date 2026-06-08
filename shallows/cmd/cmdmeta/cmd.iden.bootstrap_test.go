package cmdmeta

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/stretchr/testify/require"
)

func TestBootstrapPublicKey(t *testing.T) {
	t.Run("creates profile with admin permissions", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		_, pub, err := sshx.UnsafeNewKeyGen().Generate()
		require.NoError(t, err)

		cmd := BootstrapPublicKey{PublicKey: string(pub)}
		require.NoError(t, cmd.run(ctx, q))

		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_sso_identity_ssh"))(t))

		pid, err := sqlx.String(ctx, q, "SELECT profile_id::text FROM meta_sso_identity_ssh")
		require.NoError(t, err)

		var authz meta.Authz
		require.NoError(t, meta.AuthzFindByProfileID(ctx, q, sqlx.NewNullString(pid)).Scan(&authz))
		require.Equal(t, pid, authz.ProfileID)
		require.True(t, authz.Usermanagement)
		require.True(t, authz.LibraryRead)
		require.True(t, authz.LibraryModify)
		require.True(t, authz.BillingModify)
		require.True(t, authz.BillingRead)
	})

	t.Run("idempotent - same key does not create duplicates", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		_, pub, err := sshx.UnsafeNewKeyGen().Generate()
		require.NoError(t, err)

		cmd := BootstrapPublicKey{PublicKey: string(pub)}
		require.NoError(t, cmd.run(ctx, q))
		require.NoError(t, cmd.run(ctx, q))

		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_sso_identity_ssh"))(t))
	})
}

func TestBootstrapAuthorized(t *testing.T) {
	t.Run("creates profiles with admin permissions", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		gen := sshx.UnsafeNewKeyGen()
		buf := bytes.NewBufferString("")
		for range 3 {
			_, pub, err := gen.Generate()
			require.NoError(t, err)
			testx.Must(buf.Write(pub))(t)
		}

		path := filepath.Join(t.TempDir(), "authorized_keys")
		require.NoError(t, os.WriteFile(path, buf.Bytes(), 0600))

		cmd := BootstrapAuthorized{Path: path}
		require.NoError(t, cmd.run(ctx, q))

		require.Equal(t, 3, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
		require.Equal(t, 3, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_sso_identity_ssh"))(t))
		require.Equal(t, 3, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM authz_meta"))(t))

		rows, err := q.QueryContext(ctx, "SELECT id::text, profile_id::text FROM authz_meta")
		require.NoError(t, err)
		defer rows.Close()
		for rows.Next() {
			var authz meta.Authz
			require.NoError(t, rows.Scan(&authz.ID, &authz.ProfileID))
			require.NotEmpty(t, authz.ProfileID)

			require.NoError(t, meta.AuthzFindByProfileID(ctx, q, sqlx.NewNullString(authz.ProfileID)).Scan(&authz))
			require.True(t, authz.Usermanagement)
			require.True(t, authz.LibraryRead)
			require.True(t, authz.LibraryModify)
			require.True(t, authz.BillingModify)
			require.True(t, authz.BillingRead)
		}
		require.NoError(t, rows.Err())
	})

	t.Run("idempotent - same file does not create duplicates", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		gen := sshx.UnsafeNewKeyGen()
		buf := bytes.NewBufferString("")
		for range 3 {
			_, pub, err := gen.Generate()
			require.NoError(t, err)
			testx.Must(buf.Write(pub))(t)
		}

		path := filepath.Join(t.TempDir(), "authorized_keys")
		require.NoError(t, os.WriteFile(path, buf.Bytes(), 0600))

		cmd := BootstrapAuthorized{Path: path}
		require.NoError(t, cmd.run(ctx, q))
		require.NoError(t, cmd.run(ctx, q))

		require.Equal(t, 3, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
		require.Equal(t, 3, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_sso_identity_ssh"))(t))
	})
}
