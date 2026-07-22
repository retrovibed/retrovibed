package cmdddisc

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	ddiscapiimport "github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func newDiscoveryImportServer(t *testing.T, q *sql.DB) (*http.Client, *httptest.Server) {
	t.Helper()

	var (
		p meta.Profile
		v meta.Authz
	)
	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
	require.NoError(t, meta.ProfileInsertWithDefaults(t.Context(), q, p).Scan(&p))
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(t.Context(), q, v).Scan(&v))

	routes := mux.NewRouter()
	ddiscapi.NewHTTPDiscovery(
		q,
		searchplugin.Unimplemented{},
		ddisc.UnimplementedStrategy{},
		ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/ddisc/discovery").Subrouter())

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

func randomMagnetURI() (string, int160.T) {
	infohash := int160.Random()
	return fmt.Sprintf("magnet:?xt=urn:btih:%s", infohash.String()), infohash
}

func countDiscovery(ctx context.Context, t *testing.T, q *sql.DB, ids ...string) int {
	t.Helper()
	found := sqlx.Scan(tracking.UnknownSearch(ctx, q, tracking.UnknownSearchBuilder().Where(tracking.UnknownHashQueryByIDs(ids...))))
	rows := 0
	for range found.Iter() {
		rows++
	}
	require.NoError(t, found.Err())
	return rows
}

func TestDiscoveryImportJSONL(t *testing.T) {
	t.Run("handles empty input", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cmd := cmdDiscoveryImportJSONL{}

		c, srv := newDiscoveryImportServer(t, q)
		require.NoError(t, cmd.run(ctx, c, srv.URL, &bytes.Buffer{}))
		require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_unknown_infohashes"))(t))
	})

	t.Run("imports one discovery entry per line", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cmd := cmdDiscoveryImportJSONL{}

		var buf bytes.Buffer
		enc := jsonl.NewEncoder(&buf)

		ids := make([]string, 0, 3)
		for range 3 {
			uri, infohash := randomMagnetURI()
			require.NoError(t, enc.Encode(&ddiscapiimport.Import{Uri: uri}))
			ids = append(ids, torrentx.HashUID(&infohash))
		}

		c, srv := newDiscoveryImportServer(t, q)
		require.NoError(t, cmd.run(ctx, c, srv.URL, &buf))
		require.Equal(t, len(ids), countDiscovery(ctx, t, q, ids...))
	})

	t.Run("processes concurrently with multiple workers", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cmd := cmdDiscoveryImportJSONL{Workers: 8, Backlog: 8}

		var buf bytes.Buffer
		enc := jsonl.NewEncoder(&buf)

		ids := make([]string, 0, 20)
		for range 20 {
			uri, infohash := randomMagnetURI()
			require.NoError(t, enc.Encode(&ddiscapiimport.Import{Uri: uri}))
			ids = append(ids, torrentx.HashUID(&infohash))
		}

		c, srv := newDiscoveryImportServer(t, q)
		require.NoError(t, cmd.run(ctx, c, srv.URL, &buf))
		require.Equal(t, len(ids), countDiscovery(ctx, t, q, ids...))
	})

	t.Run("fails on malformed json line", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cmd := cmdDiscoveryImportJSONL{}

		buf := bytes.NewBufferString("not valid json\n")
		c, srv := newDiscoveryImportServer(t, q)
		err := cmd.run(ctx, c, srv.URL, buf)
		require.Error(t, err)
		require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_unknown_infohashes"))(t))
	})

	t.Run("surfaces an error on an invalid magnet uri without hanging", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cmd := cmdDiscoveryImportJSONL{}

		var buf bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&buf).Encode(&ddiscapiimport.Import{Uri: "not a magnet uri"}))

		// a bad magnet uri is a permanent failure, but backoffx.Attempt has no
		// concept of permanent vs transient errors and retries forever; bound
		// the context so the test fails fast instead of hanging until the
		// suite-level test timeout.
		bctx, bcancel := context.WithTimeout(ctx, 200*time.Millisecond)
		defer bcancel()

		c, srv := newDiscoveryImportServer(t, q)
		err := cmd.run(bctx, c, srv.URL, &buf)
		require.Error(t, err)
	})
}
