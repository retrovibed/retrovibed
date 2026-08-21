package egapplex

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// Provision decodes profile, registers it system-wide (keyed by its own UUID, so xcodebuild can
// auto-discover it), and copies it into any additional dstpaths for direct embedding into an app
// bundle.
func Provision(profile []byte, dstpaths ...string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		profilepath := egenv.WorkspaceDirectory("apple.mobileprovision")

		if err := writeFile(profilepath, profile, "provisioning profile"); err != nil {
			return err
		}

		env := shell.Runtime().Environ("APPLE_PROFILE_PATH", profilepath)

		cmds := []shell.Command{
			env.New("mkdir -p ~/Library/MobileDevice/Provisioning\\ Profiles"),
			env.New("security cms -D -i ${APPLE_PROFILE_PATH} -o /tmp/apple_profile_decoded.plist && " +
				"UUID=$(/usr/libexec/PlistBuddy -c 'Print :UUID' /tmp/apple_profile_decoded.plist) && " +
				"cp ${APPLE_PROFILE_PATH} ~/Library/MobileDevice/Provisioning\\ Profiles/${UUID}.mobileprovision"),
		}

		for _, dst := range dstpaths {
			cmds = append(cmds, env.New("cp ${APPLE_PROFILE_PATH} "+dst))
		}

		return shell.Run(ctx, cmds...)
	}
}
