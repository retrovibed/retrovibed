package media_test

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/atomicx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestWebsocket(t *testing.T) {
	t.Run("successful connection and update", func(t *testing.T) {
		const (
			torrentlen = bytesx.MiB
		)
		var (
			p  meta.Profile
			v  meta.Authz
			md tracking.Metadata
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		sdir := t.TempDir()
		info, expected, err := torrenttest.Random(sdir, torrentlen)
		require.NoError(t, err)
		smd, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(sdir)))
		require.NoError(t, err)

		lmd, err := torrent.NewFromInfo(info)
		require.NoError(t, err)

		sclient := torrenttestx.QuickClient(t)
		stor, _, err := sclient.Start(smd)
		require.NoError(t, err)
		require.NoError(t, torrent.Verify(t.Context(), stor))

		tclient := torrenttestx.QuickClient(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		md = tracking.NewMetadata(
			langx.Autoptr(smd.ID),
			tracking.MetadataOptionFromInfo(info),
			tracking.MetadataOptionAutoDescription,
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		vfs := fsx.DirVirtual(t.TempDir())
		storageClient := storage.NewFile(vfs.Path(), storage.FileOptionPathMakerInfohash)
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storageClient,
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			media.HTTPDiscoveredOptionRootStorage(vfs),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)

		server := httptest.NewServer(routes)
		defer server.Close()

		wsURL := fmt.Sprintf("ws://%s/%s/socket", server.Listener.Addr().String(), md.ID)

		c, wsResp, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
		})
		require.NoError(t, err)
		defer c.Close(websocket.StatusNormalClosure, "") //nolint: errcheck
		require.Equal(t, http.StatusSwitchingProtocols, wsResp.StatusCode)

		dl, added, err := tclient.Start(lmd, torrent.TuneClientPeer(sclient), torrent.TuneNewConns)
		require.NoError(t, err)
		require.True(t, added)
		require.NoError(t, torrent.Verify(t.Context(), dl))

		ctxTimeout, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		var (
			actual = md5.New()
			result media.Download
		)
		require.Eventually(t, func() bool {
			messageType, data, err := c.Read(ctxTimeout)
			require.NoError(t, err)
			require.Equal(t, websocket.MessageBinary, messageType)
			require.NoError(t, json.Unmarshal(data, &result))
			return result.Downloaded == torrentlen
		}, 10*time.Second, 100*time.Millisecond)

		require.Equal(t, md.ID, result.Media.Id)
		require.EqualValues(t, torrentlen, result.Bytes)
		require.EqualValues(t, torrentlen, result.Downloaded)

		n, err := torrent.DownloadInto(t.Context(), actual, dl)
		require.NoError(t, err)
		require.EqualValues(t, torrentlen, n)
		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		completedAt, err := grpcx.DecodeTime(result.CompletedAt)
		require.NoError(t, err)
		require.WithinDuration(t, time.Now(), completedAt, time.Minute)
	})
}
