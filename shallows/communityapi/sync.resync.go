package communityapi

import (
	"context"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// ResyncOne refreshes a single community's metadata from deeppool and inline-syncs
// up to 128 of its published content items.
func ResyncOne(ctx context.Context, q sqlx.Queryer, httpc *http.Client, cid string) (*community.Community, error) {
	var existing community.Community

	com, err := NewDeeppoolCommunity(httpc).Find(ctx, cid)
	if err != nil {
		return nil, errorsx.Wrapf(err, "unable to find community from deeppool - %s", cid)
	}

	if err := community.CommunityInsertWithDefaults(ctx, q, CommunityFromDeeppool(com)).Scan(&existing); err != nil {
		return nil, errorsx.Wrap(err, "unable to refresh community metadata")
	}

	if _, err := SyncContentFromDeeppool(ctx, q, NewDeeppoolPublished(httpc), cid, existing.AutoDownload != 0, 128); err != nil {
		return nil, errorsx.Wrap(err, "unable to sync published content from deeppool")
	}

	existing.LastSyncAt = time.Now()
	if err := community.CommunityUpdateLastSyncAt(ctx, q, existing).Scan(&existing); err != nil {
		return nil, errorsx.Wrap(err, "unable to update last sync time")
	}

	return &existing, nil
}
