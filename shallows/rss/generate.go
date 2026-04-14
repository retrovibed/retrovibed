package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"iter"
	"strconv"
	"strings"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

type GenOption func(*gen)

func Generator(options ...GenOption) gen {
	return langx.Clone(gen{
		timeFormat: time.RFC1123Z,
	}, options...)
}

type gen struct {
	timeFormat string
}

func (t gen) Generate(w io.Writer, c Channel, items iter.Seq[Item]) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return errorsx.Wrap(err, "failed to write XML header")
	}

	encoder := xml.NewEncoder(w)
	encoder.Indent("", "\t")

	rssStart := xml.StartElement{
		Name: xml.Name{Local: "rss"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "version"}, Value: "2.0"},
			{Name: xml.Name{Local: "xmlns:retrovibed"}, Value: "https://retrovibed.com/rss/feed"},
			{Name: xml.Name{Local: "xmlns:atom"}, Value: "http://www.w3.org/2005/Atom"},
			{Name: xml.Name{Local: "xmlns:content"}, Value: "http://purl.org/rss/1.0/modules/content/"},
			{Name: xml.Name{Local: "xmlns:media"}, Value: "http://search.yahoo.com/mrss/"},
		},
	}

	if err := encoder.EncodeToken(rssStart); err != nil {
		return errorsx.Wrap(err, "failed to encode RSS start tag")
	}

	if err := encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "channel"}}); err != nil {
		return errorsx.Wrap(err, "failed to encode channel start tag")
	}

	if err := encoder.EncodeElement(c.Title, xml.StartElement{Name: xml.Name{Local: "title"}}); err != nil {
		return errorsx.Wrap(err, "failed to encode channel title")
	}

	if err := encoder.EncodeElement(c.Link, xml.StartElement{Name: xml.Name{Local: "link"}}); err != nil {
		return errorsx.Wrap(err, "failed to encode channel link")
	}

	if err := encoder.EncodeElement(c.Description, xml.StartElement{Name: xml.Name{Local: "description"}}); err != nil {
		return errorsx.Wrap(err, "failed to encode channel description")
	}

	if err := encoder.EncodeElement(c.Copyright, xml.StartElement{Name: xml.Name{Local: "copyright"}}); err != nil {
		return errorsx.Wrap(err, "failed to encode channel copyright")
	}

	if err := encoder.EncodeElement(c.Language, xml.StartElement{Name: xml.Name{Local: "language"}}); err != nil {
		return errorsx.Wrap(err, "failed to encode channel language")
	}

	if err := encoder.EncodeElement(c.LastBuildDate.Format(t.timeFormat), xml.StartElement{Name: xml.Name{Local: "lastBuildDate"}}); err != nil {
		return errorsx.Wrap(err, "failed to encode channel lastBuildDate")
	}

	if err := encoder.EncodeElement(c.TTL, xml.StartElement{Name: xml.Name{Local: "ttl"}}); err != nil {
		return errorsx.Wrap(err, "failed to encode channel ttl")
	}

	if c.Retrovibed != nil {
		if err := encoder.EncodeElement(c.Retrovibed, xml.StartElement{Name: xml.Name{Local: "retrovibed:metadata"}}); err != nil {
			return errorsx.Wrap(err, "failed to encode channel ttl")
		}
	}

	for i := range items {
		if err := genitem(encoder, i); err != nil {
			return err
		}
	}

	if err := encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "channel"}}); err != nil {
		return errorsx.Wrap(err, "failed to encode channel end tag")
	}

	if err := encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "rss"}}); err != nil {
		return errorsx.Wrap(err, "failed to encode RSS end tag")
	}

	if err := encoder.Flush(); err != nil {
		return errorsx.Wrap(err, "failed to flush encoder")
	}

	if _, err := w.Write([]byte("\n")); err != nil {
		return errorsx.Wrap(err, "failed to include the trailing newline")
	}

	return nil
}

func genitem(encoder *xml.Encoder, item Item) error {
	itemStart := xml.StartElement{Name: xml.Name{Local: "item"}}
	if err := encoder.EncodeToken(itemStart); err != nil {
		return errorsx.Wrap(err, "failed to encode item start tag")
	}

	if err := encoder.EncodeElement(item.Guid, xml.StartElement{Name: xml.Name{Local: "guid"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "isPermaLink"}, Value: strconv.FormatBool(strings.HasPrefix(item.Guid, "http") || strings.HasPrefix(item.Guid, "https"))},
	}}); err != nil {
		return errorsx.Wrap(err, "failed to encode item GUID")
	}

	if err := encoder.EncodeElement(item.Title, xml.StartElement{Name: xml.Name{Local: "title"}}); err != nil {
		return errorsx.Wrap(err, "failed to encode item title")
	}

	if item.Source != nil {
		sourceStart := xml.StartElement{
			Name: xml.Name{Local: "source"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "url"}, Value: item.Source.URL},
			},
		}
		if err := encoder.EncodeToken(sourceStart); err != nil {
			return errorsx.Wrap(err, "failed to encode source start tag")
		}
		if err := encoder.EncodeToken(xml.CharData(item.Source.Description)); err != nil {
			return errorsx.Wrap(err, "failed to encode source character data")
		}
		if err := encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "source"}}); err != nil {
			return errorsx.Wrap(err, "failed to encode source end tag")
		}
	}

	if stringsx.Present(item.Link) {
		if err := encoder.EncodeElement(item.Link, xml.StartElement{Name: xml.Name{Local: "link"}}); err != nil {
			return errorsx.Wrap(err, "failed to encode item link")
		}
	}

	if !item.PublishDate.Equal(time.Time{}) {
		if err := encoder.EncodeElement(item.PublishDate.Format(time.RFC1123Z), xml.StartElement{Name: xml.Name{Local: "pubDate"}}); err != nil {
			return errorsx.Wrap(err, "failed to encode item pubDate")
		}
	}

	if stringsx.Present(item.Description) {
		if err := encoder.EncodeElement(item.Description, xml.StartElement{Name: xml.Name{Local: "description"}}); err != nil {
			return errorsx.Wrap(err, "failed to encode item description")
		}
	}

	if !item.Expires.IsZero() && !timex.Inf().Equal(item.Expires) && !timex.NegInf().Equal(item.Expires) {
		if err := encoder.EncodeElement(item.Expires.Format(time.RFC1123Z), xml.StartElement{Name: xml.Name{Local: "retrovibed:expires"}}); err != nil {
			return errorsx.Wrap(err, "failed to encode item expires")
		}
	}

	for _, enc := range item.Enclosures {
		enclosureStart := xml.StartElement{
			Name: xml.Name{Local: "enclosure"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "url"}, Value: enc.URL},
				{Name: xml.Name{Local: "type"}, Value: enc.Mimetype},
				{Name: xml.Name{Local: "length"}, Value: fmt.Sprintf("%d", enc.Length)}, // Convert uint64 to string
			},
		}
		if err := encoder.EncodeElement("", enclosureStart); err != nil {
			return errorsx.Wrap(err, "failed to encode enclosure element")
		}
	}

	itemEnd := xml.EndElement{Name: xml.Name{Local: "item"}}
	if err := encoder.EncodeToken(itemEnd); err != nil {
		return errorsx.Wrap(err, "failed to encode item end tag")
	}

	if err := encoder.Flush(); err != nil {
		return errorsx.Wrap(err, "failed to flush encoder after item")
	}

	return nil
}
