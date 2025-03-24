package fsx

import (
	"errors"
	"fmt"
	"io/fs"
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

func RemoveSymlink(path string) error {
	info, err := os.Stat(path)
	if IgnoreIsNotExist(err) != nil {
		return err
	}

	if info.Mode().Type()&fs.ModeSymlink != fs.ModeSymlink {
		return fmt.Errorf("unable to remove non-symlink file: %s", path)
	}

	return os.Remove(path)
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

func VirtualAsFSWithRewrite(v Virtual, rewrite func(s string) string) fs.FS {
	return vstoragefs{Virtual: v, pathrewrite: rewrite}
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

type vstoragefs struct {
	Virtual
	pathrewrite func(s string) string
}

func (t vstoragefs) Open(name string) (fs.File, error) {
	debugx.Println("opening", name, "as", t.pathrewrite(name))
	return t.Virtual.OpenFile(t.pathrewrite(name), os.O_RDONLY, 0600)
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
