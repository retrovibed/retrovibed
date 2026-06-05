package cmdtorrent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/RoaringBitmap/roaring/v2"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

// seedTorrentDir creates a random torrent and populates torrentDir with the
// .torrent metadata file, a non-zero bitmap, and the content directory that
// exportMagnets expects to find before emitting a record.
func seedTorrentDir(t *testing.T, torrentDir string) torrent.Metadata {
	t.Helper()
	require.NoError(t, os.MkdirAll(torrentDir, 0700))

	info, _, err := torrenttest.Random(t.TempDir(), 16*1024)
	require.NoError(t, err)

	md, err := torrent.NewFromInfo(info)
	require.NoError(t, err)

	require.NoError(t, torrent.NewMetadataCache(torrentDir).Write(md))

	bm := roaring.New()
	bm.Add(0)
	require.NoError(t, torrent.NewBitmapCache(torrentDir).Write(md.ID, bm))

	require.NoError(t, os.MkdirAll(filepath.Join(torrentDir, md.ID.String()), 0700))
	return md
}

// xdgTorrentDir computes the torrent directory that env.TorrentDir() will
// return for a given XDG_DATA_HOME root.
func xdgTorrentDir(xdgData string) string {
	return filepath.Join(xdgData, filepath.Base(os.Args[0]), "torrent")
}

// xdgDBPath computes the database path that cmdopts.DatabaseMeta() will
// open for a given XDG_CONFIG_HOME root.
func xdgDBPath(xdgConfig string) string {
	return filepath.Join(xdgConfig, filepath.Base(os.Args[0]), "meta.db")
}

func TestExportImportPeer(t *testing.T) {
	t.Run("roundtrip preserves infohash and encryption seed", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()
		gctx := &cmdopts.Global{Context: ctx, Shutdown: cancel, Cleanup: &sync.WaitGroup{}}

		xdgData := t.TempDir()
		xdgConfig := t.TempDir()
		t.Setenv("XDG_DATA_HOME", xdgData)
		t.Setenv("XDG_CONFIG_HOME", xdgConfig)

		torrentDir := xdgTorrentDir(xdgData)
		md := seedTorrentDir(t, torrentDir)

		const knownSeed = "550e8400-e29b-41d4-a716-446655440000"
		dbPath := xdgDBPath(xdgConfig)
		db, err := cmdopts.DatabaseCustom(ctx, dbPath)
		require.NoError(t, err)
		lmd := tracking.NewMetadata(langx.Autoptr(md.ID), tracking.MetadataOptionEncryptionSeed(knownSeed))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, db, lmd).Scan(&lmd))
		db.Close()

		magnetsFile := filepath.Join(t.TempDir(), "out.jsonl")
		require.NoError(t, exportMagnets{Path: magnetsFile}.Run(gctx, &cmdopts.SSHID{}))

		f, err := os.Open(magnetsFile)
		require.NoError(t, err)
		defer f.Close()

		dec := json.NewDecoder(f)
		var rec torrentRecord
		require.NoError(t, dec.Decode(&rec))
		require.Contains(t, rec.Magnet, md.ID.String())
		require.Equal(t, knownSeed, rec.EncryptionSeed)
		require.False(t, dec.More(), "expected exactly one record")

		importStore := fsx.DirVirtual(t.TempDir())
		importer := importPeer{Magnets: magnetsFile}
		var count int
		for importedRec, importedMd := range importer.torrents(importStore) {
			count++
			require.Equal(t, md.ID, importedMd.ID)
			require.Equal(t, knownSeed, importedRec.EncryptionSeed)
		}
		require.Equal(t, 1, count)
	})

	t.Run("skips torrents with empty bitmap", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()
		gctx := &cmdopts.Global{Context: ctx, Shutdown: cancel, Cleanup: &sync.WaitGroup{}}

		xdgData := t.TempDir()
		xdgConfig := t.TempDir()
		t.Setenv("XDG_DATA_HOME", xdgData)
		t.Setenv("XDG_CONFIG_HOME", xdgConfig)

		torrentDir := xdgTorrentDir(xdgData)
		require.NoError(t, os.MkdirAll(torrentDir, 0700))

		info, _, err := torrenttest.Random(t.TempDir(), 16*1024)
		require.NoError(t, err)
		md, err := torrent.NewFromInfo(info)
		require.NoError(t, err)
		require.NoError(t, torrent.NewMetadataCache(torrentDir).Write(md))
		require.NoError(t, os.MkdirAll(filepath.Join(torrentDir, md.ID.String()), 0700))
		// write an explicitly empty bitmap (cardinality 0)
		require.NoError(t, torrent.NewBitmapCache(torrentDir).Write(md.ID, roaring.New()))

		magnetsFile := filepath.Join(t.TempDir(), "out.jsonl")
		require.NoError(t, exportMagnets{Path: magnetsFile}.Run(gctx, &cmdopts.SSHID{}))

		info2, err := os.Stat(magnetsFile)
		require.NoError(t, err)
		require.Zero(t, info2.Size(), "expected no records for torrent with empty bitmap")
	})

	t.Run("skips torrents without content directory", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()
		gctx := &cmdopts.Global{Context: ctx, Shutdown: cancel, Cleanup: &sync.WaitGroup{}}

		xdgData := t.TempDir()
		xdgConfig := t.TempDir()
		t.Setenv("XDG_DATA_HOME", xdgData)
		t.Setenv("XDG_CONFIG_HOME", xdgConfig)

		torrentDir := xdgTorrentDir(xdgData)
		require.NoError(t, os.MkdirAll(torrentDir, 0700))

		info, _, err := torrenttest.Random(t.TempDir(), 16*1024)
		require.NoError(t, err)
		md, err := torrent.NewFromInfo(info)
		require.NoError(t, err)
		require.NoError(t, torrent.NewMetadataCache(torrentDir).Write(md))
		bm := roaring.New()
		bm.Add(0)
		require.NoError(t, torrent.NewBitmapCache(torrentDir).Write(md.ID, bm))
		// intentionally omit os.MkdirAll for the content directory

		magnetsFile := filepath.Join(t.TempDir(), "out.jsonl")
		require.NoError(t, exportMagnets{Path: magnetsFile}.Run(gctx, &cmdopts.SSHID{}))

		info2, err := os.Stat(magnetsFile)
		require.NoError(t, err)
		require.Zero(t, info2.Size(), "expected no records when content directory is absent")
	})
}
