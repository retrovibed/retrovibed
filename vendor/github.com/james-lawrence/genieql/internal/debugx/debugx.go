package debugx

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/james-lawrence/genieql/internal/contextx"
	"github.com/james-lawrence/genieql/internal/errorsx"
	"github.com/james-lawrence/genieql/internal/iox"
	"github.com/james-lawrence/genieql/internal/stringsx"
	"github.com/pkg/errors"
)

var (
	defaults = log.New(io.Discard, "DEBUG", log.LstdFlags|log.Lshortfile)
)

// Output debug output
func Output(d int, s string) error {
	return defaults.Output(d, s)
}

// Println debug output
func Println(args ...interface{}) {
	Output(2, fmt.Sprintln(args...))
}

// Printf debug output
func Printf(format string, args ...interface{}) {
	Output(2, fmt.Sprintf(format, args...))
}

func genDst() (path string, dst io.WriteCloser) {
	var (
		err error
	)

	t := time.Now()

	ts := stringsx.Reverse(strconv.Itoa(int(t.Unix())))
	path = filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s-%s.trace", filepath.Base(os.Args[0]), t.Format("2006-01-02"), ts))

	if dst, err = os.Create(path); err != nil {
		log.Println(errors.Wrapf(err, "failed to open file: %s", path))
		log.Println("routine dump falling back to stderr")
		return "stderr", iox.WriteNopCloser(os.Stderr)
	}

	return path, dst
}

// DumpRoutines writes current goroutine stack traces to a temp file.
// and returns that files path. if for some reason a file could not be opened
// it falls back to stderr
func DumpRoutines() (path string, err error) {
	var (
		dst io.WriteCloser
	)

	path, dst = genDst()
	return path, errorsx.Compact(pprof.Lookup("goroutine").WriteTo(dst, 1), dst.Close())
}

// DumpOnSignal runs the DumpRoutes method and prints to stderr whenever one of the provided
// signals is received.
func DumpOnSignal(ctx context.Context, sigs ...os.Signal) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, sigs...)

	for {
		select {
		case <-ctx.Done():
			return
		case <-signals:
			if path, err := DumpRoutines(); err == nil {
				log.Println("dump located at:", path)
			} else {
				log.Println("failed to dump routines:", err)
			}
		}
	}
}

func OnSignal(ctx context.Context, do func(ctx context.Context) error, sigs ...os.Signal) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, sigs...)

	for {
		select {
		case <-ctx.Done():
			return
		case <-signals:
			if err := do(ctx); err != nil {
				log.Println("on signal failed", err)
			}
		}
	}
}

// func CPU(dir string) func(context.Context) (err error) {
// 	return func(ctx context.Context) (err error) {
// 		return run(ctx, dir, profile.CPUProfile)
// 	}
// }

// func Memory(dir string) func(context.Context) (err error) {
// 	return func(ctx context.Context) (err error) {
// 		return run(ctx, dir, profile.MemProfile)
// 	}
// }

// func Heap(dir string) func(context.Context) (err error) {
// 	return func(ctx context.Context) (err error) {
// 		return run(ctx, dir, profile.MemProfileHeap)
// 	}
// }

// func run(ctx context.Context, dir string, strategy func(*profile.Profile)) (err error) {
// 	if err = os.MkdirAll(dir, 0700); err != nil {
// 		return errors.Wrap(err, "unable to create profiling directory")
// 	}

// 	tmpdir, err := os.MkdirTemp(dir, strings.ReplaceAll("{}.*.profile", "{}", uuid.Must(uuid.NewV7()).String()))
// 	if err != nil {
// 		return errors.Wrap(err, "unable to create profiling directory")
// 	}
// 	defer os.RemoveAll(tmpdir)

// 	p := profile.Start(
// 		strategy,
// 		profile.NoShutdownHook,
// 		profile.ProfilePath(tmpdir),
// 	)

// 	stoppable := StopFunc(func() {
// 		p.Stop()
// 		errorsx.MaybeLog(errors.Wrap(clone(path.Join(dir, "profile.pprof"), tmpdir), "unable to finalize profile"))
// 	})

// 	return errors.WithStack(Run(ctx, stoppable))
// }

type Stoppable interface {
	Stop()
}

func Run(ctx context.Context, p Stoppable) error {
	defer p.Stop()
	<-ctx.Done()
	return contextx.IgnoreDeadlineExceeded(ctx.Err())
}

type StopFunc func()

func (t StopFunc) Stop() {
	t()
}

func Noop() Stoppable {
	return StopFunc(func() {})
}

func clone(dstpath string, dir string) (err error) {
	var (
		dst, src *os.File
	)

	location := locateFirstInDir(
		dir,
		"cpu.pprof",
		"mem.pprof",
		"mutex.pprof",
		"block.pprof",
		"threadcreation.pprof",
	)

	if dst, err = os.Create(dstpath); err != nil {
		return errors.Wrap(err, "copy failed")
	}
	defer dst.Close()

	if src, err = os.Open(location); err != nil {
		return errors.Wrap(err, "copy failed")
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return errors.Wrap(err, "copy failed")
	}

	return nil
}

// LocateFirstInDir locates the first file in the given directory by name.
func locateFirstInDir(dir string, names ...string) (result string) {
	for _, name := range names {
		result = filepath.Join(dir, name)
		if _, err := os.Stat(result); err == nil {
			break
		}
	}

	return result
}
