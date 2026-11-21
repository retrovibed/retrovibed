package sshx_test

import (
	"bytes"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestKeyGenGenerate(t *testing.T) {
	t.Run("produces a parsable private key and matching authorized key", func(t *testing.T) {
		priv, pub, err := sshx.NewKeyGenSeeded("seed").Generate()
		require.NoError(t, err)

		signer, err := ssh.ParsePrivateKey(priv)
		require.NoError(t, err)

		authorized, _, _, _, err := ssh.ParseAuthorizedKey(pub)
		require.NoError(t, err)

		require.Equal(t, signer.PublicKey().Marshal(), authorized.Marshal())
	})

	t.Run("returns an error when the entropy source is exhausted", func(t *testing.T) {
		_, _, err := sshx.NewKeyGen(sshx.OptionKeyGenRand(bytes.NewReader(nil))).Generate()
		require.Error(t, err)
	})
}
