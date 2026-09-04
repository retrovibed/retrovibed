package media_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/atomicx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/james-lawrence/torrent/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredTune(t *testing.T) {
	t.Run("example 1 - successful", func(t *testing.T) {
		const (
			tlength = 128 * bytesx.KiB
		)
		var (
			p meta.Profile
			v meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		sfs := fsx.DirVirtual(t.TempDir())
		seed := torrenttestx.QuickClient(t, torrent.ClientConfigStorageDir(sfs.Path()))

		info, expected, err := torrenttest.Seeded(sfs.Path(), tlength, cryptox.NewChaCha8(t.Name()))
		require.NoError(t, err)
		require.EqualValues(t, tlength, info.TotalLength())

		md, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(sfs.Path())), torrent.OptionDisplayName(info.Name))
		require.NoError(t, err)
		_, added, err := seed.Start(md)
		require.NoError(t, err)
		require.True(t, added)
		defer seed.Stop(md)

		addrs := slicesx.MapTransform(func(v net.Addr) string { return v.String() }, seed.ListenAddrs()...)

		vfs := fsx.DirVirtual(t.TempDir())
		tclient := torrenttestx.QuickClient(t, torrent.ClientConfigStorageDir(vfs.Path()))

		md, err = torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(vfs.Path())), torrent.OptionDisplayName(info.Name))
		require.NoError(t, err)
		dl, added, err := tclient.Start(md) // add the torrent to the client but without any information to find the peer.
		require.NoError(t, err)
		require.True(t, added)

		tmd := tracking.NewMetadata(
			&md.ID,
			tracking.MetadataOptionFromInfo(info),
			tracking.MetadataOptionAutoDescription,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd))
		require.NoError(t, tracking.MetadataDownloadByID(ctx, q, tmd.ID).Scan(&tmd))

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(vfs.Path("torrent"), storage.FileOptionPathMakerInfohash),
			asyncx.NewWakeup(ctx),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			media.HTTPDiscoveredOptionRootStorage(vfs),
		).Bind(routes.PathPrefix("/").Subrouter())

		buf, err := json.Marshal(&media.DownloadTuneRequest{
			Peers: addrs,
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("/%s/tune", tmd.ID),
			bytes.NewReader(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode) // confirm the request completed successfully.

		{
			digested := md5.New()
			dctx, done := context.WithTimeout(t.Context(), 5*time.Second)
			defer done()

			n, err := torrent.DownloadInto(dctx, digested, dl)
			require.NoError(t, err)
			require.EqualValues(t, tlength, n)
			require.EqualValues(t, expected.Sum(nil), digested.Sum(nil))
		}

		path := vfs.Path(md.ID.String())
		require.FileExists(t, path)
		require.FileExists(t, fmt.Sprintf("%s.torrent", path))

		writtenMeta, err := metainfo.LoadFromFile(fmt.Sprintf("%s.torrent", path))
		require.NoError(t, err)
		require.Equal(t, md.ID.AsByteArray(), [20]byte(writtenMeta.HashInfoBytes()), "written torrent file should have matching infohash")
	})
}
