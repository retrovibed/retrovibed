package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/require"
)

func testmove(t *testing.T, routes *mux.Router, token, id, parent string) *http.Response {
	encoded, err := json.Marshal(media.MediaUpdateRequest{Media: &media.Media{ParentId: parent}})
	require.NoError(t, err)

	resp, req, err := httptestx.BuildRequestBytes(
		http.MethodPost,
		fmt.Sprintf("/%s", id),
		encoded,
		httptestx.RequestOptionAuthorization(token),
	)
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)
	return resp.Result()
}

func TestMediaMove(t *testing.T) {
	t.Run("files a row into a folder", func(t *testing.T) {
		var result media.MediaUpdateResponse
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testlibrary(t, ctx)

		dst := testfolder(t, ctx, q, uuid.Nil.String(), "photos")
		md := testfile(t, ctx, q, uuid.Nil.String(), "loose.bin")

		res := testmove(t, routes, token, md.ID, dst.ID)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, json.NewDecoder(res.Body).Decode(&result))
		require.Equal(t, dst.ID, result.Media.ParentId)

		require.NoError(t, library.MetadataFindByID(ctx, q, md.ID).Scan(&md))
		require.Equal(t, dst.ID, md.ParentID)
	})

	t.Run("rejects a move into the row's own subtree", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testlibrary(t, ctx)

		top := testfolder(t, ctx, q, uuid.Nil.String(), "photos")
		nested := testfolder(t, ctx, q, top.ID, "2026")

		require.Equal(t, http.StatusBadRequest, testmove(t, routes, token, top.ID, nested.ID).StatusCode)

		require.NoError(t, library.MetadataFindByID(ctx, q, top.ID).Scan(&top))
		require.Equal(t, uuid.Nil.String(), top.ParentID)
	})

	t.Run("an update carrying no parent leaves the row where it is", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		routes, token, q := testlibrary(t, ctx)

		dst := testfolder(t, ctx, q, uuid.Nil.String(), "photos")
		md := testfile(t, ctx, q, dst.ID, "held.bin")

		require.Equal(t, http.StatusOK, testmove(t, routes, token, md.ID, "").StatusCode)

		require.NoError(t, library.MetadataFindByID(ctx, q, md.ID).Scan(&md))
		require.Equal(t, dst.ID, md.ParentID)
	})
}
