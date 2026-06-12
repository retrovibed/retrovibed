package cmdlibrary

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type exportJSONL struct {
	Database   string   `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	MediaDir   string   `flag:"" name:"mediadir" help:"root directory of media storage" default:"${vars_media_directory}"`
	KnownMedia *bool    `flag:"" name:"known-media" help:"filter by known media presence" negatable:""`
	Torrent    *bool    `flag:"" name:"torrent" help:"filter by torrent presence" negatable:""`
	ID         []string `flag:"" name:"id" help:"export only the specified media id(s)" optional:""`
}

func (t exportJSONL) Run(gctx *cmdopts.Global) (err error) {
	db, err := cmdopts.DatabaseCustom(gctx.Context, t.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	return t.run(gctx.Context, db, fsx.DirVirtual(t.MediaDir), os.Stdout)
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
	q = q.Where(library.MetadataQueryByIDs(t.ID...))

	scanner := sqlx.Scan(library.MetadataSearch(ctx, db, q))
	var exported uint64
	for md := range scanner.Iter() {
		if err := t.export(enc, vfs, md); err != nil {
			return err
		}

		if exported += 1; exported%256 == 0 {
			log.Println("exported", exported, "records")
		}
	}
	return scanner.Err()
}

func (t exportJSONL) export(enc *json.Encoder, vfs fsx.Virtual, md library.Metadata) error {
	const chunksize = 16 * bytesx.MiB
	dcache, err := blockcache.NewDirectoryCache(vfs.Path(md.ID))
	if err != nil {
		return errorsx.Wrap(err, "open block cache")
	}

	numChunks := uint64(0)
	if md.Bytes > 0 {
		numChunks = (md.Bytes + chunksize - 1) / chunksize
	}

	if err := enc.Encode(exportHeader{Chunks: numChunks}); err != nil {
		return err
	}

	h := md5.New()
	b := bytes.NewBuffer(make([]byte, chunksize))
	w := io.MultiWriter(h, b)

	for i := uint64(0); i < numChunks; i++ {
		b.Reset()
		off := int64(md.DiskOffset) + int64(i)*chunksize

		_, err := io.CopyN(w, io.NewSectionReader(dcache, off, chunksize), chunksize)
		if errorsx.Ignore(err, io.EOF) != nil {
			failed := enc.Encode(exportTrailer{Metadata: langx.Clone(md, timex.JSONSafeEncodeOption), MD5: uuid.Nil.String()})
			return errors.Join(errorsx.Wrapf(err, "read chunk %d at offset %d", i, off), failed)
		}

		if err := enc.Encode(exportChunk{Data: b.Bytes()}); err != nil {
			return err
		}
	}

	return enc.Encode(exportTrailer{
		Metadata: langx.Clone(md, timex.JSONSafeEncodeOption),
		MD5:      hex.EncodeToString(h.Sum(nil)),
	})
}
