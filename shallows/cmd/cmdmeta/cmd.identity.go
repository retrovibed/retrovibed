package cmdmeta

type Identity struct {
	Generate  GenerateID   `cmd:"" help:"bootstrap the identity of the device itself, allows you to provide a seed for consistent generation"`
	Bootstrap Bootstrap    `cmd:"" help:"bootstrap authorized users into the system, used to initially provision the system"`
	Add       IdenAdd      `cmd:"" help:"add a user by public key via the api"`
	Show      IdenDisplay  `cmd:"" help:"display current identity"`
	Register  IdenRegister `cmd:"" help:"register the current identity with the cloud service"`
}
