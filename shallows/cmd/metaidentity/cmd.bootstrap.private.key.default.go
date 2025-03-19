package metaidentity

import (
	"os"
	"strings"

	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/x/errorsx"
	"github.com/retrovibed/retrovibed/internal/x/fsx"
	"github.com/retrovibed/retrovibed/internal/x/stringsx"

	"github.com/rymdport/portal/filechooser"
)

// private key default uses the operating system file apis to symlink the private key into
// our directory structure. Strictly speaking we should be using the flatpak document api.
func (t Bootstrap) initPrivateKey(sshdir string, id *cmdopts.SSHID) error {
	options := filechooser.OpenFileOptions{CurrentFolder: sshdir}
	files, err := filechooser.OpenFile("", "retrovibed identity - private ssh key", &options)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errorsx.Errorf("no private key was chosen")
	}

	privpath := strings.TrimPrefix(stringsx.FirstNonBlank(files...), "file://")

	return errorsx.Wrap(
		errorsx.Compact(
			fsx.IgnoreIsNotExist(os.Remove(id.PrivateKeyPath())),
			os.Symlink(privpath, id.PrivateKeyPath()),
		), "unable to symlink",
	)
}
