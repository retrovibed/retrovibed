package publishplugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"
)

// guestSSLCertDir is where each plugin invocation's TLS trust store is
// mounted inside the guest - matched by the SSL_CERT_DIR env var so the
// guest's crypto/x509 can find it.
const guestSSLCertDir = "/etc/ssl/certs"

// guestPluginConfigDir and guestPluginCacheDir are where each plugin
// invocation's per-plugin config/cache directories are mounted inside the
// guest - matched by the CONFIGURATION_DIRECTORY/CACHE_DIRECTORY env vars
// (systemd's ConfigurationDirectory=/CacheDirectory= convention) so the
// plugin can find them without hardcoding the mount point.
const guestPluginConfigDir = "/plugin/config.d"
const guestPluginCacheDir = "/plugin/cache.d"

// guestMediaDir is where the host directory containing Request.MediaPath is
// mounted inside the guest for a Publish invocation - the media's basename
// under this directory is what --media points a plugin at.
const guestMediaDir = "/plugin/media.d"

// Request is what a caller hands Publish to invoke a single named plugin.
// MediaPath, when non-empty, is a host-side flat file (the caller is
// responsible for materializing the relevant byte range of whatever storage
// backs the content into one file, since a wasm guest has no way to
// interpret block-cache internals) mounted read-only into the guest.
// Link, when non-empty, is the publicly reachable URI for the content -
// in practice the magnet URI, the same link the RSS feed advertises - so a
// plugin's post can point back at what was published.
type Request struct {
	Title       string
	Description string
	Mimetype    string
	CommunityID string
	MediaPath   string
	Link        string
}

// Result is what a plugin reports back on stdout as a single JSON object
// after a successful publish.
type Result struct {
	URL        string `json:"url"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
}

// publishWorkload is one plugin invocation dispatched onto the registry's
// shared asynccompute pool. result is buffered (size 1) and owned solely by
// the Publish call that created it, so runPublishJob's send never blocks
// regardless of whether Publish is still waiting on it.
type publishWorkload struct {
	path     string
	compiled wazero.CompiledModule
	req      Request
	result   chan<- publishOutcome
}

type publishOutcome struct {
	result *Result
	err    error
}

// Publish invokes the plugin installed at path - a WASI command run as:
//
//	<path> publish --title <t> --description <d> --mimetype <m> [--media <mounted-path>] [--community-id <id>] [--link <uri>]
//
// - decoding a single JSON object from its stdout as the *Result. Unlike
// searchplugin.T.Search, which fans a query out to every loaded plugin, this
// invokes exactly the one plugin named by path, once - callers (see
// shallows/communityapi.SyncPendingToDeeppool) are expected to loop over
// whichever plugins a community has enabled and call Publish once per
// plugin. A non-zero plugin exit is returned as an error, not silently
// swallowed - the caller decides whether that's fatal to its own loop. ctx
// alone governs how long this may run - wrap it with context.WithTimeout
// for a deadline.
func (r *Registry) Publish(ctx context.Context, path string, req Request) (*Result, error) {
	compiled, ok := r.lookup(path)
	if !ok {
		return nil, errorsx.Wrapf(ErrNotLoaded, "path: %s", path)
	}

	results := make(chan publishOutcome, 1)
	job := publishWorkload{path: path, compiled: compiled, req: req, result: results}

	if err := r.pool.Run(ctx, job); err != nil {
		return nil, err
	}

	select {
	case out := <-results:
		return out.result, out.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Environment invokes the plugin installed at path as:
//
//	<path> env
//
// returning, verbatim, the .env document it writes to stdout - the
// variables that plugin understands, with a comment per variable describing
// it (see retroapi/envfile for the exact convention). It is what lets a
// configuration UI render a form for a plugin nobody wrote a form for:
// the plugin itself is the schema.
//
// Unlike Publish this bypasses the shared worker pool. The call is short,
// rare, and interactive - queueing it behind however many multi-minute
// uploads are in flight would make the configuration screen appear hung.
func (r *Registry) Environment(ctx context.Context, path string) ([]byte, error) {
	compiled, ok := r.lookup(path)
	if !ok {
		return nil, errorsx.Wrapf(ErrNotLoaded, "path: %s", path)
	}

	// deliberately anonymous, unlike the publish invocation below: an
	// environment read can land while a publish of the same plugin is
	// mid-flight, and wazero refuses to instantiate a second module under a
	// name that is already live.
	return r.invoke(ctx, path, compiled, invocation{args: []string{"env"}})
}

// runPublishJob instantiates a single plugin, decodes its stdout as a single
// JSON object, and sends the outcome on j.result. A non-zero exit is
// reported through the outcome, not swallowed.
func (r *Registry) runPublishJob(ctx context.Context, j publishWorkload) error {
	out, err := r.publish(ctx, j.path, j.compiled, j.req)
	j.result <- publishOutcome{result: out, err: err}
	return err
}

// invocation is one command dispatched at a plugin: argv (after the binary
// path itself), the wazero module name to run it under - empty runs it
// anonymously, which is what lets concurrent invocations of the same plugin
// coexist - and the host-side media file to mount, when the command takes
// one.
type invocation struct {
	name  string
	args  []string
	media string
}

// publish runs a plugin's publish command and decodes its stdout as the
// single JSON object documented on Result.
func (r *Registry) publish(ctx context.Context, path string, compiled wazero.CompiledModule, req Request) (*Result, error) {
	args := []string{"publish", "--title", req.Title, "--description", req.Description, "--mimetype", req.Mimetype}
	if req.CommunityID != "" {
		args = append(args, "--community-id", req.CommunityID)
	}
	if req.Link != "" {
		args = append(args, "--link", req.Link)
	}
	if req.MediaPath != "" {
		args = append(args, "--media", guestMediaDir+"/"+filepath.Base(req.MediaPath))
	}

	stdout, err := r.invoke(ctx, path, compiled, invocation{name: path, args: args, media: req.MediaPath})
	if err != nil {
		return nil, err
	}

	var result Result
	if err := json.Unmarshal(stdout, &result); err != nil {
		return nil, errorsx.Wrapf(err, "publish plugin emitted invalid json: %s", path)
	}

	return &result, nil
}

// invoke runs a single command against an already-compiled plugin and
// returns whatever it wrote to stdout. Every command a plugin exposes -
// publish and env alike - reaches the guest through here, so the sandbox
// they see is identical: the host TLS trust store, the per-plugin config
// and cache directories, and the .env sidecar's variables. A non-zero exit
// is returned as an error rather than swallowed.
func (r *Registry) invoke(ctx context.Context, path string, compiled wazero.CompiledModule, inv invocation) ([]byte, error) {
	id := strings.TrimSuffix(filepath.Base(path), ".wasm")

	args := append([]string{path}, inv.args...)

	hostConfigDir := r.PluginConfigDir(id)
	hostCacheDir := r.PluginCacheDir(id)
	if err := os.MkdirAll(hostConfigDir, 0700); err != nil {
		return nil, errorsx.Wrapf(err, "unable to create plugin config directory: %s", hostConfigDir)
	}
	if err := os.MkdirAll(hostCacheDir, 0700); err != nil {
		return nil, errorsx.Wrapf(err, "unable to create plugin cache directory: %s", hostCacheDir)
	}

	wazerofs := wazero.NewFSConfig().
		WithDirMount(r.sslCertDir, guestSSLCertDir).
		WithDirMount(hostConfigDir, guestPluginConfigDir).
		WithDirMount(hostCacheDir, guestPluginCacheDir)

	if inv.media != "" {
		wazerofs = wazerofs.WithDirMount(filepath.Dir(inv.media), guestMediaDir)
	}

	var stdout bytes.Buffer
	cfg := wazero.NewModuleConfig().
		WithName(inv.name).
		WithArgs(args...).
		WithEnv("SSL_CERT_DIR", guestSSLCertDir).
		WithEnv("CONFIGURATION_DIRECTORY", guestPluginConfigDir).
		WithEnv("CACHE_DIRECTORY", guestPluginCacheDir).
		WithFSConfig(wazerofs).
		WithStdout(&stdout).
		WithStderr(os.Stderr).
		WithSysWalltime().
		WithSysNanotime()

	envpath := EnvPath(path)
	envpairs, err := readEnvFile(envpath)
	if err != nil {
		log.Println("unable to read publish plugin configuration", envpath, err)
	}
	for _, kv := range envpairs {
		k, v, _ := strings.Cut(kv, "=")
		cfg = cfg.WithEnv(k, v)
	}

	log.Println("running publish plugin", path, strings.Join(inv.args, " "))

	mod, runErr := r.runtime.InstantiateModule(ctx, compiled, cfg)
	if mod != nil {
		errorsx.Log(mod.Close(ctx))
	}

	if exit, ok := errors.AsType[*sys.ExitError](runErr); ok {
		if exit.ExitCode() != 0 {
			return nil, errorsx.Wrapf(exit, "publish plugin exited non-zero: %s", path)
		}
	} else if runErr != nil {
		return nil, errorsx.Wrapf(runErr, "unable to run publish plugin: %s", path)
	}

	return stdout.Bytes(), nil
}
