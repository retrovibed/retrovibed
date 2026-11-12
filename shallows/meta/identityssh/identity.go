package identityssh

import (
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
)

type Option func(*Identity)

func OptionTestDefaults(a *Identity) {
	a.ID = uuidx.NewV4().String()
	a.ProfileID = uuidx.Nil().String()
}
