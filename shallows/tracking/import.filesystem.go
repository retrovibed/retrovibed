package tracking

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/slicesx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/library"
)

func ImportTorrent(q sqlx.Queryer, mvfs, tvfs fsx.Virtual) library.ImportOp {
	return func(ctx context.Context, path string) (*library.Transfered, error) {
		tx, err := library.TransferedFromPath(path)
		if err != nil {
			return nil, err
		}

		md, err := torrent.NewFromMetaInfoFile(tx.Path)
		if err != nil {
			return nil, err
		}

		src, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer src.Close()

		dst, err := os.CreateTemp(tvfs.Path(), "importing.*.bin")
		if err != nil {
			return nil, err
		}
		defer os.Remove(dst.Name())
		defer dst.Close()

		if n, err := io.Copy(io.MultiWriter(tx.MD5, dst), src); err != nil {
			return nil, err
		} else {
			tx.Bytes = uint64(n)
		}

		var (
			info metainfo.Info
			tmd  Metadata
		)

		if err := bencode.NewDecoder(bytes.NewReader(md.InfoBytes)).Decode(&info); err != nil {
			return nil, errorsx.Wrap(err, "unable to extract info")
		}

		if err = MetadataInsertWithDefaults(ctx, q, NewMetadata(&md.ID, MetadataOptionFromInfo(&info))).Scan(&tmd); err != nil {
			return nil, errorsx.Wrap(err, "unable to insert metadata")
		}

		uid := md.ID.HexString()

		for _, finfo := range info.Files {
			var lmd library.Metadata

			if err := library.MetadataAssociateTorrent(ctx, q, slicesx.LastOrZero(finfo.Path...), tmd.ID).Scan(&lmd); sqlx.IgnoreNoRows(err) != nil {
				return nil, errorsx.Wrap(err, "unable to retrieve metadata")
			} else if sqlx.ErrNoRows(err) != nil {
				log.Println("unable to match", tmd.ID, finfo.Path, slicesx.LastOrZero(finfo.Path...), "with existing media")
				// ignore we can't associate
				continue
			}

			dstp := tvfs.Path(uid, filepath.Join(finfo.Path...))
			if err := os.MkdirAll(filepath.Dir(dstp), 0700); err != nil {
				return nil, errorsx.Wrap(err, "unable to create torrent file")
			}

			if err := fsx.RemoveSymlink(dstp); err != nil {
				return nil, errorsx.WithStack(err)
			}

			if err := os.Symlink(mvfs.Path(lmd.ID), dstp); err != nil {
				return nil, errorsx.Wrap(err, "unable to associate torrent file with library")
			}
		}

		if err := os.Rename(dst.Name(), tvfs.Path(fmt.Sprintf("%s.torrent", uid))); err != nil {
			return nil, errorsx.Wrap(err, "unable to symlink to original location")
		}

		if err := MetadataDownloadByID(ctx, q, tmd.ID).Scan(&tmd); err != nil {
			return nil, errorsx.Wrap(err, "unable to mark torrent for downloading")
		}

		return tx, nil
	}
}
