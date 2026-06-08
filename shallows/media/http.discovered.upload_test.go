package media_test

import (
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/atomicx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredUploadTorrent(t *testing.T) {
	var (
		p meta.Profile
		v meta.Authz
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
		media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		media.HTTPDiscoveredOptionRootStorage(vfs),
	).Bind(routes.PathPrefix("/").Subrouter())

	mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
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

	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

	resp, req, err := httptestx.BuildRequestBytes(
		http.MethodPost,
		"/",
		testx.IOBytes(buf),
		httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		httptestx.RequestOptionHeader("Content-Type", mimetype),
	)
	require.NoError(t, err)

	require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_metadata"))(t))

	routes.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Result().StatusCode)
	require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_metadata"))(t))
	require.Equal(t, "https://example.com/announce", testx.Must(sqlx.String(ctx, q, "SELECT tracker FROM torrents_metadata"))(t))
}
