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
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// PublishPluginImport reconciles the publish.d directory with the
// plugin_publishers catalog: every installed module gets a row, and rows
// whose module has gone are dropped.
//
// The registry's own fsnotify watch is what loads a plugin, but loading it
// is not enough to use it - publishing fans out over the catalog
// (CommunityPublisherFindByCommunityID -> PluginPublisherFindByID ->
// Publish), so a module with no row is never invoked and never appears as
// something a community can enable. Uploading through the HTTP API creates
// that row itself; anything installed by hand or by the CLI - in
// particular the symlinks that give one module several independently
// configured identities - relies on this.
//
// Rows are matched by path, so plugins uploaded through the API (whose ids
// are derived from their contents) are left exactly as they are, while
// anything discovered here gets an id derived from its filename. That makes
// reinstalling under the same name reuse the same id, so a community's
// existing selection survives the reinstall.
func PublishPluginImport(ctx context.Context, db sqlx.Queryer, plugindir string) error {
	log.Println("import publish plugins initiated", plugindir)
	defer log.Println("import publish plugins completed", plugindir)

	if err := os.MkdirAll(plugindir, 0o700); err != nil {
		return errorsx.Wrapf(err, "unable to create publish plugin directory: %s", plugindir)
	}

	known := make(map[string]community.PluginPublisher, 8)
	rows := sqlx.Scan(community.PluginPublisherFindAll(ctx, db))
	for pub := range rows.Iter() {
		known[pub.Path] = pub
	}
	if err := rows.Err(); err != nil {
		return errorsx.Wrap(err, "unable to list plugin publishers")
	}

	entries, err := os.ReadDir(plugindir)
	if err != nil {
		return errorsx.Wrapf(err, "unable to read publish plugin directory: %s", plugindir)
	}

	installed := make(map[string]struct{}, len(entries))
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

		installed[path] = struct{}{}

		if _, ok := known[path]; ok {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".wasm")
		pub := community.PluginPublisher{
			ID:          md5x.FormatUUID(md5x.Digest(name)),
			Path:        path,
			Description: name,
		}

		var inserted community.PluginPublisher
		if err := community.PluginPublisherInsertWithDefaults(ctx, db, pub).Scan(&inserted); err != nil {
			log.Println(errorsx.Wrapf(err, "unable to record publish plugin: %s", path))
			continue
		}

		log.Println("recorded publish plugin", inserted.ID, path)
	}

	for path, pub := range known {
		if _, ok := installed[path]; ok {
			continue
		}

		var deleted community.PluginPublisher
		if err := community.PluginPublisherDeleteByID(ctx, db, pub.ID).Scan(&deleted); err != nil {
			log.Println(errorsx.Wrapf(err, "unable to forget uninstalled publish plugin: %s", path))
			continue
		}

		log.Println("forgot uninstalled publish plugin", deleted.ID, path)
	}

	return nil
}
