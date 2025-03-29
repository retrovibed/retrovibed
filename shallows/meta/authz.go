package meta

type AuthzOption func(*Authz)

func AuthzOptionAdmin(v *Authz) {
	v.Usermanagement = true
}

func AuthzOptionProfileID(pid string) AuthzOption {
	return func(a *Authz) {
		a.ProfileID = pid
	}
}
