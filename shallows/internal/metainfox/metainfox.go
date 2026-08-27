package metainfox

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/bytesx"
	"github.com/retrovibed/retrovibed/retroapi/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

type Printer struct {
	md *metainfo.MetaInfo
}

func NewPrinter(md *metainfo.MetaInfo) Printer {
	return Printer{
		md: md,
	}
}

func (t Printer) Print(dst io.Writer) error {
	info, err := t.md.UnmarshalInfo()
	if err != nil {
		return err
	}

	c := fsx.NewWriteErrCompact()
	tw := tabwriter.NewWriter(dst, 1, 0, 2, ' ', 0)

	c.Compact(fmt.Fprintf(tw, "Name:\t%s\n", info.Name))
	c.Compact(fmt.Fprintf(tw, "InfoHash:\t%s\n", t.md.HashInfoBytes()))
	c.Compact(fmt.Fprintf(tw, "Comment:\t%s\n", t.md.Comment))
	c.Compact(fmt.Fprintf(tw, "Created By:\t%s\n", t.md.CreatedBy))
	c.Compact(fmt.Fprintf(tw, "Creation Date:\t%s\n", time.Unix(t.md.CreationDate, 0)))
	c.Compact(fmt.Fprintf(tw, "Encoding:\t%s\n", t.md.Encoding))
	c.Compact(fmt.Fprintf(tw, "Private:\t%t\n", info.Private != nil && *info.Private))
	c.Compact(fmt.Fprintf(tw, "Piece Length:\t%s\n", bytesx.Unit(info.PieceLength)))
	c.Compact(fmt.Fprintf(tw, "Pieces:\t%d\n", info.NumPieces()))
	c.Compact(fmt.Fprintf(tw, "Total Length:\t%s\n", bytesx.Unit(info.TotalLength())))

	c.Compact(fmt.Fprintf(tw, "Announce List:\t\n"))
	for _, tier := range t.md.UpvertedAnnounceList() {
		c.Compact(fmt.Fprintf(tw, "\t%s\n", strings.Join(tier, ", ")))
	}

	c.Compact(fmt.Fprintf(tw, "Files:\t\n"))
	for _, fi := range info.UpvertedFiles() {
		c.Compact(fmt.Fprintf(tw, "\t%s\t%s\n", fi.DisplayPath(&info), bytesx.Unit(fi.Length)))
	}

	return langx.FirstNonNil(tw.Flush(), c.Err())
}
