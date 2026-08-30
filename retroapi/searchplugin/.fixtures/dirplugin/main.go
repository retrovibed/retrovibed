// Command dirplugin is a search-plugin test fixture that writes a marker
// file under CACHE_DIRECTORY and reports the guest-side path it wrote to via
// the Uri field, letting registry_test.go assert the host and guest are
// looking at the same directory through the wazero mount.
package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
)

func main() {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	mimetype := fs.String("mimetype", "", "")
	fs.String("query", "", "")
	fs.Bool("adult", false, "")
	fs.Parse(os.Args[2:])

	markerPath := filepath.Join(os.Getenv("CACHE_DIRECTORY"), "marker")
	os.WriteFile(markerPath, []byte("dirplugin"), 0600)

	json.NewEncoder(os.Stdout).Encode(&ddiscapi.Import{
		Uri:      markerPath,
		Mimetype: *mimetype,
	})
}
