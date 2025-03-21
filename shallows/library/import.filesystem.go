package library

import (
	"context"
	"crypto/md5"
	"errors"
	"hash"
	"io"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/retrovibed/retrovibed/internal/asynccompute"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/md5x"
)

const ErrIterationFailed = errorsx.String("failed to yield transferred media")

type Transfered struct {
	Path     string
	Mimetype *mimetype.MIME
	MD5      hash.Hash
	Bytes    uint64
}

type ImportOp = func(ctx context.Context, path string) (*Transfered, error)

func TransferedFromPath(path string) (*Transfered, error) {
	cmimetype, err := mimetype.DetectFile(path)
	if err != nil {
		return nil, err
	}

	return &Transfered{
		Path:     path,
		MD5:      md5.New(),
		Mimetype: cmimetype,
	}, nil
}

func ImportSymlinkFile(vfs fsx.Virtual) ImportOp {
	return func(ctx context.Context, path string) (*Transfered, error) {
		tx, err := TransferedFromPath(path)
		if err != nil {
			return nil, err
		}

		src, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer src.Close()

		if n, err := io.Copy(tx.MD5, src); err != nil {
			return nil, err
		} else {
			tx.Bytes = uint64(n)
		}

		uid := md5x.FormatString(tx.MD5)

		if err := os.Remove(vfs.Path(uid)); fsx.IgnoreIsNotExist(err) != nil {
			return nil, errorsx.Wrap(err, "unable to ensure symlink destination is available")
		}

		if err := os.Symlink(tx.Path, vfs.Path(uid)); err != nil {
			return nil, errorsx.Wrap(err, "unable to symlink to original location")
		}

		return tx, nil
	}
}

func ImportCopyFile(vfs fsx.Virtual) ImportOp {
	return func(ctx context.Context, path string) (*Transfered, error) {
		tx, err := TransferedFromPath(path)
		if err != nil {
			return nil, err
		}

		src, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer src.Close()

		dst, err := os.CreateTemp(vfs.Path(), "importing.*.bin")
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
		uid := md5x.FormatString(tx.MD5)

		if err := os.Remove(vfs.Path(uid)); fsx.IgnoreIsNotExist(err) != nil {
			return nil, errorsx.Wrap(err, "unable to ensure destination is available")
		}

		if err := os.Rename(dst.Name(), vfs.Path(uid)); err != nil {
			return nil, errorsx.Wrap(err, "unable to symlink to original location")
		}

		return tx, nil
	}
}

func ImportFileDryRun(ctx context.Context, path string) (*Transfered, error) {
	return TransferedFromPath(path)
}

func NewImporter(op ImportOp, options ...asynccompute.Option[string]) *importer {
	return &importer{
		op:          op,
		computeopts: options,
	}
}

type importer struct {
	op          ImportOp
	computeopts []asynccompute.Option[string]
}

func (t importer) Import(ctx context.Context, paths ...string) iter.Seq2[*Transfered, error] {
	return func(yield func(*Transfered, error) bool) {
		var (
			capture sync.Once
			err     error
		)
		results := make(chan *Transfered)
		arena := asynccompute.New(func(ictx context.Context, path string) error {
			if info, err := os.Stat(path); err != nil {
				return err
			} else if info.IsDir() {
				return nil
			}

			tx, cause := t.op(ictx, path)
			if cause != nil {
				capture.Do(func() {
					err = errorsx.Compact(err, cause)
				})
				errorsx.Log(cause)
				return cause
			}

			select {
			case results <- tx:
				return nil
			case <-ictx.Done():
				return ctx.Err()
			}
		}, t.computeopts...)

		go func() {
			defer func() {
				ictx, done := context.WithTimeout(context.Background(), time.Minute)
				defer done()
				err = errorsx.Compact(err, asynccompute.Shutdown(ictx, arena))
				close(results)
			}()

			for _, p := range paths {
				if info, cause := os.Stat(p); errors.Is(cause, os.ErrNotExist) {
					err = errorsx.Wrap(cause, "ignoring")
					return
				} else if cause != nil {
					err = errorsx.Wrapf(cause, "failed %s", p)
					return
				} else if !info.IsDir() {
					if _, cause := arena.Run(ctx, p); cause != nil {
						err = errorsx.Wrapf(cause, "failed %s", p)
						return
					}

					continue
				}

				cause := fs.WalkDir(os.DirFS(p), ".", func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}

					if path == "." {
						return nil
					}

					_, cause := arena.Run(ctx, filepath.Join(p, path))
					return cause
				})

				if cause != nil {
					err = errorsx.Wrapf(cause, "filesystem traversal failed")
				}
			}
		}()

		for r := range results {
			if !yield(r, nil) {
				return
			}
		}

		if err != nil {
			if !yield(nil, err) {
				return
			}
		}
	}
}

func ImportFilesystem(ctx context.Context, op ImportOp, paths ...string) iter.Seq2[*Transfered, error] {
	return NewImporter(op).Import(ctx, paths...)
}
