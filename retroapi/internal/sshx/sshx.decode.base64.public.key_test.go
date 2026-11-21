package sshx_test

import (
	"encoding/base64"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestDecodeBase64PublicKey(t *testing.T) {
	t.Run("round trips an encoded public key", func(t *testing.T) {
		s, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed"))
		require.NoError(t, err)

		pub, err := sshx.DecodeBase64PublicKey(sshx.EncodeBase64PublicKey(s.PublicKey()))
		require.NoError(t, err)
		require.Equal(t, s.PublicKey().Marshal(), pub.Marshal())
	})

	t.Run("returns an error for invalid base64", func(t *testing.T) {
		_, err := sshx.DecodeBase64PublicKey("!!! not base64 !!!")
		require.Error(t, err)
	})

	t.Run("returns an error when the decoded bytes are not a public key", func(t *testing.T) {
		_, err := sshx.DecodeBase64PublicKey(base64.RawURLEncoding.EncodeToString([]byte("not a key")))
		require.Error(t, err)
	})
}
