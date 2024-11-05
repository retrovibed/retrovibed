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

	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/iox"
	"github.com/james-lawrence/deeppool/internal/x/stringsx"
	// "github.com/egdaemon/egmeta/internal/iox"
	// "github.com/egdaemon/egmeta/internal/stringsx"
)

func genDst() (path string, dst io.WriteCloser) {
	var (
		err error
	)

	t := time.Now()
	ts := stringsx.Reverse(strconv.Itoa(int(t.Unix())))
	path = filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s-%s.trace", filepath.Base(os.Args[0]), t.Format("2006-01-02"), ts))

	if dst, err = os.Create(path); err != nil {
		log.Println(errorsx.Wrapf(err, "failed to open file: %s", path))
		log.Println("routine dump falling back to stderr")
		return "stderr", iox.WriteNopCloser(os.Stderr)
	}

	return path, dst
}

// DumpRoutines writes current goroutine stack traces to a temp file.
// and returns that files path. if for some reason a file could not be opened
// it falls back to stderr
func DumpRoutines(dst io.WriteCloser) (err error) {
	return errorsx.Compact(pprof.Lookup("goroutine").WriteTo(dst, 1), dst.Close())
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
			path, dst := genDst()
			if err := DumpRoutines(dst); err == nil {
				log.Println("dump located at:", path)
			} else {
				log.Println("failed to dump routines:", err)
			}
		}
	}
}
