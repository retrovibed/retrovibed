package rss

import (
	"encoding/xml"
	"time"

	"github.com/retrovibed/retrovibed/internal/x/errorsx"
)

// rss represents the shcema of an RSS feed
type rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel channel  `xml:"channel"`
}

// in the oficial shcema channel contains more than just `item`
// but there is no need to use those fields
type channel struct {
	XMLName xml.Name `xml:"channel"`
	Title   string   `xml:"title"`
	Items   []item   `xml:"item"`
	TTL     int      `xml:"ttl"`
}

// item represent the actual feed for each news
type item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     xmlTime  `xml:"pubDate"`
	Source      *source  `xml:"source"`
}

func parseTimestamp(encoded string) (_ time.Time, err error) {
	formats := []string{
		time.Layout,
		time.RFC822,
		time.RFC850,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	for _, format := range formats {
		if ts, failed := time.Parse(format, encoded); failed == nil {
			return ts, nil
		} else {
			err = errorsx.Compact(err, failed)
		}
	}

	return time.Time{}, err
}

// this is for custom unmarshaling of date
type xmlTime struct {
	value    time.Time
	hasValue bool
}

func (t *xmlTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var encoded string
	d.DecodeElement(&encoded, &start)
	ts, err := parseTimestamp(encoded)
	if err != nil {
		*t = xmlTime{hasValue: false}
		return err
	}
	*t = xmlTime{value: ts, hasValue: true}
	return nil
}

// this represents a cource tag
type source struct {
	XMLName xml.Name `xml:"source"`
	URL     string   `xml:"url,attr"`
	Value   string   `xml:",chardata"`
}

// Item is the representation of an item
// retrieved from an RSS feed
type Item struct {
	Title       string    // Defines the title of the item
	Source      string    // Specifies a third-party source for the item
	SourceURL   string    // Specifies the link to the source
	Link        string    // Defines the hyperlink to the item
	PublishDate time.Time // Defines the last-publication date for the item
	Description string    // Describes the item
}
