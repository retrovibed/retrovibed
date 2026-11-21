package sshx_test

import (
	"bytes"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestSignerFromGenerator(t *testing.T) {
	t.Run("returns a signer for the generated key", func(t *testing.T) {
		s, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed"))
		require.NoError(t, err)
		require.NotNil(t, s)

		_, pub, err := sshx.NewKeyGenSeeded("seed").Generate()
		require.NoError(t, err)

		authorized, _, _, _, err := ssh.ParseAuthorizedKey(pub)
		require.NoError(t, err)
		require.Equal(t, authorized.Marshal(), s.PublicKey().Marshal())
	})

	t.Run("returns an error when generation fails", func(t *testing.T) {
		_, err := sshx.SignerFromGenerator(sshx.NewKeyGen(sshx.OptionKeyGenRand(bytes.NewReader(nil))))
		require.Error(t, err)
	})
}
