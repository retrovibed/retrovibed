package sshx_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestEnsurePublicKey(t *testing.T) {
	t.Run("writes the public key next to the path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "id")

		s, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed"))
		require.NoError(t, err)

		require.NoError(t, sshx.EnsurePublicKey(s, path))

		encoded, err := os.ReadFile(path + ".pub")
		require.NoError(t, err)
		require.Equal(t, ssh.MarshalAuthorizedKey(s.PublicKey()), encoded)
	})

	t.Run("overwrites an existing public key", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "id")
		require.NoError(t, os.WriteFile(path+".pub", []byte("stale"), 0600))

		s, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed"))
		require.NoError(t, err)

		require.NoError(t, sshx.EnsurePublicKey(s, path))

		encoded, err := os.ReadFile(path + ".pub")
		require.NoError(t, err)
		require.Equal(t, ssh.MarshalAuthorizedKey(s.PublicKey()), encoded)
	})

	t.Run("returns an error when the directory does not exist", func(t *testing.T) {
		s, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed"))
		require.NoError(t, err)

		require.Error(t, sshx.EnsurePublicKey(s, filepath.Join(t.TempDir(), "missing", "id")))
	})
}
