package sshx_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestNewKeyGenSeeded(t *testing.T) {
	t.Run("same string seed produces the same key", func(t *testing.T) {
		_, pub1, err := sshx.NewKeyGenSeeded("seed").Generate()
		require.NoError(t, err)

		_, pub2, err := sshx.NewKeyGenSeeded("seed").Generate()
		require.NoError(t, err)

		require.Equal(t, pub1, pub2)
	})

	t.Run("string and byte seeds with the same content produce the same key", func(t *testing.T) {
		_, pub1, err := sshx.NewKeyGenSeeded("seed").Generate()
		require.NoError(t, err)

		_, pub2, err := sshx.NewKeyGenSeeded([]byte("seed")).Generate()
		require.NoError(t, err)

		require.Equal(t, pub1, pub2)
	})

	t.Run("different seeds produce different keys", func(t *testing.T) {
		_, pub1, err := sshx.NewKeyGenSeeded("seed-a").Generate()
		require.NoError(t, err)

		_, pub2, err := sshx.NewKeyGenSeeded("seed-b").Generate()
		require.NoError(t, err)

		require.NotEqual(t, pub1, pub2)
	})
}
