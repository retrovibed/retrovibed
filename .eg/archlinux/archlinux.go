package archlinux

import (
	"context"
	"eg/compute/maintainer"
	"fmt"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

const (
	// AURContainerName is the ubuntu based container used by Publish; built
	// from .aurskel/Containerfile on top of the maintainer.Container image.
	AURContainerName = "retrovibe.archlinux.build"

	// AURRepoSSHURI is the AUR git remote for the retrovibed package.
	AURRepoSSHURI = "ssh://aur@aur.archlinux.org/retrovibed.git"
)

func Builder(name string) eg.ContainerRunner {
	return eg.Container(name)
}

// pkgver derives the PKGBUILD pkgver from the current commit, matching the
// "0.0.:autopatch:" scheme used by egdebuild: 0.0.<committer unix millisecond timestamp>.
func pkgver() string {
	return fmt.Sprintf("0.0.%d", eggit.EnvCommit().Committer.When.UnixMilli())
}

// AURRunner is the container used to publish the retrovibed package to the AUR.
func AURRunner() eg.ContainerRunner {
	return Builder(AURContainerName)
}

// Prepare builds the container image used by Publish.
func Prepare(ctx context.Context, o eg.Op) error {
	return eg.Build(AURRunner().BuildFromFile(".eg/archlinux/.aurskel/Containerfile"))(ctx, o)
}

// Publish generates the AUR .SRCINFO from the PKGBUILD via makepkg, derives a
// deterministic SSH keypair from the EG_SSH_KEY_SEED environment variable
// (via `eg ssh key`), and pushes PKGBUILD/.SRCINFO to the AUR git repository
// for the retrovibed package, committing and pushing only when something changed.
func Publish(ctx context.Context, o eg.Op) error {
	if egenv.String("", "EG_SSH_KEY_SEED") == "" {
		return fmt.Errorf("EG_SSH_KEY_SEED environment variable is required to publish to the AUR")
	}

	templatedir := egenv.WorkingDirectory(".eg", "archlinux")
	aurdir := egenv.EphemeralDirectory("aur-retrovibed")
	commitmsg := fmt.Sprintf("update retrovibed to %s", eggit.EnvCommit().Hash)

	runtime := shell.Runtime().
		Environ("GIT_AUTHOR_NAME", maintainer.Name).
		Environ("GIT_AUTHOR_EMAIL", maintainer.Email).
		Environ("GIT_COMMITTER_NAME", maintainer.Name).
		Environ("GIT_COMMITTER_EMAIL", maintainer.Email).
		Environ("PKGVER", pkgver())

	return eg.Sequential(
		shell.Op(
			// generate into a default ssh key name so ssh/git pick it up automatically.
			runtime.New(`eg ssh key --seed "$EG_SSH_KEY_SEED" --path "$HOME/.ssh/id_ed25519"`),
			runtime.Newf("if [ -d %s/.git ]; then git -C %s pull --ff-only origin master; else git clone %s %s; fi", aurdir, aurdir, AURRepoSSHURI, aurdir),
			runtime.Newf("envsubst '$PKGVER' < %s/PKGBUILD > %s/PKGBUILD", templatedir, aurdir),
			runtime.New("makepkg --printsrcinfo > .SRCINFO").Directory(aurdir),
			runtime.New("git add PKGBUILD .SRCINFO").Directory(aurdir),
			runtime.Newf("git diff --cached --quiet || (git commit -m %q && git push origin master)", commitmsg).Directory(aurdir),
		),
	)(ctx, o)
}
