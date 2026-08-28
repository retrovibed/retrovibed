package cmdmeta

type Device struct {
	Ls   DeviceLs   `cmd:"" help:"list known devices"`
	Add  DeviceAdd  `cmd:"" help:"add a device to this instance"`
	Rm   DeviceRm   `cmd:"" help:"remove a known device"`
	Edit DeviceEdit `cmd:"" help:"edit a device (default, download, description)"`
}
