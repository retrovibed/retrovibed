package sshx_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestComment(t *testing.T) {
	t.Run("appends comment and trims trailing whitespace of the key", func(t *testing.T) {
		require.Equal(t, []byte("ssh-ed25519 AAAA me@host\r\n"), sshx.Comment([]byte("ssh-ed25519 AAAA\n"), "me@host"))
	})

	t.Run("returns the key unchanged when comment is empty", func(t *testing.T) {
		require.Equal(t, []byte("ssh-ed25519 AAAA\n"), sshx.Comment([]byte("ssh-ed25519 AAAA\n"), ""))
	})

	t.Run("returns the key unchanged when comment is only whitespace", func(t *testing.T) {
		require.Equal(t, []byte("ssh-ed25519 AAAA\n"), sshx.Comment([]byte("ssh-ed25519 AAAA\n"), "  \t "))
	})
}
