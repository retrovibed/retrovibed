package authn_test

import (
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/stretchr/testify/require"
)

func TestAuthzResponseDecode(t *testing.T) {
	t.Run("ignores permissions the server has that this client does not", func(t *testing.T) {
		// device_backup is a deeppool grant with no counterpart in meta.Token.
		// decoding used to run through protojson, which rejects unknown fields
		// outright and took down every authenticated deeppool request with it.
		const body = `{"bearer":"deadbeef","token":{"sub":"a2f3c0d1","exp":1757093600,"library_read":true,"device_backup":true}}`

		var resp authn.AuthzResponse
		require.NoError(t, jsonx.UnmarshalRead(strings.NewReader(body), &resp))
		require.Equal(t, "deadbeef", resp.Bearer)
		require.Equal(t, "a2f3c0d1", resp.Token.Sub)
		require.Equal(t, int64(1757093600), resp.Token.Exp)
		require.True(t, resp.Token.LibraryRead)
	})

	t.Run("reads the registered claims by their rfc 7519 names", func(t *testing.T) {
		const body = `{"token":{"jti":"3ace4786","iss":"72bc4114","sub":"0a484599","sid":"ff0772bc","iat":1757089900,"exp":1757093600,"nbf":1757089800}}`

		var resp authn.AuthzResponse
		require.NoError(t, jsonx.UnmarshalRead(strings.NewReader(body), &resp))
		require.Equal(t, "3ace4786", resp.Token.Jti)
		require.Equal(t, "72bc4114", resp.Token.Iss)
		require.Equal(t, "0a484599", resp.Token.Sub)
		require.Equal(t, "ff0772bc", resp.Token.Sid)
		require.Equal(t, int64(1757089900), resp.Token.Iat)
		require.Equal(t, int64(1757093600), resp.Token.Exp)
		require.Equal(t, int64(1757089800), resp.Token.Nbf)
	})

	t.Run("accepts the archive quotas as either a string or a number", func(t *testing.T) {
		// protojson emitted 64 bit integers as strings, jsonx emits them as
		// numbers. deeppool has done both, so both have to decode.
		const body = `{"token":{"archive_upload":"1048576","archive_download":2097152}}`

		var resp authn.AuthzResponse
		require.NoError(t, jsonx.UnmarshalRead(strings.NewReader(body), &resp))
		require.Equal(t, uint64(1048576), resp.Token.ArchiveUpload)
		require.Equal(t, uint64(2097152), resp.Token.ArchiveDownload)
	})
}
