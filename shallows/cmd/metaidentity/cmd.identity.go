package metaidentity

import (
	"encoding/base64"
	"log"

	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/x/md5x"
	"golang.org/x/crypto/ssh"
)

type Identity struct{}

func (t Identity) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	var (
		signer ssh.Signer = id.Signer
	)

	log.Println("identity", md5x.String(ssh.FingerprintSHA256(signer.PublicKey())))
	log.Println("fingerprint", ssh.FingerprintSHA256(signer.PublicKey()))
	log.Println("base64", base64.URLEncoding.EncodeToString(signer.PublicKey().Marshal()))

	return nil
}
