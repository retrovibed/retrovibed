package ddiscapi_test

import (
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestHTTPMediaSearch(t *testing.T) {
	t.Run("search", func(t *testing.T) {
		var (
			d      ddisc.Discovered
			result ddiscapi.MediaSearchResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		d = ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPMedia(
			q,
			ddiscapi.HTTPMediaOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(formx.NewEncoder().Encode(&ddiscapi.MediaSearchRequest{
			Offset: 0,
		}))(t)

		claims = jwtx.NewJWTClaims(d.ID, jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		encoded := ddiscapi.NewMediaFromDiscovered(d)
		require.Equal(t, result.Next.Offset, uint64(0))
		require.Contains(t, result.Items, encoded)
	})

	t.Run("known media id filter", func(t *testing.T) {
		var (
			match  ddisc.Discovered
			other  ddisc.Discovered
			result ddiscapi.MediaSearchResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		knownMediaID := uuid.Must(uuid.NewV7()).String()

		matchID := int160.Random()
		match = ddisc.NewDiscovered(&matchID, ddisc.DiscoveredOptionKnownMedia(knownMediaID), ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, match).Scan(&match))

		otherID := int160.Random()
		other = ddisc.NewDiscovered(&otherID, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, other).Scan(&other))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPMedia(
			q,
			ddiscapi.HTTPMediaOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(formx.NewEncoder().Encode(&ddiscapi.MediaSearchRequest{
			KnownMediaId: knownMediaID,
		}))(t)

		claims = jwtx.NewJWTClaims(match.ID, jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		encodedmatch := ddiscapi.NewMediaFromDiscovered(match)
		require.Contains(t, result.Items, encodedmatch)
		require.Len(t, result.Items, 1)
	})

	t.Run("known media id filter disabled returns entries regardless of known_media_id", func(t *testing.T) {
		var (
			unresolved ddisc.Discovered
			resolved   ddisc.Discovered
			result     ddiscapi.MediaSearchResponse
			claims     jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		unresolvedID := int160.Random()
		unresolved = ddisc.NewDiscovered(&unresolvedID, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, unresolved).Scan(&unresolved))

		resolvedID := int160.Random()
		resolved = ddisc.NewDiscovered(&resolvedID, ddisc.DiscoveredOptionKnownMedia(uuid.Must(uuid.NewV7()).String()), ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, resolved).Scan(&resolved))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPMedia(
			q,
			ddiscapi.HTTPMediaOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(formx.NewEncoder().Encode(&ddiscapi.MediaSearchRequest{}))(t)

		claims = jwtx.NewJWTClaims(unresolved.ID, jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		encodedunresolved := ddiscapi.NewMediaFromDiscovered(unresolved)
		encodedresolved := ddiscapi.NewMediaFromDiscovered(resolved)
		require.Contains(t, result.Items, encodedunresolved)
		require.Contains(t, result.Items, encodedresolved)
		require.Len(t, result.Items, 2)
	})
}
