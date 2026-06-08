package identityssh_test

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta/identityssh"
	"github.com/stretchr/testify/require"
)

func TestImportAuthorizedKeys(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()
	q := sqltestx.Metadatabase(t)

	gen := sshx.UnsafeNewKeyGen()
	buf := bytes.NewBufferString("")

	for range 5 {
		_, pub, err := gen.Generate()
		require.NoError(t, err)
		testx.Must(buf.Write(pub))(t)
	}

	require.NoError(t, identityssh.ImportAuthorizedKeys(ctx, q, buf.Bytes()))
	require.Equal(t, 5, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
}

func TestImportPublicKey(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()
	q := sqltestx.Metadatabase(t)

	id, err := sshx.SignerFromGenerator(sshx.UnsafeNewKeyGen())
	require.NoError(t, err)

	require.NoError(t, identityssh.ImportPublicKey(ctx, q, id.PublicKey()))
	require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))

	pid, err := sqlx.String(ctx, q, "SELECT profile_id FROM meta_sso_identity_ssh")
	require.NoError(t, err)
	require.Equal(t, uuid.UUID([]byte(pid)).String(), sshx.FingerprintMD5(id.PublicKey()))
}

func TestImportPublicKeyIdempotent(t *testing.T) {
	// we want to be able to idempotent
	ctx, done := testx.Context(t)
	defer done()
	q := sqltestx.Metadatabase(t)

	id, err := sshx.SignerFromGenerator(sshx.UnsafeNewKeyGen())
	require.NoError(t, err)

	require.NoError(t, identityssh.ImportPublicKey(ctx, q, id.PublicKey()))
	require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))

	pid, err := sqlx.String(ctx, q, "SELECT profile_id FROM meta_sso_identity_ssh")
	require.NoError(t, err)
	require.Equal(t, uuid.UUID([]byte(pid)).String(), sshx.FingerprintMD5(id.PublicKey()))

	require.NoError(t, identityssh.ImportPublicKey(ctx, q, id.PublicKey()))
	require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
	pid, err = sqlx.String(ctx, q, "SELECT profile_id FROM meta_sso_identity_ssh")
	require.NoError(t, err)
	require.Equal(t, uuid.UUID([]byte(pid)).String(), sshx.FingerprintMD5(id.PublicKey()))
}
