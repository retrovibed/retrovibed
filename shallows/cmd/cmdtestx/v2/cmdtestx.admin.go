package cmdtestx

import (
	"context"
	"database/sql"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/meta/identityssh"
	"github.com/stretchr/testify/require"
)

func Admin(t *testing.T, ctx context.Context, q *sql.DB, keypath string) (admin meta.Profile, authz meta.Authz, iden identityssh.Identity) {
	t.Helper()

	signer, err := sshx.AutoCached(sshx.NewKeyGen(), keypath)
	require.NoError(t, err)

	require.NoError(t, testx.Fake(&admin, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, admin).Scan(&admin))
	require.NoError(t, meta.ProfileEnable(ctx, q, admin.ID).Scan(&admin))
	require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(admin.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))
	require.NoError(t, identityssh.IdentityInsertWithDefaults(ctx, q, identityssh.Identity{
		ID:        sshx.FingerprintMD5(signer.PublicKey()),
		ProfileID: admin.ID,
		PublicKey: sshx.EncodeBase64PublicKey(signer.PublicKey()),
	}).Scan(&iden))

	return admin, authz, iden
}
