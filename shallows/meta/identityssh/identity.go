package identityssh

import (
	"github.com/gofrs/uuid/v5"
)

type Option func(*Identity)

func OptionTestDefaults(a *Identity) {
	a.ID = uuid.Must(uuid.NewV4()).String()
	a.ProfileID = uuid.Nil.String()
}
