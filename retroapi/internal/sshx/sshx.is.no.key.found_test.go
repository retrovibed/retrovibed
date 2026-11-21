package sshx_test

import (
	"errors"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestIsNoKeyFound(t *testing.T) {
	t.Run("true for the error returned by ssh when no key is found", func(t *testing.T) {
		_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(""))
		require.Error(t, err)
		require.True(t, sshx.IsNoKeyFound(err))
	})

	t.Run("true for matching message", func(t *testing.T) {
		require.True(t, sshx.IsNoKeyFound(errors.New("ssh: no key found")))
	})

	t.Run("false for other errors", func(t *testing.T) {
		require.False(t, sshx.IsNoKeyFound(errors.New("something else")))
	})
}
