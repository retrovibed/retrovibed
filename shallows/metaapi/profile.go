package metaapi

import (
	"github.com/retrovibed/retrovibed/internal/grpcx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/meta"
)

func NewProfileFromMetaProfile(mp meta.Profile) (_ *Profile, err error) {
	var p Profile
	mp = langx.Clone(mp, meta.ProfileOptionJSONSafeEncode, meta.ProfileOptionTimezoneUTC)
	if err = grpcx.JSONDecode(mp, &p); err != nil {
		return nil, err
	}

	return &p, nil
}
