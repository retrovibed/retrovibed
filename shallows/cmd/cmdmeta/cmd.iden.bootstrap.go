package cmdmeta

import (
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"golang.org/x/crypto/ssh"

	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/huhx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/sshx"
	"github.com/retrovibed/retrovibed/internal/stringsx"
	"github.com/retrovibed/retrovibed/internal/userx"
)

type Bootstrap struct {
	Automatic  bool   `name:"automatic" help:"disable key file prompt if no path is provided" default:"false"`
	Seed       string `name:"seed" help:"generate a ssh key deterministically based on a seed value" default:"${vars_random_seed}"`
	Replace    bool   `name:"replace" help:"do not prompt and unconditionally replace the existing identity"`
	SSHKeyPath string `arg:"" name:"sshkeypath" help:"path to ssh private key to use" default:""`
}

func (t Bootstrap) replaceExistingKey(forced bool) bool {
	if forced {
		return true
	}

	return huhx.Fallback(false)(
		huhx.Bool(huh.NewInput().Prompt("y/N: ").
			Title("retrovibed generates a unique identity automatically, do you want to replace your existing identity?")))
}

func (t Bootstrap) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	var (
		s ssh.Signer
	)

	sshdir, err := userx.HomeDirectoryRel(".ssh")
	if err != nil {
		return err
	}

	if stringsx.Blank(t.SSHKeyPath) && !t.Automatic {
		if err := t.initPrivateKey(sshdir, id); err != nil {
			return err
		}
	}
	replace := t.replaceExistingKey(t.Replace)
	if stringsx.Present(t.SSHKeyPath) && replace {
		err := errorsx.Wrap(
			errorsx.Compact(
				fsx.IgnoreIsNotExist(os.Remove(env.PrivateKeyPath())),
				os.Symlink(t.SSHKeyPath, env.PrivateKeyPath()),
			), "unable to symlink",
		)
		if err != nil {
			return err
		}
	} else if replace {
		log.Println("generating credentials with seed", t.Seed)
		err := errorsx.Wrap(
			fsx.IgnoreIsNotExist(os.Remove(env.PrivateKeyPath())), "unable to remove existing key",
		)
		if err != nil {
			return err
		}
	}

	// unconditionally remove generated data from the private key, they'll be regenerated when necessary.
	if err := errorsx.Compact(
		fsx.IgnoreIsNotExist(os.Remove(env.PrivateKeyPath()+".pub")),
		fsx.IgnoreIsNotExist(os.Remove(userx.DefaultConfigDir("torrent.id"))),
	); err != nil {
		return err
	}

	if s, err = sshx.AutoCached(sshx.NewKeyGenSeeded(md5x.FormatString(md5x.Digest(t.Seed, "ssh"))), env.PrivateKeyPath()); err != nil {
		return err
	}

	return sshx.EnsurePublicKey(s, env.PrivateKeyPath())
}
