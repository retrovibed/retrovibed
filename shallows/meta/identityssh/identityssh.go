package identityssh

import (
	"context"

	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/sshx"
	"github.com/retrovibed/retrovibed/meta"
)

func ImportParsed(ctx context.Context, q sqlx.Queryer, parsed sshx.Parsed) (err error) {
	p := meta.Profile{
		Description: parsed.Comment,
	}
	if err = meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p); err != nil {
		return errorsx.Wrap(err, "unable to create profile")
	}

	authz := langx.Clone(meta.Authz{
		ProfileID: p.ID,
	}, meta.AuthzOptionAdmin)
	if err = meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz); err != nil {
		return errorsx.Wrap(err, "unable to setup authorizations")
	}

	iden := Identity{
		ID:        sshx.FingerprintMD5(parsed.PublicKey),
		PublicKey: sshx.EncodeBase64PublicKey(parsed.PublicKey),
		ProfileID: p.ID,
	}

	if err = IdentityInsertWithDefaults(ctx, q, iden).Scan(&iden); err != nil {
		return errorsx.Wrap(err, "unable create ssh identity")
	}

	return nil
}

func ImportAuthorizedKeys(ctx context.Context, q sqlx.Queryer, encoded []byte) (err error) {
	for parsed := range sshx.ParseAuthorizedKeys(encoded) {
		if err = ImportParsed(ctx, q, parsed); err != nil {
			return err
		}
	}

	return nil
}
