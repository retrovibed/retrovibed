package daemons_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttestx"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestAnnounceSeeded(t *testing.T) {
	t.Run("claims a torrent before dispatch so it cannot be re-selected while its announce is still in flight", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		started := make(chan struct{}, 1)
		release := make(chan struct{})
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case started <- struct{}{}:
			default:
			}
			<-release
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		dhts, err := dht.NewServer(32)
		require.NoError(t, err)
		t.Cleanup(dhts.Close)
		pc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, dhts.Serve(ctx, pc))

		vfs := fsx.DirVirtual(t.TempDir())
		tstore := blockcache.NewTorrentFromVirtualFS(vfs)
		tclient := torrenttestx.Client(t, autobind.NewLoopback(), torrent.NewMetadataCache(t.TempDir()), storage.NewFile(t.TempDir()))
		defer tclient.Close()

		var md tracking.Metadata
		require.NoError(t, testx.Fake(&md, tracking.MetadataOptionTestDefaults, func(m *tracking.Metadata) {
			// private skips the dht announce leg entirely, isolating this
			// test to the tracker leg, which we hold open below.
			m.Private = true
			m.Seeding = true
			m.Tracker = srv.URL
		}))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		go daemons.AnnounceSeeded(ctx, q, dhts, fsx.DirVirtual(t.TempDir()), tclient, tstore)

		select {
		case <-started:
		case <-time.After(10 * time.Second):
			t.Fatal("tracker announce never started")
		}

		// the tracker announce is now blocked mid-flight. the fix claims the
		// row (pushes next_announce_at into the future) synchronously before
		// handing it to the async worker pool, so it must already be
		// unselectable even though the announce above hasn't returned yet -
		// otherwise the next pass of the daemon's loop would re-select and
		// re-announce this same torrent while the first announce is still
		// outstanding.
		var got tracking.Metadata
		require.NoError(t, tracking.MetadataFindByID(ctx, q, md.ID).Scan(&got))
		require.True(t, got.NextAnnounceAt.After(time.Now()), "expected the torrent to be optimistically claimed before its announce completed")

		close(release)
	})

	t.Run("records the tracker's reported interval on a successful announce", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		const trackerInterval = 1800 // seconds
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, "d8:completei0e10:incompletei0e8:intervali%de5:peers0:e", trackerInterval)
		}))
		defer srv.Close()

		dhts, err := dht.NewServer(32)
		require.NoError(t, err)
		t.Cleanup(dhts.Close)
		pc, err := net.ListenPacket("udp4", "127.0.0.1:0")
		require.NoError(t, err)
		require.NoError(t, dhts.Serve(ctx, pc))

		vfs := fsx.DirVirtual(t.TempDir())
		tstore := blockcache.NewTorrentFromVirtualFS(vfs)
		tclient := torrenttestx.Client(t, autobind.NewLoopback(), torrent.NewMetadataCache(t.TempDir()), storage.NewFile(t.TempDir()))
		defer tclient.Close()

		var md tracking.Metadata
		require.NoError(t, testx.Fake(&md, tracking.MetadataOptionTestDefaults, func(m *tracking.Metadata) {
			// private skips the dht announce leg entirely, isolating this
			// test to the tracker leg's reported interval.
			m.Private = true
			m.Seeding = true
			m.Tracker = srv.URL
		}))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		start := time.Now()
		go daemons.AnnounceSeeded(ctx, q, dhts, fsx.DirVirtual(t.TempDir()), tclient, tstore)

		var got tracking.Metadata
		require.Eventually(t, func() bool {
			if err := tracking.MetadataFindByID(ctx, q, md.ID).Scan(&got); err != nil {
				return false
			}
			return got.NextAnnounceAt.After(start.Add(trackerInterval * time.Second))
		}, 10*time.Second, 50*time.Millisecond, "expected next_announce_at to reflect the tracker's reported interval")

		// worker computes next_announce_at as now + max(dht, tracker interval) + jitter in [0, 5m).
		// dht is skipped (private), so the floor is exactly the tracker's reported interval.
		diff := got.NextAnnounceAt.Sub(start)
		require.GreaterOrEqual(t, diff, trackerInterval*time.Second, "next_announce_at should be no sooner than the tracker's reported interval")
		require.Less(t, diff, trackerInterval*time.Second+6*time.Minute, "next_announce_at should be within the tracker interval plus jitter")
	})
}
