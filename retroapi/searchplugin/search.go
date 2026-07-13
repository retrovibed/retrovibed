package searchplugin

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"iter"
	"log"

	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"
)

// guestSSLCertDir is where each plugin invocation's TLS trust store is
// mounted inside the guest — matched by the SSL_CERT_DIR env var so the
// guest's crypto/x509 can find it.
const guestSSLCertDir = "/etc/ssl/certs"

// Search runs every loaded plugin as a WASI command with
// --category category --query query, decoding each line of its stdout as a
// *ddiscapi.Import and yielding it. One plugin's failure is logged and
// skipped, not fatal to the whole sequence. ctx alone governs how long this
// may run — wrap it with context.WithTimeout for a deadline.
func (r *Registry) Search(ctx context.Context, category, query string) iterx.Seq[*ddiscapi.Import] {
	return &searchSeq{r: r, category: category, query: query}
}

type searchSeq struct {
	r        *Registry
	category string
	query    string
	err      error
}

func (t *searchSeq) Each(ctx context.Context) iter.Seq[*ddiscapi.Import] {
	return func(yield func(*ddiscapi.Import) bool) {
		for path, compiled := range t.r.compiled() {
			select {
			case <-ctx.Done():
				t.err = ctx.Err()
				return
			default:
			}

			var buf bytes.Buffer
			wazerofs := wazero.NewFSConfig().WithDirMount(t.r.sslCertDir, guestSSLCertDir)
			cfg := wazero.NewModuleConfig().
				WithName(path).
				WithArgs("plugin", "--category", t.category, "--query", t.query).
				WithEnv("SSL_CERT_DIR", guestSSLCertDir).
				WithFSConfig(wazerofs).
				WithStdout(&buf)

			mod, err := t.r.runtime.InstantiateModule(ctx, compiled, cfg)
			if mod != nil {
				errorsx.Log(mod.Close(ctx))
			}

			if err != nil {
				var exit *sys.ExitError
				if !errors.As(err, &exit) || exit.ExitCode() != 0 {
					log.Println("search plugin failed", path, err)
					continue
				}
			}

			scanner := bufio.NewScanner(&buf)
			for scanner.Scan() {
				line := scanner.Bytes()
				if len(line) == 0 {
					continue
				}

				imp := &ddiscapi.Import{}
				if err := json.Unmarshal(line, imp); err != nil {
					log.Println("search plugin emitted invalid jsonl", path, err)
					continue
				}

				if !yield(imp) {
					return
				}
			}
		}
	}
}

func (t *searchSeq) Err() error {
	return t.err
}
