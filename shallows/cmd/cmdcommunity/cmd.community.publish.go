package cmdcommunity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"os"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/rss"
)

type cmdCommunityPublish struct {
	Title       string        `flag:"" name:"title" help:"title of the rss feed, defaults to the community name"`
	Description string        `flag:"" name:"description" help:"description of the rss feed, defaults to the community description"`
	Timestamp   time.Time     `flag:"" name:"publish date" help:"pub date for the feed" default:"${vars_timestamp_started}"`
	TTL         time.Duration `flag:"" name:"ttl" help:"frequency clients should check the feed in minutes" default:"1440m"`
	Copyright   string        `flag:"" name:"copyright" help:"copyright for the rss feed" default:""`
	DryRun      bool          `flag:"" name:"dry-run" help:"does not actually upload just generates" negatable:"" default:"true"`
	Name        string        `arg:"" name:"name" optional:"" help:"name of the community to upload the feed, if this is unspecified it is assumed the first message will contain the community info"`
}

func (t cmdCommunityPublish) items(c *communityapi.Community, r io.Reader) iter.Seq[rss.Item] {
	ts := time.Now()
	return func(yield func(rss.Item) bool) {
		var (
			derr error
			v    media.Published
		)
		d := jsonl.NewDecoder(r)

		for derr = d.Decode(&v); derr == nil; derr = d.Decode(&v) {
			uri := metainfo.Magnet{
				InfoHash:    metainfo.Hash(errorsx.Must(int160.FromHexEncodedString(v.Id)).AsByteArray()),
				DisplayName: v.Description,
			}

			if !yield(rss.Item{
				Guid:        v.Id,
				Title:       v.Description,
				PublishDate: ts,
				Expires:     timex.RFC3339NanoDecode(errorsx.Zero(grpcx.DecodeTime(v.ExpiresAt))),
				Link:        langx.FirstNonZero(c.Url, fmt.Sprintf("https://%s.community.retrovibe.space/%s", c.Domain, v.Id)),
				Enclosures: []rss.Enclosure{
					{URL: uri.String(), Mimetype: mimex.Bittorrent, Length: v.Bytes},
				},
			}) {
				return
			}
		}

		if errorsx.Ignore(derr, io.EOF) != nil {
			panic(derr)
		}
	}
}

func (t cmdCommunityPublish) Run(gctx *cmdopts.Global, dpc cmdopts.DeeppoolClient) (err error) {
	var (
		buf bytes.Buffer
		com = &communityapi.Community{
			Domain: t.Name,
		}
	)

	c, err := dpc.HTTPClient(gctx.Context)
	if err != nil {
		return err
	}

	if stringsx.Blank(t.Name) {
		debugx.Println("reading community from stdin")
		if err := json.NewDecoder(os.Stdin).Decode(&com); err != nil {
			return err
		}
	} else {
		info, err := metaapi.CommunityInfo(gctx.Context, c, t.Name)
		if err != nil {
			return errorsx.Wrap(err, "failed to retrieve community metadata")
		}
		com = info.Community
	}

	err = errorsx.Wrap(rss.Generator().Generate(io.MultiWriter(&buf), rss.Channel{
		Link:          langx.FirstNonZero(com.Url, fmt.Sprintf("https://%s.community.retrovibe.space", com.Domain)),
		TTL:           int(t.TTL.Minutes()),
		Title:         com.Domain,
		LastBuildDate: t.Timestamp,
		Language:      "en-us",
		Description:   com.Description,
		Copyright:     t.Copyright,
		Retrovibed:    &rss.Retrovibed{Mimetype: com.Mimetype, Entropy: com.Entropy},
	}, t.items(com, os.Stdin)), "failed to generate rss feed")
	if err != nil {
		return err
	}

	if t.DryRun {
		_, err = io.Copy(os.Stdout, &buf)
		return err
	}

	uploaded, err := metaapi.CommunityPublish(gctx.Context, c, com.Id, &buf)
	if err != nil {
		return err
	}

	if err = json.NewEncoder(os.Stdout).Encode(uploaded.Community); err != nil {
		return errorsx.Wrap(err, "unable to write to uploaded")
	}

	return nil
}
