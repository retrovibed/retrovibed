package media_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestRecentRecord(t *testing.T) {
	t.Run("records media id", func(t *testing.T) {
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

		body, err := json.Marshal(&media.RecentRecordRequest{
			Media: &media.Media{Id: md.ID},
			Query: &media.MediaSearchRequest{Query: "jazz"},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		searchreq := &media.RecentSearchRequest{Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)), Limit: 100}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp2, req2, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/?"+encoded.Encode(), nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp2, req2)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp2.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp2.Body, &result))
		require.Len(t, result.Items, 1)
		require.Equal(t, md.ID, result.Items[0].Media.Id)
	})

	t.Run("records duration in milliseconds", func(t *testing.T) {
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

		body, err := json.Marshal(&media.RecentRecordRequest{
			Media:    &media.Media{Id: md.ID},
			Query:    &media.MediaSearchRequest{Query: "jazz"},
			Duration: 7500,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		searchreq := &media.RecentSearchRequest{Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)), Limit: 100}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp2, req2, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/?"+encoded.Encode(), nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp2, req2)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp2.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp2.Body, &result))
		require.Len(t, result.Items, 1)
		require.EqualValues(t, 7500, result.Items[0].Duration)
	})

	t.Run("records position in milliseconds", func(t *testing.T) {
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

		body, err := json.Marshal(&media.RecentRecordRequest{
			Media:    &media.Media{Id: md.ID},
			Query:    &media.MediaSearchRequest{Query: "jazz"},
			Position: 4200,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		searchreq := &media.RecentSearchRequest{Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)), Limit: 100}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp2, req2, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/?"+encoded.Encode(), nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp2, req2)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp2.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp2.Body, &result))
		require.Len(t, result.Items, 1)
		require.EqualValues(t, 4200, result.Items[0].Position)
	})

	t.Run("records query fields", func(t *testing.T) {
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

		body, err := json.Marshal(&media.RecentRecordRequest{
			Media: &media.Media{Id: md.ID},
			Query: &media.MediaSearchRequest{Query: "blues", Limit: 50},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		searchreq := &media.RecentSearchRequest{Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)), Limit: 100}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp2, req2, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/?"+encoded.Encode(), nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp2, req2)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp2.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp2.Body, &result))
		require.Len(t, result.Items, 1)
		require.Equal(t, "blues", result.Items[0].Query.Query)
		require.EqualValues(t, 50, result.Items[0].Query.Limit)
	})

	t.Run("records id", func(t *testing.T) {
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

		body, err := json.Marshal(&media.RecentRecordRequest{
			Media: &media.Media{Id: md.ID},
			Query: &media.MediaSearchRequest{Query: "jazz"},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		searchreq := &media.RecentSearchRequest{Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)), Limit: 100}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp2, req2, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/?"+encoded.Encode(), nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp2, req2)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp2.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp2.Body, &result))
		require.Len(t, result.Items, 1)
		require.NotEmpty(t, result.Items[0].Id)
	})

	t.Run("records mimetype category", func(t *testing.T) {
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

		body, err := json.Marshal(&media.RecentRecordRequest{
			Media:    &media.Media{Id: md.ID},
			Query:    &media.MediaSearchRequest{Query: "jazz"},
			Mimetype: mimex.Audio,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		searchreq := &media.RecentSearchRequest{Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)), Limit: 100}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp2, req2, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/?"+encoded.Encode(), nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp2, req2)

		var result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp2.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp2.Body, &result))
		require.Len(t, result.Items, 1)
		require.Equal(t, mimex.Audio, result.Items[0].Mimetype)
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
			http.MethodPost, "/", []byte(`{not valid json`),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})

	// The session ID is md5(encoded query) with no profile component, so two profiles recording
	// the identical query for the same media previously collided on a single row owned by
	// whichever profile recorded first: the row's PRIMARY KEY was id alone, and the upsert's
	// ON CONFLICT (id) DO UPDATE never touched profile_id. The table's key is now composite
	// (profile_id, id) with a matching ON CONFLICT (id, profile_id) upsert, so identical
	// queries from different profiles land in separate rows instead of colliding.
	t.Run("same query for same media across profiles stays isolated per profile", func(t *testing.T) {
		var (
			p1, p2         meta.Profile
			authz1, authz2 meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p1, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p1).Scan(&p1))
		require.NoError(t, testx.Fake(&authz1, meta.AuthzOptionProfileID(p1.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz1).Scan(&authz1))

		require.NoError(t, testx.Fake(&p2, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p2).Scan(&p2))
		require.NoError(t, testx.Fake(&authz2, meta.AuthzOptionProfileID(p2.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz2).Scan(&authz2))

		var md library.Metadata
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		routes := mux.NewRouter()
		media.NewHTTPRecent(q, media.HTTPRecentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())

		claims1 := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p1.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz1)))
		claims2 := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p2.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz2)))

		body, err := json.Marshal(&media.RecentRecordRequest{
			Media: &media.Media{Id: md.ID},
			Query: &media.MediaSearchRequest{Query: "jazz", Limit: 25},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims1, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)

		resp2, req2, err := httptestx.BuildRequestBytes(
			http.MethodPost, "/", body,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims2, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp2, req2)
		require.Equal(t, http.StatusOK, resp2.Result().StatusCode)

		searchreq := &media.RecentSearchRequest{Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)), Limit: 100}
		encoded, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp3, req3, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/?"+encoded.Encode(), nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims1, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp3, req3)

		var p1Result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp3.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp3.Body, &p1Result))
		require.Len(t, p1Result.Items, 1, "profile 1 keeps its own session for the query")
		require.Equal(t, md.ID, p1Result.Items[0].Media.Id)

		resp4, req4, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/?"+encoded.Encode(), nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims2, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp4, req4)

		var p2Result media.RecentSearchResponse
		require.Equal(t, http.StatusOK, resp4.Result().StatusCode)
		require.NoError(t, jsonx.UnmarshalRead(resp4.Body, &p2Result))
		require.Len(t, p2Result.Items, 1, "profile 2 gets its own session for the identical query rather than losing it to profile 1")
		require.Equal(t, md.ID, p2Result.Items[0].Media.Id)

		// the content-addressed id column is identical for both (same query bytes hash to the
		// same value); ownership is distinguished by profile_id, part of the composite key.
		require.Equal(t, p1Result.Items[0].Id, p2Result.Items[0].Id)

		count, err := sqlx.Count(ctx, q, `SELECT COUNT(*) FROM library_recent_sessions`)
		require.NoError(t, err)
		require.Equal(t, 2, count, "one row per profile, not a shared/collided row")
	})
}
