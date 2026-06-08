package cmdmedia

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func jsonlBuffer(t *testing.T, records ...library.Known) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, r := range records {
		require.NoError(t, enc.Encode(r))
	}
	return &buf
}

func TestMediaEnvRun(t *testing.T) {
	cmd := knownenv{}

	t.Run("outputs RETROVIBED_ARCHIVE_START_DATE with max released date", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults, library.KnownOptionReleased(time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC))))
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults, library.KnownOptionReleased(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC))))

		knownimport{}.run(ctx, db, jsonlBuffer(t, a, b)) //nolint:errcheck

		var buf bytes.Buffer
		require.NoError(t, cmd.run(ctx, db, &buf))
		require.Contains(t, buf.String(), "RETROVIBED_ARCHIVE_START_DATE=2023-03-01")
	})

	t.Run("uses max not min released date", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b, c library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults, library.KnownOptionReleased(time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC))))
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults, library.KnownOptionReleased(time.Date(2015, 7, 4, 0, 0, 0, 0, time.UTC))))
		require.NoError(t, testx.Fake(&c, library.KnownOptionTestDefaults, library.KnownOptionReleased(time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC))))

		knownimport{}.run(ctx, db, jsonlBuffer(t, a, b, c)) //nolint:errcheck

		var buf bytes.Buffer
		require.NoError(t, cmd.run(ctx, db, &buf))
		require.Contains(t, buf.String(), "RETROVIBED_ARCHIVE_START_DATE=2022-12-31")
	})

	t.Run("empty database outputs NegInf date", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var buf bytes.Buffer
		require.NoError(t, cmd.run(ctx, db, &buf))
		require.Contains(t, buf.String(), "RETROVIBED_ARCHIVE_START_DATE="+timex.NegInf().Format(time.DateOnly))
	})

	t.Run("output is newline terminated env format", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults, library.KnownOptionReleased(time.Date(2021, 5, 20, 0, 0, 0, 0, time.UTC))))

		knownimport{}.run(ctx, db, jsonlBuffer(t, a)) //nolint:errcheck

		var buf bytes.Buffer
		require.NoError(t, cmd.run(ctx, db, &buf))

		line := strings.TrimSpace(buf.String())
		require.Equal(t, "RETROVIBED_ARCHIVE_START_DATE=2021-05-20", line)
	})
}
