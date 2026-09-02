package main

import (
	"encoding/json"
	"flag"
	"os"
)

type result struct {
	URL        string `json:"url"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
}

func main() {
	// os.Args[0] is the program name; os.Args[1] is the "publish"
	// subcommand a real kong-based plugin would consume before parsing
	// its own flags, so skip it here too.
	fs := flag.NewFlagSet("publish", flag.ExitOnError)
	title := fs.String("title", "", "")
	fs.String("description", "", "")
	fs.String("mimetype", "", "")
	fs.String("media", "", "")
	fs.String("community-id", "", "")
	fs.Parse(os.Args[2:])

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(result{
		URL:        "https://example.invalid/echo/" + *title,
		ExternalID: "echo",
		Status:     "published",
	})
}
