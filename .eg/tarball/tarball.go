// Package tarball basic functionality for creating tar balls.
// it proves the following functionality:
// - build deterministic paths to a directory for adding contents.
// - build dterministric archive names from patterns using information provided by the eg environment.
// - common patterns
// Assumptions:
// - tar/gh cli commands are available.
// - the archive patterns used are unique within the repository the workload is associated with.
// Compability guarentee: as long as you only use the functions provided by this package for accessing and generating
// the tarballs we'll ensure no breaking changes.
package tarball

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egmd5x"
)

func root(paths ...string) string {
	return egenv.CacheDirectory(".eg", "tarball", filepath.Join(paths...))
}

// Path generate a unique directory for the contents that will be inside the archive can be
// placed.
func Path(pattern string) string {
	// we dont want a deep tree in the tarball directory and we want them namespaced.
	// create a uuid from the git repository and the paths provided.
	// this will scope the paths to within a single repository in the cache.
	// longer term we'll move this into a 'run scratch pad directory'
	return root(egmd5x.String(filepath.Join(eggit.EnvCanonicalURI(), pattern)))
}

// replaces the substitution values within the pattern, resulting in the final resulting archive file's name.
func Name(pattern string) string {
	c := eggit.EnvCommit()
	return fmt.Sprintf("%s.tar.xz", c.StringReplace(pattern))
}

// simple template for naming a tarball from git commit information. see eggit.commit.StringReplace for details.
func GitPattern(prefix string) string {
	return fmt.Sprintf("%s.%%git.commit.year%%.%%git.commit.month%%.%%git.commit.day%%%%git.hash.short%%", prefix)
}

func Archive(pattern string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		dir := Path(pattern)
		name := Name(pattern)
		archive := root(name)

		return shell.Run(
			ctx,
			shell.Newf("tar -C %s -Jcvf %s .", dir, archive),
		)
	}
}

// provides the value that
func GithubRelease() string {
	c := eggit.EnvCommit()
	return c.StringReplace("r%git.commit.year%.%git.commit.month%.%git.commit.day%%git.hash.short%")
}

// Release to github, this is very experimental.
// WARNING: for local environments this assumes you've provided the token to the eg command.
// e.g.) GH_TOKEN="$(gh auth token)" eg compute local -e GH_TOKEN
// WARNING: for hosted environments: we've assumed the git auth access token for pulling the repository
// will work. this has not yet been validated.
func Github(pattern string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		c := eggit.EnvCommit()
		version := GithubRelease()
		archive := root(Name(pattern))

		runtime := shell.Runtime().Environ(
			"GH_TOKEN", egenv.String("", "EG_GIT_AUTH_ACCESS_TOKEN", "GH_TOKEN"),
		)

		log.Println("DERP DERP", c.Committer.When)
		return shell.Run(
			ctx,
			runtime.Newf("gh release create --target %s %s %s", c.Hash.String(), version, archive),
		)
	}
}
