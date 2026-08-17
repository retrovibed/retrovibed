package meta

type AuthzOption func(*Authz)

func AuthzOptionAdmin(v *Authz) {
	v.Usermanagement = true
	v.RemoteControl = true
	v.LibraryRead = true
	v.LibraryModify = true
	v.BillingModify = true
	v.BillingRead = true
}

func AuthzOptionGuest(v *Authz) {
	v.RemoteControl = true
	v.LibraryRead = true
	v.LocalOnly = true
}

func AuthzOptionNoPrivileges(v *Authz) {
	v.Usermanagement = false
	v.RemoteControl = false
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
