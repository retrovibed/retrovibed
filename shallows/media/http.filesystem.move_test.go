package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
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

		dst := testdirectory(t, ctx, q, uuid.Nil.String(), "photos")
		md := testfile(t, ctx, q, uuid.Nil.String(), "loose.bin")

		resp, req := testpost(t, token, fmt.Sprintf("/%s", md.ID), media.FilesystemMoveRequest{DirectoryId: dst.ID})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		require.Equal(t, dst.ID, result.Media.DirectoryId)

		require.NoError(t, library.MetadataFindByID(ctx, q, md.ID).Scan(&md))
		require.Equal(t, dst.ID, md.DirectoryID)
	})

	t.Run("returns a row to the root", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testfilesystem(t, ctx)

		src := testdirectory(t, ctx, q, uuid.Nil.String(), "photos")
		md := testfile(t, ctx, q, src.ID, "held.bin")

		resp, req := testpost(t, token, fmt.Sprintf("/%s", md.ID), media.FilesystemMoveRequest{})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		require.NoError(t, library.MetadataFindByID(ctx, q, md.ID).Scan(&md))
		require.Equal(t, uuid.Nil.String(), md.DirectoryID)
	})

	t.Run("rejects a move into the row's own subtree", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testfilesystem(t, ctx)

		top := testdirectory(t, ctx, q, uuid.Nil.String(), "photos")
		nested := testdirectory(t, ctx, q, top.ID, "2026")

		resp, req := testpost(t, token, fmt.Sprintf("/%s", top.ID), media.FilesystemMoveRequest{DirectoryId: nested.ID})
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)

		require.NoError(t, library.MetadataFindByID(ctx, q, top.ID).Scan(&top))
		require.Equal(t, uuid.Nil.String(), top.DirectoryID)
	})
}
