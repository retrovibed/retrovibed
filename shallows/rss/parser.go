package rss

import (
	"context"
	"encoding/xml"
	"io"
	"net/url"
)

func Parse(ctx context.Context, r io.Reader) (*channel, []Item, error) {
	return parseData(r, "")
}

func parseData(data io.Reader, originURL string) (*channel, []Item, error) {
	var r rss
	if err := xml.NewDecoder(data).Decode(&r); err != nil {
		return nil, nil, err
	}

	rssItems := make([]Item, 0, len(r.Channel.Items))
	for _, item := range r.Channel.Items {
		if item.Title == "" || item.Link == "" || item.Description == "" {
			continue
		}

		rssItem := Item{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
		}

		if item.PubDate.hasValue {
			rssItem.PublishDate = item.PubDate.value
		}

		if item.Source != nil {
			rssItem.Source = item.Source.Value
			rssItem.SourceURL = item.Source.URL
		} else {
			host := extractSource(originURL)
			rssItem.Source = host
			rssItem.SourceURL = originURL
		}

		rssItems = append(rssItems, rssItem)
	}

	return &r.Channel, rssItems, nil
}

func extractSource(urlRaw string) string {
	u, err := url.Parse(urlRaw)
	if err != nil {
		return ""
	}

	return u.Hostname()
}
