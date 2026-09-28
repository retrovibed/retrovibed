package ddiscapi_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestHTTPLocateSearch(t *testing.T) {
	t.Run("filters by state and attempts", func(t *testing.T) {
		var (
			pending   ddisc.Locate
			cooldown  ddisc.Locate
			completed ddisc.Locate
			claims    jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, ddisc.LocateInsertWithDefaults(ctx, q, ddisc.NewLocate("ubuntu", "video")).Scan(&pending))
		require.NoError(t, ddisc.LocateInsertWithDefaults(ctx, q, ddisc.NewLocate("debian", "video")).Scan(&cooldown))
		require.NoError(t, ddisc.LocateCooldown(ctx, q, cooldown.ID, time.Now().Add(time.Hour)).Scan(&cooldown))
		require.NoError(t, ddisc.LocateInsertWithDefaults(ctx, q, ddisc.NewLocate("fedora", "video")).Scan(&completed))
		require.NoError(t, ddisc.LocateCompleted(ctx, q, completed.ID).Scan(&completed))

		routes := mux.NewRouter()
		ddiscapi.NewHTTPLocate(
			q,
			asyncx.NewWakeup(t.Context()),
			ddiscapi.HTTPLocateOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims("test-subject", jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		cases := []struct {
			query    string
			expected []string
		}{
			{query: "", expected: []string{pending.ID, cooldown.ID, completed.ID}},
			{query: "pending=true", expected: []string{pending.ID}},
			{query: "completed=true", expected: []string{completed.ID}},
			{query: "attempts_min=1", expected: []string{cooldown.ID}},
			{query: "attempts_max=0", expected: []string{pending.ID, completed.ID}},
			{query: fmt.Sprintf("id=%s", completed.ID), expected: []string{completed.ID}},
			{query: "query=debian", expected: []string{cooldown.ID}},
		}

		for _, c := range cases {
			var result ddiscapi.LocateSearchResponse

			resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, fmt.Sprintf("/?%s", c.query), nil, httptestx.RequestOptionAuthorization(token))
			require.NoError(t, err)

			routes.ServeHTTP(resp, req)
			require.NoError(t, httpx.ErrorCode(resp.Result()), c.query)
			require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

			ids := []string{}
			for _, l := range result.Items {
				ids = append(ids, l.Id)
			}
			require.ElementsMatch(t, c.expected, ids, c.query)
		}
	})
}
