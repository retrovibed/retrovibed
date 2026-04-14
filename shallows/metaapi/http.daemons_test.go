package metaapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPDaemonSearch(t *testing.T) {
	var (
		p      meta.Daemon
		result metaapi.DaemonSearchResponse
		claims jwt.RegisteredClaims
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&p, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, timex.UTCEncodeOption))
	require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, p).Scan(&p))

	routes := mux.NewRouter()

	metaapi.NewHTTPDaemons(
		q,
		metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	b := testx.Must(formx.NewEncoder().Encode(&metaapi.DaemonSearchRequest{
		Offset: 0,
	}))(t)

	claims = jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration())

	resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	encoded := testx.Must(metaapi.NewDaemonFromMetaDaemon(p))(t)
	require.Equal(t, result.Next.Offset, uint64(1))
	require.Contains(t, result.Items, encoded)
}

func TestHTTPDaemonCreateNew(t *testing.T) {
	var (
		v      meta.Daemon
		result metaapi.DaemonCreateResponse
		claims jwt.RegisteredClaims
	)

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&v, meta.DaemonOptionTestDefaults, timex.UTCEncodeOption))

	routes := mux.NewRouter()

	metaapi.NewHTTPDaemons(
		q,
		metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	b := testx.Must(json.Marshal(&metaapi.DaemonCreateRequest{
		Daemon: testx.Must(metaapi.NewDaemonFromMetaDaemon(v))(t),
	}))(t)

	claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

	resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/", b, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	// a id should ge automatically generated.
	require.NotEqual(t, v.ID, result.Daemon.Id)
	// hostname should match
	require.Equal(t, v.Hostname, result.Daemon.Hostname)
}

func TestHTTPDaemonCreateUpdate(t *testing.T) {
	var (
		v      meta.Daemon
		result metaapi.DaemonCreateResponse
		claims jwt.RegisteredClaims
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&v, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, timex.UTCEncodeOption))
	require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, v).Scan(&v))

	time.Sleep(time.Millisecond)

	routes := mux.NewRouter()

	metaapi.NewHTTPDaemons(
		q,
		metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	b := testx.Must(json.Marshal(&metaapi.DaemonCreateRequest{
		Daemon: testx.Must(metaapi.NewDaemonFromMetaDaemon(v))(t),
	}))(t)

	claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

	resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/", b, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	// only timestamp updated should change.
	require.Equal(t, v.ID, result.Daemon.Id)
	require.Equal(t, v.Hostname, result.Daemon.Hostname)
	require.NotEqual(t, v.UpdatedAt, testx.Must(time.Parse(time.RFC3339Nano, result.Daemon.UpdatedAt))(t))
}

func TestHTTPDaemonUpdate(t *testing.T) {
	t.Run("standard update", func(t *testing.T) {
		var (
			v      meta.Daemon
			u      meta.Daemon
			result metaapi.DaemonCreateResponse
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&v, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, timex.UTCEncodeOption))
		require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, testx.Fake(&u, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, meta.DaemonOptionDownload, timex.UTCEncodeOption))

		require.True(t, u.Downloads)

		routes := mux.NewRouter()

		metaapi.NewHTTPDaemons(
			q,
			metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		b := testx.Must(json.Marshal(&metaapi.DaemonUpdateRequest{
			Daemon: testx.Must(metaapi.NewDaemonFromMetaDaemon(u))(t),
		}))(t)

		claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, fmt.Sprintf("/%s", v.ID), b, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.Equal(t, v.ID, result.Daemon.Id)
		require.EqualValues(t, u.Hostname, result.Daemon.Hostname)
		require.True(t, result.Daemon.Downloads)
		require.EqualValues(t, u.Description, result.Daemon.Description)
	})
}

func TestHTTPDaemonTouch(t *testing.T) {
	var (
		v, ignored meta.Daemon
		result     metaapi.DaemonLookupResponse
		claims     jwt.RegisteredClaims
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&v, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, timex.UTCEncodeOption))
	require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, v).Scan(&v))

	require.NoError(t, testx.Fake(&ignored, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, timex.UTCEncodeOption))
	require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, ignored).Scan(&ignored))

	routes := mux.NewRouter()

	metaapi.NewHTTPDaemons(
		q,
		metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

	resp, req, err := httptestx.BuildRequestBytes(http.MethodPut, fmt.Sprintf("/%s", v.ID), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.Equal(t, v.ID, result.Daemon.Id)
	require.Equal(t, v.Hostname, result.Daemon.Hostname)
	require.True(t, result.Daemon.Default)
}

func TestHTTPDaemonDelete(t *testing.T) {
	var (
		v      meta.Daemon
		result metaapi.DaemonCreateResponse
		claims jwt.RegisteredClaims
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&v, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, timex.UTCEncodeOption))
	require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, v).Scan(&v))

	routes := mux.NewRouter()

	metaapi.NewHTTPDaemons(
		q,
		metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	b := testx.Must(json.Marshal(&metaapi.DaemonCreateRequest{
		Daemon: testx.Must(metaapi.NewDaemonFromMetaDaemon(v))(t),
	}))(t)

	claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

	resp, req, err := httptestx.BuildRequestBytes(http.MethodDelete, fmt.Sprintf("/%s", v.ID), b, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.Equal(t, v.ID, result.Daemon.Id)
	require.Equal(t, v.Hostname, result.Daemon.Hostname)
	require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM meta_daemons"))(t))
}

func TestHTTPDaemonLatest(t *testing.T) {
	var (
		v, latest meta.Daemon
		result    metaapi.DaemonLookupResponse
		claims    jwt.RegisteredClaims
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	for i := 0; i < 10; i++ {
		require.NoError(t, testx.Fake(&v, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, timex.UTCEncodeOption))
		require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, v).Scan(&v))
	}

	time.Sleep(time.Millisecond)

	require.NoError(t, testx.Fake(&v, meta.DaemonOptionTestDefaults, meta.DaemonOptionMaybeID, timex.UTCEncodeOption))
	require.NoError(t, meta.DaemonInsertWithDefaults(ctx, q, v).Scan(&v))

	time.Sleep(time.Millisecond)
	require.NoError(t, meta.DaemonTouch(ctx, q, v.ID).Scan(&latest))

	routes := mux.NewRouter()

	metaapi.NewHTTPDaemons(
		q,
		metaapi.HTTPDaemonsOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	b := testx.Must(json.Marshal(&metaapi.DaemonCreateRequest{
		Daemon: testx.Must(metaapi.NewDaemonFromMetaDaemon(v))(t),
	}))(t)

	claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

	resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, "/latest", b, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.Equal(t, v.ID, result.Daemon.Id)
	require.Equal(t, v.Hostname, result.Daemon.Hostname)
	require.Equal(t, latest.UpdatedAt, testx.Must(time.Parse(time.RFC3339Nano, result.Daemon.UpdatedAt))(t))
}
