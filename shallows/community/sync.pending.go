package community

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type MetricsPublisher interface {
	Publish(ctx context.Context, communityID string, req *meta.PublishContentRequest) (*meta.PublishContentResponse, error)
}

func magnetURI(tmd tracking.Metadata, name string) string {
	ih := int160.FromBytes(tmd.Infohash)
	return metainfo.Magnet{InfoHash: metainfo.Hash(ih.AsByteArray()), DisplayName: name}.String()
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

// SyncPendingToDeeppool syncs pending published content to deeppool.
// Returns a set of community IDs that were affected.
func SyncPendingToDeeppool(ctx context.Context, q sqlx.Queryer, httpc *http.Client, metrics MetricsPublisher, mvfs, tvfs fsx.Virtual) (map[string]struct{}, error) {
	var (
		pending     = sqlx.Scan(PublishedContentFindByPendingSync(ctx, q))
		communities = make(map[string]struct{})
	)

	for pc := range pending.Iter() {
		var lmd library.Metadata

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
		if err := PublishedContentUpdateMagnetURI(ctx, q, pc.ID, pc.MagnetURI).Scan(&pc); err != nil {
			log.Println(errorsx.Wrap(err, "failed to update magnet_uri"))
			continue
		}

		var known library.Known
		if pc.KnownMediaID != "" {
			if err := library.KnownFindByID(ctx, q, pc.KnownMediaID).Scan(&known); sqlx.IgnoreNoRows(err) != nil {
				log.Println(errorsx.Wrap(err, "failed to find known media"))
			}
		}

		if pc.OAuthGoogleID != uuid.Nil.String() {
			if uerr := YouTubeUpload(ctx, q, httpc, mvfs, pc.OAuthGoogleID, lmd, stringsx.FirstNonBlank(known.Title, lmd.Description), known.Overview); uerr != nil {
				log.Println(errorsx.Wrap(uerr, "youtube cross-post failed"))
			}
		}

		if pc.PublishMode == int32(meta.PublishMode_UNLISTED) {
			if err := PublishedContentUpdatePublishedAt(ctx, q, pc.ID, time.Now()).Scan(&pc); err != nil {
				log.Println(errorsx.Wrap(err, "failed to update published_at"))
			}
			continue
		}

		if pc.PublishMode == int32(meta.PublishMode_LISTED) {
			if err := PublishedContentUpdatePublishedAt(ctx, q, pc.ID, time.Now()).Scan(&pc); err != nil {
				log.Println(errorsx.Wrap(err, "failed to update published_at"))
				continue
			}
			communities[pc.CommunityID] = struct{}{}
			log.Printf("synced listed published content %s", pc.ID)
			continue
		}

		if lmd.ArchiveID == uuid.Nil.String() || lmd.ArchiveID == uuid.Max.String() {
			continue
		}

		if _, err := metrics.Publish(ctx, pc.CommunityID, &meta.PublishContentRequest{
			PublishedContent: &meta.PublishedContent{
				Id:             pc.ID,
				KnownMediaId:   pc.KnownMediaID,
				MagnetUri:      pc.MagnetURI,
				ArchivedId:     lmd.ArchiveID,
				Title:          stringsx.FirstNonBlank(known.Title, lmd.Description),
				Description:    known.Overview,
				Mimetype:       stringsx.FirstNonBlank(known.Mimetype, lmd.Mimetype),
				EncryptionSeed: lmd.EncryptionSeed,
			},
		}); err != nil {
			log.Println(errorsx.Wrap(err, "failed to sync to deeppool"))
			continue
		}

		if err := PublishedContentUpdatePublishedAt(ctx, q, pc.ID, time.Now()).Scan(&pc); err != nil {
			log.Println(errorsx.Wrap(err, "failed to update published_at"))
			continue
		}

		communities[pc.CommunityID] = struct{}{}
		log.Printf("synced published content %s to deeppool with archive_id %s", pc.ID, lmd.ArchiveID)
	}

	if err := pending.Err(); err != nil {
		return communities, err
	}

	return communities, nil
}

// NewPendingSync creates a background worker that periodically syncs pending published content.
func NewPendingSync(ctx context.Context, q sqlx.Queryer, httpc *http.Client, metrics MetricsPublisher, published deeppool.Published, mvfs, tvfs fsx.Virtual, interval time.Duration) error {
	log.Println("pending sync worker initiated")
	defer log.Println("pending sync worker completed")

	return timex.NowAndEvery(ctx, interval, func(ctx context.Context) error {
		communities, err := SyncPendingToDeeppool(ctx, q, httpc, metrics, mvfs, tvfs)
		if err != nil {
			log.Println(errorsx.Wrap(err, "pending sync failed"))
			return nil
		}

		if len(communities) > 0 {
			log.Printf("pending sync completed: synced %d communities", len(communities))
		}

		for communityID := range communities {
			if err := RegenerateFeed(ctx, q, published, communityID); err != nil {
				log.Println(errorsx.Wrap(err, "feed regeneration failed for community "+communityID))
			}
		}

		return nil
	})
}
