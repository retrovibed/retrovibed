package backups_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/backups"
	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	const plaintext = "the metadata database, in the clear"

	encrypt := func(t *testing.T, key backups.Key) []byte {
		r, err := key.Encrypt(strings.NewReader(plaintext))
		require.NoError(t, err)
		ciphertext, err := io.ReadAll(r)
		require.NoError(t, err)
		require.NotContains(t, string(ciphertext), plaintext)
		return ciphertext
	}

	decrypt := func(t *testing.T, key backups.Key, ciphertext []byte) string {
		r, err := key.Decrypt(bytes.NewReader(ciphertext))
		require.NoError(t, err)
		decrypted, err := io.ReadAll(r)
		require.NoError(t, err)
		return string(decrypted)
	}

	key, err := backups.NewKey("seed", []byte("identity"))
	require.NoError(t, err)

	t.Run("decrypts what it encrypted", func(t *testing.T) {
		require.Equal(t, plaintext, decrypt(t, key, encrypt(t, key)))
	})

	t.Run("a different identity cannot read it", func(t *testing.T) {
		other, err := backups.NewKey("seed", []byte("a different identity"))
		require.NoError(t, err)
		require.NotEqual(t, plaintext, decrypt(t, other, encrypt(t, key)))
	})

	t.Run("a different seed cannot read it", func(t *testing.T) {
		other, err := backups.NewKey("another-seed", []byte("identity"))
		require.NoError(t, err)
		require.NotEqual(t, plaintext, decrypt(t, other, encrypt(t, key)))
	})

	t.Run("every backup runs under its own keystream", func(t *testing.T) {
		require.NotEqual(t, encrypt(t, key), encrypt(t, key))
	})

	t.Run("requires both inputs", func(t *testing.T) {
		_, err := backups.NewKey("", []byte("identity"))
		require.Error(t, err)
		_, err = backups.NewKey("seed", nil)
		require.Error(t, err)
	})
}
