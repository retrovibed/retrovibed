package cmdmeta

type CloudAccount struct {
	Edit CloudAccountEdit `cmd:"" help:"update the account's display name"`
}
