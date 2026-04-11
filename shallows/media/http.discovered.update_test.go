package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/atomicx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/grpcx"
	"github.com/retrovibed/retrovibed/internal/httptestx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/media"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/retrovibed/retrovibed/tracking"
	"github.com/stretchr/testify/require"
)

func TestHTTPDiscoveredUpdate(t *testing.T) {
	t.Run("successfully updates metadata", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			md     tracking.Metadata
			result media.DownloadUpdateResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		tclient := torrenttestx.QuickClient(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))
		require.NoError(t, testx.Fake(&md, tracking.MetadataOptionTestDefaults))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(vfs.Path(), storage.FileOptionPathMakerInfohash),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		pausedAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
		verifyAt := time.Now().Add(2 * time.Hour).UTC().Truncate(time.Second)
		completedAt := time.Now().UTC().Truncate(time.Second)
		initiatedAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPut,
			fmt.Sprintf("/%s", md.ID),
			testx.Must(json.Marshal(&media.DownloadUpdateRequest{
				Download: &media.Download{
					PausedAt:    grpcx.EncodeTime(pausedAt),
					VerifyAt:    grpcx.EncodeTime(verifyAt),
					CompletedAt: grpcx.EncodeTime(completedAt),
					InitiatedAt: grpcx.EncodeTime(initiatedAt),
				},
			}))(t),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		err = json.NewDecoder(resp.Result().Body).Decode(&result)
		require.NoError(t, err)
		require.NotNil(t, result.Download)
		require.Equal(t, md.ID, result.Download.Media.Id)

		// Verify verify_at is returned in HTTP response
		decodedVerifyAt, err := grpcx.DecodeTime(result.Download.VerifyAt)
		require.NoError(t, err)
		require.Equal(t, verifyAt, decodedVerifyAt)

		// Verify verify_at is properly updated in SQL database
		var dbMd tracking.Metadata
		require.NoError(t, tracking.MetadataFindByID(ctx, q, md.ID).Scan(&dbMd))
		require.Equal(t, verifyAt, dbMd.VerifyAt)

		// Verify paused_at is returned in HTTP response
		decodedPausedAt, err := grpcx.DecodeTime(result.Download.PausedAt)
		require.NoError(t, err)
		require.Equal(t, pausedAt, decodedPausedAt)

		// Verify paused_at is properly updated in SQL database
		require.Equal(t, pausedAt, dbMd.PausedAt)
	})

	t.Run("returns bad request for invalid JSON body", func(t *testing.T) {
		var (
			p  meta.Profile
			v  meta.Authz
			md tracking.Metadata
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		tclient := torrenttestx.QuickClient(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))
		require.NoError(t, testx.Fake(&md, tracking.MetadataOptionTestDefaults))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(vfs.Path(), storage.FileOptionPathMakerInfohash),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPut,
			fmt.Sprintf("/%s", md.ID),
			[]byte(`{"download": "invalid`),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})
}
