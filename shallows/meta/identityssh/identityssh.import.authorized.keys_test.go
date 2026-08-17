package identityssh_test

import (
	"bytes"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/meta/identityssh"
	"github.com/stretchr/testify/require"
)

func TestImportAuthorizedKeys(t *testing.T) {
	t.Run("imports every key in the blob", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)

		gen := sshx.UnsafeNewKeyGen()
		buf := bytes.NewBufferString("")

		for range 5 {
			_, pub, err := gen.Generate()
			require.NoError(t, err)
			testx.Must(buf.Write(pub))(t)
		}

		require.NoError(t, identityssh.ImportAuthorizedKeys(ctx, q, buf.Bytes()))
		require.Equal(t, 5, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
	})
}
