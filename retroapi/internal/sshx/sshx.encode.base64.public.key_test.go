package sshx_test

import (
	"encoding/base64"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestEncodeBase64PublicKey(t *testing.T) {
	t.Run("encodes the marshalled key as unpadded url safe base64", func(t *testing.T) {
		s, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed"))
		require.NoError(t, err)

		encoded := sshx.EncodeBase64PublicKey(s.PublicKey())
		require.NotContains(t, encoded, "=")

		decoded, err := base64.RawURLEncoding.DecodeString(encoded)
		require.NoError(t, err)
		require.Equal(t, s.PublicKey().Marshal(), decoded)
	})
}
