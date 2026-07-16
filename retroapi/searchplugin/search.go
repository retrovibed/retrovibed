package searchplugin

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"iter"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/langx"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"
)

// guestSSLCertDir is where each plugin invocation's TLS trust store is
// mounted inside the guest — matched by the SSL_CERT_DIR env var so the
// guest's crypto/x509 can find it.
const guestSSLCertDir = "/etc/ssl/certs"

// Search runs every loaded plugin as a WASI command with
// --mimetype mimetypes[0] [--mimetype mimetypes[1] ...] --query query,
// decoding each line of its stdout as a *ddiscapi.Import and yielding it.
// One plugin's failure is logged and skipped, not fatal to the whole
// sequence. ctx alone governs how long this may run — wrap it with
// context.WithTimeout for a deadline.
func (r *Registry) Search(ctx context.Context, mimetypes []string, query string) iterx.Seq[*ddiscapi.Import] {
	return &searchSeq{r: r, mimetypes: mimetypes, query: query}
}

type searchSeq struct {
	r         *Registry
	mimetypes []string
	query     string
	err       error
}

// workload is one plugin invocation dispatched onto the registry's shared
// asynccompute pool. results is unbuffered and shared by every job belonging
// to a single searchSeq.Each call; wg lets the dispatcher know when it is
// safe to close that channel.
type workload struct {
	path      string
	compiled  wazero.CompiledModule
	mimetypes []string
	query     string
	results   chan<- *ddiscapi.Import
	wg        *sync.WaitGroup
}

// runSearchJob instantiates a single plugin, decodes its stdout as jsonl,
// and streams the results onto j.results. Stdout is piped rather than
// buffered in full so a plugin emitting a large result set never holds its
// entire output resident in memory: the module runs on its own goroutine
// while this one scans and decodes concurrently off the pipe's reader end.
// A plugin failing (non-zero exit, bad jsonl) is logged and skipped, never
// fatal to the job. The only error returned is ctx cancellation while
// blocked handing a result to the dispatcher.
func (r *Registry) runSearchJob(ctx context.Context, j workload) error {
	defer j.wg.Done()

	stdoutr, stdoutw := io.Pipe()
	// argv[0] is the conventional program-name slot every CLI parser
	// (including kong.Parse, via os.Args[1:]) discards - "plugin" must come
	// after it to be seen as the subcommand.
	args := []string{j.path, "plugin", "--query", j.query}
	for _, m := range j.mimetypes {
		args = append(args, "--mimetype", m)
	}
	wazerofs := wazero.NewFSConfig().WithDirMount(r.sslCertDir, guestSSLCertDir)
	cfg := wazero.NewModuleConfig().
		WithName(j.path).
		WithArgs(args...).
		WithEnv("SSL_CERT_DIR", guestSSLCertDir).
		WithFSConfig(wazerofs).
		WithStdout(stdoutw).
		WithStderr(os.Stderr).
		WithSysWalltime().
		WithSysNanotime()

	envpath := strings.TrimSuffix(j.path, ".wasm") + ".env"
	envpairs, err := readEnvFile(envpath)
	if err != nil {
		log.Println("unable to read search plugin configuration", envpath, err)
	}
	for _, kv := range envpairs {
		k, v, _ := strings.Cut(kv, "=")
		cfg = cfg.WithEnv(k, v)
	}

	log.Println("running search plugin", j.path)

	go func() {
		mod, err := r.runtime.InstantiateModule(ctx, j.compiled, cfg)
		if mod != nil {
			errorsx.Log(mod.Close(ctx))
		}
		// propagates the module's outcome through the pipe itself: a
		// scan already in progress observes it as the read error once
		// buffered output is drained, instead of a side channel.
		stdoutw.CloseWithError(err)
	}()

	cause := scanResults(ctx, stdoutr, j)

	// unblocks the writer goroutine if the module is still running so it
	// can observe the failed write and exit instead of leaking; a no-op if
	// the module already finished and closed stdoutw itself.
	errorsx.Log(stdoutr.CloseWithError(cause))

	if exit, ok := errors.AsType[*sys.ExitError](cause); ok {
		if exit.ExitCode() != 0 {
			log.Println("search plugin failed", j.path, cause)
		}
		return nil
	}

	return cause
}

// scanResults decodes stdout as jsonl and streams each line onto j.results.
// A malformed line is logged and skipped. Returns non-nil only if ctx is
// done while blocked handing a result to the dispatcher.
func scanResults(ctx context.Context, stdout io.Reader, j workload) (failed error) {
	scanner := bufio.NewScanner(stdout)
	defer func() {
		failed = langx.FirstNonZero(failed, scanner.Err())
	}()

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		imp := &ddiscapi.Import{}
		if err := json.Unmarshal(line, imp); err != nil {
			log.Println("search plugin emitted invalid jsonl", j.path, err)
			continue
		}

		select {
		case j.results <- imp:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (t *searchSeq) Each(ctx context.Context) iter.Seq[*ddiscapi.Import] {
	return func(yield func(*ddiscapi.Import) bool) {
		cctx, cancel := context.WithCancel(ctx)
		defer cancel()

		results := make(chan *ddiscapi.Import)
		var wg sync.WaitGroup

		go func() {
			defer func() {
				wg.Wait()
				close(results)
			}()

			for path, compiled := range t.r.compiled() {
				select {
				case <-cctx.Done():
					t.err = cctx.Err()
					return
				default:
				}

				wg.Add(1)
				job := workload{
					path:      path,
					compiled:  compiled,
					mimetypes: t.mimetypes,
					query:     t.query,
					results:   results,
					wg:        &wg,
				}

				if err := t.r.pool.Run(cctx, job); err != nil {
					t.err = err
					wg.Done()
					return
				}
			}
		}()

		for imp := range results {
			if !yield(imp) {
				cancel()
				return
			}
		}
	}
}

func (t *searchSeq) Err() error {
	return t.err
}
