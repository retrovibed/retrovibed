package sshx_test

import (
	"slices"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/sshx"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestParseAuthorizedKeys(t *testing.T) {
	t.Run("yields nothing for empty input", func(t *testing.T) {
		require.Empty(t, slices.Collect(sshx.ParseAuthorizedKeys(nil)))
	})

	t.Run("yields nothing when no key is found", func(t *testing.T) {
		require.Empty(t, slices.Collect(sshx.ParseAuthorizedKeys([]byte("# just a comment\n"))))
	})

	t.Run("yields each key with its comment", func(t *testing.T) {
		s1, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed-a"))
		require.NoError(t, err)

		s2, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed-b"))
		require.NoError(t, err)

		encoded := append(
			sshx.Comment(ssh.MarshalAuthorizedKey(s1.PublicKey()), "first"),
			sshx.Comment(ssh.MarshalAuthorizedKey(s2.PublicKey()), "second")...,
		)

		parsed := slices.Collect(sshx.ParseAuthorizedKeys(encoded))
		require.Len(t, parsed, 2)
		require.Equal(t, s1.PublicKey().Marshal(), parsed[0].Marshal())
		require.Equal(t, "first", parsed[0].Comment)
		require.Equal(t, s2.PublicKey().Marshal(), parsed[1].Marshal())
		require.Equal(t, "second", parsed[1].Comment)
	})

	t.Run("yields key options", func(t *testing.T) {
		s, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed"))
		require.NoError(t, err)

		encoded := append([]byte(`no-pty `), ssh.MarshalAuthorizedKey(s.PublicKey())...)

		parsed := slices.Collect(sshx.ParseAuthorizedKeys(encoded))
		require.Len(t, parsed, 1)
		require.Equal(t, []string{"no-pty"}, parsed[0].Options)
	})

	t.Run("skips invalid lines", func(t *testing.T) {
		s, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed"))
		require.NoError(t, err)

		encoded := append([]byte("garbage\n"), ssh.MarshalAuthorizedKey(s.PublicKey())...)

		parsed := slices.Collect(sshx.ParseAuthorizedKeys(encoded))
		require.Len(t, parsed, 1)
		require.Equal(t, s.PublicKey().Marshal(), parsed[0].Marshal())
	})

	t.Run("stops when the consumer stops", func(t *testing.T) {
		s1, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed-a"))
		require.NoError(t, err)

		s2, err := sshx.SignerFromGenerator(sshx.NewKeyGenSeeded("seed-b"))
		require.NoError(t, err)

		encoded := append(ssh.MarshalAuthorizedKey(s1.PublicKey()), ssh.MarshalAuthorizedKey(s2.PublicKey())...)

		var count int
		for range sshx.ParseAuthorizedKeys(encoded) {
			count++
			break
		}
		require.Equal(t, 1, count)
	})
}
