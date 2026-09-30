package media_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/require"
)

func TestHTTPAutoimportSearch(t *testing.T) {
	t.Run("lists directories filtered by id", func(t *testing.T) {
		var (
			d1     library.AutoimportDirectory
			d2     library.AutoimportDirectory
			claims jwt.RegisteredClaims
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&d1, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(t.TempDir())))
		d1.Description = "movies inbox"
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, d1).Scan(&d1))
		require.NoError(t, testx.Fake(&d2, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(t.TempDir())))
		d2.Description = "music inbox"
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, d2).Scan(&d2))

		routes := mux.NewRouter()
		media.NewHTTPAutoimport(
			q,
			asyncx.NewWakeup(t.Context()),
			media.HTTPAutoimportOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims = jwtx.NewJWTClaims("test-subject", jwtx.ClaimsOptionAuthnExpiration())
		token := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		cases := []struct {
			query    string
			expected []string
		}{
			{query: "", expected: []string{d1.ID, d2.ID}},
			{query: fmt.Sprintf("id=%s", d1.ID), expected: []string{d1.ID}},
			{query: "query=movies", expected: []string{d1.ID}},
			{query: "query=inbox", expected: []string{d1.ID, d2.ID}},
			{query: "query=podcasts", expected: []string{}},
		}

		for _, c := range cases {
			var result media.AutoimportDirectorySearchResponse

			resp, req, err := httptestx.BuildRequestBytes(http.MethodGet, fmt.Sprintf("/?%s", c.query), nil, httptestx.RequestOptionAuthorization(token))
			require.NoError(t, err)

			routes.ServeHTTP(resp, req)
			require.NoError(t, httpx.ErrorCode(resp.Result()))
			require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

			ids := make([]string, 0, len(result.Items))
			for _, d := range result.Items {
				ids = append(ids, d.Id)
			}

			require.ElementsMatch(t, c.expected, ids, c.query)
		}
	})
}
