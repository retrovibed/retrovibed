package main

import (
	"encoding/json"
	"flag"
	"os"
)

type result struct {
	Uri      string `json:"uri"`
	Health   uint32 `json:"health"`
	Mimetype string `json:"mimetype"`
}

// environment is what the "env" subcommand reports: the variables this
// plugin understands, in the .env-with-comments form retroapi/envfile parses.
const environment = `# appended to every result's magnet as its display name
PLUGIN_TOKEN="" # token echoed back in results
`

func main() {
	if len(os.Args) > 1 && os.Args[1] == "env" {
		os.Stdout.WriteString(environment)
		return
	}

	// os.Args[0] is the program name; os.Args[1] is the "search"
	// subcommand a real kong-based plugin would consume before parsing
	// its own flags, so skip it here too.
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	mimetype := fs.String("mimetype", "", "")
	fs.String("query", "", "")
	fs.Bool("adult", false, "")
	fs.Parse(os.Args[2:])

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(result{
		Uri:      "magnet:?xt=urn:btih:1111111111111111111111111111111111111111&dn=" + os.Getenv("PLUGIN_TOKEN"),
		Health:   42,
		Mimetype: *mimetype,
	})
}
