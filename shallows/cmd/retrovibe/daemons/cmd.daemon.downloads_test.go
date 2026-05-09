package daemons

import (
	"context"
	"crypto/tls"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/downloads"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestDownloadWatcher(t *testing.T) {
	t.Run("inaccessible directory does not fatal", func(t *testing.T) {
		db := sqltestx.Metadatabase(t)
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		dwatcher, err := downloads.NewDirectoryWatcher(ctx, &tls.Config{}, db)
		require.NoError(t, err)

		err = dwatcher.Add("/nonexistent/path/that/does/not/exist")
		require.Error(t, err, "Add should return an error for inaccessible directories")
	})
}
