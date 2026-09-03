package backups_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/backups"
	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	seed := "5c6d3f2e-1a0b-4c9d-8e7f-0123456789ab"
	privatekey := []byte("-----BEGIN PRIVATE KEY-----\nMC4CAQAwBQYDK2VwBCIEIA==\n-----END PRIVATE KEY-----\n")

	t.Run("is deterministic", func(t *testing.T) {
		a, err := backups.Key(seed, privatekey)
		require.NoError(t, err)
		b, err := backups.Key(seed, privatekey)
		require.NoError(t, err)
		require.Equal(t, a, b)
		require.Len(t, a, 64)
	})

	t.Run("changes with the seed", func(t *testing.T) {
		a, err := backups.Key(seed, privatekey)
		require.NoError(t, err)
		b, err := backups.Key("another-seed", privatekey)
		require.NoError(t, err)
		require.NotEqual(t, a, b)
	})

	t.Run("changes with the private key", func(t *testing.T) {
		a, err := backups.Key(seed, privatekey)
		require.NoError(t, err)
		b, err := backups.Key(seed, []byte("a different identity"))
		require.NoError(t, err)
		require.NotEqual(t, a, b)
	})

	t.Run("requires both inputs", func(t *testing.T) {
		_, err := backups.Key("", privatekey)
		require.Error(t, err)
		_, err = backups.Key(seed, nil)
		require.Error(t, err)
	})
}
