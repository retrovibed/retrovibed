// Command publishplugin-noop is a minimal, worked example of a publish
// plugin for retroapi/publishplugin's Registry. It performs no real
// publish - it exists to show the protocol a plugin must satisfy, with
// nothing else competing for attention. Copy this package as the starting
// point for a real plugin (e.g. one that posts to Mastodon, Discord, or
// another platform's API).
//
// build it for the registry with:
//
//	GOOS=wasip1 GOARCH=wasm go build -o noop.wasm ./retroapi/examples/publishplugin-noop
//
// and drop noop.wasm in the well-known publish.d plugin directory (see
// publishplugin.PublishPluginDir) to have Registry load and run it.
//
// Two kong subcommands are exposed. "publish", invoked by
// publishplugin.Registry.Publish as:
//
//	<binary> publish --title <t> --description <d> --mimetype <m> [--media <path>] [--community-id <id>] [--link <uri>]
//
// and "env", invoked by publishplugin.Registry.Environment as:
//
//	<binary> env
//
// which prints the variables this plugin understands as a .env document,
// one comment per variable describing it. That output is the plugin's
// configuration schema: it is what the console renders a settings form
// from, so a plugin nobody wrote a UI for still gets one.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"

	// autohijack points net.DefaultResolver and http.DefaultTransport at
	// wasinet's virtual sockets when this is built for wasip1 (a no-op on
	// every other platform). Registry.Publish mounts the host's real TLS
	// trust store into the sandbox and sets SSL_CERT_DIR before running a
	// plugin, so a plain net/http.Post - TLS verification included - works
	// unmodified once this import is in place. Every plugin that talks to
	// the network needs it; without it a wasip1 build has no network at
	// all. This noop plugin never calls net/http, but the import is kept
	// here because it's the one piece of wiring every real plugin repeats.
	_ "github.com/egdaemon/wasinet/wasinet/autohijack"
)

// cli is parsed straight from argv[1:] by kong.Parse. Registry.Publish
// always invokes a loaded plugin as exactly:
//
//	<binary> publish --title <t> --description <d> --mimetype <m> [--media <mounted-path>] [--community-id <id>] [--link <uri>]
//
// (see retroapi/publishplugin/publish.go's invoke) - argv[1] is literally
// the subcommand name ("publish"), which kong resolves from the field name
// below.
var cli struct {
	Publish publishCmd `cmd:"" help:"invoked by publishplugin.Registry to publish content"`
	Env     envCmd     `cmd:"" help:"invoked by publishplugin.Registry to report the variables this plugin understands"`
}

// publishCmd's flags match Registry's invocation 1:1: Media is only ever
// passed when the caller supplied Request.MediaPath, so a plugin that
// doesn't need the file can simply not read it, and Link only when the
// content has a publicly reachable URI (in practice its magnet link).
type publishCmd struct {
	Title       string `flag:"" name:"title" help:"title of the content being published"`
	Description string `flag:"" name:"description" help:"description of the content being published"`
	Mimetype    string `flag:"" name:"mimetype" help:"mimetype of the content being published"`
	Media       string `flag:"" name:"media" help:"guest path to the mounted media file, when the caller provided one"`
	CommunityID string `flag:"" name:"community-id" help:"id of the community the content is being published on behalf of"`
	Link        string `flag:"" name:"link" help:"publicly reachable uri for the content, when it has one"`
}

// envCmd takes no flags - it is a pure declaration of what this plugin can
// be configured with.
type envCmd struct{}

// environment is this plugin's configuration schema, in the format
// retroapi/envfile parses: a comment block immediately preceding a
// KEY=VALUE line (or a trailing "# ..." on the line itself, which wins)
// becomes that variable's help text. Values given here are defaults shown
// to whoever configures the plugin, not secrets - a real plugin declares
// its API token here with an empty value and lets the operator fill it in.
const environment = `# the noop plugin has nothing to configure; a real plugin
# declares one commented KEY=VALUE line per variable it reads.
NOOP_EXAMPLE="" # placeholder, ignored by this plugin
`

// Run prints the declaration verbatim on stdout. Registry.Environment
// returns these bytes untouched, so anything else written here would end up
// in the configuration form.
func (cmd *envCmd) Run(ctx context.Context) error {
	_, err := fmt.Fprint(os.Stdout, environment)
	return err
}

// Run reports the mounted directories and the media path on stderr - proof,
// when reading plugin logs, that the mounts and env vars documented above
// actually reach the guest, and that the media file (if any) is readable -
// then emits a single deterministic Result, just to prove the rest of the
// contract end to end. A real plugin performs the actual platform API call
// here instead, via net/http reaching the real network through the
// autohijack wiring above.
//
// Registry.Publish reads stdout as a single JSON object; anything written
// to stderr is logged, not parsed.
func (cmd *publishCmd) Run(ctx context.Context) error {
	fmt.Fprintln(os.Stderr, "publishplugin-noop: CONFIGURATION_DIRECTORY =", os.Getenv("CONFIGURATION_DIRECTORY"))
	fmt.Fprintln(os.Stderr, "publishplugin-noop: CACHE_DIRECTORY =", os.Getenv("CACHE_DIRECTORY"))

	fmt.Fprintln(os.Stderr, "publishplugin-noop: link =", cmd.Link)

	if cmd.Media != "" {
		info, err := os.Stat(cmd.Media)
		if err != nil {
			return fmt.Errorf("unable to stat mounted media: %w", err)
		}
		fmt.Fprintln(os.Stderr, "publishplugin-noop: media size =", info.Size())
	}

	result := struct {
		URL        string `json:"url"`
		ExternalID string `json:"external_id"`
		Status     string `json:"status"`
	}{
		URL:        "https://example.invalid/noop/" + cmd.CommunityID,
		ExternalID: "noop",
		Status:     "published",
	}

	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		return fmt.Errorf("failed to encode result: %w", err)
	}

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	kctx := kong.Parse(&cli, kong.BindTo(ctx, (*context.Context)(nil)))
	kctx.FatalIfErrorf(kctx.Run())
}
