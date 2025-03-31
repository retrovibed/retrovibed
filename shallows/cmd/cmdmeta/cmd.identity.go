package cmdmeta

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"golang.org/x/crypto/ssh"
)

type Identity struct {
	Bootstrap Bootstrap   `cmd:"" help:"bootstrap authorized users into the system, used to initially provision the system"`
	Show      IdenDisplay `cmd:"" help:"display current identity"`
}

type IdenDisplay struct{}

func (t IdenDisplay) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	var (
		signer ssh.Signer = id.Signer
	)

	log.Println("identity", md5x.String(ssh.FingerprintSHA256(signer.PublicKey())))
	log.Println("fingerprint", ssh.FingerprintSHA256(signer.PublicKey()))
	log.Println("public", strings.TrimSpace(string(ssh.MarshalAuthorizedKey(signer.PublicKey()))))
	log.Println("base64", base64.URLEncoding.EncodeToString(signer.PublicKey().Marshal()))

	return nil
}
