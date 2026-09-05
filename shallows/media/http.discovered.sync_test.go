package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttestx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/atomicx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestDownloadSyncMetadata(t *testing.T) {
	var (
		p     meta.Profile
		authz meta.Authz
		m     library.Known
		v0    library.Metadata
		v1    library.Metadata
		v2    library.Metadata
		tmd   tracking.Metadata
	)
	ctx, done := testx.Context(t)
	defer done()

	tclient := torrenttestx.QuickClient(t)
	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
	require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

	require.NoError(t, testx.Fake(&m, library.KnownOptionTestDefaults, library.KnownOptionRandomID))
	require.NoError(t, library.KnownInsertWithDefaults(ctx, q, m).Scan(&m))

	require.NoError(t, testx.Fake(&tmd, tracking.MetadataOptionTestDefaults))
	require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd))

	require.NoError(t, testx.Fake(&v0, library.MetadataOptionTestDefaults, library.MetadataOptionTestID(uuidx.WithSuffix(1)), library.MetadataOptionTorrentID(tmd.ID)))
	require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, v0).Scan(&v0))
	require.NoError(t, testx.Fake(&v1, library.MetadataOptionTestDefaults, library.MetadataOptionTestID(uuidx.WithSuffix(2)), library.MetadataOptionTorrentID(tmd.ID)))
	require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, v1).Scan(&v1))
	require.NoError(t, testx.Fake(&v2, library.MetadataOptionTestDefaults, library.MetadataOptionTestID(uuidx.WithSuffix(3))))
	require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, v2).Scan(&v2))

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

	tmd.KnownMediaID = m.ID
	encoded, err := json.Marshal(media.MetadataSyncRequest{
		Media: new(langx.Clone(media.Media{}, media.MediaOptionFromTorrentMetadata(tmd))),
	})
	require.NoError(t, err)

	resp, req, err := httptestx.BuildRequestBytes(
		http.MethodPost,
		fmt.Sprintf("/%s/metadatasync", tmd.ID),
		encoded,
		httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
	)
	require.NoError(t, err)

	require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata WHERE known_media_id = ?", m.ID))(t))
	require.Equal(t, 2, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata WHERE torrent_id = ?", tmd.ID))(t))
	require.Equal(t, uuid.Max.String(), testx.Must(sqlx.String(ctx, q, "SELECT known_media_id::text FROM torrents_metadata WHERE id = ?", tmd.ID))(t))

	routes.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Result().StatusCode)
	require.Equal(t, m.ID, testx.Must(sqlx.String(ctx, q, "SELECT known_media_id::text FROM torrents_metadata WHERE id = ?", tmd.ID))(t))
	require.Equal(t, 2, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata WHERE known_media_id = ?", m.ID))(t))

}
