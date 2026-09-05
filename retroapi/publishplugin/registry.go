// Package publishplugin loads and runs sandboxed WASM publisher plugins -
// the publishing-side counterpart to retroapi/searchplugin. Where a search
// plugin is fanned a query out to every loaded module, a publisher plugin is
// invoked by name, once, on behalf of whichever community enabled it (see
// shallows/community's plugin_publishers/community_publisher tables) - see
// publish.go for the invocation contract.
package publishplugin

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/egdaemon/wasinet/wasinet/wnetruntime"
	"github.com/egdaemon/wasinet/wazeronet"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/fsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/langx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// T is the interface *Registry satisfies for Publish - the seam every
// consumer of a publish-plugin registry should depend on instead of the
// concrete type, so tests can substitute a fake.
type T interface {
	Publish(ctx context.Context, path string, req Request) (*Result, error)
}

// E is the interface *Registry satisfies for Environment, kept separate
// from T because the two have disjoint consumers: the publishing daemon
// only ever publishes, while the configuration API and CLI only ever read a
// plugin's declared variables. Depending on the narrower of the two keeps a
// fake down to the one method its caller actually exercises.
type E interface {
	Environment(ctx context.Context, path string) ([]byte, error)
}

// Unimplemented is a safe default T and E: every call fails with
// errors.ErrUnsupported instead of silently succeeding, so callers that
// haven't wired up a real registry get a clear signal rather than a bare nil
// interface passed around.
type Unimplemented struct{}

func (Unimplemented) Publish(ctx context.Context, path string, req Request) (*Result, error) {
	return nil, errors.ErrUnsupported
}

func (Unimplemented) Environment(ctx context.Context, path string) ([]byte, error) {
	return nil, errors.ErrUnsupported
}

// ErrNotLoaded is returned by Publish when path has not been (or is no
// longer) loaded into the registry.
const ErrNotLoaded = errorsx.String("publish plugin not loaded")

// Registry loads and holds compiled publish-plugin wasm modules, keeping
// them resident so repeated publishes never re-read or re-compile from disk
// - that only happens once per plugin, when the directory watcher reports
// it.
type Registry struct {
	runtime wazero.Runtime
	// sslCertDir is the host's real (symlink-resolved) TLS trust store
	// directory, mounted read-only into every plugin invocation so guest
	// TLS verification (system cert lookup via SSL_CERT_DIR) works -
	// wasip1 has no trust store of its own.
	sslCertDir string

	// configDir is the config root (e.g. userx.DefaultConfigDir(...)) this
	// registry's plugins read per-plugin config from
	// (PluginConfigDir(configDir, name), under PublishPluginDir(configDir)).
	configDir string
	// cacheDir is the cache root (e.g. userx.DefaultCacheDirectory(...))
	// this registry's plugins write per-plugin cache state to
	// (PluginCacheDir(cacheDir, name), under PublishPluginDir(cacheDir)).
	cacheDir string

	mu      sync.RWMutex
	modules map[string]wazero.CompiledModule

	// pool runs plugin invocations for every Publish call; it is shared and
	// long-lived so concurrent publishes reuse the same worker goroutines
	// instead of spinning up new ones per call.
	pool *asynccompute.Pool[publishWorkload]
}

// Option customizes a Registry built by NewRegistry/NewRegistryWithSocket.
type Option func(*Registry)

// OptionConfigDir overrides the config root a Registry's plugins read
// per-plugin config from (default: userx.DefaultConfigDir(userx.DefaultRelRoot())).
// Also where the publish.d plugin directory itself is watched from.
func OptionConfigDir(dir string) Option {
	return func(r *Registry) { r.configDir = dir }
}

// OptionCacheDir overrides the cache root a Registry's plugins write
// per-plugin cache state to (default: userx.DefaultCacheDirectory(userx.DefaultRelRoot())).
func OptionCacheDir(dir string) Option {
	return func(r *Registry) { r.cacheDir = dir }
}

// NewRegistry builds a Registry using wasinet's default Virtual socket
// (public addresses only - plugins never open real kernel sockets, every
// dial/lookup goes through the host's *net.Dialer/*net.Resolver instead).
// See NewRegistryWithSocket to substitute a different wnetruntime.Socket.
func NewRegistry(ctx context.Context, options ...Option) (*Registry, error) {
	return NewRegistryWithSocket(ctx, defaultSocket(), options...)
}

// NewRegistryWithSocket builds a Registry using sock for wasinet's guest
// networking instead of the default Virtual/PublicFirewall socket - e.g. to
// route plugin traffic through a different Dialer/Resolver (wireguard,
// dnscache, a test double) - and starts watching the well-known publish.d
// plugin directory (PublishPluginDir(${vars_user_configuration_directory}))
// for changes. Pass OptionConfigDir/OptionCacheDir to point a Registry at
// different roots (e.g. a t.TempDir() in tests) instead of the real,
// hardcoded userx directories.
func NewRegistryWithSocket(ctx context.Context, sock wnetruntime.Socket, options ...Option) (*Registry, error) {
	r, err := newRegistry(ctx, sock, options...)
	if err != nil {
		return nil, err
	}

	if err := watch(ctx, r, PublishPluginDir(r.configDir)); err != nil {
		return nil, err
	}

	return r, nil
}

func defaultSocket() wnetruntime.Socket {
	return wnetruntime.Virtual(&net.Dialer{}, &net.ListenConfig{}, net.DefaultResolver, wnetruntime.PublicFirewall())
}

// newRegistry builds the wazero runtime + WASI + wasinet wiring for sock
// without starting the publish.d directory watch, so tests can Load plugins
// directly without touching the real, hardcoded plugin directory.
func newRegistry(ctx context.Context, sock wnetruntime.Socket, options ...Option) (*Registry, error) {
	cachedir := userx.DefaultCacheDirectory(userx.DefaultRelRoot())

	// Plugins are large (tens of MB) wasm binaries; without a persistent
	// compilation cache wazero AOT-compiles every one of them from scratch
	// on every daemon start, which can add tens of seconds to startup.
	// Caching compiled code here means only new/changed plugins pay that
	// cost.
	compiledir := filepath.Join(cachedir, "wazero.compiled")
	if err := os.MkdirAll(compiledir, 0700); err != nil {
		return nil, errorsx.Wrap(err, "unable to create wazero compilation cache directory")
	}

	cache, err := wazero.NewCompilationCacheWithDir(compiledir)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to open wazero compilation cache")
	}

	runtime := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCompilationCache(cache))

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
		configDir: userx.DefaultConfigDir(userx.DefaultRelRoot()),
		cacheDir:  cachedir,
		modules:   map[string]wazero.CompiledModule{},
	}
	for _, opt := range options {
		opt(r)
	}
	r.pool = asynccompute.New(r.runPublishJob)

	return r, nil
}

// PublishPluginDir is the well-known publish.d directory nested under root -
// the config-rooted form (root = userx.DefaultConfigDir(userx.DefaultRelRoot()))
// is where plugin .wasm/.env files live; the cache-rooted form
// (root = userx.DefaultCacheDirectory(userx.DefaultRelRoot())) is its
// cache-side analogue.
func PublishPluginDir(root string) string {
	return filepath.Join(root, "publish.d")
}

// EnvPath is the .env sidecar belonging to the plugin installed at wasm - a
// sibling file with the .wasm suffix swapped for .env. Every consumer of a
// plugin's configuration (the runner in publish.go, the HTTP environment
// service, the install/config CLI) resolves it through here so the rule
// cannot drift between them. Because it is derived from the path rather
// than the file's contents, each symlink to a shared module gets its own
// configuration.
func EnvPath(wasm string) string {
	return strings.TrimSuffix(wasm, ".wasm") + ".env"
}

// PluginConfigDir is the well-known per-plugin configuration directory,
// nested under root's publish.d: {root}/publish.d/{name}.config.d.
func PluginConfigDir(root string, name string) string {
	return filepath.Join(PublishPluginDir(root), name+".config.d")
}

// PluginCacheDir is the well-known per-plugin cache directory, nested under
// root's publish.d: {root}/publish.d/{name}.cache.d.
func PluginCacheDir(root string, name string) string {
	return filepath.Join(PublishPluginDir(root), name+".cache.d")
}

// PluginConfigDir returns this registry's per-plugin config directory for name.
func (r *Registry) PluginConfigDir(name string) string {
	return PluginConfigDir(r.configDir, name)
}

// PluginCacheDir returns this registry's per-plugin cache directory for name.
func (r *Registry) PluginCacheDir(name string) string {
	return PluginCacheDir(r.cacheDir, name)
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

// lookup returns the compiled module registered for path, if any.
func (r *Registry) lookup(path string) (wazero.CompiledModule, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.modules[path]
	return m, ok
}
