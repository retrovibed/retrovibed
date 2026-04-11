package cryptox_test

import (
	"bytes"
	"crypto/cipher"
	"crypto/md5"
	"io"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/internal/bytesx"
	"github.com/retrovibed/retrovibed/internal/cryptox"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/chacha20"
)

func EncryptionTest[T ~[]byte | string](t *testing.T, seed T, limit int64, r func(prng io.Reader, src io.Reader) (*cipher.StreamReader, error), w func(prng io.Reader, src io.Writer) (*cipher.StreamWriter, error)) {
	var (
		expected = md5.New()
		actual   = md5.New()
	)

	encrypt, err := cryptox.NewReaderChaCha20(
		cryptox.NewChaCha8(seed),
		io.LimitReader(
			io.TeeReader(cryptox.NewChaCha8(seed), expected),
			limit,
		),
	)
	require.NoError(t, err)

	decrypt, err := cryptox.NewWriterChaCha20(
		cryptox.NewChaCha8(seed),
		actual,
	)
	require.NoError(t, err)

	n, err := io.Copy(decrypt, encrypt)
	require.NoError(t, err)
	require.Equal(t, limit, n)
	require.Equal(t, expected.Sum(nil), actual.Sum(nil))
}

func TestNewReaderChaCha20(t *testing.T) {
	t.Run("should return an error if PRNG provides insufficient data for key", func(t *testing.T) {
		cipher, err := cryptox.NewReaderChaCha20(io.LimitReader(cryptox.NewChaCha8("derp"), chacha20.KeySize-1), cryptox.NewChaCha8("derp"))
		require.Error(t, err)
		require.Nil(t, cipher)
		require.Contains(t, err.Error(), "unexpected EOF")
	})

	t.Run("should return an error if PRNG provides insufficient data for nonce", func(t *testing.T) {
		cipher, err := cryptox.NewReaderChaCha20(io.LimitReader(cryptox.NewChaCha8("derp"), chacha20.KeySize), cryptox.NewChaCha8("derp"))

		require.Error(t, err)
		require.Nil(t, cipher)
		require.Contains(t, err.Error(), "EOF")
	})

	t.Run("returned cipher should be functional for encryption/decryption", func(t *testing.T) {
		EncryptionTest(t, uuid.Must(uuid.NewV4()).Bytes(), bytesx.MiB, cryptox.NewReaderChaCha20, cryptox.NewWriterChaCha20)
	})
}

func TestChaCha20Offset(t *testing.T) {
	t.Run("return ciphers should allow for random indexing into data", func(t *testing.T) {
		var (
			encrypted bytes.Buffer
			decrypted bytes.Buffer
		)

		original := testx.IOBytes(io.LimitReader(cryptox.NewChaCha8(t.Name()), 129))

		w, err := cryptox.NewWriterChaCha20(cryptox.NewChaCha8(t.Name()), &encrypted)
		require.NoError(t, err)
		n, err := io.Copy(w, bytes.NewBuffer(original))
		require.NoError(t, err)
		require.EqualValues(t, len(original), n)

		r, err := cryptox.NewOffsetReaderChaCha20(cryptox.NewChaCha8(t.Name()), io.NewSectionReader(bytes.NewReader(encrypted.Bytes()), 64, int64(len(original))), 64)
		require.NoError(t, err)

		n, err = io.Copy(&decrypted, r)
		require.NoError(t, err)
		require.EqualValues(t, len(original)-64, n)
		require.Equal(t, original[64:], decrypted.Bytes())
	})

	t.Run("encrypt/decrypt using just an writer", func(t *testing.T) {
		var (
			encrypted bytes.Buffer
			decrypted bytes.Buffer
		)

		original := testx.IOBytes(io.LimitReader(cryptox.NewChaCha8(t.Name()), 129))

		w, err := cryptox.NewWriterChaCha20(cryptox.NewChaCha8(t.Name()), &encrypted)
		require.NoError(t, err)
		n, err := io.Copy(w, bytes.NewBuffer(original))
		require.NoError(t, err)
		require.EqualValues(t, len(original), n)

		w, err = cryptox.NewOffsetWriterChaCha20(cryptox.NewChaCha8(t.Name()), &decrypted, 64)
		require.NoError(t, err)

		n, err = io.Copy(w, io.NewSectionReader(bytes.NewReader(encrypted.Bytes()), 64, int64(len(original))))
		require.NoError(t, err)
		require.EqualValues(t, len(original)-64, n)
		require.Equal(t, original[64:], decrypted.Bytes())
	})
}
