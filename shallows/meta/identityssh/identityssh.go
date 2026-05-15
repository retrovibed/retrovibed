package identityssh

import (
	"context"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"golang.org/x/crypto/ssh"
)

func InitializeAdmin(ctx context.Context, q sqlx.Queryer, pub ssh.PublicKey) (err error) {
	var (
		p      meta.Profile
		parsed = sshx.Parsed{
			PublicKey: pub,
		}
	)

	if p, err = importParsed(ctx, q, parsed); err != nil {
		return err
	}

	authz := langx.Clone(meta.Authz{
		ID:        sshx.FingerprintMD5(parsed.PublicKey),
		ProfileID: p.ID,
	}, meta.AuthzOptionAdmin)
	if err = meta.AuthzInsertWithIDDefaults(ctx, q, authz).Scan(&authz); err != nil {
		return errorsx.Wrap(err, "unable to setup authorizations")
	}

	if err := meta.ProfileAutoEnable(ctx, q, &p); err != nil {
		return errorsx.Wrap(err, "unable to enable profile")
	}

	return nil
}

func ImportPublicKey(ctx context.Context, q sqlx.Queryer, pub ssh.PublicKey) (err error) {
	var (
		parsed = sshx.Parsed{
			PublicKey: pub,
		}
	)

	return ImportParsed(ctx, q, parsed)
}

func ImportParsed(ctx context.Context, q sqlx.Queryer, parsed sshx.Parsed) (err error) {
	_, err = importParsed(ctx, q, parsed)
	return err
}

func ImportAuthorizedKeys(ctx context.Context, q sqlx.Queryer, encoded []byte) (err error) {
	for parsed := range sshx.ParseAuthorizedKeys(encoded) {
		if err = ImportParsed(ctx, q, parsed); err != nil {
			return err
		}
	}

	return nil
}

func importParsed(ctx context.Context, q sqlx.Queryer, parsed sshx.Parsed) (_zero meta.Profile, err error) {
	p := meta.Profile{
		ID:      sshx.FingerprintMD5(parsed.PublicKey),
		Display: parsed.Comment,
	}

	if err = meta.ProfileInsertWithID(ctx, q, p).Scan(&p); err != nil {
		return _zero, errorsx.Wrap(err, "unable to create profile")
	}

	iden := Identity{
		ID:        sshx.FingerprintMD5(parsed.PublicKey),
		PublicKey: sshx.EncodeBase64PublicKey(parsed.PublicKey),
		ProfileID: p.ID,
	}

	if err = IdentityInsertWithDefaults(ctx, q, iden).Scan(&iden); err != nil {
		return _zero, errorsx.Wrap(err, "unable create ssh identity")
	}

	return p, nil
}
