package sshx_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestFingerprintMD5(t *testing.T) {
	t.Run("is stable for the same key", func(t *testing.T) {
		s1, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed"))
		require.NoError(t, err)

		s2, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed"))
		require.NoError(t, err)

		fp := sshx.FingerprintMD5(s1.PublicKey())
		require.NotEmpty(t, fp)
		require.Equal(t, fp, sshx.FingerprintMD5(s2.PublicKey()))
	})

	t.Run("differs for different keys", func(t *testing.T) {
		s1, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed-a"))
		require.NoError(t, err)

		s2, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed-b"))
		require.NoError(t, err)

		require.NotEqual(t, sshx.FingerprintMD5(s1.PublicKey()), sshx.FingerprintMD5(s2.PublicKey()))
	})
}
