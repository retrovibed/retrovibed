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

// environment is what the "env" subcommand reports: the variables this
// plugin understands, in the .env-with-comments form retroapi/envfile
// parses. The echo fixture understands none, but still declares one so
// tests have something with a hint to assert against.
const environment = `# echoed back verbatim as the published status
ECHO_STATUS="published" # status reported for every echo publish
`

func main() {
	// os.Args[0] is the program name; os.Args[1] is the subcommand a real
	// kong-based plugin would consume before parsing its own flags, so
	// dispatch on it here too.
	if len(os.Args) > 1 && os.Args[1] == "env" {
		os.Stdout.WriteString(environment)
		return
	}

	fs := flag.NewFlagSet("publish", flag.ExitOnError)
	title := fs.String("title", "", "")
	fs.String("description", "", "")
	fs.String("mimetype", "", "")
	fs.String("media", "", "")
	fs.String("community-id", "", "")
	link := fs.String("link", "", "")
	fs.Parse(os.Args[2:])

	// the link is echoed back as the external id when the caller supplied
	// one, so a test can prove Request.Link actually reaches the guest;
	// with no link this stays "echo" and the flag is invisible.
	external := "echo"
	if *link != "" {
		external = *link
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(result{
		URL:        "https://example.invalid/echo/" + *title,
		ExternalID: external,
		Status:     "published",
	})
}
