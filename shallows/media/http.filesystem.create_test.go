package media_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

// an authorized admin, which every request here needs before it reaches a handler.
func testauthz(t *testing.T, ctx context.Context, q sqlx.Queryer) string {
	var (
		p meta.Profile
		v meta.Authz
	)

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
	return httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)
}

func testlibrary(t *testing.T, ctx context.Context) (*mux.Router, string, *sql.DB) {
	q := sqltestx.Metadatabase(t)
	token := testauthz(t, ctx, q)

	routes := mux.NewRouter()
	media.NewHTTPLibrary(
		q,
		asyncx.NewWakeup(t.Context()),
		asyncx.NewWakeup(t.Context()),
		fsx.DirVirtual(t.TempDir()),
		nil,
		media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	return routes, token, q
}

func testfilesystem(t *testing.T, ctx context.Context) (*mux.Router, string, *sql.DB) {
	q := sqltestx.Metadatabase(t)
	token := testauthz(t, ctx, q)

	routes := mux.NewRouter()
	media.NewHTTPFilesystem(
		q,
		media.HTTPFilesystemOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	return routes, token, q
}

// a directory row below parent. directories carry a generated id because they have no
// content to hash.
func testdirectory(t *testing.T, ctx context.Context, q sqlx.Queryer, parent, description string) library.Metadata {
	var md library.Metadata
	require.NoError(t, testx.Fake(
		&md,
		library.MetadataOptionTestDefaults,
		library.MetadataOptionTestID(uuid.Must(uuid.NewV7()).String()),
		library.MetadataOptionMimetype(mimex.Directory),
		library.MetadataOptionDescription(description),
		library.MetadataOptionAutoDescription(description),
		library.MetadataOptionDirectoryID(parent),
		library.MetadataOptionBytes(0),
	))
	require.NoError(t, library.DirectoryUpsert(ctx, q, md).Scan(&md))
	return md
}

// a content row below parent.
func testfile(t *testing.T, ctx context.Context, q sqlx.Queryer, parent, description string) library.Metadata {
	var md library.Metadata
	require.NoError(t, testx.Fake(
		&md,
		library.MetadataOptionTestDefaults,
		library.MetadataOptionTestRandomID,
		library.MetadataOptionMimetype(mimex.Binary),
		library.MetadataOptionDescription(description),
		library.MetadataOptionAutoDescription(description),
		library.MetadataOptionDirectoryID(parent),
	))
	require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))
	return md
}

func testpost(t *testing.T, token, path string, body any) (*httptest.ResponseRecorder, *http.Request) {
	encoded, err := json.Marshal(body)
	require.NoError(t, err)

	resp, req, err := httptestx.BuildRequestBytes(
		http.MethodPost,
		path,
		encoded,
		httptestx.RequestOptionAuthorization(token),
	)
	require.NoError(t, err)
	return resp, req
}

func TestFilesystemCreate(t *testing.T) {
	t.Run("creates a directory at the root", func(t *testing.T) {
		var result media.FilesystemCreateResponse
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testfilesystem(t, ctx)

		resp, req := testpost(t, token, "/", media.FilesystemCreateRequest{Name: "photos"})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.Equal(t, mimex.Directory, result.Media.Mimetype)
		require.Equal(t, "photos", result.Media.Description)
		require.Equal(t, uuid.Nil.String(), result.Media.DirectoryId)

		// the id cannot be the md5 of the content the way every other row's is, because a
		// directory has no content.
		require.NotEqual(t, uuid.Nil.String(), result.Media.Id)

		var md library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&md))
		require.Equal(t, uint64(0), md.Bytes)

		// a directory that reached the identification daemon would be matched against
		// known media and then published to the discovery network.
		require.Equal(t, uuid.Nil.String(), md.KnownMediaID)
	})

	t.Run("creates a directory inside another", func(t *testing.T) {
		var parent, child media.FilesystemCreateResponse
		ctx, done := testx.Context(t)
		defer done()

		routes, token, _ := testfilesystem(t, ctx)

		resp, req := testpost(t, token, "/", media.FilesystemCreateRequest{Name: "photos"})
		routes.ServeHTTP(resp, req)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&parent))

		resp, req = testpost(t, token, "/", media.FilesystemCreateRequest{Name: "2026", DirectoryId: parent.Media.Id})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&child))

		require.Equal(t, parent.Media.Id, child.Media.DirectoryId)
	})

	t.Run("rejects a directory with no name", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, _ := testfilesystem(t, ctx)

		resp, req := testpost(t, token, "/", media.FilesystemCreateRequest{})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
