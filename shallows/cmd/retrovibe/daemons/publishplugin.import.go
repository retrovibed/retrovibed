package daemons

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// PublishPluginImport reconciles the publish.d directory with the
// plugin_publishers catalog: every installed module gets a row, and rows
// that no longer describe the module at their path are dropped.
//
// The registry's own fsnotify watch is what loads a plugin, but loading it
// is not enough to use it - publishing fans out over the catalog
// (CommunityPublisherFindByCommunityID -> PluginPublisherFindByID ->
// Publish), so a module with no row is never invoked and never appears as
// something a community can enable. Uploading through the HTTP API creates
// that row itself, as does the community importer; anything installed by
// hand or by the CLI - in particular the symlinks that give one module
// several independently configured identities - relies on this.
//
// Rows are keyed by publishplugin.Identity, so the id is the plugin rather
// than its filename. Identical bytes under the same identity always hash to
// the same id, which is what makes a reinstall, a re-import of the same
// module from another community, and a duplicate install all land on the row
// that already exists, leaving a community's selection intact. Different
// bytes are a different plugin, so replacing a module in place retires its
// row and records the replacement - the same thing that already happens when
// a new build is uploaded.
func PublishPluginImport(ctx context.Context, db sqlx.Queryer, plugindir string) error {
	log.Println("import publish plugins initiated", plugindir)
	defer log.Println("import publish plugins completed", plugindir)

	if err := os.MkdirAll(plugindir, 0o700); err != nil {
		return errorsx.Wrapf(err, "unable to create publish plugin directory: %s", plugindir)
	}

	entries, err := os.ReadDir(plugindir)
	if err != nil {
		return errorsx.Wrapf(err, "unable to read publish plugin directory: %s", plugindir)
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".wasm") {
			continue
		}

		path := filepath.Join(plugindir, entry.Name())

		// Stat rather than the DirEntry's own type: these are routinely
		// symlinks, and following them is what distinguishes a live link
		// from one whose target was uninstalled out from under it.
		if info, serr := os.Stat(path); serr != nil || !info.Mode().IsRegular() {
			log.Println("skipping unusable publish plugin", path, errorsx.Compact(serr))
			continue
		}

		if err := publishplugin.VerifyWasmMagicPath(path); err != nil {
			log.Println(errorsx.Wrap(err, "skipping publish plugin"))
			continue
		}

		id, err := publishplugin.Identity(path)
		if err != nil {
			log.Println(errorsx.Wrap(err, "skipping publish plugin"))
			continue
		}

		// a module named after its own digest carries no label, only an id.
		// the upsert leaves a blank description alone, so whatever the upload
		// endpoint or the community importer recorded survives this pass.
		name := strings.TrimSuffix(entry.Name(), ".wasm")
		if name == id {
			name = ""
		}

		var inserted community.PluginPublisher
		if err := community.PluginPublisherInsertWithDefaults(ctx, db, community.PluginPublisher{
			ID:          id,
			Path:        path,
			Description: name,
		}).Scan(&inserted); err != nil {
			log.Println(errorsx.Wrapf(err, "unable to record publish plugin: %s", path))
			continue
		}

		log.Println("recorded publish plugin", inserted.ID, path)
	}

	// whatever the pass above did not write is a row whose module is gone,
	// was renamed, or was replaced by a different build - its path no longer
	// hashes to the id it was recorded under.
	rows := sqlx.Scan(community.PluginPublisherFindAll(ctx, db))
	for pub := range rows.Iter() {
		if id, err := publishplugin.Identity(pub.Path); err == nil && id == pub.ID {
			continue
		}

		var deleted community.PluginPublisher
		if err := community.PluginPublisherDeleteByID(ctx, db, pub.ID).Scan(&deleted); err != nil {
			log.Println(errorsx.Wrapf(err, "unable to forget uninstalled publish plugin: %s", pub.Path))
			continue
		}

		log.Println("forgot uninstalled publish plugin", deleted.ID, pub.Path)
	}

	if err := rows.Err(); err != nil {
		return errorsx.Wrap(err, "unable to list plugin publishers")
	}

	return nil
}
