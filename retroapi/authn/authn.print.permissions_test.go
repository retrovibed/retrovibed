package authn_test

import (
	"bytes"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/stretchr/testify/require"
)

func TestPrintPermissions(t *testing.T) {
	t.Run("lists granted permissions", func(t *testing.T) {
		var out bytes.Buffer
		token := &authn.Token{LibraryRead: true, RemoteControl: true, ArchiveUpload: 1024}
		require.NoError(t, authn.PrintPermissions(&out, token))
		require.Equal(t, "permissions library_read,remote_control,archive_upload=1024\n", out.String())
	})

	t.Run("includes the media metadata upload permission", func(t *testing.T) {
		var out bytes.Buffer
		require.NoError(t, authn.PrintPermissions(&out, &authn.Token{DeviceBackup: true, MediaMetadataUpload: true}))
		require.Equal(t, "permissions device_backup,media_metadata_upload\n", out.String())
	})

	t.Run("no permissions granted", func(t *testing.T) {
		var out bytes.Buffer
		require.NoError(t, authn.PrintPermissions(&out, &authn.Token{}))
		require.Equal(t, "permissions \n", out.String())
	})
}
