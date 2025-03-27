package meta

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/squirrelx"
	"github.com/retrovibed/retrovibed/internal/timex"
)

// Option for a profile
type ProfileOption func(*Profile)

func ProfileOptionTestDefaults(p *Profile) {
	p.DisabledAt = timex.Inf()
	p.DisabledManuallyAt = timex.Inf()
	p.DisabledPendingApprovalAt = timex.Inf()
}

func ProfileOptionTimezoneUTC(p *Profile) {
	p.CreatedAt = p.CreatedAt.UTC()
	p.UpdatedAt = p.UpdatedAt.UTC()
	p.DisabledAt = p.DisabledAt.UTC()
	p.DisabledManuallyAt = p.DisabledManuallyAt.UTC()
	p.DisabledPendingApprovalAt = p.DisabledPendingApprovalAt.UTC()
}

func ProfileOptionJSONSafeEncode(p *Profile) {
	p.CreatedAt = timex.RFC3339NanoEncode(p.CreatedAt)
	p.UpdatedAt = timex.RFC3339NanoEncode(p.UpdatedAt)
	p.DisabledAt = timex.RFC3339NanoEncode(p.DisabledAt)
	p.DisabledManuallyAt = timex.RFC3339NanoEncode(p.DisabledManuallyAt)
	p.DisabledPendingApprovalAt = timex.RFC3339NanoEncode(p.DisabledPendingApprovalAt)
}

func ProfileOptionJSONSafeDecode(p *Profile) {
	p.CreatedAt = timex.RFC3339NanoDecode(p.CreatedAt)
	p.UpdatedAt = timex.RFC3339NanoDecode(p.UpdatedAt)
	p.DisabledAt = timex.RFC3339NanoDecode(p.DisabledAt)
	p.DisabledManuallyAt = timex.RFC3339NanoDecode(p.DisabledManuallyAt)
	p.DisabledPendingApprovalAt = timex.RFC3339NanoDecode(p.DisabledPendingApprovalAt)
}

// ProfileSearch scan results of the query
func ProfileSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) ProfileScanner {
	return NewProfileScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func ProfileSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(ProfileScannerStaticColumns)...).From("meta_profiles")
}
