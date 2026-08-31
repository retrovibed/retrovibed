package media_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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

// stands up the library routes behind an authorized admin, which every request here needs
// before it reaches a handler.
func testlibrary(t *testing.T, ctx context.Context) (*mux.Router, string, *sql.DB) {
	var (
		p meta.Profile
		v meta.Authz
	)

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
	require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

	routes := mux.NewRouter()
	media.NewHTTPLibrary(
		q,
		asyncx.NewWakeup(t.Context()),
		asyncx.NewWakeup(t.Context()),
		fsx.DirVirtual(t.TempDir()),
		nil,
		media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

	return routes, httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource), q
}

// a directory row below parent. folders carry a generated id because they have no content
// to hash.
func testfolder(t *testing.T, ctx context.Context, q sqlx.Queryer, parent, description string) library.Metadata {
	var md library.Metadata
	require.NoError(t, testx.Fake(
		&md,
		library.MetadataOptionTestDefaults,
		library.MetadataOptionTestID(uuid.Must(uuid.NewV7()).String()),
		library.MetadataOptionMimetype(mimex.Directory),
		library.MetadataOptionDescription(description),
		library.MetadataOptionAutoDescription(description),
		library.MetadataOptionParentID(parent),
		library.MetadataOptionBytes(0),
	))
	require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))
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
		library.MetadataOptionParentID(parent),
	))
	require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))
	return md
}

func testmkdir(t *testing.T, token string, form url.Values) (*httptest.ResponseRecorder, *http.Request) {
	resp, req, err := httptestx.BuildRequestBytes(
		http.MethodPost,
		"/",
		[]byte(form.Encode()),
		httptestx.RequestOptionAuthorization(token),
		httptestx.RequestOptionHeader("Content-Type", "application/x-www-form-urlencoded"),
	)
	require.NoError(t, err)
	return resp, req
}

func TestLibraryMkdir(t *testing.T) {
	t.Run("creates a folder at the root", func(t *testing.T) {
		var result media.MediaUploadResponse
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testlibrary(t, ctx)

		resp, req := testmkdir(t, token, url.Values{
			"mimetype":    {mimex.Directory},
			"description": {"photos"},
		})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.Equal(t, mimex.Directory, result.Media.Mimetype)
		require.Equal(t, "photos", result.Media.Description)
		require.Equal(t, uuid.Nil.String(), result.Media.ParentId)

		// the id cannot be the md5 of the content the way every other row's is, because a
		// folder has no content.
		require.NotEqual(t, uuid.Nil.String(), result.Media.Id)

		var md library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&md))
		require.Equal(t, uint64(0), md.Bytes)

		// a folder that reached the identification daemon would be matched against known
		// media and then published to the discovery network.
		require.Equal(t, uuid.Nil.String(), md.KnownMediaID)
	})

	t.Run("creates a folder inside another folder", func(t *testing.T) {
		var parent, child media.MediaUploadResponse
		ctx, done := testx.Context(t)
		defer done()

		routes, token, _ := testlibrary(t, ctx)

		resp, req := testmkdir(t, token, url.Values{
			"mimetype":    {mimex.Directory},
			"description": {"photos"},
		})
		routes.ServeHTTP(resp, req)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&parent))

		resp, req = testmkdir(t, token, url.Values{
			"mimetype":    {mimex.Directory},
			"description": {"2026"},
			"parent_id":   {parent.Media.Id},
		})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&child))

		require.Equal(t, parent.Media.Id, child.Media.ParentId)
	})

	t.Run("rejects a folder with no name", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, _ := testlibrary(t, ctx)

		resp, req := testmkdir(t, token, url.Values{"mimetype": {mimex.Directory}})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("a request without the directory mimetype still requires content", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, _ := testlibrary(t, ctx)

		resp, req := testmkdir(t, token, url.Values{"description": {"photos"}})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
