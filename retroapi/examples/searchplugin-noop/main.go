// Command searchplugin-noop is a minimal, worked example of a search
// plugin for retroapi/searchplugin's Registry. It performs no real search -
// it exists to show the protocol a plugin must satisfy, with nothing else
// competing for attention. Copy this package as the starting point for a
// real plugin (see ../../retrodscrape for one that actually queries real
// sites: 1337x.to and apibay.org).
//
// build it for the registry with:
//
//	GOOS=wasip1 GOARCH=wasm go build -o noop.wasm ./retroapi/examples/searchplugin-noop
//
// and drop noop.wasm in the well-known search.d plugin directory (see
// searchplugin.SearchPluginDir) to have Registry load and run it.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"

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

// mimetypeFlag collects every --mimetype occurrence into a slice, since a
// search request can carry several candidate discovery mimetypes (see
// ddisc.Category) and stdlib flag has no built-in repeatable-flag type.
type mimetypeFlag []string

func (m *mimetypeFlag) String() string { return strings.Join(*m, ",") }
func (m *mimetypeFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

func main() {
	// Registry.Search invokes every plugin as exactly:
	//   <binary> plugin --mimetype <m> [--mimetype <m> ...] --query <query>
	// argv[0] is the conventional program-name slot, discarded like any
	// CLI; argv[1] is always the literal subcommand "plugin" (the registry
	// has no other entrypoint), so flag parsing starts at argv[2:] rather
	// than argv[1:] - a plugin that flag.Parse(os.Args[1:]) instead will
	// see "plugin" as its first non-flag argument and stop parsing before
	// it ever reaches --mimetype/--query.
	if len(os.Args) < 2 || os.Args[1] != "plugin" {
		fmt.Fprintln(os.Stderr, "searchplugin-noop: expected a \"plugin\" subcommand")
		os.Exit(1)
	}

	fs := flag.NewFlagSet("plugin", flag.ExitOnError)
	var mimetypes mimetypeFlag
	fs.Var(&mimetypes, "mimetype", "discovery mimetype to search within (repeatable)")
	query := fs.String("query", "", "search text to query")
	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "searchplugin-noop:", err)
		os.Exit(1)
	}

	if *query == "" {
		fmt.Fprintln(os.Stderr, "searchplugin-noop: --query is required")
		os.Exit(1)
	}

	var mimetype string
	if len(mimetypes) > 0 {
		mimetype = mimetypes[0]
	}

	// A real plugin fetches results here - via net/http, reaching the real
	// network through the autohijack wiring above - and emits one
	// *ddiscapi.Import per result found, same as retrodscrape's leetx.go
	// and piratebay.go do. This noop plugin skips the request and
	// fabricates a single deterministic result instead, just to prove the
	// rest of the contract end to end.
	imp := &ddiscapi.Import{
		Magnet:   "magnet:?xt=urn:btih:0000000000000000000000000000000000000000&dn=" + *query,
		Health:   0,
		Mimetype: mimetype,
	}

	// Registry.Search reads stdout as newline-delimited JSON, one
	// *ddiscapi.Import per line; anything written to stderr is logged, not
	// parsed. A plugin emitting many results should reuse a single
	// json.Encoder/writer across all of them rather than reopening one per
	// result.
	if err := json.NewEncoder(os.Stdout).Encode(imp); err != nil {
		fmt.Fprintln(os.Stderr, "searchplugin-noop: failed to encode result:", err)
		os.Exit(1)
	}
}
