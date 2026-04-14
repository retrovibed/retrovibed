package metaapi

import (
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

func NewProfileFromMetaProfile(mp meta.Profile) (_ *Profile, err error) {
	var p Profile
	mp = langx.Clone(mp, timex.JSONSafeEncodeOption, timex.UTCEncodeOption)
	if err = grpcx.JSONDecode(mp, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func NewMetaProfileFromProfile(p *Profile, options ...meta.ProfileOption) (mp meta.Profile, err error) {
	if err = grpcx.JSONEncode(p, &mp); err != nil {
		return mp, err
	}

	mp = langx.Clone(mp, timex.JSONSafeDecodeOption, timex.UTCEncodeOption, langx.Compose(options...))
	return mp, nil
}

func ProfileStatusOf(p meta.Profile) ProfileStatus {
	now := time.Now()
	disabledAt := timex.Min(p.DisabledAt, p.DisabledManuallyAt)
	pendingAt := p.DisabledPendingApprovalAt

	if disabledAt.Before(now) && !pendingAt.Before(now) {
		return ProfileStatus_DISABLED
	}
	if pendingAt.Before(now) && !disabledAt.Before(now) {
		return ProfileStatus_PENDING
	}
	if disabledAt.After(now) && pendingAt.After(now) {
		return ProfileStatus_ENABLED
	}
	return ProfileStatus_NONE
}
