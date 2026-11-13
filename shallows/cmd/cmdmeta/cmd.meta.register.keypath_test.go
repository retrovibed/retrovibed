package cmdmeta_test

import (
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	cmdtestx "github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx/v2"
	"github.com/stretchr/testify/require"
)

func TestCloudRegisterPrivateKeyPath(t *testing.T) {
	harness := cmdtestx.New(t, &cmdmeta.Cloud{}).Options(cmdtestx.Standard(t))

	// a missing key fails while loading the signer, before any request is made, so the error
	// shows which path the command used.
	t.Run("loads the identity from the provided private-key-path", func(t *testing.T) {
		keypath := filepath.Join(t.TempDir(), "custom-id")

		err := harness.Run("register", "--private-key-path", keypath)

		require.Error(t, err)
		require.ErrorContains(t, err, keypath)
	})

	t.Run("each provided private-key-path is used", func(t *testing.T) {
		dir := t.TempDir()
		first := filepath.Join(dir, "first")
		second := filepath.Join(dir, "second")

		err := harness.Run("register", "--private-key-path", first)
		require.ErrorContains(t, err, first)
		require.NotContains(t, err.Error(), second)

		err = harness.Run("register", "--private-key-path", second)
		require.ErrorContains(t, err, second)
		require.NotContains(t, err.Error(), first)
	})
}
