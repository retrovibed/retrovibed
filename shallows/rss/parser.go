package rss

import (
	"context"
	"crypto/md5"
	"encoding/xml"
	"hash"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

func Parse(ctx context.Context, r io.Reader) (hash.Hash, *channel, []Item, error) {
	return parseData(r, "")
}

func parseData(data io.Reader, originURL string) (hash.Hash, *channel, []Item, error) {
	var (
		r    Root
		hash = md5.New()
	)

	if err := xml.NewDecoder(data).Decode(&r); err != nil {
		return nil, nil, nil, err
	}

	r.Channel.Title = strings.TrimSpace(r.Channel.Title)

	rssItems := make([]Item, 0, len(r.Channel.Items))
	for _, item := range r.Channel.Items {
		if item.Title == "" && item.Link == "" && item.Description == "" {
			log.Println("skipping item", spew.Sdump(item))
			continue
		}

		if _, err := hash.Write([]byte(item.GUID)); err != nil {
			return nil, nil, nil, errorsx.Wrap(err, "failed to write guid into digest")
		}

		rssItem := Item{
			Title:       strings.TrimSpace(item.Title),
			Description: item.Description,
			Link:        item.Link,
			Enclosures:  item.Enclosures,
			PublishDate: item.PubDate.Timestamp(timex.Inf()), // use inf to mark as unknown published dates. allows for item.Published.Before(time.Now()).
			Expires:     item.Expires.Timestamp(timex.Inf()), // use inf to mark as unknown expires dates. allows for item.Published.Before(time.Now()).
		}

		if item.Source != nil {
			rssItem.Source = &Source{
				Description: item.Source.Value,
				URL:         item.Source.URL,
			}
		} else {
			host := extractSource(originURL)
			rssItem.Source = &Source{
				Description: host,
				URL:         originURL,
			}
		}

		rssItems = append(rssItems, rssItem)
	}

	return hash, &r.Channel, rssItems, nil
}

func extractSource(urlRaw string) string {
	u, err := url.Parse(urlRaw)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

// FindBittorrentEnclosures returns enclosures that represent bittorrent-transportable
// content: magnet URIs (whose type attribute commonly describes the content
// rather than the transport) or enclosures explicitly tagged with the
// bittorrent mimetype.
func FindBittorrentEnclosures(items ...Item) (enc []Enclosure) {
	for _, i := range items {
		for _, e := range i.Enclosures {
			u, err := url.Parse(e.URL)
			if err != nil {
				continue
			}

			if u.Scheme == "magnet" || e.Mimetype == mimex.Bittorrent {
				enc = append(enc, e)
			}
		}
	}

	return enc
}

func ItemToEnclosure(i Item, mimetype string) Enclosure {
	return Enclosure{
		URL:      i.Link,
		Mimetype: mimetype,
	}
}
