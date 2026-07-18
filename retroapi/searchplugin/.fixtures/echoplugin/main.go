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

func main() {
	// os.Args[0] is the program name; os.Args[1] is the "plugin"
	// subcommand a real kong-based plugin would consume before parsing
	// its own flags, so skip it here too.
	fs := flag.NewFlagSet("plugin", flag.ExitOnError)
	mimetype := fs.String("mimetype", "", "")
	query := fs.String("query", "", "")
	fs.Bool("adult", false, "")
	fs.Parse(os.Args[2:])

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(result{
		Uri:      "magnet:?xt=urn:btih:1111111111111111111111111111111111111111&dn=" + *query,
		Health:   42,
		Mimetype: *mimetype,
	})
}
