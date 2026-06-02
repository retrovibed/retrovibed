package community

import (
	"context"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// SyncContentFromDeeppool fetches published content for a community from deeppool
// and imports it into the local database. Returns the number of items synced.
func SyncContentFromDeeppool(ctx context.Context, q sqlx.Queryer, client deeppool.Published, communityID string, autodownload bool) (int, error) {
	resp, err := client.List(ctx, communityID)
	if err != nil {
		return 0, errorsx.Wrap(err, "failed to fetch published content from deeppool")
	}

	synced := 0
	for _, pc := range resp.Items {
		if err := SyncPublishedContentItem(ctx, q, pc, autodownload); err != nil {
			log.Println("failed to sync published content", pc.Id, err)
			continue
		}
		synced++
	}

	return synced, nil
}

// SyncPublishedContentItem syncs a single published content item from deeppool into the local database.
func SyncPublishedContentItem(ctx context.Context, q sqlx.Queryer, pc *meta.PublishedContent, autodownload bool) error {
	md, err := metainfo.ParseMagnetURI(pc.MagnetUri)
	if err != nil {
		return errorsx.Wrap(err, "failed to parse magnet URI")
	}

	tmeta := tracking.NewMetadata(
		langx.Autoptr(int160.FromByteArray(md.InfoHash)),
		tracking.MetadataOptionFromMagnet(&md),
		tracking.MetadataOptionDescription(md.DisplayName),
		tracking.MetadataOptionKnownMediaID(stringsx.FirstNonBlank(pc.KnownMediaId, uuid.Max.String())),
		tracking.MetadataOptionEntropySeed(md.InfoHash.Bytes()),
		tracking.MetadataOptionAutoDescription,
		tracking.MetadataOptionAutoHidden,
	)

	if err = tracking.MetadataInsertWithDefaults(ctx, q, tmeta).Scan(&tmeta); err != nil {
		return errorsx.Wrap(err, "failed to insert torrent metadata")
	}

	if autodownload {
		if err = tracking.MetadataAutoDownloadByID(ctx, q, tmeta.ID).Scan(&tmeta); err != nil {
			return errorsx.Wrap(err, "failed to mark torrent for automatic download")
		}
	}

	dbpc := NewPublishedContent(PublishedContent{
		ID:            pc.Id,
		CommunityID:   pc.CommunityId,
		MagnetURI:     pc.MagnetUri,
		LibraryID:     stringsx.FirstNonBlank(pc.LibraryId, uuid.Nil.String()),
		OAuthGoogleID: pc.OauthGoogleId,
		KnownMediaID:  tmeta.KnownMediaID,
	})

	if err = PublishedContentInsertWithDefaults(ctx, q, dbpc).Scan(&dbpc); err != nil {
		return errorsx.Wrap(err, "failed to insert published content")
	}

	return nil
}

// NewSubscriptionSync creates a background worker that periodically syncs
// content from all subscribed communities.
func NewSubscriptionSync(ctx context.Context, q sqlx.Queryer, client deeppool.Published, interval time.Duration) error {
	log.Println("subscription sync worker initiated")
	defer log.Println("subscription sync worker completed")

	return timex.NowAndEvery(ctx, interval, func(ctx context.Context) error {
		subs := sqlx.Scan(CommunitySubscriptionFindAll(ctx, q))
		for sub := range subs.Iter() {
			autodownload := sub.AutoDownload != 0
			synced, err := SyncContentFromDeeppool(ctx, q, client, sub.CommunityID, autodownload)
			if err != nil {
				log.Println(errorsx.Wrap(err, "subscription sync failed for "+sub.CommunityID))
				continue
			}

			if synced > 0 {
				log.Printf("subscription sync: imported %d items for %s", synced, sub.CommunityID)
			}

			if err = CommunitySubscriptionUpdateLastSyncAt(ctx, q, sub.CommunityID, time.Now()).Scan(&sub); err != nil {
				log.Println(errorsx.Wrap(err, "failed to update last_sync_at for "+sub.CommunityID))
			}
		}

		if err := subs.Err(); err != nil {
			log.Println(errorsx.Wrap(err, "subscription sync iteration failed"))
		}

		return nil
	})
}
