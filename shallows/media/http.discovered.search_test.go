package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/atomicx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredSearch(t *testing.T) {
	t.Run("standard search", func(t *testing.T) {
		// ensure that search doesnt return an error.
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

		encoder := formx.NewEncoder()
		query, err := encoder.Encode(media.MediaSearchRequest{
			Limit: 100,
			Query: "derp",
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/available?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
	})

	t.Run("filter by query", func(t *testing.T) {
		var (
			p         meta.Profile
			v         meta.Authz
			mdMatch   tracking.Metadata
			mdNoMatch tracking.Metadata
			result    media.DownloadSearchResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		// Metadata that should match the query
		require.NoError(t, testx.Fake(&mdMatch,
			tracking.MetadataOptionTestDefaults,
			tracking.MetadataOptionBytes(100),
			tracking.MetadataOptionDownloaded(50),
			tracking.MetadataOptionDescription("unique_query_term_video"),
			tracking.MetadataOptionAutoDescription,
			func(t *tracking.Metadata) { t.InitiatedAt = time.Now() },
		))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, mdMatch).Scan(&mdMatch))
		require.NoError(t, tracking.MetadataDownloadByID(ctx, q, mdMatch.ID).Scan(&mdMatch))

		// Metadata that should NOT match the query
		require.NoError(t, testx.Fake(&mdNoMatch,
			tracking.MetadataOptionTestDefaults,
			tracking.MetadataOptionBytes(100),
			tracking.MetadataOptionDownloaded(50),
			tracking.MetadataOptionDescription("another_item_audio"),
			tracking.MetadataOptionAutoDescription,
			func(t *tracking.Metadata) { t.InitiatedAt = time.Now() },
		))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, mdNoMatch).Scan(&mdNoMatch))
		require.NoError(t, tracking.MetadataDownloadByID(ctx, q, mdNoMatch.ID).Scan(&mdNoMatch))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(vfs.Path(), storage.FileOptionPathMakerInfohash),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()

		// Request with a specific query
		query, err := encoder.Encode(media.DownloadSearchRequest{
			Limit: 100,
			Query: "unique_query_term",
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/available?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		err = json.NewDecoder(resp.Result().Body).Decode(&result)
		require.NoError(t, err)

		require.Len(t, result.Items, 1, "Expected 1 items in the search results")
		require.Equal(t, mdMatch.ID, result.Items[0].Media.Id, "Expected item to match the query")
		initiatedAt, err := grpcx.DecodeTime(result.Items[0].InitiatedAt)
		require.NoError(t, err)
		require.WithinDuration(t, time.Now(), initiatedAt, time.Second)
	})

	t.Run("include records that dont have a known size", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			md     tracking.Metadata
			result media.DownloadSearchResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, testx.Fake(&md,
			tracking.MetadataOptionTestDefaults,
			tracking.MetadataOptionBytes(0),
			tracking.MetadataOptionDownloaded(0),
		))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))
		require.NoError(t, tracking.MetadataDownloadByID(ctx, q, md.ID).Scan(&md))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(vfs.Path(), storage.FileOptionPathMakerInfohash),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()

		query, err := encoder.Encode(media.DownloadSearchRequest{
			Limit: 100,
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/available?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		err = json.NewDecoder(resp.Result().Body).Decode(&result)
		require.NoError(t, err)

		require.Len(t, result.Items, 1, "Expected only 1 incomplete metadata item in the search results")
		require.Equal(t, md.ID, result.Items[0].Media.Id, "Expected the completed item to be returned")
		initiatedAt, err := grpcx.DecodeTime(result.Items[0].InitiatedAt)
		require.NoError(t, err)
		require.WithinDuration(t, time.Now(), initiatedAt, time.Second)
	})

	t.Run("hidden torrents not returned by default", func(t *testing.T) {
		var (
			p         meta.Profile
			v         meta.Authz
			mdVisible tracking.Metadata
			mdHidden  tracking.Metadata
			result    media.DownloadSearchResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, testx.Fake(&mdVisible, tracking.MetadataOptionTestDefaults))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, mdVisible).Scan(&mdVisible))

		require.NoError(t, testx.Fake(&mdHidden,
			tracking.MetadataOptionTestDefaults,
			func(m *tracking.Metadata) { m.HiddenAt = time.Now() },
		))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, mdHidden).Scan(&mdHidden))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(vfs.Path(), storage.FileOptionPathMakerInfohash),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()

		query, err := encoder.Encode(media.DownloadSearchRequest{
			Limit: 100,
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/available?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		err = json.NewDecoder(resp.Result().Body).Decode(&result)
		require.NoError(t, err)

		require.Len(t, result.Items, 1, "Expected only the visible item in search results")
		require.Equal(t, mdVisible.ID, result.Items[0].Media.Id, "Expected the hidden item to be excluded")
	})

	t.Run("filter by completed=true", func(t *testing.T) {
		var (
			p            meta.Profile
			v            meta.Authz
			mdIncomplete tracking.Metadata
			mdCompleted  tracking.Metadata
			result       media.DownloadSearchResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, testx.Fake(&mdIncomplete,
			tracking.MetadataOptionTestDefaults,
			tracking.MetadataOptionBytes(100),
			tracking.MetadataOptionDownloaded(50),
		))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, mdIncomplete).Scan(&mdIncomplete))
		require.NoError(t, tracking.MetadataDownloadByID(ctx, q, mdIncomplete.ID).Scan(&mdIncomplete))

		require.NoError(t, testx.Fake(&mdCompleted,
			tracking.MetadataOptionTestDefaults,
			tracking.MetadataOptionBytes(100),
		))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, mdCompleted).Scan(&mdCompleted))
		require.NoError(t, tracking.MetadataCompleteByID(ctx, q, mdCompleted.ID, 0, mdCompleted.Bytes, mdCompleted.Bytes, 0).Scan(&mdCompleted))
		require.EqualValues(t, mdCompleted.Bytes, mdCompleted.Downloaded)
		require.WithinDuration(t, time.Now(), mdCompleted.CompletedAt, time.Second)

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(vfs.Path(), storage.FileOptionPathMakerInfohash),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()

		query, err := encoder.Encode(media.DownloadSearchRequest{
			Limit:     100,
			Completed: true,
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/available?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		err = json.NewDecoder(resp.Result().Body).Decode(&result)
		require.NoError(t, err)

		require.Len(t, result.Items, 1, "Expected only 1 incomplete metadata item in the search results")
		require.Equal(t, mdCompleted.ID, result.Items[0].Media.Id, "Expected the completed item to be returned")
		require.Equal(t, grpcx.EncodeTime(timex.RFC3339NanoEncode(timex.Inf())), result.Items[0].InitiatedAt, "Expected completed item without download to have unset initiated_at")
	})

	t.Run("filter by hidden=true", func(t *testing.T) {
		var (
			p         meta.Profile
			v         meta.Authz
			mdVisible tracking.Metadata
			mdHidden  tracking.Metadata
			result    media.DownloadSearchResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, testx.Fake(&mdVisible, tracking.MetadataOptionTestDefaults))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, mdVisible).Scan(&mdVisible))

		require.NoError(t, testx.Fake(&mdHidden,
			tracking.MetadataOptionTestDefaults,
			func(m *tracking.Metadata) { m.HiddenAt = time.Now() },
		))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, mdHidden).Scan(&mdHidden))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPDiscovered(
			q,
			atomicx.PointerPtr(tclient),
			storage.NewFile(vfs.Path(), storage.FileOptionPathMakerInfohash),
			media.HTTPDiscoveredOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		encoder := formx.NewEncoder()

		query, err := encoder.Encode(media.DownloadSearchRequest{
			Limit:  100,
			Hidden: true,
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/available?%s", query.Encode()),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		err = json.NewDecoder(resp.Result().Body).Decode(&result)
		require.NoError(t, err)

		require.Len(t, result.Items, 1, "Expected only the hidden item in search results")
		require.Equal(t, mdHidden.ID, result.Items[0].Media.Id, "Expected the hidden item to be returned")
	})
}
