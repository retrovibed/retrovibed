package meta

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

func ProfileCreated(p *Profile) bool {
	return p.CreatedAt.Equal(p.UpdatedAt)
}

func ProfileAutoEnable(ctx context.Context, q sqlx.Queryer, p *Profile) error {
	// if !ProfileCreated(p) {
	// 	return nil
	// }

	return ProfileEnable(ctx, q, p.ID).Scan(p)
}

// Option for a profile
type ProfileOption func(*Profile)

func ProfileOptionTestDefaults(p *Profile) {
	p.DisabledAt = timex.Inf()
	p.DisabledManuallyAt = timex.Inf()
	p.DisabledPendingApprovalAt = timex.Inf()
}

func ProfileOptionDisplay(s string) ProfileOption {
	return func(p *Profile) {
		p.Display = s
	}
}

func ProfileOptionAutoDisplay(s string) ProfileOption {
	return func(p *Profile) {
		p.Display = langx.FirstNonZero(p.Display, s)
	}
}

// ProfileSearch scan results of the query
func ProfileSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) ProfileScanner {
	return NewProfileScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func ProfileSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(ProfileScannerStaticColumns)...).From("meta_profiles")
}

func QueryIsEnabled(e uint32) squirrel.Sqlizer {
	switch e {
	case 0:
		return squirrel.Expr("'t'")
	case 1: // disabled
		return squirrel.Expr("LEAST(meta_profiles.disabled_at, meta_profiles.disabled_manually_at) < NOW() AND NOT LEAST(meta_profiles.disabled_pending_approval_at) < NOW()")
	case 2: // pending
		return squirrel.Expr("LEAST(meta_profiles.disabled_pending_approval_at) < NOW() AND NOT LEAST(meta_profiles.disabled_at, meta_profiles.disabled_manually_at) < NOW()")
	case 3: // enabled
		return squirrel.Expr("LEAST(meta_profiles.disabled_at, meta_profiles.disabled_manually_at, meta_profiles.disabled_pending_approval_at) > NOW()")
	default: // default, everything but pending
		return squirrel.Expr("NOT LEAST(meta_profiles.disabled_pending_approval_at) < NOW()")
	}
}
