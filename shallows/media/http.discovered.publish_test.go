package media_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/james-lawrence/torrent/torrenttestx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/atomicx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredPublishTorrent(t *testing.T) {
	t.Run("example 1 - successful case", func(t *testing.T) {
		var (
			p     meta.Profile
			v     meta.Authz
			r     media.PublishedUploadResponse
			found tracking.Metadata
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
		dvfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(dvfs.Path("torrent"), storage.FileOptionPathMakerInfohash),
			asyncx.NewWakeup(ctx),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			media.HTTPDiscoveredOptionRootStorage(dvfs),
		).Bind(routes.PathPrefix("/").Subrouter())

		info, err := torrenttest.RandomMulti(vfs.Path(), 5, 128*bytesx.KiB, bytesx.MiB, metainfo.OptionPieceLength(bytesx.MiB))
		require.NoError(t, err)

		md, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(vfs.Path())), torrent.OptionDisplayName(info.Name))
		require.NoError(t, err)

		mimetype, buf, err := media.PublishRequest(ctx, md, &media.PublishedUploadRequest{
			Entropy:  uuidx.WithSuffix(16),
			Mimetype: mimex.RetrovibedMediaArchive,
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestContext(
			ctx,
			http.MethodPost,
			"/publish",
			buf,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_metadata"))(t))

		routes.ServeHTTP(resp, req)
		require.NoError(t, buf.Close())

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &r))
		require.Equal(t, info.Name, r.Published.Description)
		require.Equal(t, md.ID.String(), r.Published.Id)
		require.Equal(t, grpcx.EncodeTime(langx.Autoderef(timex.JSONSafeEncode(new(timex.Inf())))), r.Published.ExpiresAt)
		require.Equal(t, mimex.RetrovibedMediaArchive, r.Published.Mimetype)
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_metadata"))(t))
		require.Equal(t, md5x.FormatUUID(md5x.Digest(md.ID.Bytes(), uuid.FromStringOrNil(uuidx.WithSuffix(16)).Bytes())), testx.Must(sqlx.String(ctx, q, "SELECT encryption_seed::text FROM torrents_metadata"))(t))
		require.Equal(t, "", testx.Must(sqlx.String(ctx, q, "SELECT tracker FROM torrents_metadata"))(t))
		require.Equal(t, mimex.RetrovibedMediaArchive, testx.Must(sqlx.String(ctx, q, "SELECT mimetype FROM torrents_metadata"))(t))
		require.EqualValues(t, info.TotalLength(), testx.Must(sqlx.Value[int](ctx, q, "SELECT bytes FROM torrents_metadata"))(t))
		require.EqualValues(t, info.TotalLength(), testx.Must(sqlx.Value[int](ctx, q, "SELECT available FROM torrents_metadata"))(t))
		require.EqualValues(t, 0, testx.Must(sqlx.Value[int](ctx, q, "SELECT downloaded FROM torrents_metadata"))(t), "published content was never fetched from peers")
		require.NoError(t, tracking.MetadataFindByInfohash(t.Context(), q, r.Published.Id).Scan(&found))
		require.True(t, found.Seeding)
		require.WithinDuration(t, timex.Inf(), found.ExpiresAt, 0)
		require.WithinDuration(t, timex.NegInf(), found.NextAnnounceAt, 0)
		require.WithinDuration(t, time.Now(), found.CompletedAt, time.Second)

		path := dvfs.Path("torrent", md.ID.String())
		require.DirExists(t, path)
		require.FileExists(t, fmt.Sprintf("%s.torrent", path))

		writtenMeta, err := metainfo.LoadFromFile(fmt.Sprintf("%s.torrent", path))
		require.NoError(t, err)
		require.Equal(t, md.ID.AsByteArray(), [20]byte(writtenMeta.HashInfoBytes()), "written torrent file should have matching infohash")
	})

	t.Run("example 2 - ttl is set", func(t *testing.T) {
		var (
			p     meta.Profile
			v     meta.Authz
			r     media.PublishedUploadResponse
			found tracking.Metadata
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
		dvfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(dvfs.Path("torrent"), storage.FileOptionPathMakerInfohash),
			asyncx.NewWakeup(ctx),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			media.HTTPDiscoveredOptionRootStorage(dvfs),
		).Bind(routes.PathPrefix("/").Subrouter())

		info, err := torrenttest.RandomMulti(vfs.Path(), 5, 128*bytesx.KiB, bytesx.MiB, metainfo.OptionPieceLength(bytesx.MiB))
		require.NoError(t, err)

		md, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(vfs.Path())), torrent.OptionDisplayName(info.Name))
		require.NoError(t, err)

		mimetype, buf, err := media.PublishRequest(ctx, md, &media.PublishedUploadRequest{
			Entropy:  uuidx.WithSuffix(16),
			Mimetype: mimex.RetrovibedMediaArchive,
			Ttl:      uint64(time.Minute),
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestContext(
			ctx,
			http.MethodPost,
			"/publish",
			buf,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_metadata"))(t))

		routes.ServeHTTP(resp, req)
		require.NoError(t, buf.Close())

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &r))
		require.Equal(t, info.Name, r.Published.Description)
		require.Equal(t, md.ID.String(), r.Published.Id)
		require.NotEqual(t, grpcx.EncodeTime(langx.Autoderef(timex.JSONSafeEncode(new(timex.Inf())))), r.Published.ExpiresAt)
		require.Equal(t, mimex.RetrovibedMediaArchive, r.Published.Mimetype)
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_metadata"))(t))
		require.Equal(t, md5x.FormatUUID(md5x.Digest(md.ID.Bytes(), uuid.FromStringOrNil(uuidx.WithSuffix(16)).Bytes())), testx.Must(sqlx.String(ctx, q, "SELECT encryption_seed::text FROM torrents_metadata"))(t))
		require.Equal(t, "", testx.Must(sqlx.String(ctx, q, "SELECT tracker FROM torrents_metadata"))(t))
		require.Equal(t, mimex.RetrovibedMediaArchive, testx.Must(sqlx.String(ctx, q, "SELECT mimetype FROM torrents_metadata"))(t))
		require.EqualValues(t, info.TotalLength(), testx.Must(sqlx.Value[int](ctx, q, "SELECT bytes FROM torrents_metadata"))(t))
		require.EqualValues(t, info.TotalLength(), testx.Must(sqlx.Value[int](ctx, q, "SELECT available FROM torrents_metadata"))(t))
		require.EqualValues(t, 0, testx.Must(sqlx.Value[int](ctx, q, "SELECT downloaded FROM torrents_metadata"))(t), "published content was never fetched from peers")
		require.NoError(t, tracking.MetadataFindByInfohash(t.Context(), q, r.Published.Id).Scan(&found))
		require.True(t, found.Seeding)
		require.WithinDuration(t, time.Now().Add(time.Minute), found.ExpiresAt, 500*time.Millisecond)
		require.WithinDuration(t, time.Now(), found.CompletedAt, time.Second)

		path := dvfs.Path("torrent", md.ID.String())
		require.DirExists(t, path)
		require.FileExists(t, fmt.Sprintf("%s.torrent", path))

		writtenMeta, err := metainfo.LoadFromFile(fmt.Sprintf("%s.torrent", path))
		require.NoError(t, err)
		require.Equal(t, md.ID.AsByteArray(), [20]byte(writtenMeta.HashInfoBytes()), "written torrent file should have matching infohash")

		// ensure we can update the ttl when necessary.
		mimetype, buf, err = media.PublishRequest(ctx, md, &media.PublishedUploadRequest{
			Entropy:  uuidx.WithSuffix(16),
			Mimetype: mimex.RetrovibedMediaArchive,
			Ttl:      uint64(time.Hour),
		})
		require.NoError(t, err)

		resp, req, err = httptestx.BuildRequestContext(
			ctx,
			http.MethodPost,
			"/publish",
			buf,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, buf.Close())

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &r))

		require.NoError(t, tracking.MetadataFindByInfohash(t.Context(), q, r.Published.Id).Scan(&found))
		require.True(t, found.Seeding)
		require.WithinDuration(t, time.Now().Add(time.Hour), found.ExpiresAt, 500*time.Millisecond)
	})
}
