package sshx_test

import (
	"bytes"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestOptionKeyGenRand(t *testing.T) {
	t.Run("same entropy source produces the same key", func(t *testing.T) {
		var kg1, kg2 sshx.KeyGen
		sshx.OptionKeyGenRand(bytes.NewReader(bytes.Repeat([]byte{1}, 32)))(&kg1)
		sshx.OptionKeyGenRand(bytes.NewReader(bytes.Repeat([]byte{1}, 32)))(&kg2)

		priv1, pub1, err := kg1.Generate()
		require.NoError(t, err)

		priv2, pub2, err := kg2.Generate()
		require.NoError(t, err)

		require.Equal(t, priv1, priv2)
		require.Equal(t, pub1, pub2)
	})

	t.Run("different entropy sources produce different keys", func(t *testing.T) {
		var kg1, kg2 sshx.KeyGen
		sshx.OptionKeyGenRand(bytes.NewReader(bytes.Repeat([]byte{1}, 32)))(&kg1)
		sshx.OptionKeyGenRand(bytes.NewReader(bytes.Repeat([]byte{2}, 32)))(&kg2)

		_, pub1, err := kg1.Generate()
		require.NoError(t, err)

		_, pub2, err := kg2.Generate()
		require.NoError(t, err)

		require.NotEqual(t, pub1, pub2)
	})
}
