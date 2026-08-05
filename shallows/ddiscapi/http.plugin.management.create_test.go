package ddiscapi_test

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPPluginManagementCreate(t *testing.T) {
	configDir := t.TempDir()

	ctx, done := testx.Context(t)
	defer done()

	reg := testx.Must(searchplugin.NewRegistry(ctx, searchplugin.OptionConfigDir(configDir), searchplugin.OptionCacheDir(t.TempDir())))(t)

	routes := mux.NewRouter()
	ddiscapi.NewHTTPPluginManagement(
		reg,
		ddiscapi.HTTPPluginManagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		ddiscapi.HTTPPluginManagementOptionDir(searchplugin.SearchPluginDir(configDir)),
	).Bind(routes.PathPrefix("/").Subrouter())

	wasmPath := filepath.Join(t.TempDir(), "noop.wasm")
	build := exec.Command("go", "build", "-o", wasmPath, "./.fixtures/noopplugin")
	build.Env = append(build.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))
	wasm := errorsx.Must(os.ReadFile(wasmPath))

	var v meta.Authz
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionAdmin))
	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
	token := httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)

	t.Run("valid upload is installed and loaded", func(t *testing.T) {
		mimetype, body, err := httpx.Multipart(func(w *multipart.Writer) error {
			if err := w.WriteField("name", "noop"); err != nil {
				return err
			}

			part, err := w.CreateFormFile("content", "noop.wasm")
			if err != nil {
				return err
			}

			_, err = io.Copy(part, testx.Read(wasmPath))
			return err
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/",
			testx.IOBytes(body),
			httptestx.RequestOptionAuthorization(token),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result ddiscapi.PluginCreateResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, "noop", result.Plugin.Name)
		require.EqualValues(t, len(wasm), result.Plugin.Size)
		require.FileExists(t, filepath.Join(searchplugin.SearchPluginDir(configDir), "noop.wasm"))
	})

	t.Run("invalid wasm magic rejected", func(t *testing.T) {
		mimetype, body, err := httpx.Multipart(func(w *multipart.Writer) error {
			if err := w.WriteField("name", "bogus"); err != nil {
				return err
			}

			part, err := w.CreateFormFile("content", "bogus.wasm")
			if err != nil {
				return err
			}

			_, err = part.Write([]byte("not a wasm module"))
			return err
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/",
			testx.IOBytes(body),
			httptestx.RequestOptionAuthorization(token),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.NoFileExists(t, filepath.Join(searchplugin.SearchPluginDir(configDir), "bogus.wasm"))
	})

	t.Run("path traversal via name rejected", func(t *testing.T) {
		mimetype, body, err := httpx.Multipart(func(w *multipart.Writer) error {
			if err := w.WriteField("name", "../evil"); err != nil {
				return err
			}

			part, err := w.CreateFormFile("content", "evil.wasm")
			if err != nil {
				return err
			}

			_, err = io.Copy(part, testx.Read(wasmPath))
			return err
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/",
			testx.IOBytes(body),
			httptestx.RequestOptionAuthorization(token),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.NoFileExists(t, filepath.Join(searchplugin.SearchPluginDir(configDir), "..", "evil.wasm"))
	})

	t.Run("requires privileged token", func(t *testing.T) {
		claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		unprivileged := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)

		mimetype, body, err := httpx.Multipart(func(w *multipart.Writer) error {
			if err := w.WriteField("name", "noop2"); err != nil {
				return err
			}

			part, err := w.CreateFormFile("content", "noop.wasm")
			if err != nil {
				return err
			}

			_, err = io.Copy(part, testx.Read(wasmPath))
			return err
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/",
			testx.IOBytes(body),
			httptestx.RequestOptionAuthorization(unprivileged),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.NoFileExists(t, filepath.Join(searchplugin.SearchPluginDir(configDir), "noop2.wasm"))
	})
}
