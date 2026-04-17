package media_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/rss"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestRSSFeedSearch(t *testing.T) {
	t.Run("should properly search feeds", func(t *testing.T) {
		var (
			p      meta.Profile
			authz  meta.Authz
			result rss.FeedSearchResponse
			feed   tracking.RSS
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		require.NoError(t, testx.Fake(&feed, tracking.RSSOptionTestDefaults, tracking.RSSOptionDefaultEncryptionSeed))
		require.NoError(t, tracking.RSSInsertWithDefaults(ctx, q, feed).Scan(&feed))
		routes := mux.NewRouter()

		media.NewHTTPRSSFeed(
			q,
			media.HTTPRSSFeedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			"/",
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Len(t, result.Items, 1)
		require.NotEmpty(t, result.Items[0].EncryptionSeed)
	})
}

func TestRSSFeedCreate(t *testing.T) {
	t.Run("should properly create a feed", func(t *testing.T) {
		var (
			p      meta.Profile
			authz  meta.Authz
			result rss.FeedCreateResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()

		media.NewHTTPRSSFeed(
			q,
			media.HTTPRSSFeedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		encoded, err := json.Marshal(rss.FeedCreateRequest{
			Feed: &rss.Feed{
				Url:          "https://example.com/71b9dbed-e985-42df-ba18-82dbf5779358",
				Description:  "hello world",
				Autodownload: true,
				Autoarchive:  true,
				Contributing: false,
			},
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			encoded,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "d81f2387-d17a-cb19-cfb5-a1f3830589bd", result.Feed.Id)
		require.Equal(t, "hello world", result.Feed.Description)
		require.WithinDuration(t, time.Now(), errorsx.Must(grpcx.DecodeTime(result.Feed.NextCheck)), time.Second)
		require.True(t, result.Feed.Autodownload)
		require.True(t, result.Feed.Autoarchive)
		require.False(t, result.Feed.Contributing)
	})

	t.Run("should properly update a feed", func(t *testing.T) {
		var (
			p      meta.Profile
			authz  meta.Authz
			result rss.FeedCreateResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()

		media.NewHTTPRSSFeed(
			q,
			media.HTTPRSSFeedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		encoded, err := json.Marshal(rss.FeedCreateRequest{
			Feed: &rss.Feed{
				Url:          "https://example.com/71b9dbed-e985-42df-ba18-82dbf5779358",
				Description:  "hello world",
				Autodownload: true,
				Autoarchive:  true,
				Contributing: false,
				Digest:       uuid.Nil.String(),
			},
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			encoded,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "d81f2387-d17a-cb19-cfb5-a1f3830589bd", result.Feed.Id)
		require.Equal(t, "hello world", result.Feed.Description)
		require.True(t, result.Feed.Autodownload)
		require.True(t, result.Feed.Autoarchive)
		require.False(t, result.Feed.Contributing)
		require.Equal(t, uuid.Nil.String(), result.Feed.Digest)

		// clone and update the next check timestamp
		ts := time.Now()
		tmp := proto.CloneOf(&result)
		tmp.Feed.Description = tmp.Feed.Description + " 2"
		tmp.Feed.NextCheck = grpcx.EncodeTime(timex.RFC3339NanoEncode(ts.Truncate(time.Second).UTC()))
		tmp.Feed.Digest = errorsx.Must(uuid.NewV7()).String()

		encoded, err = json.Marshal(tmp)
		require.NoError(t, err)

		resp, req, err = httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			encoded,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "d81f2387-d17a-cb19-cfb5-a1f3830589bd", result.Feed.Id)
		require.Equal(t, "hello world 2", result.Feed.Description)
		require.Equal(t, tmp.Feed.NextCheck, result.Feed.NextCheck)
		require.True(t, result.Feed.Autodownload)
		require.True(t, result.Feed.Autoarchive)
		require.False(t, result.Feed.Contributing)
		require.Equal(t, tmp.Feed.Digest, result.Feed.Digest)
	})

	t.Run("should generate a consistent EncryptionSeed based on the url if its not provided", func(t *testing.T) {
		var (
			p      meta.Profile
			authz  meta.Authz
			result rss.FeedUpdateResponse
			feed   tracking.RSS
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()

		media.NewHTTPRSSFeed(
			q,
			media.HTTPRSSFeedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		encoded, err := json.Marshal(rss.FeedCreateRequest{
			Feed: &rss.Feed{
				Url:          "https://example.com/71b9dbed-e985-42df-ba18-82dbf5779358",
				Description:  "hello world",
				Autodownload: true,
				Autoarchive:  true,
				Contributing: false,
			},
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			encoded,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "d81f2387-d17a-cb19-cfb5-a1f3830589bd", result.Feed.Id)
		require.Equal(t, "hello world", result.Feed.Description)
		require.True(t, result.Feed.Autodownload)
		require.True(t, result.Feed.Autoarchive)
		require.False(t, result.Feed.Contributing)
		require.NoError(t, tracking.RSSFindByID(ctx, q, result.Feed.Id).Scan(&feed))
		require.Equal(t, "5ababd60-3b22-7803-02dd-8d83498e5172", feed.EncryptionSeed)
	})

	t.Run("should use the encryption seed if provided", func(t *testing.T) {
		var (
			p      meta.Profile
			authz  meta.Authz
			result rss.FeedUpdateResponse
			feed   tracking.RSS
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		routes := mux.NewRouter()

		media.NewHTTPRSSFeed(
			q,
			media.HTTPRSSFeedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		encoded, err := json.Marshal(rss.FeedCreateRequest{
			Feed: &rss.Feed{
				Url:            "https://example.com/71b9dbed-e985-42df-ba18-82dbf5779358",
				Description:    "hello world",
				Autodownload:   true,
				Autoarchive:    true,
				Contributing:   false,
				EncryptionSeed: uuidx.WithSuffix(16),
			},
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			encoded,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "d81f2387-d17a-cb19-cfb5-a1f3830589bd", result.Feed.Id)
		require.Equal(t, "hello world", result.Feed.Description)
		require.True(t, result.Feed.Autodownload)
		require.True(t, result.Feed.Autoarchive)
		require.False(t, result.Feed.Contributing)
		require.NoError(t, tracking.RSSFindByID(ctx, q, result.Feed.Id).Scan(&feed))
		require.Equal(t, uuidx.WithSuffix(16), feed.EncryptionSeed)
	})
}
