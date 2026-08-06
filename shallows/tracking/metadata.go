package tracking

import (
	"context"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/anacrolix/missinggo/pubsub"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	rootenv "github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"golang.org/x/exp/constraints"
	"golang.org/x/time/rate"
)

func MetadataOptionNoop(*Metadata) {}

func MetadataOptionInitiate(md *Metadata) {
	md.InitiatedAt = time.Now()
}

func MetadataOptionFromInfo(i *metainfo.Info) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = strings.ToValidUTF8(i.Name, "\uFFFD")
		m.Bytes = uint64(i.TotalLength())
		m.Private = langx.Autoderef(i.Private)
	}
}

func MetadataOptionFromMagnet(i *metainfo.Magnet) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = strings.ToValidUTF8(i.DisplayName, "\uFFFD")
	}
}

func MetadataOptionMimetype(d string) func(*Metadata) {
	return func(m *Metadata) {
		m.Mimetype = stringsx.FirstNonBlank(d, mimex.Bittorrent)
	}
}

func MetadataOptionDescription(d string) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = d
	}
}

// Currently will select just the first tracker due to poor list support in duckdb.
func MetadataOptionTrackers(d ...string) func(*Metadata) {
	return func(m *Metadata) {
		m.Tracker = slicesx.FirstOrZero(d...)
	}
}

func MetadataOptionKnownMediaID(d string) func(*Metadata) {
	return func(m *Metadata) {
		m.KnownMediaID = d
	}
}

func MetadataOptionDisableKnownMedia(m *Metadata) {
	m.KnownMediaID = uuid.Nil.String()
}

func MetadataOptionAutoEntropySeed(d []byte) func(*Metadata) {
	return func(m *Metadata) {
		m.EncryptionSeed = md5x.FormatUUID(md5x.Digest(m.Infohash, d))
	}
}

func MetadataOptionEntropySeed(d ...[]byte) func(*Metadata) {
	return func(m *Metadata) {
		m.EncryptionSeed = md5x.FormatUUID(md5x.Digest(d...))
	}
}

func MetadataOptionEncryptionSeed(seed string) func(*Metadata) {
	return func(m *Metadata) {
		if stringsx.Blank(seed) {
			return
		}
		m.EncryptionSeed = seed
	}
}

func MetadataOptionAutoArchive(b bool) func(*Metadata) {
	return func(m *Metadata) {
		m.Archivable = b
	}
}

func MetadataOptionBytes[T constraints.Integer](b T) func(*Metadata) {
	return func(m *Metadata) {
		m.Bytes = uint64(b)
	}
}

// MetadataOptionAutoBytes sets Bytes only when it hasn't already been
// determined from an authoritative source (e.g. a downloaded torrent's info
// dict). Useful for magnet uris, where the real size is unknown until the
// info dict is fetched and an approximate size (e.g. from an rss enclosure)
// is the best available hint.
func MetadataOptionAutoBytes[T constraints.Integer](b T) func(*Metadata) {
	return func(m *Metadata) {
		m.Bytes = langx.FirstNonZero(m.Bytes, uint64(b))
	}
}

func MetadataOptionDownloaded[T constraints.Integer](b T) func(*Metadata) {
	return func(m *Metadata) {
		m.Downloaded = uint64(b)
	}
}

func MetadataOptionUploaded[T constraints.Integer](b T) func(*Metadata) {
	return func(m *Metadata) {
		m.Uploaded = uint64(b)
	}
}

// mark the torrent as seeding if the downloaded field is == to the bytes.
func MetadataOptionAutoSeeding(m *Metadata) {
	m.Seeding = m.Downloaded == m.Bytes
}

func MetadataOptionTestDefaults(m *Metadata) {
	*m = NewMetadata(
		new(int160.Random()),
		MetadataOptionBytes(max(m.Bytes, 1)),
		MetadataOptionDownloaded(max(m.Downloaded, 1)),
		MetadataOptionUploaded(max(m.Uploaded, 1)),
	)
}

func MetadataOptionAutoDescription(m *Metadata) {
	m.AutoDescription = library.NormalizedDescription(m.Description)
}

func MetadataOptionAutoHidden(m *Metadata) {
	_, ok := slicesx.Find(func(mime string) bool {
		return m.Mimetype == mime
	}, mimex.RetrovibedMediaArchive, mimex.RetrovibedNeural, mimex.RetrovibedDiscoverySearch)

	if !ok {
		return
	}

	m.HiddenAt = time.Now()
}

func MetadataOptionExpiresAt(ts time.Time) func(*Metadata) {
	return func(m *Metadata) {
		m.ExpiresAt = langx.FirstNonZero(ts, timex.Inf())
	}
}

// sets expired at to be now().add(duration).
func MetadataOptionExpiresAtTTL(d time.Duration) func(*Metadata) {
	return func(m *Metadata) {
		if d == 0 {
			return
		}

		m.ExpiresAt = time.Now().Add(d)
	}
}

func MetadataOptionCompleted(m *Metadata) {
	m.CompletedAt = time.Now()
}

func NewMetadata(md *int160.T, options ...func(*Metadata)) (m Metadata) {
	r := langx.Clone(Metadata{
		ID:             torrentx.HashUID(md),
		Infohash:       md.Bytes(),
		InitiatedAt:    timex.Inf(),
		CompletedAt:    timex.Inf(),
		ImportedAt:     timex.Inf(),
		HiddenAt:       timex.Inf(),
		PausedAt:       timex.Inf(),
		VerifyAt:       timex.Inf(),
		ExpiresAt:      timex.Inf(),
		TombstonedAt:   timex.Inf(),
		NextAnnounceAt: timex.NegInf(),
		KnownMediaID:   uuid.Max.String(),
		EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		Mimetype:       mimex.Bittorrent,
	}, options...)
	return r
}

func MetadataQueryNotInitiated() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.initiated_at = 'infinity'")
}

func MetadataQueryInitiated() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.initiated_at < NOW()")
}

func MetadataQueryIncomplete() squirrel.Sqlizer {
	return squirrel.Expr("(torrents_metadata.completed_at == 'infinity')")
}

func MetadataQueryCompleted(b bool) squirrel.Sqlizer {
	if b {
		return squirrel.Expr("torrents_metadata.completed_at < 'infinity'")
	}

	return squirrel.Expr("torrents_metadata.completed_at == 'infinity'")
}

func MetadataQueryHidden(b bool) squirrel.Sqlizer {
	if b {
		return squirrel.Expr("torrents_metadata.hidden_at < 'infinity'")
	}

	return squirrel.Expr("torrents_metadata.hidden_at = 'infinity'")
}

func MetadataQueryNotTombstoned() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.tombstoned_at = 'infinity'")
}

func MetadataQueryNotPaused() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.paused_at = 'infinity'")
}

func MetadataQueryNeedsVerification() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.verify_at < NOW()")
}

func MetadataQuerySeeding() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.seeding")
}

func MetadataQueryHasTracker() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.tracker != ''")
}

func MetadataQueryNeedsAnnounce() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.next_announce_at < NOW()")
}

func MetadataQueryAnnounceable() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.next_announce_at < 'infinity'")
}

func MetadataQueryNeedsKnownMediaID() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.known_media_id = 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'")
}

func MetadataQueryMediaArchive() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.mimetype = ?", mimex.RetrovibedMediaArchive)
}

func MetadataQueryNotMediaArchive() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.mimetype != ?", mimex.RetrovibedMediaArchive)
}

func MetadataQueryNeural() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.mimetype = ?", mimex.RetrovibedNeural)
}

func MetadataQueryNotNeural() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.mimetype != ?", mimex.RetrovibedNeural)
}

func MetadataQueryDiscoverySearch() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.mimetype = ?", mimex.RetrovibedDiscoverySearch)
}

func MetadataQueryNotDiscoverySearch() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.mimetype != ?", mimex.RetrovibedDiscoverySearch)
}

func MetadataQueryCreatedAfter(ts time.Time) squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.created_at BETWEEN ? AND NOW()", ts)
}

func MetadataQueryNotImported() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.imported_at = 'infinity'")
}

func MetadataQueryNotHidden() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.hidden_at = 'infinity'")
}

func MetadataSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) MetadataScanner {
	return NewMetadataScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func MetadataSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(MetadataScannerStaticColumns)...).From("torrents_metadata")
}

func Verify(ctx context.Context, t torrent.Torrent) error {
	log.Println("verify initiated", t.Metadata().DisplayName)
	defer log.Println("verify completed", t.Metadata().DisplayName, spew.Sdump(t.Stats()))
	return torrent.Verify(ctx, t)
}

// clears out data from storage and resets any torrent data to zero.
func Reset(ctx context.Context, q sqlx.Queryer, vfs fsx.Virtual, md *Metadata) (err error) {
	mediavfs := fsx.DirVirtual(vfs.Path(rootenv.MediaDirName))
	torrentvfs := fsx.DirVirtual(vfs.Path(rootenv.TorrentDirName))

	log.Println("removing torrent initiated", md.ID, int160.FromBytes(md.Infohash).String())
	defer log.Println("removing torrent completed", md.ID, int160.FromBytes(md.Infohash).String())

	if err = sqlx.Discard(sqlx.Scan(library.MetadataTombstoneByTorrentID(ctx, q, md.ID))); err != nil {
		return err
	}

	if err = MetadataResetByID(ctx, q, md.ID).Scan(md); err != nil {
		return err
	}

	ls := sqlx.Scan(library.MetadataSearch(ctx, q, library.MetadataSearchBuilder().Where(
		library.MetadataQueryByTorrentID(md.ID),
	)))

	for lmd := range ls.Iter() {
		log.Println("removing library content", lmd.ID)
		if err = os.RemoveAll(mediavfs.Path(lmd.ID)); err != nil {
			return err
		}
	}

	matches, err := fs.Glob(os.DirFS(torrentvfs.Path()), int160.FromBytes(md.Infohash).String()+"*")
	if err != nil {
		return err
	}

	for _, m := range matches {
		p := torrentvfs.Path(m)
		log.Println("removing torrent content", m, "->", p)
		if err = os.RemoveAll(p); err != nil {
			return err
		}
	}

	return nil
}

func DownloadInto(ctx context.Context, q sqlx.Queryer, vfs fsx.Virtual, mc library.QueryCleaner, md *Metadata, t torrent.Torrent, dst io.Writer, pub *asyncx.Wakeup, options ...torrent.Tuner) (err error) {
	var (
		downloaded int64
	)

	pctx, done := context.WithCancel(ctx)
	defer done()

	// update the progress.
	go DownloadProgress(pctx, q, md, t)

	mediavfs := fsx.DirVirtual(vfs.Path(rootenv.MediaDirName))
	torrentvfs := fsx.DirVirtual(vfs.Path(rootenv.TorrentDirName))
	bcache, err := blockcache.NewDirectoryCache(torrentvfs.Path(int160.FromBytes(md.Infohash).String()))
	if err != nil {
		return err
	}

	// just copying as we receive data to block until done.
	if downloaded, err = torrent.DownloadInto(ctx, dst, t, torrent.TuneAnnounceUntilComplete, torrent.TuneNewConns, langx.ComposeErr(options...)); err != nil {
		return errorsx.Wrap(err, "download failed")
	}

	log.Println("content transfer to library initiated", t.Metadata().ID.String())
	defer log.Println("content transfer to library completed", t.Metadata().ID.String())

	archive, op := detectImportStrategy(t, bcache, torrentvfs, mediavfs)

	for tx, cause := range library.ImportFilesystem(ctx, op, archive, ".") {
		if cause != nil {
			log.Println("import failed", cause)
			err = errorsx.Compact(err, cause)
			continue
		}

		desc := stringsx.Join(" ", md.Description, DescriptionFromPath(md, tx.Path))
		log.Println("------------------------------------------- cleaned", md.Description, tx.Path, "->", desc)

		lmd := library.NewMetadata(
			md5x.FormatUUID(tx.MD5),
			library.MetadataOptionDescription(strings.TrimSpace(desc)),
			library.MetadataOptionAutoDescription(library.NormalizedDescription(desc)),
			library.MetadataOptionBytes(tx.Bytes),
			library.MetadataOptionOffset(tx.Offset),
			library.MetadataOptionTorrentID(md.ID),
			library.MetadataOptionKnownMediaID(md.KnownMediaID),
			library.MetadataOptionMimetype(tx.Mimetype.String()),
			library.MetadataOptionEncryptionSeed(md.EncryptionSeed),
			library.MetadataOptionArchivable(md.Archivable),
			library.MetadataOptionHidden(md.HiddenAt.Before(time.Now())),
		)

		if err := library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd); err != nil {
			return errorsx.Wrap(err, "unable to record library metadata")
		}

		log.Println("new library content", lmd.ID, lmd.Description)
	}

	if err != nil {
		return errorsx.Wrap(err, "failed to transfer files into library")
	}

	stats := t.Stats()
	bytes := uint64(0)
	if i := t.Info(); i != nil {
		bytes = uint64(i.TotalLength())
	}

	if err := MetadataCompleteByID(ctx, q, md.ID, 0, bytes, uint64(downloaded), stats.BytesWrittenData.Uint64()).Scan(md); err != nil {
		return errorsx.Wrap(err, "unable to mark completed")
	}

	pub.Broadcast()

	return nil
}

func DescriptionFromPath(md *Metadata, path string) string {
	tmp := filepath.Base(path)
	if tmp == hex.EncodeToString(md.Infohash) {
		return ""
	}

	if tmp == md.Description {
		return ""
	}

	return tmp
}

func DownloadProgress(ctx context.Context, q sqlx.Queryer, md *Metadata, dl torrent.Torrent) {
	var (
		statsfreq = envx.Duration(1*time.Minute, env.TorrentDownloadStats)
		sub       pubsub.Subscription
	)

	log.Println("monitoring download progress initiated", md.ID, md.Description, md.Tracker, statsfreq)
	defer log.Println("monitoring download progress completed", md.ID, md.Description, md.Tracker)

	// Revisit once resume is working.
	if err := dl.Tune(torrent.TuneSubscribe(&sub)); err != nil {
		log.Println("unable to subscribe", err)
		return
	}
	defer sub.Close()

	statst := time.NewTicker(statsfreq)
	l := rate.NewLimiter(rate.Every(time.Second), 1)

	for {
		select {
		case <-statst.C:
			stats := dl.Stats()
			info := dl.Info()

			log.Printf(
				"%s - %s - %s: info(%t) %s\n", md.ID, hex.EncodeToString(md.Infohash), md.Description, info != nil, stats,
			)

			if err := dl.Tune(torrent.TuneNewConns); err != nil {
				log.Println("unable to request new connections", err)
				continue
			}

			current := uint64(dl.BytesCompleted())
			if md.Downloaded == current || info == nil {
				continue
			}

			uctx, done := context.WithTimeout(context.Background(), time.Second)
			if err := MetadataProgressByID(uctx, q, md.ID, uint16(stats.ActivePeers), uint64(info.TotalLength()), current).Scan(md); err != nil {
				done()
				log.Println("failed to update progress", err)
			}
			done()
		case evt := <-sub.Values:
			switch evt.(type) {
			case torrent.TorrentComplete:
				// when torrent is complete should should trigger emit a final event.
			default:
				if !l.Allow() {
					continue
				}
			}

			stats := dl.Stats()
			info := dl.Info()
			current := uint64(dl.BytesCompleted())
			if md.Downloaded == current || info == nil {
				continue
			}

			statst.Reset(statsfreq)

			log.Printf(
				"%s - %s - %s: info(%t) %s\n", md.ID, hex.EncodeToString(md.Infohash), md.Description, true, stats,
			)

			if err := MetadataProgressByID(ctx, q, md.ID, uint16(stats.ActivePeers), uint64(info.TotalLength()), current).Scan(md); err != nil {
				log.Println("failed to update progress", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func detectBluray(archive fs.StatFS) error {
	_, err := archive.Stat("BDMV/index.bdmv")
	if err == nil {
		return nil
	}

	return errors.ErrUnsupported
}

func detectDVD(archive fs.StatFS) error {
	_, err := archive.Stat("VIDEO_TS/VIDEO_TS.IFO")
	if err == nil {
		return nil
	}

	return errors.ErrUnsupported
}

func detectImportStrategy(md torrent.Torrent, bcache *blockcache.DirCache, srcvfs, vfs fsx.Virtual) (fs.StatFS, library.ImportOp) {
	id := md.Metadata().ID
	archive := blockcache.TorrentFilesystem(bcache, md.Info())

	if detectBluray(archive) == nil {
		return blockcache.TorrentSingleFile(bcache, md.Info()), ImportSymlink(id, srcvfs, vfs)
	}

	if detectDVD(archive) == nil {
		return blockcache.TorrentSingleFile(bcache, md.Info()), ImportSymlink(id, srcvfs, vfs)
	}

	return archive, ImportSymlink(id, srcvfs, vfs)
}

func ImportSymlink(id int160.T, srcvfs, vfs fsx.Virtual) library.ImportOp {
	critical := &sync.Mutex{}
	return func(ctx context.Context, root fs.StatFS, path string) (*library.Transfered, error) {
		tx, err := library.TransferedFromPath(root, path)
		if err != nil {
			return nil, err
		}

		src, err := root.Open(path)
		if err != nil {
			return nil, err
		}
		defer src.Close()

		if n, ok := src.(*blockcache.File); ok {
			tx.Offset = n.Offset
		}

		if n, err := io.Copy(tx.MD5, src); err != nil {
			return nil, err
		} else {
			tx.Bytes = uint64(n)
		}

		uid := md5x.FormatUUID(tx.MD5)

		critical.Lock()
		defer critical.Unlock()
		if err := os.Remove(vfs.Path(uid)); fsx.IgnoreIsNotExist(err) != nil {
			return nil, errorsx.Wrap(err, "unable to ensure symlink destination is available")
		}

		if err := os.Symlink(srcvfs.Path(id.String()), vfs.Path(uid)); err != nil {
			return nil, errorsx.Wrapf(err, "unable to symlink to original location: %s -> %s", srcvfs.Path(id.String()), vfs.Path(uid))
		}

		log.Printf("symlinked: %s -> %s\n", srcvfs.Path(id.String()), vfs.Path(uid))

		return tx, nil
	}
}
