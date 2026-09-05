package cmdmedia

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestReindexRun(t *testing.T) {
	t.Run("resolves each file's real path so distinct files get distinct descriptions", func(t *testing.T) {
		ctx := t.Context()
		db := sqltestx.Metadatabase(t)

		info := testx.Must(torrenttest.Tree(t.TempDir(), rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{"file1.mkv", "file2.mkv"}))(t)
		md, err := torrent.NewFromInfo(info)
		require.NoError(t, err)

		torrentdir := t.TempDir()
		raw := metainfo.MetaInfo{InfoBytes: md.InfoBytes}
		require.NoError(t, os.WriteFile(filepath.Join(torrentdir, md.ID.String()+tracking.TorrentSuffix), testx.Must(metainfo.Encode(raw))(t), 0600))

		tmd := tracking.NewMetadata(new(md.ID), tracking.MetadataOptionFromInfo(info), tracking.MetadataOptionAutoDescription)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, db, tmd).Scan(&tmd))

		var files []metainfo.File
		for f := range metainfo.Files(info) {
			files = append(files, f)
		}
		require.Len(t, files, 2)

		mediastore := fsx.DirVirtual(t.TempDir())

		lmd0ID := uuid.Must(uuid.NewV4()).String()
		require.NoError(t, os.Symlink(filepath.Join(torrentdir, md.ID.String()), mediastore.Path(lmd0ID)))
		lmd0 := library.NewMetadata(lmd0ID,
			library.MetadataOptionDescription(files[0].Path),
			library.MetadataOptionBytes(files[0].Length),
			library.MetadataOptionOffset(files[0].Offset),
			library.MetadataOptionTorrentID(tmd.ID),
		)
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, lmd0).Scan(&lmd0))

		lmd1ID := uuid.Must(uuid.NewV4()).String()
		require.NoError(t, os.Symlink(filepath.Join(torrentdir, md.ID.String()), mediastore.Path(lmd1ID)))
		lmd1 := library.NewMetadata(lmd1ID,
			library.MetadataOptionDescription(files[1].Path),
			library.MetadataOptionBytes(files[1].Length),
			library.MetadataOptionOffset(files[1].Offset),
			library.MetadataOptionTorrentID(tmd.ID),
		)
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, lmd1).Scan(&lmd1))

		require.NoError(t, reindex{DryRun: false}.run(ctx, db, library.QueryCleanerNoop(), mediastore))

		var got0, got1 library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, db, lmd0.ID).Scan(&got0))
		require.NoError(t, library.MetadataFindByID(ctx, db, lmd1.ID).Scan(&got1))

		require.Contains(t, got0.Description, files[0].Path)
		require.Contains(t, got1.Description, files[1].Path)
		require.NotEqual(t, got0.Description, got1.Description)
	})

	t.Run("dry run leaves records unchanged", func(t *testing.T) {
		ctx := t.Context()
		db := sqltestx.Metadatabase(t)

		info := testx.Must(torrenttest.Tree(t.TempDir(), rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{"file1.mkv"}))(t)
		md, err := torrent.NewFromInfo(info)
		require.NoError(t, err)

		torrentdir := t.TempDir()
		raw := metainfo.MetaInfo{InfoBytes: md.InfoBytes}
		require.NoError(t, os.WriteFile(filepath.Join(torrentdir, md.ID.String()+tracking.TorrentSuffix), testx.Must(metainfo.Encode(raw))(t), 0600))

		tmd := tracking.NewMetadata(new(md.ID), tracking.MetadataOptionFromInfo(info), tracking.MetadataOptionAutoDescription)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, db, tmd).Scan(&tmd))

		var files []metainfo.File
		for f := range metainfo.Files(info) {
			files = append(files, f)
		}
		require.Len(t, files, 1)

		mediastore := fsx.DirVirtual(t.TempDir())

		lmdID := uuid.Must(uuid.NewV4()).String()
		require.NoError(t, os.Symlink(filepath.Join(torrentdir, md.ID.String()), mediastore.Path(lmdID)))
		lmd := library.NewMetadata(lmdID,
			library.MetadataOptionDescription("original description"),
			library.MetadataOptionBytes(files[0].Length),
			library.MetadataOptionOffset(files[0].Offset),
			library.MetadataOptionTorrentID(tmd.ID),
		)
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, lmd).Scan(&lmd))

		require.NoError(t, reindex{DryRun: true}.run(ctx, db, library.QueryCleanerNoop(), mediastore))

		var got library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, db, lmd.ID).Scan(&got))
		require.Equal(t, "original description", got.Description)
	})

	t.Run("unindexed flag only updates records missing an auto description", func(t *testing.T) {
		ctx := t.Context()
		db := sqltestx.Metadatabase(t)

		info := testx.Must(torrenttest.Tree(t.TempDir(), rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{"file1.mkv", "file2.mkv"}))(t)
		md, err := torrent.NewFromInfo(info)
		require.NoError(t, err)

		torrentdir := t.TempDir()
		raw := metainfo.MetaInfo{InfoBytes: md.InfoBytes}
		require.NoError(t, os.WriteFile(filepath.Join(torrentdir, md.ID.String()+tracking.TorrentSuffix), testx.Must(metainfo.Encode(raw))(t), 0600))

		tmd := tracking.NewMetadata(new(md.ID), tracking.MetadataOptionFromInfo(info), tracking.MetadataOptionAutoDescription)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, db, tmd).Scan(&tmd))

		var files []metainfo.File
		for f := range metainfo.Files(info) {
			files = append(files, f)
		}
		require.Len(t, files, 2)

		mediastore := fsx.DirVirtual(t.TempDir())

		indexedID := uuid.Must(uuid.NewV4()).String()
		require.NoError(t, os.Symlink(filepath.Join(torrentdir, md.ID.String()), mediastore.Path(indexedID)))
		indexed := library.NewMetadata(indexedID,
			library.MetadataOptionDescription("already indexed"),
			library.MetadataOptionAutoDescription("already indexed"),
			library.MetadataOptionBytes(files[0].Length),
			library.MetadataOptionOffset(files[0].Offset),
			library.MetadataOptionTorrentID(tmd.ID),
		)
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, indexed).Scan(&indexed))

		unindexedID := uuid.Must(uuid.NewV4()).String()
		require.NoError(t, os.Symlink(filepath.Join(torrentdir, md.ID.String()), mediastore.Path(unindexedID)))
		unindexed := library.NewMetadata(unindexedID,
			library.MetadataOptionDescription("not yet indexed"),
			library.MetadataOptionBytes(files[1].Length),
			library.MetadataOptionOffset(files[1].Offset),
			library.MetadataOptionTorrentID(tmd.ID),
		)
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, unindexed).Scan(&unindexed))

		require.NoError(t, reindex{Unindexed: true, DryRun: false}.run(ctx, db, library.QueryCleanerNoop(), mediastore))

		var gotIndexed, gotUnindexed library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, db, indexed.ID).Scan(&gotIndexed))
		require.NoError(t, library.MetadataFindByID(ctx, db, unindexed.ID).Scan(&gotUnindexed))

		require.Equal(t, "already indexed", gotIndexed.Description)
		require.Contains(t, gotUnindexed.Description, files[1].Path)
	})

	t.Run("record with no matching tracking metadata is skipped without failing", func(t *testing.T) {
		ctx := t.Context()
		db := sqltestx.Metadatabase(t)

		mediastore := fsx.DirVirtual(t.TempDir())

		lmdID := uuid.Must(uuid.NewV4()).String()
		orphanTorrentID := uuid.Must(uuid.NewV4()).String()
		lmd := library.NewMetadata(lmdID,
			library.MetadataOptionDescription("orphaned"),
			library.MetadataOptionTorrentID(orphanTorrentID),
		)
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, lmd).Scan(&lmd))

		require.NoError(t, reindex{DryRun: false}.run(ctx, db, library.QueryCleanerNoop(), mediastore))

		var got library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, db, lmd.ID).Scan(&got))
		require.Equal(t, "orphaned", got.Description)
	})

	t.Run("record whose torrent file is missing from disk is skipped without failing", func(t *testing.T) {
		ctx := t.Context()
		db := sqltestx.Metadatabase(t)

		info := testx.Must(torrenttest.Tree(t.TempDir(), rand.Reader, 16*bytesx.KiB, 64*bytesx.KiB, []string{"file1.mkv"}))(t)
		md, err := torrent.NewFromInfo(info)
		require.NoError(t, err)

		tmd := tracking.NewMetadata(new(md.ID), tracking.MetadataOptionFromInfo(info), tracking.MetadataOptionAutoDescription)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, db, tmd).Scan(&tmd))

		mediastore := fsx.DirVirtual(t.TempDir())

		// symlink points at a torrent directory that has no <id>.torrent file written to disk.
		missingTorrentDir := t.TempDir()
		lmdID := uuid.Must(uuid.NewV4()).String()
		require.NoError(t, os.Symlink(filepath.Join(missingTorrentDir, md.ID.String()), mediastore.Path(lmdID)))
		lmd := library.NewMetadata(lmdID,
			library.MetadataOptionDescription("unreachable"),
			library.MetadataOptionTorrentID(tmd.ID),
		)
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, lmd).Scan(&lmd))

		require.NoError(t, reindex{DryRun: false}.run(ctx, db, library.QueryCleanerNoop(), mediastore))

		var got library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, db, lmd.ID).Scan(&got))
		require.Equal(t, "unreachable", got.Description)
	})

	t.Run("folders are not media and never reach the reindex loop", func(t *testing.T) {
		ctx := t.Context()
		db := sqltestx.Metadatabase(t)
		mediastore := fsx.DirVirtual(t.TempDir())

		// the loop treats a nil torrent id as an unexpected state and logs it. a folder has
		// no torrent by construction, so every folder would report as a defect.
		dir := library.NewMetadata(
			uuid.Must(uuid.NewV7()).String(),
			library.MetadataOptionDescription("photos"),
			library.MetadataOptionMimetype(mimex.Directory),
		)
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, dir).Scan(&dir))

		require.NoError(t, reindex{DryRun: false}.run(ctx, db, library.QueryCleanerNoop(), mediastore))

		var got library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, db, dir.ID).Scan(&got))
		require.Equal(t, "photos", got.Description)
	})
}
