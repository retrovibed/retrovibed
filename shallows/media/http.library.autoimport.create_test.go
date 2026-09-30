package media_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
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

func TestHTTPAutoimportCreate(t *testing.T) {
	t.Run("creates a monitored directory", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			d      library.AutoimportDirectory
			result media.AutoimportDirectoryCreateResponse
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		root := t.TempDir()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		media.NewHTTPAutoimport(
			q,
			asyncx.NewWakeup(t.Context()),
			media.HTTPAutoimportOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

		encoded, err := jsonx.Marshal(&media.AutoimportDirectoryCreateRequest{
			Directory: &media.AutoimportDirectory{
				Path:        root,
				Description: "inbox",
				Debounce:    60,
				Mode:        library.AutoimportModeMove,
			},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/", encoded, httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))

		require.NoError(t, library.AutoimportDirectoryFindByID(ctx, q, result.Directory.Id).Scan(&d))
		require.Equal(t, root, d.Path)
		require.Equal(t, "inbox", d.Description)
		require.Equal(t, library.AutoimportModeMove, d.Mode)
		require.Equal(t, uuid.Nil.String(), d.LibraryDirectoryID)
		require.Equal(t, uint64(60), result.Directory.Debounce)
	})

	t.Run("rejects invalid requests", func(t *testing.T) {
		var (
			p meta.Profile
			v meta.Authz
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		root := t.TempDir()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		media.NewHTTPAutoimport(
			q,
			asyncx.NewWakeup(t.Context()),
			media.HTTPAutoimportOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

		cases := []struct {
			name      string
			directory *media.AutoimportDirectory
		}{
			{name: "relative path", directory: &media.AutoimportDirectory{Path: "inbox", Debounce: 60}},
			{name: "missing directory", directory: &media.AutoimportDirectory{Path: filepath.Join(root, "missing"), Debounce: 60}},
			{name: "zero debounce", directory: &media.AutoimportDirectory{Path: root}},
			{name: "unknown mode", directory: &media.AutoimportDirectory{Path: root, Debounce: 60, Mode: 2}},
		}

		for _, c := range cases {
			encoded, err := jsonx.Marshal(&media.AutoimportDirectoryCreateRequest{Directory: c.directory})
			require.NoError(t, err)

			resp, req, err := httptestx.BuildRequestBytes(http.MethodPost, "/", encoded, httptestx.RequestOptionAuthorization(token))
			require.NoError(t, err)

			routes.ServeHTTP(resp, req)
			require.Equal(t, http.StatusBadRequest, resp.Code, c.name)
		}
	})
}
