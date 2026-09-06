package daemons_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/stretchr/testify/require"
)

func TestPublishPluginImport(t *testing.T) {
	// the magic header publishplugin.VerifyWasmMagicPath insists on,
	// followed by enough filler that a truncated read would be obvious.
	wasm := append([]byte{0x00, 0x61, 0x73, 0x6D}, bytes.Repeat([]byte{0x42}, 64)...)

	t.Run("records an installed plugin under a content derived id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		plugindir := filepath.Join(t.TempDir(), "publish.d")
		require.NoError(t, os.MkdirAll(plugindir, 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(plugindir, "lemmy.wasm"), wasm, 0o600))

		require.NoError(t, daemons.PublishPluginImport(ctx, q, plugindir))

		id, err := publishplugin.Identity(filepath.Join(plugindir, "lemmy.wasm"))
		require.NoError(t, err)

		var row community.PluginPublisher
		require.NoError(t, community.PluginPublisherFindByID(ctx, q, id).Scan(&row))
		require.Equal(t, filepath.Join(plugindir, "lemmy.wasm"), row.Path)
		require.Equal(t, "lemmy", row.Description)
	})

	t.Run("reinstalling the same bytes keeps the row", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		plugindir := filepath.Join(t.TempDir(), "publish.d")
		require.NoError(t, os.MkdirAll(plugindir, 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(plugindir, "lemmy.wasm"), wasm, 0o600))

		require.NoError(t, daemons.PublishPluginImport(ctx, q, plugindir))

		id, err := publishplugin.Identity(filepath.Join(plugindir, "lemmy.wasm"))
		require.NoError(t, err)

		var before community.PluginPublisher
		require.NoError(t, community.PluginPublisherFindByID(ctx, q, id).Scan(&before))

		// what a community has selected hangs off the id, so it has to
		// survive reinstalling the identical module.
		var selected community.CommunityPublisher
		require.NoError(t, community.CommunityPublisherInsertWithDefaults(ctx, q, community.CommunityPublisher{
			ID:          uuid.Must(uuid.NewV7()).String(),
			CommunityID: uuid.Must(uuid.NewV7()).String(),
			PublisherID: id,
		}).Scan(&selected))

		require.NoError(t, os.WriteFile(filepath.Join(plugindir, "lemmy.wasm"), wasm, 0o600))
		require.NoError(t, daemons.PublishPluginImport(ctx, q, plugindir))

		var after community.PluginPublisher
		require.NoError(t, community.PluginPublisherFindByID(ctx, q, id).Scan(&after))
		require.Equal(t, before.ID, after.ID)
		require.Equal(t, before.CreatedAt, after.CreatedAt)
		require.Equal(t, "lemmy", after.Description)

		var stillselected community.CommunityPublisher
		require.NoError(t, community.CommunityPublisherFindByID(ctx, q, selected.ID).Scan(&stillselected))
		require.Equal(t, id, stillselected.PublisherID)
	})

	t.Run("replaces the row when a module is upgraded in place", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		plugindir := filepath.Join(t.TempDir(), "publish.d")
		require.NoError(t, os.MkdirAll(plugindir, 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(plugindir, "lemmy.wasm"), wasm, 0o600))

		require.NoError(t, daemons.PublishPluginImport(ctx, q, plugindir))

		before, err := publishplugin.Identity(filepath.Join(plugindir, "lemmy.wasm"))
		require.NoError(t, err)

		upgraded := append([]byte{0x00, 0x61, 0x73, 0x6D}, bytes.Repeat([]byte{0x43}, 128)...)
		require.NoError(t, os.WriteFile(filepath.Join(plugindir, "lemmy.wasm"), upgraded, 0o600))

		require.NoError(t, daemons.PublishPluginImport(ctx, q, plugindir))

		after, err := publishplugin.Identity(filepath.Join(plugindir, "lemmy.wasm"))
		require.NoError(t, err)
		require.NotEqual(t, before, after)

		var row community.PluginPublisher
		require.Error(t, community.PluginPublisherFindByID(ctx, q, before).Scan(&row))
		require.NoError(t, community.PluginPublisherFindByID(ctx, q, after).Scan(&row))
		require.Equal(t, filepath.Join(plugindir, "lemmy.wasm"), row.Path)
	})

	t.Run("gives each symlink its own row", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		plugindir := filepath.Join(t.TempDir(), "publish.d")
		require.NoError(t, os.MkdirAll(plugindir, 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(plugindir, "lemmy.wasm"), wasm, 0o600))
		require.NoError(t, os.Symlink("lemmy.wasm", filepath.Join(plugindir, "lemmy-movies.wasm")))
		require.NoError(t, os.Symlink("lemmy.wasm", filepath.Join(plugindir, "lemmy-music.wasm")))

		require.NoError(t, daemons.PublishPluginImport(ctx, q, plugindir))

		paths := make([]string, 0, 3)
		rows := sqlx.Scan(community.PluginPublisherFindAll(ctx, q))
		for pub := range rows.Iter() {
			paths = append(paths, pub.Path)
		}
		require.NoError(t, rows.Err())
		require.ElementsMatch(t, []string{
			filepath.Join(plugindir, "lemmy.wasm"),
			filepath.Join(plugindir, "lemmy-movies.wasm"),
			filepath.Join(plugindir, "lemmy-music.wasm"),
		}, paths)
	})

	t.Run("is idempotent and leaves uploaded rows alone", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		plugindir := filepath.Join(t.TempDir(), "publish.d")
		require.NoError(t, os.MkdirAll(plugindir, 0o700))

		// what the upload endpoint produces: the row's id is the content
		// digest and the filename is that same id.
		uploaded := md5x.FormatUUID(md5x.Digest(string(wasm)))
		uploadedPath := filepath.Join(plugindir, uploaded+".wasm")
		require.NoError(t, os.WriteFile(uploadedPath, wasm, 0o600))

		var seeded community.PluginPublisher
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: uploaded, Path: uploadedPath, Description: "uploaded", Mimetype: "application/vnd.retrovibe.publisher.test",
		}).Scan(&seeded))

		require.NoError(t, daemons.PublishPluginImport(ctx, q, plugindir))
		require.NoError(t, daemons.PublishPluginImport(ctx, q, plugindir))

		rows := sqlx.Scan(community.PluginPublisherFindAll(ctx, q))
		found := make([]community.PluginPublisher, 0, 1)
		for pub := range rows.Iter() {
			found = append(found, pub)
		}
		require.NoError(t, rows.Err())
		require.Len(t, found, 1)
		require.Equal(t, uploaded, found[0].ID)
		require.Equal(t, "uploaded", found[0].Description)
	})

	t.Run("forgets a plugin whose module is gone", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		plugindir := filepath.Join(t.TempDir(), "publish.d")
		require.NoError(t, os.MkdirAll(plugindir, 0o700))

		id := uuid.Must(uuid.NewV7()).String()
		var seeded community.PluginPublisher
		require.NoError(t, community.PluginPublisherInsertWithDefaults(ctx, q, community.PluginPublisher{
			ID: id, Path: filepath.Join(plugindir, "uninstalled.wasm"), Description: "uninstalled",
		}).Scan(&seeded))

		require.NoError(t, daemons.PublishPluginImport(ctx, q, plugindir))

		var row community.PluginPublisher
		require.Error(t, community.PluginPublisherFindByID(ctx, q, id).Scan(&row))
	})

	t.Run("skips a dangling symlink and anything that is not a wasm module", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		plugindir := filepath.Join(t.TempDir(), "publish.d")
		require.NoError(t, os.MkdirAll(plugindir, 0o700))
		require.NoError(t, os.Symlink("gone.wasm", filepath.Join(plugindir, "dangling.wasm")))
		require.NoError(t, os.WriteFile(filepath.Join(plugindir, "junk.wasm"), []byte("not a wasm module"), 0o600))
		require.NoError(t, os.WriteFile(filepath.Join(plugindir, "lemmy.env"), []byte("LEMMY_INSTANCE=https://lemmy.ml\n"), 0o600))

		require.NoError(t, daemons.PublishPluginImport(ctx, q, plugindir))

		rows := sqlx.Scan(community.PluginPublisherFindAll(ctx, q))
		for pub := range rows.Iter() {
			t.Fatalf("recorded an unusable plugin: %s", pub.Path)
		}
		require.NoError(t, rows.Err())
	})
}
