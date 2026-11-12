package media_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/require"
)

func TestFilesystemMove(t *testing.T) {
	t.Run("files a row into a directory", func(t *testing.T) {
		var result media.FilesystemMoveResponse
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testfilesystem(t, ctx)

		dst := testdirectory(t, ctx, q, uuidx.Nil().String(), "photos")
		md := testfile(t, ctx, q, uuidx.Nil().String(), "loose.bin")

		resp, req := testpost(t, token, fmt.Sprintf("/%s", md.ID), media.FilesystemMoveRequest{DirectoryId: dst.ID})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.NoError(t, jsonx.UnmarshalRead(resp.Body, &result))
		require.Equal(t, dst.ID, result.Media.DirectoryId)

		require.NoError(t, library.MetadataFindByID(ctx, q, md.ID).Scan(&md))
		require.Equal(t, dst.ID, md.DirectoryID)
	})

	t.Run("returns a row to the root", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testfilesystem(t, ctx)

		src := testdirectory(t, ctx, q, uuidx.Nil().String(), "photos")
		md := testfile(t, ctx, q, src.ID, "held.bin")

		resp, req := testpost(t, token, fmt.Sprintf("/%s", md.ID), media.FilesystemMoveRequest{})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		require.NoError(t, library.MetadataFindByID(ctx, q, md.ID).Scan(&md))
		require.Equal(t, uuidx.Nil().String(), md.DirectoryID)
	})

	t.Run("rejects a move into the row's own subtree", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testfilesystem(t, ctx)

		top := testdirectory(t, ctx, q, uuidx.Nil().String(), "photos")
		nested := testdirectory(t, ctx, q, top.ID, "2026")

		resp, req := testpost(t, token, fmt.Sprintf("/%s", top.ID), media.FilesystemMoveRequest{DirectoryId: nested.ID})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)

		require.NoError(t, library.MetadataFindByID(ctx, q, top.ID).Scan(&top))
		require.Equal(t, uuidx.Nil().String(), top.DirectoryID)
	})
}
