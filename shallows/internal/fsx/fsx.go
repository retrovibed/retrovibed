package fsx

import (
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"log"
	"os"
	"path/filepath"

	"github.com/retrovibed/retrovibed/internal/debugx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
)

func ErrIsNotExist(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func IgnoreIsNotExist(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

func IgnoreIsExist(err error) error {
	if errors.Is(err, os.ErrExist) {
		return nil
	}

	return err
}

func AutoCached(path string, gen func() ([]byte, error)) (s []byte, err error) {
	if s, err = os.ReadFile(path); err == nil {
		return s, nil
	}

	if s, err = gen(); err != nil {
		return nil, err
	}

	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	if err = os.WriteFile(path, s, 0600); err != nil {
		return nil, err
	}

	return s, err
}

// IsRegularFile returns true IFF a non-directory file exists at the provided path.
func IsRegularFile(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

func IsSymlink(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return info.Mode().Type()&fs.ModeSymlink == fs.ModeSymlink
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return ErrIsNotExist(err) == nil
}

func RemoveSymlink(path string) error {
	info, err := os.Stat(path)
	if ErrIsNotExist(err) != nil {
		return nil
	} else if err != nil {
		return err
	}

	if info.Mode().Type()&fs.ModeSymlink != fs.ModeSymlink {
		return fmt.Errorf("refusing to remove non-symlink file: %s", path)
	}

	return os.Remove(path)
}

// resolves a symlink to its true path if an error occurs
// it'll log to debug and return the original path provided.
func ResolveSymlink(path string) string {
	p, err := os.Readlink(path)
	if err != nil {
		debugx.Println(errorsx.Wrapf(err, "unable to resolve symlink: %s", path))
		return path
	}

	return p
}

type Virtual interface {
	// returns the path rooted at the virtual fs from the fragments.
	Path(rel ...string) string
	MkDirAll(path string, perm os.FileMode) error
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
}

func VirtualAsFS(v Virtual) fs.FS {
	return vstoragefs{Virtual: v, pathrewrite: func(s string) string { return s }}
}

func DirVirtual(dir string) Virtual {
	return dirvirt{root: dir}
}

type dirvirt struct {
	root string
}

func (t dirvirt) Path(rel ...string) string {
	return filepath.Join(t.root, filepath.Join(rel...))
}

func (t dirvirt) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(filepath.Join(t.root, name), flag, perm)
}

func (t dirvirt) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, filepath.Join(t.root, newpath))
}

func (t dirvirt) MkDirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(filepath.Join(t.root, path), perm)
}

func VirtualAsFSWithRewrite(v Virtual, rewrite func(s string) string) fs.FS {
	return vstoragefs{Virtual: v, pathrewrite: rewrite}
}

type vstoragefs struct {
	Virtual
	pathrewrite func(s string) string
}

func (t vstoragefs) Open(name string) (fs.File, error) {
	debugx.Println("opening", name, "as", t.pathrewrite(name))
	return t.OpenFile(t.pathrewrite(name), os.O_RDONLY, 0600)
}

func MkDirs(perm fs.FileMode, paths ...string) (err error) {
	for _, p := range paths {
		if err = os.MkdirAll(p, perm); err != nil {
			return errorsx.Wrapf(err, "unable to create directory: %s", p)
		}

		if err = os.Chmod(p, perm); err != nil {
			return errorsx.Wrapf(err, "unable to set directory mod: %s", p)
		}
	}

	return nil
}

func Touch(perm fs.FileMode, paths ...string) (err error) {
	for _, p := range paths {
		touched, cause := os.OpenFile(p, os.O_CREATE|os.O_RDONLY, perm)
		if cause != nil {
			err = errorsx.Compact(err, errorsx.Wrapf(cause, "unable to create directory: %s", p))
			continue
		}

		err = errorsx.Compact(err, errorsx.Wrapf(touched.Close(), "unable to create directory: %s", p))
	}

	return err
}

func PrintFS(d fs.FS) {
	errorsx.Log(log.Output(2, fmt.Sprintln("--------- FS WALK INITIATED ---------")))
	defer func() { errorsx.Log(log.Output(3, fmt.Sprintln("--------- FS WALK COMPLETED ---------"))) }()

	err := fs.WalkDir(d, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info := errorsx.Zero(d.Info())
		errorsx.Log(log.Output(7, fmt.Sprintf("%v %4d %s\n", info.Mode(), info.Size(), path)))

		return nil
	})
	if err != nil {
		errorsx.Log(log.Output(2, fmt.Sprintln("fs walk failed", err)))
	}
}

func Walk(d fs.FS) *Walker {
	return &Walker{
		root: d,
	}
}

type Walker struct {
	root   fs.FS
	failed error
}

func (t *Walker) Walk() iter.Seq2[string, fs.DirEntry] {
	return func(yield func(string, fs.DirEntry) bool) {
		t.failed = fs.WalkDir(t.root, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !yield(path, d) {
				return fmt.Errorf("unable to yield path: %s", path)
			}

			if d.IsDir() && path != "." {
				return fs.SkipDir
			}

			return nil
		})
	}
}

func (t *Walker) Err() error {
	return t.failed
}

func AppendTo(path string, perm os.FileMode, block []byte, blocks ...[]byte) error {
	d, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer d.Close()

	for _, encoded := range blocks {
		block = append(block, encoded...)

	}

	if _, err = d.Write(block); err != nil {
		return err
	}

	return nil
}
