package searchplugin

import (
	"context"
	"os"
	"sync"

	"github.com/egdaemon/wasinet/wazeronet"
	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Registry loads and holds compiled search-plugin wasm modules, keeping them
// resident so repeated searches never re-read or re-compile from disk — that
// only happens once per plugin, when the directory watcher reports it.
type Registry struct {
	runtime wazero.Runtime

	mu      sync.RWMutex
	modules map[string]wazero.CompiledModule
}

// NewRegistry builds a Registry, wires up WASI + wasinet networking
// (public addresses only), and starts watching the well-known search.d
// plugin directory (${vars_user_configuration_directory}/search.d) for
// changes. There is no configuration surface for any of this.
func NewRegistry(ctx context.Context) (*Registry, error) {
	r, err := newRegistry(ctx)
	if err != nil {
		return nil, err
	}

	if err := watch(ctx, r, searchPluginDir()); err != nil {
		return nil, err
	}

	return r, nil
}

// newRegistry builds the wazero runtime + WASI + wasinet wiring without
// starting the search.d directory watch, so tests can Load plugins directly
// without touching the real, hardcoded plugin directory.
func newRegistry(ctx context.Context) (*Registry, error) {
	runtime := wazero.NewRuntime(ctx)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, runtime); err != nil {
		return nil, errorsx.Wrap(err, "unable to instantiate wasi")
	}

	if _, err := wazeronet.Module(runtime, newPublicOnlySocket()).Instantiate(ctx); err != nil {
		return nil, errorsx.Wrap(err, "unable to instantiate wasinet")
	}

	return &Registry{
		runtime: runtime,
		modules: map[string]wazero.CompiledModule{},
	}, nil
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
	for k, v := range r.modules {
		out[k] = v
	}
	return out
}
