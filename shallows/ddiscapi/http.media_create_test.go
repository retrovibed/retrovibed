package ddiscapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPMediaCreate(t *testing.T) {
	var (
		p      meta.Profile
		v      meta.Authz
		result ddiscapi.MediaCreateResponse
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

	id := int160.Random()
	infohash := id.Bytes()

	routes := mux.NewRouter()
	ddiscapi.NewHTTPMedia(
		q,
		ddiscapi.HTTPMediaOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

	body := testx.Must(json.Marshal(ddiscapi.MediaCreateRequest{
		Media: &ddiscapi.Media{
			Infohash:    infohash,
			Title:       "test title",
			Description: "test description",
			Mimetype:    "video/mp4",
		},
	}))(t)

	resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/", body, httptestx.RequestOptionAuthorization(token))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.Equal(t, "test title", result.Media.GetTitle())
	require.Equal(t, "test description", result.Media.GetDescription())
	require.Equal(t, "video/mp4", result.Media.GetMimetype())
	require.Equal(t, infohash, result.Media.GetInfohash())

	var d ddisc.Discovered
	require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, result.Media.GetId()).Scan(&d))
	require.Equal(t, "test title", d.Title)
}
