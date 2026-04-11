package meta

type AuthzOption func(*Authz)

func AuthzOptionAdmin(v *Authz) {
	v.Usermanagement = true
	v.LibraryRead = true
	v.LibraryModify = true
}

func AuthzOptionNoPrivileges(v *Authz) {
	v.Usermanagement = false
	v.LibraryRead = false
	v.LibraryModify = false
	v.BillingModify = false
	v.BillingRead = false
	v.CommunityModify = false
}

func AuthzOptionProfileID(pid string) AuthzOption {
	return func(a *Authz) {
		a.ProfileID = pid
	}
}
