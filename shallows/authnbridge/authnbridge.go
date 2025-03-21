// Package authnbridge provides functionality for the console to read its identity using the daemons codebase.
package authnbridge

import (
	"log"

	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/sshx"
)

func Bearer() string {
	signer, err := sshx.AutoCached(sshx.NewKeyGen(), env.PrivateKeyPath())
	if err != nil {
		log.Println("unable to read identity:", err)
		return ""
	}
	_ = signer

	return ""
}
