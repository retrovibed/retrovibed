package archlinux

import (
	"context"
	"eg/compute/maintainer"
	"fmt"
	"path/filepath"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

const (
	// ContainerName is the arch based container used by Generate, Build and Publish; built from
	// .buildskel/Containerfile. makepkg needs a real arch environment to resolve dependencies.
	ContainerName = "retrovibe.archlinux.build"

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

// Runner is the container used to generate, build and publish the retrovibed package.
func Runner() eg.ContainerRunner {
	return Builder(ContainerName)
}

// Prepare builds the container image used by Generate, Build and Publish.
func Prepare(ctx context.Context, o eg.Op) error {
	return eg.Build(Runner().BuildFromFile(".eg/archlinux/.buildskel/Containerfile"))(ctx, o)
}

// PackageDirectory is where Build places the resulting *.pkg.tar.* files.
func PackageDirectory() string {
	return egenv.DirectoryCache("pacman")
}

// Build runs makepkg against the PKGBUILD rendered by Generate (at pkgbuild),
// building the checked out commit instead of the upstream default branch. the
// resulting packages are written to PackageDirectory.
func Build(pkgbuild string) eg.OpFn {
	pkgdir := PackageDirectory()
	workdir := egenv.DirectoryTemp("makepkg-retrovibed")

	runtime := shell.Runtime().
		Environ("PKGDEST", pkgdir).
		Environ("BUILDDIR", filepath.Join("/", "tmp", "build")).
		Environ("SRCDEST", filepath.Join("/", "tmp", "src")).
		Environ("PACKAGER", fmt.Sprintf("%s <%s>", maintainer.Name, maintainer.Email))

	return shell.Op(
		runtime.Newf("mkdir -p %s %s", pkgdir, workdir),
		runtime.Newf("cp %s %s/PKGBUILD", pkgbuild, workdir),
		// build the checked out commit instead of the upstream default branch.
		runtime.Newf("sed -i 's|git+https://github.com/retrovibed/retrovibed.git|git+file://%s|' %s/PKGBUILD", egenv.DirectoryWorking(), workdir),
		runtime.New("makepkg -f --noconfirm --syncdeps").Directory(workdir),
	)
}

// Generate renders the PKGBUILD template (substituting $PKGVER) into path.
// It has no dependency on the AUR git repository, so the rendered file is
// available (e.g. for the GitHub release upload) even when the AUR git
// remote is unreachable.
func Generate(path string) eg.OpFn {
	templatedir := egenv.DirectoryWorking(".eg", "archlinux")

	runtime := shell.Runtime().Environ("PKGVER", pkgver())

	return shell.Op(
		runtime.Newf("id; sudo ls -ldn /workload /workload/eg /workload/eg/.eg"),
		runtime.Newf("ls -lha %s", templatedir),
		runtime.Newf("ls -lha %s", path),
		runtime.Newf("cat %s/PKGBUILD | envsubst '$PKGVER' | tee %s > /dev/null", templatedir, path),
	)
}

// Publish generates the AUR .SRCINFO from the PKGBUILD via makepkg, derives a
// deterministic SSH keypair from the EG_SSH_KEY_SEED environment variable
// (via `eg ssh key`), and pushes PKGBUILD/.SRCINFO to the AUR git repository
// for the retrovibed package, committing and pushing only when something changed.
// It consumes the PKGBUILD rendered by Generate from path.
func Publish(path string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		if egenv.String("", "EG_SSH_KEY_SEED") == "" {
			return fmt.Errorf("EG_SSH_KEY_SEED environment variable is required to publish to the AUR")
		}

		aurdir := egenv.DirectoryTemp("aur-retrovibed")
		commitmsg := fmt.Sprintf("update retrovibed to %s", eggit.EnvCommit().Hash)

		runtime := shell.Runtime().
			Environ("GIT_AUTHOR_NAME", maintainer.Name).
			Environ("GIT_AUTHOR_EMAIL", maintainer.Email).
			Environ("GIT_COMMITTER_NAME", maintainer.Name).
			Environ("GIT_COMMITTER_EMAIL", maintainer.Email)

		return eg.Sequential(
			shell.Op(
				// generate into a default ssh key name so ssh/git pick it up automatically.
				runtime.New(`eg ssh key --seed "$EG_SSH_KEY_SEED" --path "$HOME/.ssh/id_ed25519"`),
				runtime.Newf("if [ -d %s/.git ]; then git -C %s pull --ff-only origin master; else git clone %s %s; fi", aurdir, aurdir, AURRepoSSHURI, aurdir),
				runtime.Newf("cp %s %s/PKGBUILD", path, aurdir),
				runtime.New("makepkg --printsrcinfo > .SRCINFO").Directory(aurdir),
				runtime.New("git add PKGBUILD .SRCINFO").Directory(aurdir),
				runtime.Newf("git diff --cached --quiet || (git commit -m %q && git push origin master)", commitmsg).Directory(aurdir),
			),
		)(ctx, o)
	}
}
