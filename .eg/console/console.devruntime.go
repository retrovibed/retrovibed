package console

import "github.com/egdaemon/eg/runtime/wasi/shell"

// devRuntime applies flutterRuntimev2's defaults and the environment shared by every
// local development run of the console.
func devRuntime(rt shell.Command) shell.Command {
	return flutterRuntimev2(rt).
		// Environ("RETROVIBED_TORRENT_DEBUG", "true").
		Environ("RETROVIBED_TORRENT_AUTO_DISCOVERY", "true").
		Environ("RETROVIBED_METADATA_AUTODOWNLOAD", "false").
		Environ("RETROVIBED_DDISC_INDEX_RATIO", "99").
		Environ("RETROVIBED_DDISC_BACKGROUND_FREQUENCY", "20s").
		Environ("RETROVIBED_DDISC_BACKGROUND_WORKERS", "1")
}
