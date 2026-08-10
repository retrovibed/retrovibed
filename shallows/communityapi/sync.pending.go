package communityapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type MetricsPublisher interface {
	Publish(ctx context.Context, req *PublishContentRequest, torrent io.Reader) (*PublishContentResponse, error)
}

func magnetURI(tmd tracking.Metadata, name string) string {
	return metainfo.NewMagnetFromInfohash(tmd.Infohash, metainfo.MagnetOptionDisplayName(name)).String()
}

func ensureTorrent(ctx context.Context, q sqlx.Queryer, mvfs, tvfs fsx.Virtual, lmd *library.Metadata) (tmd tracking.Metadata, err error) {
	if lmd.TorrentID != uuid.Nil.String() {
		if err = tracking.MetadataFindByID(ctx, q, lmd.TorrentID).Scan(&tmd); err != nil {
			return tmd, err
		}

		if !tmd.Seeding || tmd.CompletedAt.Equal(timex.Inf()) {
			if err = media.ValidateTorrent(ctx, q, tvfs, &tmd); err != nil {
				return tmd, err
			}
		}

		return tmd, err
	}

	return media.GenerateTorrent(ctx, q, mvfs, tvfs, lmd)
}

// SyncPendingToDeeppool syncs pending published content to deeppool and regenerates affected feeds.
func SyncPendingToDeeppool(ctx context.Context, q sqlx.Queryer, httpc *http.Client, metrics MetricsPublisher, publisher FeedPublisher, archiver library.Archiver, mvfs, tvfs fsx.Virtual) error {
	pending := sqlx.Scan(community.PublishedContentFindByPendingSync(ctx, q))

	for pc := range pending.Iter() {
		var (
			lmd   library.Metadata
			known library.Known
		)

		if err := library.MetadataFindByID(ctx, q, pc.LibraryID).Scan(&lmd); err != nil {
			log.Println(errorsx.Wrap(err, "failed to find library metadata"))
			continue
		}

		tmd, err := ensureTorrent(ctx, q, mvfs, tvfs, &lmd)
		if err != nil {
			log.Println(errorsx.Wrap(err, "failed to ensure torrent"))
			continue
		}

		pc.MagnetURI = magnetURI(tmd, lmd.Description)
		if err := community.PublishedContentUpdateMagnetURI(ctx, q, pc.ID, pc.MagnetURI).Scan(&pc); err != nil {
			log.Println(errorsx.Wrap(err, "failed to update magnet_uri"))
			continue
		}

		if pc.KnownMediaID != "" {
			if err := library.KnownFindByID(ctx, q, pc.KnownMediaID).Scan(&known); sqlx.IgnoreNoRows(err) != nil {
				log.Println(errorsx.Wrap(err, "failed to find known media"))
			}
		}

		if pc.OAuthGoogleID != uuid.Nil.String() {
			if uerr := community.YouTubeUpload(ctx, q, httpc, mvfs, pc.OAuthGoogleID, lmd, stringsx.FirstNonBlank(known.Title, lmd.Description), known.Overview); uerr != nil {
				log.Println(errorsx.Wrap(uerr, "youtube cross-post failed"))
			}
		}

		if pc.PublishMode == int32(PublishMode_UNLISTED) {
			if err := community.PublishedContentUpdatePublishedAt(ctx, q, pc.ID, time.Now()).Scan(&pc); err != nil {
				log.Println(errorsx.Wrap(err, "failed to update published_at"))
			}
			continue
		}

		if pc.PublishMode == int32(PublishMode_LISTED) {
			if err := community.PublishedContentUpdatePublishedAt(ctx, q, pc.ID, time.Now()).Scan(&pc); err != nil {
				log.Println(errorsx.Wrap(err, "failed to update published_at"))
				continue
			}

			if err := RegenerateFeed(ctx, q, publisher, pc.CommunityID); err != nil {
				log.Println(errorsx.Wrap(err, "feed regeneration failed for community "+pc.CommunityID))
				continue
			}

			log.Printf("synced listed published content %s", pc.ID)
			continue
		}

		if lmd.ArchiveID == uuid.Max.String() {
			continue // archival in progress, waiting
		} else if lmd.ArchiveID == uuid.Nil.String() {
			if err := library.Archive(ctx, q, &lmd, mvfs, archiver); err != nil {
				log.Println(errorsx.Wrap(err, "failed to archive media"))
				continue
			}
		}

		encoded, err := os.ReadFile(tvfs.Path(fmt.Sprintf("%s.torrent", int160.FromBytes(tmd.Infohash).String())))
		if err != nil {
			log.Println(errorsx.Wrap(err, "failed to open torrent file"))
			continue
		}

		req := PublishContentRequest{
			PublishedContent: &PublishedContent{
				Id:             pc.ID,
				CommunityId:    pc.CommunityID,
				KnownMediaId:   pc.KnownMediaID,
				MagnetUri:      pc.MagnetURI,
				Title:          stringsx.FirstNonBlank(known.Title, lmd.Description),
				Description:    known.Overview,
				ArchivedId:     lmd.ArchiveID,
				Mimetype:       stringsx.FirstNonBlank(known.Mimetype, lmd.Mimetype),
				EncryptionSeed: lmd.EncryptionSeed,
				Bytes:          lmd.Bytes,
			},
		}
		if _, err = metrics.Publish(ctx, &req, io.NopCloser(bytes.NewReader(encoded))); err != nil {
			log.Println(errorsx.Wrap(err, "failed to sync to deeppool"))
			continue
		}

		if err := community.PublishedContentUpdatePublishedAt(ctx, q, pc.ID, time.Now()).Scan(&pc); err != nil {
			log.Println(errorsx.Wrap(err, "failed to update published_at"))
			continue
		}

		if err := RegenerateFeed(ctx, q, publisher, pc.CommunityID); err != nil {
			log.Println(errorsx.Wrapf(err, "feed regeneration failed for community: %s", pc.CommunityID))
			continue
		}

		log.Printf("synced published content %s to deeppool with archive_id %s", pc.ID, lmd.ArchiveID)
	}

	return pending.Err()
}
