package media_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

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

func TestDiscoveredPause(t *testing.T) {
	t.Run("successful pause - stops the torrent", func(t *testing.T) {
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
			fmt.Sprintf("/%s/pause", tmd.ID),
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
		assert.EqualValues(t, tmd.Downloaded, latest.Downloaded)
		assert.EqualValues(t, tmd.Available, latest.Available)
		assert.EqualValues(t, tmd.Uploaded, latest.Uploaded)
		assert.WithinDuration(t, tmd.CreatedAt, latest.CreatedAt, 0)
		assert.LessOrEqual(t, tmd.UpdatedAt, latest.UpdatedAt)
		assert.WithinDuration(t, timex.Inf(), latest.InitiatedAt, 0)
		assert.WithinDuration(t, time.Now(), latest.PausedAt, 5*time.Second)
	})

	t.Run("pause non-existent returns not found", func(t *testing.T) {
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
			"/non-existent-id/pause",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("should require authentication", func(t *testing.T) {
		var (
			p   meta.Profile
			tmd tracking.Metadata
		)
		ctx, done := testx.Context(t)
		defer done()

		tclient := torrenttestx.QuickClient(t)
		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

		require.NoError(t, testx.Fake(&tmd, tracking.MetadataOptionTestDefaults))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd))

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

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodDelete,
			fmt.Sprintf("/%s/pause", tmd.ID),
			nil,
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode)
	})
}
