package metaapi_test

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPFileConfig(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		t.Run("file does not exist", func(t *testing.T) {
			filepath := filepath.Join(t.TempDir(), "nonexistent.json")
			resp, req, err := httptestx.BuildRequestContext(
				t.Context(),
				http.MethodGet,
				"/",
				nil,
				httptestx.RequestOptionAuthorization(httpauthtest.UnsafeTokenAuto(t)),
			)
			require.NoError(t, err)

			routes := mux.NewRouter()
			metaapi.NewHTTPFileConfig(
				filepath,
				metaapi.HTTPFileConfigOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			).Bind(routes)

			routes.ServeHTTP(resp, req)

			require.NoError(t, httpx.ErrorCode(resp.Result()))
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, "null", string(bytes.TrimSpace(body)))
		})

		t.Run("empty file", func(t *testing.T) {
			filepath := filepath.Join(t.TempDir(), "nonexistent.json")
			require.NoError(t, os.WriteFile(filepath, nil, 0600))

			resp, req, err := httptestx.BuildRequestContext(
				t.Context(),
				http.MethodGet,
				"/",
				nil,
				httptestx.RequestOptionAuthorization(httpauthtest.UnsafeTokenAuto(t)),
			)
			require.NoError(t, err)

			routes := mux.NewRouter()
			metaapi.NewHTTPFileConfig(
				filepath,
				metaapi.HTTPFileConfigOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			).Bind(routes)

			routes.ServeHTTP(resp, req)

			require.NoError(t, httpx.ErrorCode(resp.Result()))
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, "null", string(bytes.TrimSpace(body)))
		})

		t.Run("file exists", func(t *testing.T) {
			filepath := filepath.Join(t.TempDir(), "config.json")
			content := []byte(`{"version":"1.0"}`)
			require.NoError(t, os.WriteFile(filepath, content, 0600))

			resp, req, err := httptestx.BuildRequestContext(
				t.Context(),
				http.MethodGet,
				"/",
				nil,
				httptestx.RequestOptionAuthorization(httpauthtest.UnsafeTokenAuto(t)),
			)
			require.NoError(t, err)

			routes := mux.NewRouter()
			metaapi.NewHTTPFileConfig(
				filepath,
				metaapi.HTTPFileConfigOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			).Bind(routes)

			routes.ServeHTTP(resp, req)

			require.NoError(t, httpx.ErrorCode(resp.Result()))
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, content, bytes.TrimSpace(body))
		})

		t.Run("unauthenticated request", func(t *testing.T) {
			filepath := filepath.Join(t.TempDir(), "config.json")

			resp, req, err := httptestx.BuildRequestContext(
				t.Context(),
				http.MethodGet,
				"/",
				nil,
			)
			require.NoError(t, err)

			routes := mux.NewRouter()
			metaapi.NewHTTPFileConfig(
				filepath,
				metaapi.HTTPFileConfigOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			).Bind(routes)

			routes.ServeHTTP(resp, req)

			require.ErrorContains(t, httpx.ErrorCode(resp.Result()), "Unauthorized")
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("successful update", func(t *testing.T) {
			filepath := filepath.Join(t.TempDir(), "config.json")
			content := []byte(`{"version":"2.0"}`)

			resp, req, err := httptestx.BuildRequestContext(
				t.Context(),
				http.MethodPost,
				"/",
				bytes.NewReader(content),
				httptestx.RequestOptionAuthorization(httpauthtest.UnsafeTokenAuto(t)),
			)
			require.NoError(t, err)

			routes := mux.NewRouter()
			metaapi.NewHTTPFileConfig(
				filepath,
				metaapi.HTTPFileConfigOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			).Bind(routes)

			routes.ServeHTTP(resp, req)

			require.NoError(t, httpx.ErrorCode(resp.Result()))
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, content, bytes.TrimSpace(body))

			fileContent, err := os.ReadFile(filepath)
			require.NoError(t, err)
			require.Equal(t, content, fileContent)
		})

		t.Run("unauthenticated request", func(t *testing.T) {
			filepath := filepath.Join(t.TempDir(), "config.json")
			content := []byte(`{"version":"2.0"}`)

			resp, req, err := httptestx.BuildRequestContext(
				t.Context(),
				http.MethodPost,
				"/",
				bytes.NewReader(content),
			)
			require.NoError(t, err)

			routes := mux.NewRouter()
			metaapi.NewHTTPFileConfig(
				filepath,
				metaapi.HTTPFileConfigOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			).Bind(routes)

			routes.ServeHTTP(resp, req)

			require.ErrorContains(t, httpx.ErrorCode(resp.Result()), "Unauthorized")
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("successful delete", func(t *testing.T) {
			var (
				empty = []byte(`{}`)
			)
			filepath := filepath.Join(t.TempDir(), "config.json")
			require.NoError(t, os.WriteFile(filepath, empty, 0600))

			resp, req, err := httptestx.BuildRequestContext(
				t.Context(),
				http.MethodDelete,
				"/",
				nil,
				httptestx.RequestOptionAuthorization(httpauthtest.UnsafeTokenAuto(t)),
			)
			require.NoError(t, err)

			routes := mux.NewRouter()
			metaapi.NewHTTPFileConfig(
				filepath,
				metaapi.HTTPFileConfigOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			).Bind(routes)

			routes.ServeHTTP(resp, req)

			require.NoError(t, httpx.ErrorCode(resp.Result()))
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, empty, bytes.TrimSpace(body))

			_, err = os.Stat(filepath)
			require.True(t, os.IsNotExist(err))
		})

		t.Run("file does not exist", func(t *testing.T) {
			filepath := filepath.Join(t.TempDir(), "nonexistent.json")

			resp, req, err := httptestx.BuildRequestContext(
				t.Context(),
				http.MethodDelete,
				"/",
				nil,
				httptestx.RequestOptionAuthorization(httpauthtest.UnsafeTokenAuto(t)),
			)
			require.NoError(t, err)

			routes := mux.NewRouter()
			metaapi.NewHTTPFileConfig(
				filepath,
				metaapi.HTTPFileConfigOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			).Bind(routes)

			routes.ServeHTTP(resp, req)

			require.NoError(t, httpx.ErrorCode(resp.Result()))
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, "null", string(bytes.TrimSpace(body)))
		})

		t.Run("unauthenticated request", func(t *testing.T) {
			filepath := filepath.Join(t.TempDir(), "config.json")

			resp, req, err := httptestx.BuildRequestContext(
				t.Context(),
				http.MethodDelete,
				"/",
				nil,
			)
			require.NoError(t, err)

			routes := mux.NewRouter()
			metaapi.NewHTTPFileConfig(
				filepath,
				metaapi.HTTPFileConfigOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			).Bind(routes)

			routes.ServeHTTP(resp, req)

			require.ErrorContains(t, httpx.ErrorCode(resp.Result()), "Unauthorized")
		})
	})
}
