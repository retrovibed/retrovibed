package cmdddisc

import (
	"bytes"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func newLocateServer(t *testing.T, q *sql.DB) (*http.Client, *httptest.Server) {
	t.Helper()

	var (
		p meta.Profile
		v meta.Authz
	)
	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(t.Context(), q, p).Scan(&p))
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(t.Context(), q, v).Scan(&v))

	routes := mux.NewRouter()
	ddiscapi.NewHTTPLocate(
		q,
		asyncx.NewWakeup(t.Context()),
		ddiscapi.HTTPLocateOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/l").Subrouter())

	srv := httptest.NewServer(routes)
	t.Cleanup(srv.Close)

	token := httpauthtest.UnsafeClaimsToken(
		metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))),
		httpauthtest.UnsafeJWTSecretSource,
	)

	headers := http.Header{"Authorization": []string{token}}
	return &http.Client{
		Transport: httpx.NewHeadersTransport(headers, httpx.HTORoundTripper(
			httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil),
		)),
	}, srv
}

func TestMediaLocateJSONL(t *testing.T) {
	t.Run("handles empty input", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cmd := cmdMediaLocateJSONL{}

		c, srv := newLocateServer(t, q)
		require.NoError(t, cmd.run(ctx, c, srv.URL, &bytes.Buffer{}))
		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_locate"))
	})

	t.Run("submits one locate request per known media record", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cmd := cmdMediaLocateJSONL{}

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "My Title"
		known.Mimetype = "video"

		var buf bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&buf).Encode(known))

		c, srv := newLocateServer(t, q)
		require.NoError(t, cmd.run(ctx, c, srv.URL, &buf))

		expected := ddisc.NewLocate(known.Title, known.Mimetype)
		var found ddisc.Locate
		require.NoError(t, ddisc.LocateFindByID(ctx, q, expected.ID).Scan(&found))
		require.Equal(t, known.Title, found.Query)
		require.Equal(t, known.Mimetype, found.Mimetype)
		require.Equal(t, known.UID, found.KnownMediaID)
	})

	t.Run("processes concurrently with multiple workers", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cmd := cmdMediaLocateJSONL{Workers: 8, Backlog: 8}

		var buf bytes.Buffer
		enc := jsonl.NewEncoder(&buf)

		for i := range 20 {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionRandomID))
			known.Mimetype = "video"
			require.NoError(t, enc.Encode(known))
			_ = i
		}

		c, srv := newLocateServer(t, q)
		require.NoError(t, cmd.run(ctx, c, srv.URL, &buf))
		require.Equal(t, 20, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_locate"))
	})

	t.Run("fails on malformed json line", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cmd := cmdMediaLocateJSONL{}

		buf := bytes.NewBufferString("not valid json\n")

		c, srv := newLocateServer(t, q)
		err := cmd.run(ctx, c, srv.URL, buf)
		require.Error(t, err)
		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_locate"))
	})

	t.Run("skips a locate request with a blank mimetype without hanging", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cmd := cmdMediaLocateJSONL{}

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Mimetype = ""

		var buf bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&buf).Encode(known))

		// a blank mimetype is skipped up front rather than submitted to the
		// server, so this should return immediately rather than hang.
		bctx, bcancel := context.WithTimeout(ctx, 10*time.Second)
		defer bcancel()

		c, srv := newLocateServer(t, q)
		err := cmd.run(bctx, c, srv.URL, &buf)
		require.NoError(t, err)
		require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_locate"))
	})
}
