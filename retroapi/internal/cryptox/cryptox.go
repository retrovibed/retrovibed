package cryptox

import (
	"crypto/cipher"
	"crypto/sha256"
	"io"
	"math/rand/v2"

	"golang.org/x/crypto/chacha20"
)

func NewChaCha8[T ~[]byte | string](seed T) *rand.ChaCha8 {
	return rand.NewChaCha8(sha256.Sum256([]byte(seed)))
}

func NewReaderChaCha20(prng, src io.Reader) (*cipher.StreamReader, error) {
	key := make([]byte, chacha20.KeySize)
	if _, err := io.ReadFull(prng, key); err != nil {
		return nil, err
	}

	nonce := make([]byte, chacha20.NonceSize)
	if _, err := io.ReadFull(prng, nonce); err != nil {
		return nil, err
	}

	s, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}

	return &cipher.StreamReader{
		S: s,
		R: src,
	}, nil
}

func NewWriterChaCha20(prng io.Reader, dst io.Writer) (*cipher.StreamWriter, error) {
	key := make([]byte, chacha20.KeySize)
	if _, err := io.ReadFull(prng, key); err != nil {
		return nil, err
	}

	nonce := make([]byte, chacha20.NonceSize)
	if _, err := io.ReadFull(prng, nonce); err != nil {
		return nil, err
	}

	s, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}
	return &cipher.StreamWriter{
		S: s,
		W: dst,
	}, nil
}

func NewOffsetWriterChaCha20(prng io.Reader, dst io.Writer, offset uint32) (*cipher.StreamWriter, error) {
	key := make([]byte, chacha20.KeySize)
	if _, err := io.ReadFull(prng, key); err != nil {
		return nil, err
	}

	nonce := make([]byte, chacha20.NonceSize)
	if _, err := io.ReadFull(prng, nonce); err != nil {
		return nil, err
	}

	s, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}

	s.SetCounter(offset / 64)

	return &cipher.StreamWriter{
		S: s,
		W: dst,
	}, nil
}

func NewOffsetReaderChaCha20(prng io.Reader, src io.Reader, offset uint32) (*cipher.StreamReader, error) {
	key := make([]byte, chacha20.KeySize)
	if _, err := io.ReadFull(prng, key); err != nil {
		return nil, err
	}

	nonce := make([]byte, chacha20.NonceSize)
	if _, err := io.ReadFull(prng, nonce); err != nil {
		return nil, err
	}

	s, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}

	s.SetCounter(offset / 64)

	return &cipher.StreamReader{
		S: s,
		R: src,
	}, nil
}
