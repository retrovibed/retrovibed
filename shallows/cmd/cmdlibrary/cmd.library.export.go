package cmdlibrary

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type exportJSONL struct {
	Database   string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	KnownMedia *bool  `flag:"" name:"known-media" help:"filter by known media presence" negatable:""`
	Torrent    *bool  `flag:"" name:"torrent" help:"filter by torrent presence" negatable:""`
}

func (t exportJSONL) Run(gctx *cmdopts.Global) (err error) {
	db, err := cmdopts.DatabaseCustom(gctx.Context, t.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	return t.run(gctx.Context, db, fsx.DirVirtual(env.TorrentDir()), os.Stdout)
}

func (t exportJSONL) run(ctx context.Context, db sqlx.Queryer, vfs fsx.Virtual, w io.Writer) error {
	enc := json.NewEncoder(w)

	q := library.MetadataSearchBuilder().Where(library.MetadataQueryNotTombstoned()).OrderBy("id ASC")
	if t.KnownMedia != nil {
		q = q.Where(library.MetadataQueryHasKnownMedia(*t.KnownMedia))
	}
	if t.Torrent != nil {
		q = q.Where(library.MetadataQueryHasTorrent(*t.Torrent))
	}

	scanner := sqlx.Scan(library.MetadataSearch(ctx, db, q))
	var exported uint64
	for md := range scanner.Iter() {
		if err := exportItem(enc, vfs, md); err != nil {
			return err
		}
		exported++
		if exported%256 == 0 {
			log.Println("exported", exported, "records")
		}
	}
	return scanner.Err()
}

func exportItem(enc *json.Encoder, vfs fsx.Virtual, md library.Metadata) error {
	dcache, err := blockcache.NewDirectoryCache(vfs.Path(md.ID))
	if err != nil {
		return errorsx.Wrap(err, "open block cache")
	}

	numChunks := uint64(0)
	if md.Bytes > 0 {
		numChunks = (md.Bytes + exportChunkSize - 1) / exportChunkSize
	}

	if err := enc.Encode(exportHeader{Chunks: numChunks}); err != nil {
		return err
	}

	h := md5.New()
	buf := make([]byte, exportChunkSize)
	for i := uint64(0); i < numChunks; i++ {
		off := int64(md.DiskOffset) + int64(i)*exportChunkSize
		n, readErr := dcache.ReadAt(buf, off)
		if readErr != nil && readErr != io.EOF {
			timex.JSONSafeEncodeOption(&md)
			_ = enc.Encode(exportTrailer{Metadata: md, MD5: ""})
			return errorsx.Wrapf(readErr, "read chunk %d at offset %d", i, off)
		}
		h.Write(buf[:n])
		if err := enc.Encode(exportChunk{Data: buf[:n]}); err != nil {
			return err
		}
	}

	timex.JSONSafeEncodeOption(&md)
	return enc.Encode(exportTrailer{
		Metadata: md,
		MD5:      hex.EncodeToString(h.Sum(nil)),
	})
}
