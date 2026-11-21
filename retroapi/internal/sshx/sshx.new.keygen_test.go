package sshx_test

import (
	"bytes"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestNewKeyGen(t *testing.T) {
	t.Run("defaults to a random source", func(t *testing.T) {
		_, pub1, err := sshx.NewKeyGen().Generate()
		require.NoError(t, err)

		_, pub2, err := sshx.NewKeyGen().Generate()
		require.NoError(t, err)

		require.NotEqual(t, pub1, pub2)
	})

	t.Run("applies options", func(t *testing.T) {
		_, pub1, err := sshx.NewKeyGen(sshx.OptionKeyGenRand(bytes.NewReader(bytes.Repeat([]byte{7}, 32)))).Generate()
		require.NoError(t, err)

		_, pub2, err := sshx.NewKeyGen(sshx.OptionKeyGenRand(bytes.NewReader(bytes.Repeat([]byte{7}, 32)))).Generate()
		require.NoError(t, err)

		require.Equal(t, pub1, pub2)
	})
}
