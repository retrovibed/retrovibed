package console

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// RunDevWatch runs the console as systemd user services that hot reload on dart changes
// and restart when the native libraries change, see .eg/dev/.systemd.
func RunDevWatch(ctx context.Context, o eg.Op) error {
	runtime := devRuntime(
		shell.Runtime().
			Environ("RETROVIBED_META_ENDPOINT", "api.retrovibe.space").
			Environ("RETROVIBED_CONSOLE_ENDPOINT", "console.retrovibe.space"),
	)

	return shell.Run(
		ctx,
		runtime.New("mkdir -p ~/.config/environment.d"),
		// the services read this file, see EnvironmentFile in .eg/dev/.systemd.
		runtime.New("env | tee ~/.config/environment.d/retro.env"),
		shell.Newf("ls %s/* | xargs -I {} machinectl shell egd@.host /usr/bin/systemctl --user link --force {}", egenv.DirectoryWorking(".eg", "dev", ".systemd")).Privileged(),
		shell.New("machinectl shell egd@.host /usr/bin/systemctl --user daemon-reload").Privileged(),
		shell.New("machinectl shell egd@.host /usr/bin/systemctl --user enable --now retro-console.target retro-console.service retro-console-watch.service retro-console-watch-native.service").Privileged(),
		shell.New("machinectl shell egd@.host /usr/bin/journalctl --user -f -n 100 -u retro-console.service -u retro-console-build.service").Privileged().Timeout(egenv.TTL()),
	)
}
