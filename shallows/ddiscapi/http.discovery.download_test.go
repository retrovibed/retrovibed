package ddiscapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/bytesx"
	pluginapi "github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestHTTPDiscoveryDownload(t *testing.T) {
	t.Run("downloads an already persisted candidate by id, filling in the rest via find", func(t *testing.T) {
		var result ddiscapi.DiscoveryDownloadResponse

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		disc := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, disc).Scan(&disc))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			nil,
			tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir())),
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(disc.ID, jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		body := testx.Must(json.Marshal(ddiscapi.DiscoveryDownloadRequest{
			Discovery:    &ddiscapi.Discovery{Id: disc.ID},
			Autodownload: true,
		}))(t)
		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/download", body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		require.Equal(t, disc.ID, result.Discovery.Id)
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_metadata WHERE initiated_at <= NOW()"))(t))
	})

	t.Run("proto3-json encoded uint64 fields (bytes as string) decode successfully", func(t *testing.T) {
		var result ddiscapi.DiscoveryDownloadResponse

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		disc := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionAutoMagnet)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, disc).Scan(&disc))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			nil,
			tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir())),
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(disc.ID, jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		body := fmt.Appendf(nil, `{"discovery":{"id":%q,"bytes":"12345"},"autodownload":true}`, disc.ID)
		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/download", body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		require.Equal(t, disc.ID, result.Discovery.Id)
	})

	t.Run("downloads an ephemeral candidate never persisted by the websocket search", func(t *testing.T) {
		var result ddiscapi.DiscoveryDownloadResponse

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		uri := (&metainfo.Magnet{InfoHash: metainfo.Hash(id.Bytes())}).String()
		ephemeralID := uuid.Must(uuid.NewV7()).String()

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			nil,
			tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir())),
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(ephemeralID, jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		body := testx.Must(json.Marshal(ddiscapi.DiscoveryDownloadRequest{
			Discovery: &ddiscapi.Discovery{
				Id:  ephemeralID,
				Uri: uri,
			},
			Autodownload: true,
		}))(t)
		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/download", body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		var disc ddisc.Discovered
		require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, ephemeralID).Scan(&disc))
		require.Equal(t, id.Bytes(), disc.Infohash)
		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_metadata WHERE initiated_at <= NOW()"))(t))
	})

	t.Run("missing discovery payload returns not found", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			nil,
			tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir())),
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims("missing", jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/download", nil, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("blank uri for an unknown id returns not found", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			nil,
			tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir())),
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims("missing", jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		body := testx.Must(json.Marshal(ddiscapi.DiscoveryDownloadRequest{
			Discovery: &ddiscapi.Discovery{Id: "missing"},
		}))(t)
		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/download", body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("unfetchable uri returns internal server error", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			nil,
			tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir())),
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		ephemeralID := uuid.Must(uuid.NewV7()).String()
		claims := jwtx.NewJWTClaims(ephemeralID, jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		body := testx.Must(json.Marshal(ddiscapi.DiscoveryDownloadRequest{
			Discovery: &ddiscapi.Discovery{
				Id:  ephemeralID,
				Uri: "https://127.0.0.1:1/nonexistent.torrent",
			},
			Autodownload: true,
		}))(t)
		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/download", body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)
	})

	t.Run("resolving a queue-persisted candidate corrects its placeholder infohash and default-private flag", func(t *testing.T) {
		var result ddiscapi.DiscoveryDownloadResponse

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		// build a real, public .torrent and serve it - same shape a
		// tracker-uri candidate from a search plugin would carry.
		info, _, err := torrenttest.Random(t.TempDir(), 128*bytesx.KiB)
		require.NoError(t, err)
		mi := metainfo.MetaInfo{InfoBytes: testx.Must(metainfo.Encode(info))(t)}
		encoded := testx.Must(metainfo.Encode(mi))(t)
		mux2 := http.NewServeMux()
		mux2.HandleFunc("/x.torrent", httptestx.HandleIO(bytes.NewReader(encoded)))
		torrentSrv := httptest.NewServer(mux2)
		defer torrentSrv.Close()

		// exactly what ddisc.search.queue.background.go persists for an
		// external candidate: placeholder infohash, defaulted private,
		// never fetched yet.
		disc := ddisc.NewDiscoveredFromImport(&pluginapi.Import{Uri: torrentSrv.URL + "/x.torrent"})
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, disc).Scan(&disc))
		require.True(t, disc.Private)

		routes := mux.NewRouter()
		ddiscapi.NewHTTPDiscovery(
			q,
			searchplugin.Unimplemented{},
			nil,
			tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir())),
			ddiscapi.HTTPDiscoveryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := jwtx.NewJWTClaims(disc.ID, jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		// only the id is sent, same as KnownMediaLocator downloading an
		// already-persisted row it never streamed itself - the handler
		// looks the rest up by id.
		body := testx.Must(json.Marshal(ddiscapi.DiscoveryDownloadRequest{
			Discovery:    &ddiscapi.Discovery{Id: disc.ID},
			Autodownload: true,
		}))(t)
		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/download", body, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		var resolved ddisc.Discovered
		require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, disc.ID).Scan(&resolved))
		require.NotEqual(t, disc.Infohash, resolved.Infohash, "the placeholder infohash must be replaced by the real one")
		require.Equal(t, mi.ID().Bytes(), resolved.Infohash)
		require.False(t, resolved.Private, "a public torrent must have its default-private flag cleared once resolved")
		require.EqualValues(t, 1, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM torrents_metadata WHERE initiated_at <= NOW()"))(t))
	})
}
