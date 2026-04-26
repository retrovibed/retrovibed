package flathub

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"eg/compute/tarballs"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

//go:embed .flathubskel
var flathubskel embed.FS

//go:embed .metainfoskel
var metainfoskel embed.FS

func cloneAssets(dir string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		skel, err := fs.Sub(flathubskel, ".flathubskel")
		if err != nil {
			return err
		}
		return egfs.CloneFS(ctx, dir, ".", skel)
	}
}

// Metainfo generates the AppStream metainfo file into the tarball staging directory.
func Metainfo(b *tarballs.Build) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		c := eggit.EnvCommit()
		skel, err := fs.Sub(metainfoskel, ".metainfoskel")
		if err != nil {
			return err
		}

		tmpl, err := template.ParseFS(skel, "space.retrovibe.Console.metainfo.xml")
		if err != nil {
			return err
		}

		dir := filepath.Join(egtarball.Path(tarballs.Retrovibed(b)), "usr", "share", "metainfo")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(dir, "space.retrovibe.Console.metainfo.xml"))
		if err != nil {
			return err
		}
		defer f.Close()

		return tmpl.Execute(f, struct {
			Version string
			Date    string
		}{
			Version: fmt.Sprintf("%d.%d.%d", c.Committer.When.Year(), int(c.Committer.When.Month()), c.Committer.When.Day()),
			Date:    c.Committer.When.Format("2006-01-02"),
		})
	}
}

// Release updates the dedicated Flathub app repo with the latest manifest.
// Used for ongoing releases after the app has been accepted to Flathub.
func Release(ctx context.Context, o eg.Op) error {
	return eg.Sequential(
		shell.Op(
			shell.Newf("gh repo clone flathub/space.retrovibe.Console %s", egenv.EphemeralDirectory("flathub")),
			shell.Newf("gh release download --repo retrovibed/retrovibed --pattern 'space.retrovibe.Console.yml' --dir %s --clobber", egenv.EphemeralDirectory("flathub")),
		),
		cloneAssets(egenv.EphemeralDirectory("flathub")),
		shell.Op(
			shell.Newf("git -C %s add -A", egenv.EphemeralDirectory("flathub")),
			shell.Newf("git -C %s -c user.name='Retrovibed Engineering' -c user.email='engineering@retrovibe.space' commit -m 'Update to latest retrovibed release'", egenv.EphemeralDirectory("flathub")),
			shell.Newf("git -C %s push origin HEAD", egenv.EphemeralDirectory("flathub")),
		),
	)(ctx, o)
}

// Submit opens the initial PR to flathub/flathub for new app submission.
// One-time use only. See: https://docs.flathub.org/docs/for-app-authors/submission
func Submit(ctx context.Context, o eg.Op) error {
	return eg.Sequential(
		shell.Op(
			shell.New("mkdir -m 0700 -p /home/egd/.ssh"),
			shell.New("eg ssh key --path /home/egd/.ssh/id_ed25519"),
			shell.New("git config --global url.\"git@github.com:\".insteadOf \"https://github.com/\""),
			shell.Newf("gh repo fork flathub/flathub --org retrovibed --fork-name flathub --remote-name origin --clone -- %s", egenv.EphemeralDirectory("flathub")),
			shell.Newf("git -C %s remote -v", egenv.EphemeralDirectory("flathub")),
			shell.Newf("git -C %s checkout -b space.retrovibe.Console origin/new-pr", egenv.EphemeralDirectory("flathub")),
			shell.Newf("gh release download --repo retrovibed/retrovibed --pattern 'space.retrovibe.Console.yml' --dir %s --clobber", egenv.EphemeralDirectory("flathub")),
		),
		cloneAssets(egenv.EphemeralDirectory("flathub")),
		shell.Op(
			shell.Newf("chown -R egd:egd %s", egenv.EphemeralDirectory("flathub")).Privileged(),
			shell.Newf("git -C %s add -A .", egenv.EphemeralDirectory("flathub")),
			shell.Newf("git -C %s -c user.name='Retrovibed Engineering' -c user.email='engineering@retrovibe.space' commit -m 'Add space.retrovibe.Console'", egenv.EphemeralDirectory("flathub")),
			shell.Newf("git -C %s remote -v", egenv.EphemeralDirectory("flathub")),
			shell.Newf("git -C %s push origin HEAD", egenv.EphemeralDirectory("flathub")),
		),
	)(ctx, o)
}
