package media_test

import (
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttestx"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/atomicx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPDiscoveredTorrent(t *testing.T) {
	t.Run("returns torrent info for a known id", func(t *testing.T) {
		var (
			p    meta.Profile
			v    meta.Authz
			resp metaapi.TorrentInfoResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(vfs.Path(), storage.FileOptionPathMakerInfohash),
			asyncx.NewWakeup(ctx),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			media.HTTPDiscoveredOptionRootStorage(vfs),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)

		uploadMimetype, uploadBuf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.Binary, "content", "example.bin"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(testx.Read(".fixtures", "example.1.torrent"), 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		uploadResp, uploadReq, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(uploadBuf),
			httptestx.RequestOptionAuthorization(token),
			httptestx.RequestOptionHeader("Content-Type", uploadMimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(uploadResp, uploadReq)
		require.Equal(t, http.StatusOK, uploadResp.Result().StatusCode)

		id := testx.Must(sqlx.String(ctx, q, "SELECT id::text FROM torrents_metadata"))(t)

		torrentResp, torrentReq, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/"+id+"/torrent",
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(torrentResp, torrentReq)

		require.Equal(t, http.StatusOK, torrentResp.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(torrentResp.Body, &resp))
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
		var (
			p     meta.Profile
			authz meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		tclient := torrenttestx.QuickClient(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(vfs.Path(), storage.FileOptionPathMakerInfohash),
			asyncx.NewWakeup(ctx),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			media.HTTPDiscoveredOptionRootStorage(vfs),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/non-existent-id/torrent",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})
}
