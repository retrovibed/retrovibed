package daemons_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverFromRSSFeeds(t *testing.T) {
	t.Run("should be able to locate and create torrent metadata from an rss feed", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)
		defer tclient.Close()

		vfs := fsx.DirVirtual(t.TempDir())
		tstore := blockcache.NewTorrentFromVirtualFS(vfs)

		mux := http.NewServeMux()

		mux.HandleFunc("/index.xml", func(w http.ResponseWriter, r *http.Request) {
			v := strings.ReplaceAll(testx.ReadString(testx.Fixture("torrent.rss", "arch.linux", "index.xml")), "https://archlinux.org", fmt.Sprintf("http://%s", r.Host))
			httptestx.HandleIO(strings.NewReader(v))(w, r)
		})
		mux.HandleFunc("/releng/releases/{id}/torrent/", func(w http.ResponseWriter, r *http.Request) {
			httptestx.HandleIO(testx.Read(testx.Fixture("torrent.rss", "arch.linux"), fmt.Sprintf("%s.torrent", r.PathValue("id"))))(w, r)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		require.NoError(t, fsx.MkDirs(0700, vfs.Path("torrent")))

		feed := langx.Clone(tracking.RSS{}, tracking.RSSOptionDefaultFeeds(tracking.RSS{
			Description:    "Arch Linux - iso",
			URL:            fmt.Sprintf("%s/index.xml", srv.URL),
			Contributing:   true,
			Autodownload:   true,
			LastBuiltAt:    time.Date(2025, time.June, 01, 0, 0, 0, 0, time.UTC),
			EncryptionSeed: uuidx.WithSuffix(16),
		}), tracking.RSSOptionDefaultEncryptionSeed)

		require.NoError(t, tracking.RSSInsertDefaultFeed(t.Context(), q, feed).Scan(&feed))

		require.NotEqual(t, time.Date(2025, time.July, 01, 17, 47, 32, 0, time.UTC), feed.LastBuiltAt)
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")))
		require.NoError(t, daemons.DiscoverFromRSSFeedsOnce(t.Context(), q, vfs, library.QueryCleanerNoop(), tclient, tstore))
		require.Equal(t, 0, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")))
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_metadata")))

		var (
			actual  tracking.Metadata
			updfeed tracking.RSS
		)

		require.NoError(t, tracking.RSSFindByID(t.Context(), q, feed.ID).Scan(&updfeed))
		require.Equal(t, time.Date(2025, time.July, 01, 17, 47, 32, 0, time.UTC), updfeed.LastBuiltAt)

		require.NoError(t, tracking.MetadataFindByID(t.Context(), q, errorsx.Must(sqlx.String(t.Context(), q, "SELECT id::text FROM torrents_metadata"))).Scan(&actual))
		// these values should all be generated consistently
		require.Equal(t, "9f676b73-25ef-674d-6443-c90e562c28db", actual.ID)
		require.Equal(t, "3ae42d96-ac70-58a7-c9f2-71ecb1c36232", actual.EncryptionSeed)
		require.EqualValues(t, 0x50ea8000, actual.Bytes)
		require.False(t, actual.Private)
		require.False(t, actual.Archivable)
		require.True(t, actual.InitiatedAt.Before(time.Now().Add(time.Millisecond)))
	})

	t.Run("should download feeds when digests differ", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)
		vfs := fsx.DirVirtual(t.TempDir())
		tstore := blockcache.NewTorrentFromVirtualFS(vfs)

		mux := http.NewServeMux()

		mux.HandleFunc("/index.xml", func(w http.ResponseWriter, r *http.Request) {
			v := strings.ReplaceAll(testx.ReadString(testx.Fixture("torrent.rss", "arch.linux", "index.xml")), "https://archlinux.org", fmt.Sprintf("http://%s", r.Host))
			httptestx.HandleIO(strings.NewReader(v))(w, r)
		})
		mux.HandleFunc("/releng/releases/{id}/torrent/", func(w http.ResponseWriter, r *http.Request) {
			httptestx.HandleIO(testx.Read(testx.Fixture("torrent.rss", "arch.linux"), fmt.Sprintf("%s.torrent", r.PathValue("id"))))(w, r)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		require.NoError(t, fsx.MkDirs(0700, vfs.Path("torrent")))

		feed := langx.Clone(tracking.RSS{}, tracking.RSSOptionDefaultFeeds(tracking.RSS{
			Description:  "Arch Linux - iso",
			URL:          fmt.Sprintf("%s/index.xml", srv.URL),
			Contributing: true,
			LastBuiltAt:  time.Now(),
		}), tracking.RSSOptionDefaultEncryptionSeed)

		require.NoError(t, tracking.RSSInsertDefaultFeed(t.Context(), q, feed).Scan(&feed))

		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")))
		require.NoError(t, daemons.DiscoverFromRSSFeedsOnce(t.Context(), q, vfs, library.QueryCleanerNoop(), tclient, tstore))
		require.Equal(t, 0, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")))
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_metadata")))
	})

	t.Run("should skip feeds with equal digests", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)
		vfs := fsx.DirVirtual(t.TempDir())
		tstore := blockcache.NewTorrentFromVirtualFS(vfs)

		mux := http.NewServeMux()

		mux.HandleFunc("/index.xml", func(w http.ResponseWriter, r *http.Request) {
			v := strings.ReplaceAll(testx.ReadString(testx.Fixture("torrent.rss", "arch.linux", "index.xml")), "https://archlinux.org", fmt.Sprintf("http://%s", r.Host))
			httptestx.HandleIO(strings.NewReader(v))(w, r)
		})
		mux.HandleFunc("/releng/releases/{id}/torrent/", func(w http.ResponseWriter, r *http.Request) {
			httptestx.HandleIO(testx.Read(testx.Fixture("torrent.rss", "arch.linux"), fmt.Sprintf("%s.torrent", r.PathValue("id"))))(w, r)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		require.NoError(t, fsx.MkDirs(0700, vfs.Path("torrent")))

		feed := langx.Clone(tracking.RSS{}, tracking.RSSOptionDefaultFeeds(tracking.RSS{
			Description:  "Arch Linux - iso",
			URL:          fmt.Sprintf("%s/index.xml", srv.URL),
			Contributing: true,
			LastBuiltAt:  time.Date(2025, time.July, 01, 17, 47, 32, 0, time.UTC),
		}), tracking.RSSOptionDigest("356ecab8-91df-beff-5eb5-a9ba480ff213"), tracking.RSSOptionDefaultEncryptionSeed)

		require.NoError(t, tracking.RSSInsertDefaultFeed(t.Context(), q, feed).Scan(&feed))

		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")))
		require.NoError(t, daemons.DiscoverFromRSSFeedsOnce(t.Context(), q, vfs, library.QueryCleanerNoop(), tclient, tstore))
		require.Equal(t, 0, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")))
		require.Equal(t, 0, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_metadata")))
	})

	t.Run("should skip items from before feeds last sync", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)
		vfs := fsx.DirVirtual(t.TempDir())
		tstore := blockcache.NewTorrentFromVirtualFS(vfs)

		mux := http.NewServeMux()

		mux.HandleFunc("/index.xml", func(w http.ResponseWriter, r *http.Request) {
			v := strings.ReplaceAll(testx.ReadString(testx.Fixture("torrent.rss", "example.2", "index.xml")), "https://archlinux.org", fmt.Sprintf("http://%s", r.Host))
			httptestx.HandleIO(strings.NewReader(v))(w, r)
		})
		mux.HandleFunc("/releng/releases/{id}/torrent/", func(w http.ResponseWriter, r *http.Request) {
			httptestx.HandleIO(testx.Read(testx.Fixture("torrent.rss", "arch.linux"), fmt.Sprintf("%s.torrent", r.PathValue("id"))))(w, r)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		require.NoError(t, fsx.MkDirs(0700, vfs.Path("torrent")))

		feed := langx.Clone(tracking.RSS{}, tracking.RSSOptionDefaultFeeds(tracking.RSS{
			Description:  "Arch Linux - iso",
			URL:          fmt.Sprintf("%s/index.xml", srv.URL),
			Contributing: true,
			LastBuiltAt:  time.Date(2025, time.July, 01, 17, 0, 0, 0, time.UTC),
		}), tracking.RSSOptionDefaultEncryptionSeed)

		require.NoError(t, tracking.RSSInsertDefaultFeed(t.Context(), q, feed).Scan(&feed))

		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")))
		require.NoError(t, daemons.DiscoverFromRSSFeedsOnce(t.Context(), q, vfs, library.QueryCleanerNoop(), tclient, tstore))
		require.Equal(t, 0, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")))
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_metadata")))
		require.True(t, sqltestx.Bool(t, q, "SELECT hidden_at == 'infinity' FROM torrents_metadata"))
	})

	t.Run("should handle magnet uri", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)
		vfs := fsx.DirVirtual(t.TempDir())
		tstore := blockcache.NewTorrentFromVirtualFS(vfs)

		mux := http.NewServeMux()

		mux.HandleFunc("/index.xml", func(w http.ResponseWriter, r *http.Request) {
			httptestx.HandleIO(testx.Read(testx.Fixture("torrent.rss", "example.3", "index.xml")))(w, r)
		})

		srv := httptest.NewServer(mux)
		defer srv.Close()

		require.NoError(t, fsx.MkDirs(0700, vfs.Path("torrent")))

		feed := langx.Clone(tracking.RSS{}, tracking.RSSOptionDefaultFeeds(tracking.RSS{
			Description:  "Arch Linux - iso",
			URL:          fmt.Sprintf("%s/index.xml", srv.URL),
			Contributing: true,
			LastBuiltAt:  time.Date(2025, time.July, 01, 17, 0, 0, 0, time.UTC),
		}), tracking.RSSOptionDefaultEncryptionSeed)

		require.NoError(t, tracking.RSSInsertDefaultFeed(t.Context(), q, feed).Scan(&feed))

		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")))
		require.NoError(t, daemons.DiscoverFromRSSFeedsOnce(t.Context(), q, vfs, library.QueryCleanerNoop(), tclient, tstore))
		require.Equal(t, 0, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_feed_rss WHERE next_check < NOW()")))
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_metadata")))
	})

	t.Run("should handle an item with only a link", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		tclient := torrenttestx.QuickClient(t)
		vfs := fsx.DirVirtual(t.TempDir())
		tstore := blockcache.NewTorrentFromVirtualFS(vfs)

		mux := http.NewServeMux()

		mux.HandleFunc("/index.xml", func(w http.ResponseWriter, r *http.Request) {
			v := strings.ReplaceAll(testx.ReadString(testx.Fixture("torrent.rss", "example.6", "index.xml")), "https://example.com", fmt.Sprintf("http://%s", r.Host))
			httptestx.HandleIO(strings.NewReader(v))(w, r)
		})
		mux.HandleFunc("/retrovibed.media.metadata.archive.torrent", func(w http.ResponseWriter, r *http.Request) {
			httptestx.HandleIO(testx.Read(testx.Fixture("torrent.rss", "example.6"), "example.torrent"))(w, r)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		require.NoError(t, fsx.MkDirs(0700, vfs.Path("torrent")))

		feed := langx.Clone(tracking.RSS{}, tracking.RSSOptionDefaultFeeds(tracking.RSS{
			Description:  t.Name(),
			URL:          fmt.Sprintf("%s/index.xml", srv.URL),
			Contributing: true,
			LastBuiltAt:  time.Date(2025, time.July, 01, 17, 0, 0, 0, time.UTC),
		}), tracking.RSSOptionDefaultEncryptionSeed)

		require.NoError(t, tracking.RSSInsertDefaultFeed(t.Context(), q, feed).Scan(&feed))

		require.NoError(t, daemons.DiscoverFromRSSFeedsOnce(t.Context(), q, vfs, library.QueryCleanerNoop(), tclient, tstore))
		require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_metadata")))
	})

	t.Run("metadata", func(t *testing.T) {
		t.Run("should be marked as hidden", func(t *testing.T) {
			q := sqltestx.Metadatabase(t)

			tclient := torrenttestx.QuickClient(t)
			vfs := fsx.DirVirtual(t.TempDir())
			tstore := blockcache.NewTorrentFromVirtualFS(vfs)

			mux := http.NewServeMux()

			mux.HandleFunc("/index.xml", func(w http.ResponseWriter, r *http.Request) {
				httptestx.HandleIO(testx.Read(testx.Fixture("torrent.rss", "example.4", "index.xml")))(w, r)
			})

			srv := httptest.NewServer(mux)
			defer srv.Close()

			require.NoError(t, fsx.MkDirs(0700, vfs.Path("torrent")))

			feed := langx.Clone(tracking.RSS{}, tracking.RSSOptionDefaultFeeds(tracking.RSS{
				Description:  t.Name(),
				URL:          fmt.Sprintf("%s/index.xml", srv.URL),
				Contributing: true,
				LastBuiltAt:  time.Date(2025, time.July, 01, 17, 0, 0, 0, time.UTC),
			}), tracking.RSSOptionDefaultEncryptionSeed)

			require.NoError(t, tracking.RSSInsertDefaultFeed(t.Context(), q, feed).Scan(&feed))

			require.Equal(t, 0, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_metadata")))
			require.NoError(t, daemons.DiscoverFromRSSFeedsOnce(t.Context(), q, vfs, library.QueryCleanerNoop(), tclient, tstore))
			require.Equal(t, 1, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_metadata")))
			assert.Equal(t, mimex.RetrovibedMediaArchive, sqltestx.String(t, q, "SELECT mimetype FROM torrents_metadata"))
			assert.WithinDuration(t, time.Now(), sqltestx.Timestamp(t, q, "SELECT hidden_at FROM torrents_metadata"), 100*time.Millisecond)
		})

		t.Run("should set mimetype, hidden, and description from a feed of magnet links", func(t *testing.T) {
			q := sqltestx.Metadatabase(t)

			tclient := torrenttestx.QuickClient(t)
			vfs := fsx.DirVirtual(t.TempDir())
			tstore := blockcache.NewTorrentFromVirtualFS(vfs)

			mux := http.NewServeMux()

			mux.HandleFunc("/index.xml", func(w http.ResponseWriter, r *http.Request) {
				httptestx.HandleIO(testx.Read(testx.Fixture("torrent.rss", "example.5", "index.xml")))(w, r)
			})

			srv := httptest.NewServer(mux)
			defer srv.Close()

			require.NoError(t, fsx.MkDirs(0700, vfs.Path("torrent")))

			feed := langx.Clone(tracking.RSS{}, tracking.RSSOptionDefaultFeeds(tracking.RSS{
				Description:  t.Name(),
				URL:          fmt.Sprintf("%s/index.xml", srv.URL),
				Contributing: true,
				LastBuiltAt:  time.Date(2024, time.July, 01, 17, 0, 0, 0, time.UTC),
			}), tracking.RSSOptionDefaultEncryptionSeed)

			require.NoError(t, tracking.RSSInsertDefaultFeed(t.Context(), q, feed).Scan(&feed))

			require.NoError(t, daemons.DiscoverFromRSSFeedsOnce(t.Context(), q, vfs, library.QueryCleanerNoop(), tclient, tstore))
			require.Equal(t, 30, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_metadata")))
			require.Equal(t, 30, errorsx.Zero(sqlx.Count(t.Context(), q, "SELECT COUNT (*) FROM torrents_metadata WHERE mimetype = ?", mimex.RetrovibedMediaArchive)))

			const query = "SELECT mimetype, description, hidden_at, bytes FROM torrents_metadata WHERE description = 'retrovibed.media.metadata.archive.05.tar.gz'"

			var (
				mimetype    string
				description string
				hiddenAt    time.Time
				bytes       uint64
			)
			require.NoError(t, q.QueryRowContext(t.Context(), query).Scan(&mimetype, &description, &hiddenAt, &bytes))

			assert.Equal(t, mimex.RetrovibedMediaArchive, mimetype)
			assert.Equal(t, "retrovibed.media.metadata.archive.05.tar.gz", description)
			assert.WithinDuration(t, time.Now(), hiddenAt, time.Second)
			assert.EqualValues(t, 14266525, bytes)
		})
	})
}
