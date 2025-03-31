package identityssh_test

import (
	"bytes"
	"testing"

	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/sshx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/meta/identityssh"
	"github.com/stretchr/testify/require"
)

func TestImportAuthorizedKeys(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()
	q := sqltestx.Metadatabase(t)
	defer q.Close()

	gen := sshx.UnsafeNewKeyGen()
	buf := bytes.NewBufferString("")

	for i := 0; i < 5; i++ {
		_, pub, err := gen.Generate()
		require.NoError(t, err)
		testx.Must(buf.Write(pub))(t)
	}

	require.NoError(t, identityssh.ImportAuthorizedKeys(ctx, q, buf.Bytes()))
	require.Equal(t, 5, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_profiles"))(t))
}
