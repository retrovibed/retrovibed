package cmdmeta

type Usermanagement struct {
	Ls     U12TLs     `cmd:"" help:"list profiles"`
	Grant  U12TGrant  `cmd:"" help:"grant access to a profile"`
	Revoke U12TRevoke `cmd:"" help:"revoke access from a profile"`
}
