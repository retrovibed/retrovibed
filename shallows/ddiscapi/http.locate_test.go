package ddiscapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestHTTPLocateCreate(t *testing.T) {
	var (
		result ddiscapi.LocateCreateResponse
		claims jwt.RegisteredClaims
	)

	q := sqltestx.Metadatabase(t)

	routes := mux.NewRouter()
	ddiscapi.NewHTTPLocate(
		q,
		asyncx.NewWakeup(t.Context()),
		ddiscapi.HTTPLocateOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims = jwtx.NewJWTClaims("test-subject", jwtx.ClaimsOptionAuthnExpiration())

	body := testx.Must(json.Marshal(ddiscapi.LocateCreateRequest{
		Locate: &ddiscapi.Locate{
			Query:    "ubuntu",
			Mimetype: "video",
		},
	}))(t)

	resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/", body, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.Equal(t, "ubuntu", result.Locate.Query)
	require.Equal(t, "video", result.Locate.Mimetype)
	require.NotEmpty(t, result.Locate.Id)
	require.NotEmpty(t, result.Locate.KnownMediaId)
}

func TestHTTPLocateCreateRejectsEmpty(t *testing.T) {
	var claims jwt.RegisteredClaims

	q := sqltestx.Metadatabase(t)

	routes := mux.NewRouter()
	ddiscapi.NewHTTPLocate(
		q,
		asyncx.NewWakeup(t.Context()),
		ddiscapi.HTTPLocateOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims = jwtx.NewJWTClaims("test-subject", jwtx.ClaimsOptionAuthnExpiration())

	body := testx.Must(json.Marshal(ddiscapi.LocateCreateRequest{
		Locate: &ddiscapi.Locate{},
	}))(t)

	resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/", body, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
}
