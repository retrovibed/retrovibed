package sshx_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestUnsafeNewKeyGen(t *testing.T) {
	t.Run("is deterministic", func(t *testing.T) {
		_, pub1, err := sshx.UnsafeNewKeyGen().Generate()
		require.NoError(t, err)

		_, pub2, err := sshx.UnsafeNewKeyGen().Generate()
		require.NoError(t, err)

		require.Equal(t, pub1, pub2)
	})

	t.Run("is seeded with the unsafe seed", func(t *testing.T) {
		_, pub1, err := sshx.UnsafeNewKeyGen().Generate()
		require.NoError(t, err)

		_, pub2, err := sshx.NewKeyGenSeeded("unsafe").Generate()
		require.NoError(t, err)

		require.Equal(t, pub1, pub2)
	})
}
