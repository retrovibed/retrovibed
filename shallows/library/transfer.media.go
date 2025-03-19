package library

import (
	"context"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"iter"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/retrovibed/retrovibed/internal/x/fsx"
)

type Transfered struct {
	Path     string
	Mimetype *mimetype.MIME
	MD5      hash.Hash
	Bytes    uint64
}

// used to import files from a given directory tree into the library.
// it'll walk the tree, create a copy of each file into the media based on the contents md5.
func ImportDirectory(ctx context.Context, rootstore fsx.Virtual, subtree string) iter.Seq2[Transfered, error] {
	ErrIterationFailed := fmt.Errorf("failed to yield transferred media")
	fsi := os.DirFS(rootstore.Path(subtree))

	return func(yield func(Transfered, error) bool) {
		err := fs.WalkDir(fsi, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			cmimetype, err := mimetype.DetectFile(rootstore.Path(subtree, path))
			if err != nil {
				return err
			}

			var (
				tx = Transfered{
					Path:     rootstore.Path(subtree, path),
					MD5:      md5.New(),
					Mimetype: cmimetype,
				}
			)

			// log.Println("initiated copy", path, "to", tmp.Name())
			// defer log.Println("completed copy", path, "to", tmp.Name())

			src, err := os.Open(rootstore.Path(subtree, path))
			if err != nil {
				return err
			}
			defer src.Close()

			if n, err := io.Copy(tx.MD5, src); err != nil {
				// if n, err := io.Copy(io.MultiWriter(tmp, tx.MD5), src); err != nil {
				return err
			} else {
				tx.Bytes = uint64(n)
			}

			if !yield(tx, nil) {
				return ErrIterationFailed
			}

			return nil
		})

		if err != ErrIterationFailed && err != nil {
			yield(Transfered{}, err)
		}
	}
}
