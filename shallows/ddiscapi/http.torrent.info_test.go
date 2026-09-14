package ddiscapi_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	rootenv "github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestHTTPTorrentInfo(t *testing.T) {
	t.Run("returns torrent info for a known id", func(t *testing.T) {
		var (
			md   tracking.Metadata
			resp metaapi.TorrentInfoResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		fixture := testx.Fixture("example.1.torrent")
		mi, err := metainfo.LoadFromFile(fixture)
		require.NoError(t, err)
		hash := mi.HashInfoBytes()

		md = tracking.NewMetadata(new(int160.FromByteArray(hash)), tracking.MetadataOptionAutoDescription)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		require.NoError(t, fsx.MkDirs(0700, vfs.Path(rootenv.TorrentDirName)))
		require.NoError(t, os.WriteFile(
			vfs.Path(rootenv.TorrentDirName, metainfo.Hash(md.Infohash).String()+tracking.TorrentSuffix),
			testx.IOBytes(testx.Read(".fixtures", "example.1.torrent")),
			0600,
		))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPTorrentInfo(
			q,
			ddiscapi.HTTPTorrentInfoOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			ddiscapi.HTTPTorrentInfoOptionRootStorage(vfs),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp2, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/"+md.ID+"/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp2, req)

		require.Equal(t, http.StatusOK, resp2.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp2.Body, &resp))
		require.Equal(t, "example.1", resp.Details.Name)
		require.EqualValues(t, 132, resp.Details.Length)
		require.True(t, resp.Details.Private)
		require.Equal(t, "mktorrent 1.1", resp.Meta.CreatedBy)
		require.Equal(t, []string{"https://example.com/announce"}, resp.Meta.AnnounceList)
		require.Len(t, resp.Files, 1)
		require.Equal(t, "example.1", resp.Files[0].Name)
		require.EqualValues(t, 132, resp.Files[0].Length)
	})

	t.Run("unknown id returns 404", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		routes := mux.NewRouter()
		ddiscapi.NewHTTPTorrentInfo(
			q,
			ddiscapi.HTTPTorrentInfoOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			ddiscapi.HTTPTorrentInfoOptionRootStorage(vfs),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/non-existent-id/", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("unauthenticated", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		vfs := fsx.DirVirtual(t.TempDir())

		routes := mux.NewRouter()
		ddiscapi.NewHTTPTorrentInfo(
			q,
			ddiscapi.HTTPTorrentInfoOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			ddiscapi.HTTPTorrentInfoOptionRootStorage(vfs),
		).Bind(routes.PathPrefix("/").Subrouter())

		resp, req, err := httptestx.BuildRequestContextBytes(ctx, http.MethodGet, "/non-existent-id/", nil)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
