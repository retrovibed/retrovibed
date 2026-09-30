package media_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPAutoimportUpdate(t *testing.T) {
	t.Run("updates the editable fields but not the path", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			d      library.AutoimportDirectory
			result media.AutoimportDirectoryUpdateResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		root := t.TempDir()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, testx.Fake(&d, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(root)))
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, d).Scan(&d))

		routes := mux.NewRouter()
		media.NewHTTPAutoimport(
			q,
			asyncx.NewWakeup(t.Context()),
			media.HTTPAutoimportOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

		encoded, err := jsonx.Marshal(&media.AutoimportDirectoryUpdateRequest{
			Directory: &media.AutoimportDirectory{
				Path:        "/ignored",
				Description: "updated",
				Debounce:    120,
				Mode:        library.AutoimportModeMove,
			},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, fmt.Sprintf("/%s", d.ID), encoded, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		require.NoError(t, library.AutoimportDirectoryFindByID(ctx, q, d.ID).Scan(&d))
		require.Equal(t, root, d.Path)
		require.Equal(t, "updated", d.Description)
		require.Equal(t, 2*time.Minute, d.Debounce)
		require.Equal(t, library.AutoimportModeMove, d.Mode)
		require.Equal(t, d.ID, result.Directory.Id)
	})
}
