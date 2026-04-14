package library

import (
	"context"
	"crypto/md5"
	"errors"
	"hash"
	"io"
	"io/fs"
	"iter"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/retrovibed/retrovibed/shallows/internal/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
)

const ErrIterationFailed = errorsx.String("failed to yield transferred media")

type Transfered struct {
	Path     string
	Mimetype *mimetype.MIME
	MD5      hash.Hash
	Offset   uint64
	Bytes    uint64
}

type ImportOp = func(ctx context.Context, root fs.StatFS, path string) (*Transfered, error)

func TransferedFromPath(root fs.StatFS, path string) (*Transfered, error) {
	src, err := root.Open(path)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	cmimetype, err := mimetype.DetectReader(src)
	if err != nil {
		return nil, err
	}

	return &Transfered{
		Path:     path,
		MD5:      md5.New(),
		Mimetype: cmimetype,
	}, nil
}

func ImportSymlinkFile(srcvfs fsx.Virtual, vfs fsx.Virtual) ImportOp {
	if err := os.MkdirAll(vfs.Path(), 0700); err != nil {
		log.Println(errorsx.Wrap(err, "unable to ensure directory"))
	}

	return func(ctx context.Context, root fs.StatFS, path string) (*Transfered, error) {
		tx, err := TransferedFromPath(root, path)
		if err != nil {
			return nil, err
		}

		src, err := root.Open(path)
		if err != nil {
			return nil, err
		}
		defer src.Close()

		if n, err := io.Copy(tx.MD5, src); err != nil {
			return nil, err
		} else {
			tx.Bytes = uint64(n)
		}

		uid := md5x.FormatUUID(tx.MD5)

		if err := os.Remove(vfs.Path(uid)); fsx.IgnoreIsNotExist(err) != nil {
			return nil, errorsx.Wrap(err, "unable to ensure symlink destination is available")
		}

		if err := os.Symlink(srcvfs.Path(tx.Path), vfs.Path(uid)); err != nil {
			return nil, errorsx.Wrapf(err, "unable to symlink to original location: %s -> %s", srcvfs.Path(tx.Path), vfs.Path(uid))
		}

		return tx, nil
	}
}

func ImportCopyFile(vfs fsx.Virtual) ImportOp {
	return func(ctx context.Context, root fs.StatFS, path string) (*Transfered, error) {
		tx, err := TransferedFromPath(root, path)
		if err != nil {
			return nil, err
		}

		src, err := root.Open(path)
		if err != nil {
			return nil, err
		}
		defer src.Close()

		dst, err := os.CreateTemp(vfs.Path(), "importing.*.bin")
		if err != nil {
			return nil, err
		}
		defer os.Remove(dst.Name()) //nolint: errcheck
		defer dst.Close()

		if n, err := io.Copy(io.MultiWriter(tx.MD5, dst), src); err != nil {
			return nil, err
		} else {
			tx.Bytes = uint64(n)
		}
		uid := md5x.FormatUUID(tx.MD5)

		if err := os.Remove(vfs.Path(uid)); fsx.IgnoreIsNotExist(err) != nil {
			return nil, errorsx.Wrap(err, "unable to ensure destination is available")
		}

		if err := os.Rename(dst.Name(), vfs.Path(uid)); err != nil {
			return nil, errorsx.Wrap(err, "unable to symlink to original location")
		}

		return tx, nil
	}
}

func ImportFileDryRun(ctx context.Context, root fs.StatFS, path string) (*Transfered, error) {
	return TransferedFromPath(root, path)
}

func NewImporter(op ImportOp, root fs.StatFS, options ...asynccompute.Option[string]) *importer {
	return &importer{
		op:          op,
		root:        root,
		computeopts: options,
	}
}

type importer struct {
	op          ImportOp
	root        fs.StatFS
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
			if info, err := t.root.Stat(path); err != nil {
				return err
			} else if info.IsDir() {
				return nil
			}

			tx, cause := t.op(ictx, t.root, path)
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
				err = errorsx.Compact(err, asynccompute.Shutdown(ctx, arena))
				close(results)
			}()

			for _, p := range paths {
				if info, cause := t.root.Stat(p); errors.Is(cause, fs.ErrNotExist) {
					err = errorsx.Wrapf(cause, "ignoring %s", p)
					return
				} else if cause != nil {
					err = errorsx.Wrapf(cause, "failed %s", p)
					return
				} else if !info.IsDir() {
					if cause := arena.Run(ctx, p); cause != nil {
						err = errorsx.Wrapf(cause, "failed %s", p)
						return
					}

					continue
				}
				subroot, cause := fs.Sub(t.root, p)
				if cause != nil {
					err = errorsx.Wrapf(cause, "failed %s", p)
					return
				}
				cause = fs.WalkDir(subroot, ".", func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}

					if path == "." {
						return nil
					}

					return arena.Run(ctx, filepath.Join(p, path))
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

func ImportFilesystem(ctx context.Context, op ImportOp, root fs.StatFS, paths ...string) iter.Seq2[*Transfered, error] {
	return NewImporter(op, root).Import(ctx, paths...)
}
