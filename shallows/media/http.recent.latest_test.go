package media_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestRecentLatest(t *testing.T) {
	t.Run("zero state", func(t *testing.T) {
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

		searchreq := &media.RecentSearchRequest{
			Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)),
			Limit:   100,
		}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+encoded.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Empty(t, result.Items)
	})

	t.Run("returns results", func(t *testing.T) {
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

		for idx := range 3 {
			var md library.Metadata
			require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

			var session library.RecentSession
			require.NoError(t, testx.Fake(&session, library.RecentSessionOptionTestDefaults, library.RecentSessionOptionID(uuidx.WithSuffix(idx)), library.RecentSessionOptionMediaID(md.ID)))
			require.NoError(t, library.RecentSessionInsertWithDefaults(ctx, q, session).Scan(&session))
		}

		routes := mux.NewRouter()
		media.NewHTTPRecent(q, media.HTTPRecentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		searchreq := &media.RecentSearchRequest{
			Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)),
			Limit:   100,
		}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+encoded.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 3)
	})

	t.Run("returns Duration Position Query", func(t *testing.T) {
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

		encodedQuery, err := proto.Marshal(&media.MediaSearchRequest{Query: "rock", Limit: 25})
		require.NoError(t, err)

		var session library.RecentSession
		require.NoError(t, testx.Fake(&session, library.RecentSessionOptionTestDefaults, library.RecentSessionOptionMediaID(md.ID)))
		session.Duration = 5000 * time.Millisecond
		session.Position = 3000 * time.Millisecond
		session.Query = encodedQuery
		require.NoError(t, library.RecentSessionInsertWithDefaults(ctx, q, session).Scan(&session))

		routes := mux.NewRouter()
		media.NewHTTPRecent(q, media.HTTPRecentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		searchreq := &media.RecentSearchRequest{
			Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)),
			Limit:   100,
		}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+encoded.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
		require.Equal(t, md.ID, result.Items[0].Media.Id)
		require.EqualValues(t, 5000, result.Items[0].Duration)
		require.EqualValues(t, 3000, result.Items[0].Position)
		require.Equal(t, "rock", result.Items[0].Query.Query)
		require.EqualValues(t, 25, result.Items[0].Query.Limit)
	})

	t.Run("no duplicates for same query", func(t *testing.T) {
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

		// record two sessions with the same query but different media IDs
		for range 2 {
			var md library.Metadata
			require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

			body, err := json.Marshal(&media.RecentRecordRequest{
				Media: &media.Media{Id: md.ID},
				Query: &media.MediaSearchRequest{Query: "rock", Limit: 25},
			})
			require.NoError(t, err)

			resp, req, err := httptestx.BuildRequestBytes(
				http.MethodPost,
				"/",
				body,
				httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			)
			require.NoError(t, err)
			routes.ServeHTTP(resp, req)
			require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		}

		searchreq := &media.RecentSearchRequest{
			Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)),
			Limit:   100,
		}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+encoded.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
		require.Equal(t, "rock", result.Items[0].Query.Query)

		count, err := sqlx.Count(ctx, q, `SELECT COUNT(*) FROM library_recent_sessions`)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("filters by mimetype", func(t *testing.T) {
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

		for idx := range 3 {
			var md library.Metadata
			require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

			var session library.RecentSession
			require.NoError(t, testx.Fake(&session, library.RecentSessionOptionTestDefaults, library.RecentSessionOptionID(uuidx.WithSuffix(idx)), library.RecentSessionOptionMediaID(md.ID)))
			session.Mimetype = mimex.Video
			require.NoError(t, library.RecentSessionInsertWithDefaults(ctx, q, session).Scan(&session))
		}

		for idx := range 3 {
			var md library.Metadata
			require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

			var session library.RecentSession
			require.NoError(t, testx.Fake(&session, library.RecentSessionOptionTestDefaults, library.RecentSessionOptionID(uuidx.WithSuffix(10+idx)), library.RecentSessionOptionMediaID(md.ID)))
			session.Mimetype = mimex.Audio
			require.NoError(t, library.RecentSessionInsertWithDefaults(ctx, q, session).Scan(&session))
		}

		routes := mux.NewRouter()
		media.NewHTTPRecent(q, media.HTTPRecentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		searchreq := &media.RecentSearchRequest{
			Created:  meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)),
			Mimetype: mimex.Video,
			Limit:    100,
		}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+encoded.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 3)
		for _, item := range result.Items {
			require.Equal(t, mimex.Video, item.Mimetype)
		}
	})

	t.Run("empty mimetype returns all", func(t *testing.T) {
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

		for idx := range 3 {
			var md library.Metadata
			require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

			var session library.RecentSession
			require.NoError(t, testx.Fake(&session, library.RecentSessionOptionTestDefaults, library.RecentSessionOptionID(uuidx.WithSuffix(idx)), library.RecentSessionOptionMediaID(md.ID)))
			session.Mimetype = mimex.Video
			require.NoError(t, library.RecentSessionInsertWithDefaults(ctx, q, session).Scan(&session))
		}

		for idx := range 3 {
			var md library.Metadata
			require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

			var session library.RecentSession
			require.NoError(t, testx.Fake(&session, library.RecentSessionOptionTestDefaults, library.RecentSessionOptionID(uuidx.WithSuffix(10+idx)), library.RecentSessionOptionMediaID(md.ID)))
			session.Mimetype = mimex.Audio
			require.NoError(t, library.RecentSessionInsertWithDefaults(ctx, q, session).Scan(&session))
		}

		routes := mux.NewRouter()
		media.NewHTTPRecent(q, media.HTTPRecentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		searchreq := &media.RecentSearchRequest{
			Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)),
			Limit:   100,
		}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+encoded.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 6)
	})

	t.Run("tombstoned excluded", func(t *testing.T) {
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

		var alive, tombstoned library.Metadata
		require.NoError(t, testx.Fake(&alive, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, alive).Scan(&alive))
		var aliveSession library.RecentSession
		require.NoError(t, testx.Fake(&aliveSession, library.RecentSessionOptionTestDefaults, library.RecentSessionOptionMediaID(alive.ID)))
		require.NoError(t, library.RecentSessionInsertWithDefaults(ctx, q, aliveSession).Scan(new(library.RecentSession)))

		require.NoError(t, testx.Fake(&tombstoned, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, tombstoned).Scan(&tombstoned))
		var tombstonedSession library.RecentSession
		require.NoError(t, testx.Fake(&tombstonedSession, library.RecentSessionOptionTestDefaults, library.RecentSessionOptionMediaID(tombstoned.ID)))
		require.NoError(t, library.RecentSessionInsertWithDefaults(ctx, q, tombstonedSession).Scan(new(library.RecentSession)))
		require.NoError(t, library.MetadataTombstoneByID(ctx, q, tombstoned.ID).Scan(&tombstoned))

		routes := mux.NewRouter()
		media.NewHTTPRecent(q, media.HTTPRecentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		searchreq := &media.RecentSearchRequest{
			Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)),
			Limit:   100,
		}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/?"+encoded.Encode(),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
		require.Equal(t, alive.ID, result.Items[0].Media.Id)
	})
}
