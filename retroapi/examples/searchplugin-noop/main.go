// Command searchplugin-noop is a minimal, worked example of a search
// plugin for retroapi/searchplugin's Registry. It performs no real search -
// it exists to show the protocol a plugin must satisfy, with nothing else
// competing for attention. Copy this package as the starting point for a
// real plugin (see ../../retrodscrape for one that actually queries real).
//
// build it for the registry with:
//
//	GOOS=wasip1 GOARCH=wasm go build -o noop.wasm ./retroapi/examples/searchplugin-noop
//
// and drop noop.wasm in the well-known search.d plugin directory (see
// searchplugin.SearchPluginDir) to have Registry load and run it. A build
// can also bake a default --source tag in at compile time via
// `-ldflags "-X main.source=mysite"` (see the source var below).
//
// Two kong subcommands are exposed: "plugin", the only one
// searchplugin.Registry ever invokes, and "recommendations", a second
// worked example showing the same *ddiscapi.Import output shape used for
// content a plugin recommends independent of any search query. Registry
// has no contract for "recommendations" today - it's here purely to
// demonstrate the shape a second plugin command could take.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/mimex"

	// autohijack points net.DefaultResolver and http.DefaultTransport at
	// wasinet's virtual sockets when this is built for wasip1 (a no-op on
	// every other platform). Registry.Search mounts the host's real TLS
	// trust store into the sandbox and sets SSL_CERT_DIR before running a
	// plugin, so a plain net/http.Get - TLS verification included - works
	// unmodified once this import is in place. Every plugin that talks to
	// the network needs it; without it a wasip1 build has no network at
	// all. This noop plugin never calls net/http, but the import is kept
	// here because it's the one piece of wiring every real plugin repeats.
	_ "github.com/egdaemon/wasinet/wasinet/autohijack"
)

// Registry mounts two per-plugin host directories into the wasm guest,
// following the systemd ConfigurationDirectory=/CacheDirectory= convention:
// a config directory at guest path /plugin/config.d, host-side
// searchplugin.PluginConfigDir(root, name) ({root}/search.d/{name}.config.d),
// and a cache directory at guest path /plugin/cache.d, host-side
// searchplugin.PluginCacheDir(root, name) ({root}/search.d/{name}.cache.d).
// A plugin reads these paths from the CONFIGURATION_DIRECTORY and
// CACHE_DIRECTORY env vars rather than hardcoding them - e.g. an API key or
// site-specific settings file would live under CONFIGURATION_DIRECTORY, and
// anything safe to lose (an HTTP response cache, a session cookie jar) under
// CACHE_DIRECTORY. This noop plugin needs neither and doesn't read them, but
// a real plugin typically does.

// source is this build's default provenance tag, baked in via
// `go build -ldflags "-X main.source=mysite"`. This is the general pattern
// for values a specific deployment of a plugin needs but that can't be a
// literal constant in shared source - a site identifier, a default
// category, an embedded API key. Blank (unconfigured) by default; a
// --source flag at runtime still overrides whatever was baked in, same as
// any other kong default.
var source = ""

// cli is parsed straight from argv[1:] by kong.Parse. Registry.Search
// always invokes a loaded plugin as exactly:
//
//	<binary> plugin --mimetype <m> [--mimetype <m> ...] --query <query> [--adult] [--public]
//
// (see retroapi/searchplugin/search.go's runSearchJob) - argv[1] is
// literally the subcommand name "plugin", which kong resolves the same way
// it resolves "recommendations" below. Unlike retrodscrape's real plugins
// (unit3d, leetx, piratebay), this example has no separate bare-argument
// dev mode to preserve alongside that contract, so there's no need to
// hand-parse os.Args before handing off to kong.
var cli struct {
	Plugin          pluginCmd          `cmd:"" help:"invoked by searchplugin.Registry to run a search"`
	Recommendations recommendationsCmd `cmd:"" help:"emit recommended content, independent of any search query"`
}

// pluginCmd's flags match Registry's invocation 1:1: repeatable --mimetype,
// required --query, and --adult/--public only ever appended when true
// (Registry never passes --adult=false/--public=false), so a plugin
// predating either flag still works for ordinary searches.
type pluginCmd struct {
	Mimetype []string `flag:"" name:"mimetype" help:"discovery mimetype to search within (repeatable)"`
	Query    string   `flag:"" name:"query" help:"search text to query" required:""`
	Source   string   `flag:"" name:"source" help:"provenance tag for every result (bakeable via -ldflags -X main.source=...)" default:"${source}"`
	Adult    bool     `flag:"" name:"adult" help:"allow adult content in results (Registry passes this whenever the caller allows it)"`
	Public   bool     `flag:"" name:"public" help:"only return results a public source could produce (Registry passes this when the caller wants public-only results)"`
}

// Run reports the mounted directories on stderr - proof, when reading
// plugin logs, that the mounts and env vars documented above actually
// reach the guest - then fabricates a single deterministic result, just to
// prove the rest of the contract end to end. A real plugin fetches results
// here instead, via net/http reaching the real network through the
// autohijack wiring above, emitting one *ddiscapi.Import per result found.
//
// This noop plugin stands in for a private-tracker source, so Public true
// means zero results - a plugin backed by an already-public source would
// instead just ignore the flag and search normally.
//
// Registry.Search reads stdout as newline-delimited JSON, one
// *ddiscapi.Import per line; anything written to stderr is logged, not
// parsed. A plugin emitting many results should reuse a single
// json.Encoder/writer across all of them rather than reopening one per
// result.
func (cmd *pluginCmd) Run(ctx context.Context) error {
	fmt.Fprintln(os.Stderr, "searchplugin-noop: CONFIGURATION_DIRECTORY =", os.Getenv("CONFIGURATION_DIRECTORY"))
	fmt.Fprintln(os.Stderr, "searchplugin-noop: CACHE_DIRECTORY =", os.Getenv("CACHE_DIRECTORY"))

	if cmd.Public {
		fmt.Fprintln(os.Stderr, "searchplugin-noop: standing in for a private-tracker source, --public requested - returning zero results")
		return nil
	}

	var mimetype string
	if len(cmd.Mimetype) > 0 {
		mimetype = cmd.Mimetype[0]
	}

	imp := &ddiscapi.Import{
		Uri:      "magnet:?xt=urn:btih:0000000000000000000000000000000000000000&dn=" + cmd.Query,
		Uritype:  mimex.Magnet,
		Health:   0,
		Mimetype: mimetype,
		Source:   cmd.Source,
	}

	if err := json.NewEncoder(os.Stdout).Encode(imp); err != nil {
		return fmt.Errorf("failed to encode result: %w", err)
	}

	return nil
}

// recommendationsCmd has no Registry contract behind it - it's a second
// worked example of the *ddiscapi.Import output shape, for content a
// plugin recommends independent of any search query. Unlike pluginCmd
// there's no --query: recommendations aren't search-driven by definition.
type recommendationsCmd struct {
	Mimetype []string `flag:"" name:"mimetype" help:"discovery mimetype to recommend content within (repeatable)"`
	Limit    uint     `flag:"" name:"limit" help:"number of recommended results to emit" default:"5"`
	Source   string   `flag:"" name:"source" help:"provenance tag for every result (bakeable via -ldflags -X main.source=...)" default:"${source}"`
	Public   bool     `flag:"" name:"public" help:"only return results a public source could produce"`
}

// licenseCycle is the fixed order recommendationsCmd.Run cycles fabricated
// entries' Licensed field through, so a worked example touches every value
// of ddiscapi.Import_LicenseStatus rather than always leaving it at its
// zero value the way pluginCmd's single fabricated result does.
var licenseCycle = [...]ddiscapi.Import_LicenseStatus{
	ddiscapi.Import_Unknown,
	ddiscapi.Import_Unlicensed,
	ddiscapi.Import_Licensed,
}

// Run fabricates Limit deterministic results, descending Popularity to
// look recommendation-like and cycling Licensed through every
// ddiscapi.Import_LicenseStatus value, and encodes them the same way
// pluginCmd does. Like pluginCmd, this noop stands in for a
// private-tracker source, so Public true means zero results.
func (cmd *recommendationsCmd) Run(ctx context.Context) error {
	if cmd.Public {
		fmt.Fprintln(os.Stderr, "searchplugin-noop: standing in for a private-tracker source, --public requested - returning zero results")
		return nil
	}

	var mimetype string
	if len(cmd.Mimetype) > 0 {
		mimetype = cmd.Mimetype[0]
	}

	enc := json.NewEncoder(os.Stdout)
	for i := uint(0); i < cmd.Limit; i++ {
		imp := &ddiscapi.Import{
			Uri:        fmt.Sprintf("magnet:?xt=urn:btih:%040d&dn=recommended-%d", i, i),
			Uritype:    mimex.Magnet,
			Health:     0,
			Mimetype:   mimetype,
			Source:     cmd.Source,
			Title:      fmt.Sprintf("recommended #%d", i),
			Popularity: float64(cmd.Limit - i),
			Licensed:   licenseCycle[i%uint(len(licenseCycle))],
		}

		if err := enc.Encode(imp); err != nil {
			return fmt.Errorf("failed to encode result: %w", err)
		}
	}

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	kctx := kong.Parse(&cli, kong.BindTo(ctx, (*context.Context)(nil)), kong.Vars{"source": source})
	kctx.FatalIfErrorf(kctx.Run())
}
