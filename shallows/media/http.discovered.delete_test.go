package media_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttestx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/atomicx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredDelete(t *testing.T) {
	t.Run("successful delete - resets the torrent", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
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

		require.NoError(t, testx.Fake(&tmd, tracking.MetadataOptionTestDefaults))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd))
		require.Greater(t, tmd.Bytes, uint64(0))
		require.Greater(t, tmd.Downloaded, uint64(0))
		require.Greater(t, tmd.Uploaded, uint64(0))

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
			http.MethodDelete,
			fmt.Sprintf("/%s", tmd.ID),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var (
			latest tracking.Metadata
		)

		require.NoError(t, tracking.MetadataFindByID(t.Context(), q, tmd.ID).Scan(&latest))
		assert.EqualValues(t, tmd.Bytes, latest.Bytes)
		assert.EqualValues(t, 0, latest.Downloaded)
		assert.EqualValues(t, 0, latest.Available)
		assert.EqualValues(t, 0, latest.Uploaded)
		assert.EqualValues(t, tmd.Archivable, tmd.Archivable)
		assert.WithinDuration(t, tmd.CreatedAt, latest.CreatedAt, 0)
		assert.LessOrEqual(t, tmd.UpdatedAt, latest.UpdatedAt)
		assert.LessOrEqual(t, tmd.InitiatedAt, latest.InitiatedAt)
		assert.WithinDuration(t, timex.Inf(), latest.InitiatedAt, 0)
		assert.WithinDuration(t, timex.Inf(), latest.PausedAt, 0)
		assert.WithinDuration(t, timex.Inf(), latest.CompletedAt, 0)
		assert.WithinDuration(t, timex.Inf(), latest.HiddenAt, 0)
		assert.WithinDuration(t, timex.Inf(), latest.VerifyAt, 0)
		assert.WithinDuration(t, timex.NegInf(), latest.NextAnnounceAt, 0)
	})

	t.Run("delete non-existent returns not found", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		tclient := torrenttestx.QuickClient(t)
		q := sqltestx.Metadatabase(t)

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
			http.MethodDelete,
			"/non-existent-id",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})
}
