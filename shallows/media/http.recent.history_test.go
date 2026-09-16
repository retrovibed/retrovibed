package media_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestRecentHistory(t *testing.T) {
	t.Run("upserts a new watch history row", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		routes := mux.NewRouter()
		media.NewHTTPRecent(q, media.HTTPRecentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		playbackID := uuidx.WithSuffix(1)
		body, err := json.Marshal(&library.WatchHistoryRecordRequest{
			Record: &library.WatchHistoryRecord{
				Id:       playbackID,
				MediaId:  md.ID,
				Duration: 5000,
			},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/history", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var row library.WatchHistory
		scanner := library.WatchHistorySearch(ctx, q, library.WatchHistorySearchBuilder().Where(library.WatchHistoryQueryProfileID(p.ID)))
		require.True(t, scanner.Next())
		require.NoError(t, scanner.Scan(&row))
		require.NoError(t, scanner.Err())
		require.NoError(t, scanner.Close())

		require.Equal(t, playbackID, row.ID)
		require.Equal(t, md.ID, row.MediaID)
		require.Equal(t, p.ID, row.ProfileID)
		require.Equal(t, 5*time.Second, row.Duration)
	})

	t.Run("same record id updates duration in place", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		routes := mux.NewRouter()
		media.NewHTTPRecent(q, media.HTTPRecentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		playbackID := uuidx.WithSuffix(1)

		heartbeat := func(duration uint64) {
			body, err := json.Marshal(&library.WatchHistoryRecordRequest{
				Record: &library.WatchHistoryRecord{
					Id:       playbackID,
					MediaId:  md.ID,
					Duration: duration,
				},
			})
			require.NoError(t, err)

			resp, req, err := httptestx.BuildRequestBytes(
				http.MethodPost, "/history", body,
				httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			)
			require.NoError(t, err)
			routes.ServeHTTP(resp, req)
			require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		}

		heartbeat(1000)
		heartbeat(9000)

		var row library.WatchHistory
		scanner := library.WatchHistorySearch(ctx, q, library.WatchHistorySearchBuilder().Where(library.WatchHistoryQueryProfileID(p.ID)))
		require.True(t, scanner.Next())
		require.NoError(t, scanner.Scan(&row))
		require.False(t, scanner.Next(), "repeat heartbeats for the same playback id must upsert a single row")
		require.NoError(t, scanner.Err())
		require.NoError(t, scanner.Close())

		require.Equal(t, 9*time.Second, row.Duration)
	})

	t.Run("different record id creates a separate row", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		routes := mux.NewRouter()
		media.NewHTTPRecent(q, media.HTTPRecentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		for _, id := range []string{uuidx.WithSuffix(1), uuidx.WithSuffix(2)} {
			body, err := json.Marshal(&library.WatchHistoryRecordRequest{
				Record: &library.WatchHistoryRecord{
					Id:       id,
					MediaId:  md.ID,
					Duration: 1000,
				},
			})
			require.NoError(t, err)

			resp, req, err := httptestx.BuildRequestBytes(
				http.MethodPost, "/history", body,
				httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			)
			require.NoError(t, err)
			routes.ServeHTTP(resp, req)
			require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		}

		count := 0
		scanner := library.WatchHistorySearch(ctx, q, library.WatchHistorySearchBuilder().Where(library.WatchHistoryQueryProfileID(p.ID)))
		for scanner.Next() {
			var row library.WatchHistory
			require.NoError(t, scanner.Scan(&row))
			count++
		}
		require.NoError(t, scanner.Err())
		require.NoError(t, scanner.Close())

		require.Equal(t, 2, count, "a rewatch (new playback id) must not overwrite the earlier row")
	})

	t.Run("malformed json returns 400", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()
		media.NewHTTPRecent(q, media.HTTPRecentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/history", []byte(`{not valid json`),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})
}
