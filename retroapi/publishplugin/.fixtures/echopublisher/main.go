package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/langx"
	"github.com/retrovibed/retrovibed/retroapi/internal/md5x"
	"github.com/retrovibed/retrovibed/retroapi/internal/stringsx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
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
	// stderr, never stdout: the registry decodes stdout as the single
	// result object, so anything else written there breaks the publish.
	fmt.Fprintln(os.Stderr, "echopublisher:", strings.Join(os.Args[1:], " "))

	// os.Args[0] is the program name; os.Args[1] is the subcommand a real
	// kong-based plugin would consume before parsing its own flags, so
	// dispatch on it here too.
	if len(os.Args) > 1 && os.Args[1] == "env" {
		os.Stdout.WriteString(environment)
		return
	}

	fs := flag.NewFlagSet("publish", flag.ExitOnError)
	title := fs.String("title", "", "")
	media := fs.String("media", "", "")
	description := fs.String("description", "", "")
	mimetype := fs.String("mimetype", "", "")
	cid := fs.String("community-id", "", "")
	magnet := fs.String("magnet", "", "")
	fs.Bool("adult", false, "")
	fs.Parse(os.Args[2:])

	contentpath := userx.DefaultRuntimeDirectory(langx.Zero(media))

	content, err := os.Open(contentpath)
	if err != nil {
		log.Fatalln(errorsx.Wrap(err, "content missing"))
	}
	defer content.Close()

	if stringsx.Blank(langx.FirstNonZero(*description, *mimetype, *cid, *magnet)) {
		log.Println("description", description)
		log.Println("mimetype", mimetype)
		log.Println("community id", cid)
		log.Println("magnet", magnet)
		log.Fatalln("missing cli arguments")
	}
	enc := json.NewEncoder(os.Stdout)

	// the content is digested as the external id when the caller supplied
	// one, so a test can prove the content actually reaches the guest;
	enc.Encode(result{
		URL:        "https://example.invalid/echo/" + *title,
		ExternalID: md5x.FormatUUID(md5x.IO(content)),
		Status:     "published",
	})
}
