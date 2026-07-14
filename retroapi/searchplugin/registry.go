package searchplugin

import (
	"context"
	"maps"
	"net"
	"os"
	"sync"

	"github.com/egdaemon/wasinet/wasinet/wnetruntime"
	"github.com/egdaemon/wasinet/wazeronet"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/fsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/langx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Registry loads and holds compiled search-plugin wasm modules, keeping them
// resident so repeated searches never re-read or re-compile from disk — that
// only happens once per plugin, when the directory watcher reports it.
type Registry struct {
	runtime wazero.Runtime
	// sslCertDir is the host's real (symlink-resolved) TLS trust store
	// directory, mounted read-only into every plugin invocation so guest
	// TLS verification (system cert lookup via SSL_CERT_DIR) works —
	// wasip1 has no trust store of its own.
	sslCertDir string

	mu      sync.RWMutex
	modules map[string]wazero.CompiledModule

	// pool runs plugin invocations for every Search call; it is shared and
	// long-lived so concurrent searches reuse the same worker goroutines
	// instead of spinning up new ones per call.
	pool *asynccompute.Pool[workload]
}

// NewRegistry builds a Registry using wasinet's default Virtual socket
// (public addresses only — plugins never open real kernel sockets, every
// dial/lookup goes through the host's *net.Dialer/*net.Resolver instead).
// See NewRegistryWithSocket to substitute a different wnetruntime.Socket.
func NewRegistry(ctx context.Context) (*Registry, error) {
	return NewRegistryWithSocket(ctx, defaultSocket())
}

// NewRegistryWithSocket builds a Registry using sock for wasinet's guest
// networking instead of the default Virtual/PublicFirewall socket — e.g. to
// route plugin traffic through a different Dialer/Resolver (wireguard,
// dnscache, a test double) — and starts watching the well-known search.d
// plugin directory (${vars_user_configuration_directory}/search.d) for
// changes. There is no other configuration surface.
func NewRegistryWithSocket(ctx context.Context, sock wnetruntime.Socket) (*Registry, error) {
	r, err := newRegistry(ctx, sock)
	if err != nil {
		return nil, err
	}

	if err := watch(ctx, r, searchPluginDir()); err != nil {
		return nil, err
	}

	return r, nil
}

func defaultSocket() wnetruntime.Socket {
	return wnetruntime.Virtual(&net.Dialer{}, &net.ListenConfig{}, net.DefaultResolver, wnetruntime.PublicFirewall())
}

// newRegistry builds the wazero runtime + WASI + wasinet wiring for sock
// without starting the search.d directory watch, so tests can Load plugins
// directly without touching the real, hardcoded plugin directory.
func newRegistry(ctx context.Context, sock wnetruntime.Socket) (*Registry, error) {
	runtime := wazero.NewRuntime(ctx)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, runtime); err != nil {
		return nil, errorsx.Wrap(err, "unable to instantiate wasi")
	}

	if _, err := wazeronet.Module(runtime, sock).Instantiate(ctx); err != nil {
		return nil, errorsx.Wrap(err, "unable to instantiate wasinet")
	}

	r := &Registry{
		runtime: runtime,
		sslCertDir: langx.FirstNonZero(
			fsx.LocatePhysicalPath(
				"/etc/ssl/certs",
				"/etc/pki/tls/certs",
				"/usr/share/ca-certificates",
			),
			"/etc/ssl/certs",
		),
		modules: map[string]wazero.CompiledModule{},
	}
	r.pool = asynccompute.New(r.runSearchJob)

	return r, nil
}

func searchPluginDir() string {
	return userx.DefaultConfigDir(userx.DefaultRelRoot(), "search.d")
}

// Load compiles path (a .wasm file) and adds it to the registry. Compiling
// is the only disk IO this package performs after startup, and happens
// exactly once per plugin file.
func (r *Registry) Load(ctx context.Context, path string) error {
	bin, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	compiled, err := r.runtime.CompileModule(ctx, bin)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.modules[path] = compiled

	return nil
}

// Unload removes a previously loaded plugin.
func (r *Registry) Unload(path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.modules, path)
}

func (r *Registry) compiled() map[string]wazero.CompiledModule {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[string]wazero.CompiledModule, len(r.modules))
	maps.Copy(out, r.modules)
	return out
}
