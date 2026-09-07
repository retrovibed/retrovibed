package media_test

import (
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/require"
)

func testfssearch(t *testing.T, routes *mux.Router, token string, req *media.FilesystemSearchRequest) (result *media.FilesystemSearchResponse) {
	result = &media.FilesystemSearchResponse{}
	query, err := formx.NewEncoder().Encode(req)
	require.NoError(t, err)

	resp, hreq, err := httptestx.BuildRequestBytes(
		http.MethodGet,
		"/?"+query.Encode(),
		nil,
		httptestx.RequestOptionAuthorization(token),
	)
	require.NoError(t, err)

	routes.ServeHTTP(resp, hreq)
	require.Equal(t, http.StatusOK, resp.Result().StatusCode)
	require.NoError(t, jsonx.UnmarshalRead(resp.Body, result))

	return result
}

func testfsids(result *media.FilesystemSearchResponse) []string {
	return slicesx.MapTransform(func(m *media.Media) string { return m.Id }, result.Items...)
}

func TestFilesystemSearch(t *testing.T) {
	t.Run("an absent directory lists the root", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testfilesystem(t, ctx)

		top := testdirectory(t, ctx, q, uuid.Nil.String(), "photos")
		testfile(t, ctx, q, top.ID, "held.bin")
		loose := testfile(t, ctx, q, uuid.Nil.String(), "loose.bin")

		result := testfssearch(t, routes, token, &media.FilesystemSearchRequest{Limit: 10})
		require.ElementsMatch(t, []string{top.ID, loose.ID}, testfsids(result))
		require.Empty(t, result.Breadcrumb)
	})

	t.Run("lists only the contents of the named directory", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testfilesystem(t, ctx)

		top := testdirectory(t, ctx, q, uuid.Nil.String(), "photos")
		nested := testdirectory(t, ctx, q, top.ID, "2026")
		held := testfile(t, ctx, q, top.ID, "held.bin")
		testfile(t, ctx, q, uuid.Nil.String(), "loose.bin")

		result := testfssearch(t, routes, token, &media.FilesystemSearchRequest{Limit: 10, DirectoryId: top.ID})
		require.ElementsMatch(t, []string{nested.ID, held.ID}, testfsids(result))
	})

	t.Run("sorts directories ahead of files", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testfilesystem(t, ctx)

		top := testdirectory(t, ctx, q, uuid.Nil.String(), "photos")
		zzz := testdirectory(t, ctx, q, top.ID, "zzz")
		aaa := testfile(t, ctx, q, top.ID, "aaa.bin")

		result := testfssearch(t, routes, token, &media.FilesystemSearchRequest{Limit: 10, DirectoryId: top.ID})
		require.Equal(t, []string{zzz.ID, aaa.ID}, testfsids(result))
	})

	t.Run("carries the path root first", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testfilesystem(t, ctx)

		top := testdirectory(t, ctx, q, uuid.Nil.String(), "photos")
		nested := testdirectory(t, ctx, q, top.ID, "2026")
		held := testfile(t, ctx, q, nested.ID, "held.bin")

		result := testfssearch(t, routes, token, &media.FilesystemSearchRequest{Limit: 10, DirectoryId: nested.ID})
		require.Equal(t, []string{held.ID}, testfsids(result))
		require.Equal(t, []string{top.ID, nested.ID}, slicesx.MapTransform(func(m *media.Media) string { return m.Id }, result.Breadcrumb...))
	})
}

func TestFilesystemDelete(t *testing.T) {
	t.Run("deleting a directory deletes its contents", func(t *testing.T) {
		result := &media.FilesystemDeleteResponse{}
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testfilesystem(t, ctx)

		top := testdirectory(t, ctx, q, uuid.Nil.String(), "photos")
		nested := testdirectory(t, ctx, q, top.ID, "2026")
		testfile(t, ctx, q, nested.ID, "held.bin")
		outside := testfile(t, ctx, q, uuid.Nil.String(), "loose.bin")

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodDelete,
			"/"+top.ID,
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, result))

		require.Equal(t, top.ID, result.Media.Id)
		// the client is told how much went with it, which is what the console warns about.
		require.Equal(t, uint64(2), result.Removed)

		// nothing below the deleted directory may survive as an orphan.
		remaining := testfssearch(t, routes, token, &media.FilesystemSearchRequest{Limit: 10})
		require.Equal(t, []string{outside.ID}, testfsids(remaining))
	})
}
