package metaaccounts

import (
	"encoding/base64"
	"log"

	"github.com/james-lawrence/deeppool/cmd/cmdopts"
	"github.com/james-lawrence/deeppool/internal/x/md5x"
	"golang.org/x/crypto/ssh"

	"github.com/rymdport/portal/filechooser"
)

type Identity struct {
	SSHKeyPath string `name:"sshkeypath" help:"path to ssh key to use" default:"${vars_ssh_key_path}" hidden:"true"`
}

func (t Identity) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	var (
		signer ssh.Signer = id.Signer
	)

	log.Println("identity", md5x.String(ssh.FingerprintSHA256(signer.PublicKey())))
	log.Println("fingerprint", ssh.FingerprintSHA256(signer.PublicKey()))
	log.Println("base64", base64.URLEncoding.EncodeToString(signer.PublicKey().Marshal()))

	return nil
}

type Bootstrap struct {
	SSHKeyPath string `arg:"" name:"sshkeypath" help:"path to ssh private key to use, will default to well known common paths" default:"${vars_ssh_key_path}"`
}

func (t Bootstrap) Run(gctx *cmdopts.Global) (err error) {
	var (
		signer ssh.Signer
	)

	options := filechooser.OpenFileOptions{Multiple: true}
	files, err := filechooser.OpenFile("", "Select files", &options)
	if err != nil {
		return err
	}

	return nil
}
