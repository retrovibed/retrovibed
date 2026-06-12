package cmdmeta

import (
	"bytes"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/stretchr/testify/require"
)

func TestCloudRegister(t *testing.T) {
	signer, err := sshx.SignerFromGenerator(sshx.UnsafeNewKeyGen())
	require.NoError(t, err)

	t.Run("prints identity and account information", func(t *testing.T) {
		var buf bytes.Buffer
		session := &authn.Session{Account: &authn.Account{Id: "account-id"}}

		require.NoError(t, CloudRegister{}.run(&buf, signer, session))

		out := buf.String()
		require.Contains(t, out, "fingerprint")
		require.Contains(t, out, "account     account-id")
		require.Contains(t, out, "public")
	})

	t.Run("nil session omits account line", func(t *testing.T) {
		var buf bytes.Buffer

		require.NoError(t, CloudRegister{}.run(&buf, signer, nil))

		require.NotContains(t, buf.String(), "account")
	})
}
