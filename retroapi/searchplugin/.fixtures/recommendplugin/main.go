// Command recommendplugin is a search-plugin test fixture standing in for
// the "recommendations" subcommand - registry_test.go uses it to assert
// Registry.Recommend invokes that subcommand (not "search") and decodes its
// output the same way Search's echoplugin fixture does.
package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
)

func main() {
	// os.Args[0] is the program name; os.Args[1] is the "recommendations"
	// subcommand a real kong-based plugin would consume before parsing its
	// own flags, so skip it here too.
	fs := flag.NewFlagSet("recommendations", flag.ExitOnError)
	mimetype := fs.String("mimetype", "", "")
	limit := fs.String("limit", "", "")
	fs.Parse(os.Args[2:])

	json.NewEncoder(os.Stdout).Encode(&ddiscapi.Import{
		Uri:      "magnet:?xt=urn:btih:2222222222222222222222222222222222222222&dn=recommended",
		Mimetype: *mimetype,
		Title:    "limit=" + *limit,
	})
}
