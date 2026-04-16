package community

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"slices"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/rss"
)

type FeedPublisher interface {
	Find(ctx context.Context, communityID string) (*meta.Community, error)
	UploadFeed(ctx context.Context, communityID string, feed io.Reader) error
}

const (
	feedDefaultTTL      = 24 * 60
	feedDefaultLanguage = "en"
)

// RegenerateFeed generates an RSS feed for a community and uploads it to deeppool.
func RegenerateFeed(ctx context.Context, q sqlx.Queryer, published FeedPublisher, communityID string) error {
	c, err := published.Find(ctx, communityID)
	if err != nil {
		return errorsx.Wrap(err, "failed to find community")
	}

	items, err := buildFeedItems(ctx, q, c)
	if err != nil {
		return errorsx.Wrap(err, "failed to build feed items")
	}

	buf := new(bytes.Buffer)
	channel := rss.Channel{
		Title:         c.Domain,
		Link:          fmt.Sprintf("https://%s.community.retrovibe.space", c.Domain), // TODO: return the fully qualified url from the find function. aka add a uri field.
		Description:   c.Description,
		TTL:           feedDefaultTTL,
		LastBuildDate: time.Now().UTC(),
		Language:      feedDefaultLanguage,
		Retrovibed:    &rss.Retrovibed{Entropy: c.Entropy, Mimetype: c.Mimetype},
	}
	if err := rss.Generator().Generate(buf, channel, slices.Values(items)); err != nil {
		return errorsx.Wrap(err, "failed to generate RSS feed")
	}

	if err := published.UploadFeed(ctx, communityID, buf); err != nil {
		return errorsx.Wrap(err, "failed to upload RSS feed")
	}

	log.Printf("regenerated feed for community %s", communityID)
	return nil
}

func buildFeedItems(ctx context.Context, q sqlx.Queryer, community *meta.Community) ([]rss.Item, error) {
	var items []rss.Item

	scanner := sqlx.Scan(PublishedContentFindByCommunityIDForFeed(ctx, q, community.Id))
	for pc := range scanner.Iter() {
		var (
			known library.Known
			lmd   library.Metadata
		)

		if err := library.KnownFindByID(ctx, q, pc.KnownMediaID).Scan(&known); sqlx.IgnoreNoRows(err) != nil {
			continue
		}

		if err := library.MetadataFindByID(ctx, q, pc.LibraryID).Scan(&lmd); sqlx.IgnoreNoRows(err) != nil {
			continue
		}

		items = append(items, rss.Item{
			Guid:        pc.ID,
			Title:       stringsx.FirstNonBlank(known.Title, known.OriginalTitle, lmd.Description),
			Link:        pc.MagnetURI,
			PublishDate: pc.PublishedAt,
			Description: known.Overview,
			Enclosures: []rss.Enclosure{
				{
					URL:      pc.MagnetURI,
					Mimetype: lmd.Mimetype,
					Length:   lmd.Bytes,
				},
			},
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
