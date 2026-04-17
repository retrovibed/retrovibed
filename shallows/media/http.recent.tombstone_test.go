package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestRecentTombstone(t *testing.T) {
	t.Run("removes session from latest", func(t *testing.T) {
		var (
			p     meta.Profile
			authz meta.Authz
			md    library.Metadata
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))

		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		routes := mux.NewRouter()
		media.NewHTTPRecent(q, media.HTTPRecentOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource)).Bind(routes.PathPrefix("/").Subrouter())
		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

		query := &media.MediaSearchRequest{Query: "jazz"}
		encoded, err := proto.Marshal(query)
		require.NoError(t, err)

		var rs library.RecentSession
		require.NoError(t, library.RecentSessionInsertWithDefaults(ctx, q, library.RecentSession{
			ID:      uuid.Must(uuid.NewV4()).String(),
			MediaID: md.ID,
			Query:   encoded,
		}).Scan(&rs))

		resp2, req2, err := httptestx.BuildRequestBytes(
			http.MethodDelete, fmt.Sprintf("/%s", rs.ID), nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp2, req2)
		require.Equal(t, http.StatusOK, resp2.Result().StatusCode)

		searchreq := &media.RecentSearchRequest{Created: meta.NewDateRange(timex.NewRangeDuration(24 * time.Hour)), Limit: 100}
		searchenc, err := formx.NewEncoder().Encode(&searchreq)
		require.NoError(t, err)

		resp3, req3, err := httptestx.BuildRequestBytes(
			http.MethodGet, "/?"+searchenc.Encode(), nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp3, req3)
		require.Equal(t, http.StatusOK, resp3.Result().StatusCode)

		var result media.RecentSearchResponse
		require.NoError(t, json.NewDecoder(resp3.Body).Decode(&result))
		require.Empty(t, result.Items)
	})

	t.Run("returns 404 for unknown id", func(t *testing.T) {
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
			http.MethodDelete, "/00000000-0000-0000-0000-000000000000", nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})
}
