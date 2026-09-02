package communityapi_test

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPCommunityPublisherCreate(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	dir := t.TempDir()

	reg := testx.Must(publishplugin.NewRegistry(ctx, publishplugin.OptionConfigDir(t.TempDir()), publishplugin.OptionCacheDir(t.TempDir())))(t)

	routes := mux.NewRouter()
	communityapi.NewHTTPCommunityPublisher(
		q,
		reg,
		communityapi.HTTPCommunityPublisherOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		communityapi.HTTPCommunityPublisherOptionDir(dir),
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

	upload := func(t *testing.T, authorization string, fields map[string]string, content []byte) (*httptest.ResponseRecorder, *http.Request) {
		mimetype, body, err := httpx.Multipart(func(w *multipart.Writer) error {
			for k, v := range fields {
				if err := w.WriteField(k, v); err != nil {
					return err
				}
			}
			if content == nil {
				return nil
			}
			part, err := w.CreateFormFile("content", "publisher.wasm")
			if err != nil {
				return err
			}
			_, err = part.Write(content)
			return err
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/",
			testx.IOBytes(body),
			httptestx.RequestOptionAuthorization(authorization),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		return resp, req
	}

	t.Run("valid upload is installed and loaded", func(t *testing.T) {
		expectedID := md5x.FormatUUID(md5x.Digest(wasm))

		resp, req := upload(t, token, map[string]string{"description": "YouTube", "mimetype": "application/vnd.retrovibe.publisher.youtube"}, wasm)
		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))

		var result communityapi.PluginPublisherCreateResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, expectedID, result.Publisher.Id)
		require.Equal(t, "YouTube", result.Publisher.Description)
		require.Equal(t, "application/vnd.retrovibe.publisher.youtube", result.Publisher.Mimetype)

		installed := filepath.Join(dir, expectedID+".wasm")
		require.FileExists(t, installed)
		got, err := os.ReadFile(installed)
		require.NoError(t, err)
		require.Equal(t, wasm, got)

		var row community.PluginPublisher
		require.NoError(t, community.PluginPublisherFindByID(ctx, q, expectedID).Scan(&row))
		require.Equal(t, "YouTube", row.Description)
		require.Equal(t, installed, row.Path)
	})

	t.Run("re-uploading identical content upserts rather than duplicates", func(t *testing.T) {
		fields := map[string]string{"description": "Spotify", "mimetype": "application/vnd.retrovibe.publisher.spotify"}
		expectedID := md5x.FormatUUID(md5x.Digest(wasm))

		for range 2 {
			resp, req := upload(t, token, fields, wasm)
			routes.ServeHTTP(resp, req)
			require.NoError(t, httpx.ErrorCode(resp.Result()))
		}

		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM plugin_publishers WHERE id = '"+expectedID+"'"))
	})

	t.Run("invalid wasm magic rejected", func(t *testing.T) {
		content := []byte("not a wasm module")

		resp, req := upload(t, token, map[string]string{"description": "bogus", "mimetype": "application/vnd.retrovibe.publisher.bogus"}, content)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.NoFileExists(t, filepath.Join(dir, md5x.FormatUUID(md5x.Digest(content))+".wasm"))
	})

	t.Run("missing mimetype rejected", func(t *testing.T) {
		resp, req := upload(t, token, map[string]string{"description": "no mimetype"}, wasm)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("missing content file rejected", func(t *testing.T) {
		resp, req := upload(t, token, map[string]string{"mimetype": "application/vnd.retrovibe.publisher.instagram"}, nil)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("requires a privileged token", func(t *testing.T) {
		claims := jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())
		unprivileged := httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)
		content := []byte("x-publisher-binary")

		resp, req := upload(t, unprivileged, map[string]string{"mimetype": "application/vnd.retrovibe.publisher.x"}, content)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.NoFileExists(t, filepath.Join(dir, md5x.FormatUUID(md5x.Digest(content))+".wasm"))
	})
}
